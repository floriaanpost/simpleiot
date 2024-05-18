package modbus

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const rtuFrameStartLen = 2

// RtuADU defines an ADU for RTU packets
type RtuADU struct {
	PDU
	Address byte
	CRC     uint16
}

// RTU defines an RTU connection
type RTU struct {
	port              io.ReadWriteCloser
	framebuf          *bytes.Buffer
	incomingFrameType frameType
}

// NewRTU creates a new RTU transport
func NewRTU(port io.ReadWriteCloser) *RTU {
	return &RTU{
		port:     port,
		framebuf: &bytes.Buffer{},
	}
}

func (r *RTU) Read(p []byte) (int, error) {
	// the layer above always wants full frames, so make sure to reset any data left in the buffer
	r.framebuf.Reset()

	n, err := io.CopyN(r.framebuf, r.port, rtuFrameStartLen)
	if err != nil {
		return int(n), err
	}

	// TODO: Implement timeout

	frame := r.framebuf.Bytes()
	funcCode := FunctionCode(frame[1])

	if funcCode&0x80 != 0 { // this is an exception
		n, err := io.CopyN(r.framebuf, r.port, 5-rtuFrameStartLen) // exception frames are 5 bytes long
		if err != nil {
			return int(n), err
		}

		return r.framebuf.Read(p)
	}

	// the frame length depends on if it is a request or a response
	switch r.incomingFrameType {
	case request:
		n, err := r.bufferRequestFrame(funcCode)
		if err != nil {
			return n, err
		}
	case response:
		n, err := r.bufferResponseFrame(funcCode)
		if err != nil {
			return n, err
		}
	default:
		// TODO: should not be possible, how to handle this?
	}

	return r.framebuf.Read(p)
}

func (r *RTU) Write(p []byte) (int, error) {
	return r.port.Write(p)
}

// Close closes the serial port
func (r *RTU) Close() error {
	return r.port.Close()
}

// Encode encodes a RTU packet
func (r *RTU) Encode(id byte, pdu PDU) ([]byte, error) {
	ret := make([]byte, len(pdu.Data)+2+2)
	ret[0] = id
	ret[1] = byte(pdu.FunctionCode)
	copy(ret[2:], pdu.Data)
	crc := RtuCrc(ret[:len(ret)-2])
	binary.BigEndian.PutUint16(ret[len(ret)-2:], crc)
	return ret, nil
}

// Decode decodes a RTU packet
func (r *RTU) Decode(packet []byte) (byte, PDU, error) {
	err := CheckRtuCrc(packet)
	if err != nil {
		return 0, PDU{}, err
	}

	ret := PDU{}

	ret.FunctionCode = FunctionCode(packet[1])

	if len(packet) < 4 {
		return 0, PDU{}, fmt.Errorf("short packet, got %d bytes", len(packet))
	}

	id := packet[0]

	ret.Data = packet[2 : len(packet)-2]

	return id, ret, nil
}

// Type returns TransportType
func (r *RTU) Type() TransportType {
	return TransportTypeRTU
}

func (r *RTU) setIncomingFrameType(ft frameType) {
	r.incomingFrameType = ft
}

func (r *RTU) bufferRequestFrame(funcCode FunctionCode) (int, error) {
	// get minimum amount of bytes for this function code
	frameLen, err := r.bufferMinimalLength(minRequestLen, funcCode)
	if err != nil {
		return frameLen, err
	}

	// get remaining bytes if needed
	switch funcCode {
	case FuncCodeReadDiscreteInputs, FuncCodeReadCoils, FuncCodeReadInputRegisters, FuncCodeReadHoldingRegisters, FuncCodeWriteSingleCoil, FuncCodeWriteSingleRegister:
		// nothing to do
	case FuncCodeWriteMultipleCoils, FuncCodeWriteMultipleRegisters:
		dataLen := int(r.framebuf.Bytes()[6])
		// copy next bytes to buffer
		n, err := io.CopyN(r.framebuf, r.port, int64(dataLen))
		if err != nil {
			return int(n), err
		}
		frameLen += int(n)
	case FuncCodeReadWriteMultipleRegisters, FuncCodeMaskWriteRegister, FuncCodeReadFIFOQueue:
		fallthrough
	default:
		// not implemented, keep reading until timeout // TODO: fix this!
		n, err := io.CopyN(r.framebuf, r.port, 256)
		return int(n), err
	}

	return frameLen, nil
}

func (r *RTU) bufferResponseFrame(funcCode FunctionCode) (int, error) {
	// get minimum amount of bytes for this function code
	frameLen, err := r.bufferMinimalLength(minResponseLen, funcCode)
	if err != nil {
		return frameLen, err
	}

	// get remaining bytes if needed
	switch funcCode {
	case FuncCodeWriteSingleCoil, FuncCodeWriteSingleRegister, FuncCodeWriteMultipleCoils, FuncCodeWriteMultipleRegisters:
		// nothing to do
	case FuncCodeReadCoils, FuncCodeReadDiscreteInputs, FuncCodeReadHoldingRegisters, FuncCodeReadInputRegisters:
		dataLen := int(r.framebuf.Bytes()[2])

		// copy next bytes to buffer
		n, err := io.CopyN(r.framebuf, r.port, int64(dataLen))
		if err != nil {
			return int(n), err
		}
		frameLen += int(n)
	case FuncCodeReadWriteMultipleRegisters, FuncCodeMaskWriteRegister, FuncCodeReadFIFOQueue:
		fallthrough
	default:
		// not implemented, keep reading until timeout // TODO: fix this!
		n, err := io.CopyN(r.framebuf, r.port, 256)
		return int(n), err
	}

	return frameLen, nil
}

func (r *RTU) bufferMinimalLength(minLen map[FunctionCode]int, funcCode FunctionCode) (int, error) {
	frameLen, ok := minLen[funcCode]
	if !ok {
		// unknown function code, keep reading until timeout // TODO: Do we want this?
		n, err := io.CopyN(r.framebuf, r.port, 256)
		return int(n), err
	}

	frameLen += 1 + 2 // include address and checksum

	// copy next bytes to buffer
	n, err := io.CopyN(r.framebuf, r.port, int64(frameLen-rtuFrameStartLen))
	if err != nil {
		return int(n), err
	}

	return frameLen, nil
}
