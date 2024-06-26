package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/nats-io/nats.go"
	"github.com/simpleiot/simpleiot/data"
)

// InfluxMeasurement is the Influx measurement to which all points are written
const InfluxMeasurement = "points"

// Db represents the configuration for a SIOT DB client
type Db struct {
	ID            string   `node:"id"`
	Parent        string   `node:"parent"`
	Description   string   `point:"description"`
	URI           string   `point:"uri"`
	Org           string   `point:"org"`
	Bucket        string   `point:"bucket"`
	AuthToken     string   `point:"authToken"`
	TagPointTypes []string `point:"tagPointType"`
}

// DbClient is a SIOT database client
type DbClient struct {
	nc            *nats.Conn
	config        Db
	stop          chan struct{}
	newPoints     chan NewPoints
	newEdgePoints chan NewPoints
	newDbPoints   chan NewPoints
	upSub         *nats.Subscription
	upSubHr       *nats.Subscription
	historySub    *nats.Subscription
	nodeCache     nodeCache
	client        influxdb2.Client
	writeAPI      api.WriteAPI
}

// NewDbClient ...
func NewDbClient(nc *nats.Conn, config Db) Client {
	return &DbClient{
		nc:            nc,
		config:        config,
		stop:          make(chan struct{}),
		newPoints:     make(chan NewPoints),
		newEdgePoints: make(chan NewPoints),
		newDbPoints:   make(chan NewPoints),
		nodeCache:     newNodeCache(config.TagPointTypes),
	}
}

