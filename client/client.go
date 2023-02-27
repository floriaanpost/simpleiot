package client

import (
	"github.com/nats-io/nats.go"
	"github.com/simpleiot/simpleiot/data"
)

// Client interface describes methods a Simple Iot client must implement.
// This is to be kept as simple as possible, and the ClientManager does all
// the heavy lifting of interacting with the rest of the SIOT system.
// Run should block until Stop is called.
// Start MUST return when Stop is called.
// Stop does not block -- wait until Run returns if you need to know the client
// is stopped.
// Points and EdgePoints are called when there are updates to the client node.
// The client Manager filters out all points with Origin set to "" because it
// assumes the point was generated by the client and the client already knows about it.
// Thus, if you want points to get to a client, Origin must be set.
type Client interface {
	RunStop

	Points(string, []data.Point)
	EdgePoints(string, string, []data.Point)
}

// DefaultClients returns an actor for the default group of built in clients
func DefaultClients(nc *nats.Conn) (*Group, error) {
	g := NewGroup("Default clients")

	sc := NewManager(nc, NewSerialDevClient)
	g.Add(sc)

	cb := NewManager(nc, NewCanBusClient)
	g.Add(cb)

	rc := NewManager(nc, NewRuleClient)
	g.Add(rc)

	db := NewManager(nc, NewDbClient)
	g.Add(db)

	sg := NewManager(nc, NewSignalGeneratorClient)
	g.Add(sg)

	sync := NewManager(nc, NewSyncClient)
	g.Add(sync)

	metrics := NewManager(nc, NewMetricsClient)
	g.Add(metrics)

	particle := NewManager(nc, NewParticleClient)
	g.Add(particle)

	return g, nil
}
