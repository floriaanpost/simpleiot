package modbus

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/simpleiot/simpleiot/test"
)

// Client defines a Modbus client (master)
type Client struct {
	transport Transport
	debug     int
}

// NewClient is used to create a new modbus client
// port must return an entire packet for each Read().
// github.com/simpleiot/simpleiot/respreader is a good
// way to do this.
func NewClient(transport Transport, debug int) *Client {
	return &Client{
		transport: transport,
		debug:     debug,
	}
}

// SetDebugLevel allows you to change debug level on the fly
func (c *Client) SetDebugLevel(debug int) {
	c.debug = debug
}

// Close closes the client transport
func (c *Client) Close() error {
	return c.transport.Close()
}

// ReadCoils is used to read modbus coils
func (c *Client) ReadCoils(id byte, coil, count uint16) ([]bool, error) {
	ret := []bool{}
	req := ReadCoils(coil, count)

	resp, err := c.executeRequest(id, req)
	if err != nil {
		return ret, err
	}

	return resp.RespReadBits()
}

// WriteSingleCoil is used to read modbus coils
func (c *Client) WriteSingleCoil(id byte, coil uint16, v bool) error {
	req := WriteSingleCoil(coil, v)

	resp, err := c.executeRequest(id, req)
	if err != nil {
		return err
	}

	if !bytes.Equal(req.Data, resp.Data) {
		return errors.New("did not get the correct response data")
	}

	return nil
}

// ReadDiscreteInputs is used to read modbus discrete inputs
func (c *Client) ReadDiscreteInputs(id byte, input, count uint16) ([]bool, error) {
	ret := []bool{}
	req := ReadDiscreteInputs(input, count)

	resp, err := c.executeRequest(id, req)
	if err != nil {
		return ret, err
	}

	return resp.RespReadBits()
}

// ReadHoldingRegs is used to read modbus coils
func (c *Client) ReadHoldingRegs(id byte, reg, count uint16) ([]uint16, error) {
	ret := []uint16{}
	req := ReadHoldingRegs(reg, count)

	resp, err := c.executeRequest(id, req)
	if err != nil {
		return ret, err
	}

	return resp.RespReadRegs()
}

// ReadInputRegs is used to read modbus coils
func (c *Client) ReadInputRegs(id byte, reg, count uint16) ([]uint16, error) {
	ret := []uint16{}
	req := ReadInputRegs(reg, count)

	resp, err := c.executeRequest(id, req)
	if err != nil {
		return ret, err
	}

	return resp.RespReadRegs()
}

// WriteSingleReg writes to a single holding register
func (c *Client) WriteSingleReg(id byte, reg, value uint16) error {
	req := WriteSingleReg(reg, value)

	resp, err := c.executeRequest(id, req)
	if err != nil {
		return err
	}

	if !bytes.Equal(req.Data, resp.Data) {
		return errors.New("did not get the correct response data")
	}

	return nil
}

// executeRequests sends data to the server and parses the response. It also does
// basic error checking:
// - Did the server return an exception?
// - Did the server respond with the same function code?
func (c *Client) executeRequest(id byte, req PDU) (PDU, error) {
	if c.debug >= 1 {
		fmt.Printf("Modbus client %s ID:0x%x req:%v\n", req.FunctionCode, id, req)
	}

	packet, err := c.transport.Encode(id, req)
	if err != nil {
		return PDU{}, err
	}

	if c.debug >= 9 {
		fmt.Printf("Modbus client %s tx: %s\n", req.FunctionCode, test.HexDump(packet))
	}

	_, err = c.transport.Write(packet)
	if err != nil {
		return PDU{}, err
	}

	// FIXME, what is max modbus packet size?
	buf := make([]byte, 200)
	cnt, err := c.transport.Read(buf)
	if err != nil {
		return PDU{}, err
	}

	buf = buf[:cnt]

	if c.debug >= 9 {
		fmt.Printf("Modbus client %s rx: %s\n", req.FunctionCode, test.HexDump(packet))
	}

	_, resp, err := c.transport.Decode(buf)
	if err != nil {
		return PDU{}, err
	}

	if c.debug >= 1 {
		fmt.Printf("Modbus client %s ID:0x%x resp:%v\n", req.FunctionCode, id, resp)
	}

	funcCode := resp.FunctionCode
	hasException := resp.FunctionCode&0x80 != 0
	if hasException {
		funcCode = resp.FunctionCode & (0x80 ^ 0xFF) // reset the bit
	}

	// first check function code before checking exception
	if funcCode != req.FunctionCode {
		return PDU{}, errors.New("resp contains wrong function code")
	}

	if hasException {
		return resp, ExceptionCode(resp.Data[0])
	}

	return resp, err
}