// Run runs the main logic for this client and blocks until stopped
func (dbc *DbClient) Run() error {
	log.Println("Starting db client:", dbc.config.Description)
	var err error

	// FIXME, we probably want to store edge points too ...

	subject := fmt.Sprintf("up.%v.*", dbc.config.Parent)
	dbc.upSub, err = dbc.nc.Subscribe(subject, func(msg *nats.Msg) {
		points, err := data.PbDecodePoints(msg.Data)
		if err != nil {
			log.Println("Error decoding points in db upSub:", err)
			return
		}

		// find node ID for points
		chunks := strings.Split(msg.Subject, ".")
		if len(chunks) != 3 {
			log.Println("rule client up sub, malformed subject:", msg.Subject)
			return
		}

		dbc.newDbPoints <- NewPoints{chunks[2], "", points}
	})

	if err != nil {
		return fmt.Errorf("subscribing to %v: %w", subject, err)
	}

	subjectHR := fmt.Sprintf("phrup.%v.*", dbc.config.Parent)
	dbc.upSubHr, err = dbc.nc.Subscribe(subjectHR, func(msg *nats.Msg) {
		// find node ID for points
		chunks := strings.Split(msg.Subject, ".")
		if len(chunks) != 3 {
			log.Println("rule client up hr sub, malformed subject:", msg.Subject)
			return
		}

		nodeID := chunks[2]

		// Update nodeCache with no points
		err := dbc.nodeCache.Update(dbc.nc, NewPoints{
			ID: nodeID,
		})
		if err != nil {
			log.Printf("error updating cache: %v", err)
		}

		err = data.DecodeSerialHrPayload(msg.Data, func(pt data.Point) {
			tags := map[string]string{
				"type": pt.Type,
				"key":  pt.Key,
			}
			dbc.nodeCache.CopyTags(nodeID, tags)
			p := influxdb2.NewPoint(InfluxMeasurement,
				tags,
				map[string]interface{}{
					"value": pt.Value,
				},
				pt.Time)
			dbc.writeAPI.WritePoint(p)
		})

		if err != nil {
			log.Println("DB: error decoding HR data:", err)
		}
	})

	if err != nil {
		return fmt.Errorf("subscribing to %v: %w", subjectHR, err)
	}

	subjectHistory := fmt.Sprintf("history.%v", dbc.config.ID)
	dbc.historySub, err = dbc.nc.Subscribe(subjectHistory, func(msg *nats.Msg) {
		query := new(data.HistoryQuery)
		results := new(data.HistoryResults)
		ctx := context.Background()

		// Defer encoding and sending response
		defer func() {
			res, err := json.Marshal(results)
			if err != nil {
				err = msg.Respond([]byte(`{"error":"error encoding response"}`))
				if err != nil {
					log.Printf("error responding to history query: %v", err)
				}
			} else {
				err = msg.Respond(res)
				if err != nil {
					log.Printf("error responding to history query: %v", err)
				}
			}
		}()

		// Parse query
		err = json.Unmarshal(msg.Data, query)
		if err != nil {
			results.ErrorMessage = "parsing query: " + err.Error()
			return
		}
		log.Printf("received history query: %+v", query)

		// Execute query
		query.Execute(
			ctx,
			dbc.client.QueryAPI(dbc.config.Org),
			dbc.config.Bucket,
			InfluxMeasurement,
			results,
		)
	})

	if err != nil {
		return fmt.Errorf("subscribing to %v: %w", subjectHistory, err)
	}

	setupAPI := func() {
		log.Println("Setting up Influx API")
		// you can set things like retries, batching, precision, etc in client options.
		dbc.client = influxdb2.NewClientWithOptions(dbc.config.URI,
			dbc.config.AuthToken, influxdb2.DefaultOptions())
		dbc.writeAPI = dbc.client.WriteAPI(dbc.config.Org, dbc.config.Bucket)

		influxErrors := dbc.writeAPI.Errors()

		go func() {
			for err := range influxErrors {
				if err != nil {
					log.Println("Influx write error:", err)
				}

			}
			log.Println("Influxdb write api closed")
		}()
	}

	setupAPI()

done:
	for {
		select {
		case <-dbc.stop:
			log.Println("Stopping db client:", dbc.config.Description)
			break done
		case pts := <-dbc.newPoints:
			err := data.MergePoints(pts.ID, pts.Points, &dbc.config)
			if err != nil {
				log.Println("error merging new points:", err)
			}

			for _, p := range pts.Points {
				switch p.Type {
				case data.PointTypeURI,
					data.PointTypeOrg,
					data.PointTypeBucket,
					data.PointTypeAuthToken:
					// we need to restart the influx write API
					dbc.client.Close()
					setupAPI()
				case data.PointTypeTagPointType:
					dbc.nodeCache = newNodeCache(dbc.config.TagPointTypes)
				}
			}

		case pts := <-dbc.newEdgePoints:
			err := data.MergeEdgePoints(pts.ID, pts.Parent, pts.Points, &dbc.config)
			if err != nil {
				log.Println("error merging new points:", err)
			}
		case pts := <-dbc.newDbPoints:
			// Update nodeCache if needed
			err := dbc.nodeCache.Update(dbc.nc, pts)
			if err != nil {
				log.Printf("error updating cache: %v", err)
			}
			// Add points to InfluxDB
			for _, point := range pts.Points {
				tags := map[string]string{
					"type": point.Type,
					"key":  point.Key,
				}
				dbc.nodeCache.CopyTags(pts.ID, tags)
				p := influxdb2.NewPoint(InfluxMeasurement,
					tags,
					map[string]interface{}{
						"value": point.Value,
						"text":  point.Text,
					},
					point.Time)
				dbc.writeAPI.WritePoint(p)
			}
		}
	}

	// clean up
	_ = dbc.upSub.Unsubscribe()
	_ = dbc.upSubHr.Unsubscribe()
	_ = dbc.historySub.Unsubscribe()
	dbc.client.Close()
	return nil
}

// Stop sends a signal to the Run function to exit
func (dbc *DbClient) Stop(_ error) {
	close(dbc.stop)
}

// Points is called by the Manager when new points for this
// node are received.
func (dbc *DbClient) Points(nodeID string, points []data.Point) {
	dbc.newPoints <- NewPoints{nodeID, "", points}
}

// EdgePoints is called by the Manager when new edge points for this
// node are received.
func (dbc *DbClient) EdgePoints(nodeID, parentID string, points []data.Point) {
	dbc.newEdgePoints <- NewPoints{nodeID, parentID, points}
}
