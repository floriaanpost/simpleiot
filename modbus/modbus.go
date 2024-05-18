package modbus

import "fmt"

// FunctionCode represents a modbus function code
type FunctionCode byte

// Defined valid function codes
const (
	// Bit access
	FuncCodeReadDiscreteInputs FunctionCode = 2
	FuncCodeReadCoils          FunctionCode = 1
	FuncCodeWriteSingleCoil    FunctionCode = 5
	FuncCodeWriteMultipleCoils FunctionCode = 15

	// 16-bit access
	FuncCodeReadInputRegisters         FunctionCode = 4
	FuncCodeReadHoldingRegisters       FunctionCode = 3
	FuncCodeWriteSingleRegister        FunctionCode = 6
	FuncCodeWriteMultipleRegisters     FunctionCode = 16
	FuncCodeReadWriteMultipleRegisters FunctionCode = 23
	FuncCodeMaskWriteRegister          FunctionCode = 22
	FuncCodeReadFIFOQueue              FunctionCode = 24
)

// ExceptionCode represents a modbus exception code
type ExceptionCode byte

// Defined valid exception codes
const (
	ExcIllegalFunction              ExceptionCode = 1
	ExcIllegalAddress               ExceptionCode = 2
	ExcIllegalValue                 ExceptionCode = 3
	ExcServerDeviceFailure          ExceptionCode = 4
	ExcAcknowledge                  ExceptionCode = 5
	ExcServerDeviceBusy             ExceptionCode = 6
	ExcMemoryParityError            ExceptionCode = 8
	ExcGatewayPathUnavilable        ExceptionCode = 0x0a
	ExcGatewayTargetFailedToRespond ExceptionCode = 0x0b
)

// define valid values for write coil
const (
	WriteCoilValueOn  uint16 = 0xff00
	WriteCoilValueOff uint16 = 0
)

// minRequestLen is the minimum number of PDU bytes for a request with
// the given function code (not including slave address or checksum,
// which are part of the ADU).
var minRequestLen = map[FunctionCode]int{
	FuncCodeReadDiscreteInputs:         5,
	FuncCodeReadCoils:                  5,
	FuncCodeWriteSingleCoil:            5,
	FuncCodeWriteMultipleCoils:         7,
	FuncCodeReadInputRegisters:         5,
	FuncCodeReadHoldingRegisters:       5,
	FuncCodeWriteSingleRegister:        5,
	FuncCodeWriteMultipleRegisters:     8,
	FuncCodeReadWriteMultipleRegisters: 12,
	FuncCodeMaskWriteRegister:          7,
	FuncCodeReadFIFOQueue:              3,
}

// minResponseLen is the minimum number of PDU bytes for a resppnse with
// the given function code (not including slave address or checksum,
// which are part of the ADU).
var minResponseLen = map[FunctionCode]int{
	FuncCodeReadDiscreteInputs:     2,
	FuncCodeReadCoils:              2,
	FuncCodeWriteSingleCoil:        5,
	FuncCodeWriteMultipleCoils:     5,
	FuncCodeReadInputRegisters:     2,
	FuncCodeReadHoldingRegisters:   2,
	FuncCodeWriteSingleRegister:    5,
	FuncCodeWriteMultipleRegisters: 5,
	// TODO: find out what these frames look like
	FuncCodeReadWriteMultipleRegisters: 0,
	FuncCodeMaskWriteRegister:          0,
	FuncCodeReadFIFOQueue:              0,
}

func (e ExceptionCode) Error() string {
	switch e {
	case ExcIllegalFunction:
		return "ILLEGAL FUNCTION"
	case ExcIllegalAddress:
		return "ILLEGAL DATA ADDRESS"
	case ExcIllegalValue:
		return "ILLEGAL DATA VALUE"
	case ExcServerDeviceFailure:
		return "SERVER DEVICE FAILURE"
	case ExcAcknowledge:
		return "ACKNOWLEDGE"
	case ExcServerDeviceBusy:
		return "SERVER DEVICE BUSY"
	case ExcMemoryParityError:
		return "MEMORY PARITY ERROR"
	case ExcGatewayPathUnavilable:
		return "GATEWAY PATH UNAVAILABLE"
	case ExcGatewayTargetFailedToRespond:
		return "GATEWAY TARGET DEVICE FAILED TO RESPOND"
	}
	return fmt.Sprintf("unknown exception code %x", int(e))
}

// This is used to set the incoming type of the transport.
// It is needed because to find the end of the frame, we need to know if
// we are dealing with a request or a response frame.
type frameType int

const (
	response frameType = iota
	request
)
