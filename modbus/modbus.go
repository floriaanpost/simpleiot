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

func (fc FunctionCode) String() string {
	suffix := ""
	if fc&0x80 != 0 {
		fc &= (0x80 ^ 0xFF) // reset fault bit
		suffix = "Fault"
	}
	switch fc {
	case FuncCodeReadDiscreteInputs:
		return "ReadDiscreteInputs" + suffix
	case FuncCodeReadCoils:
		return "ReadCoils" + suffix
	case FuncCodeWriteSingleCoil:
		return "WriteSingleCoil" + suffix
	case FuncCodeWriteMultipleCoils:
		return "WriteMultipleCoils" + suffix
	case FuncCodeReadInputRegisters:
		return "ReadInputRegisters" + suffix
	case FuncCodeReadHoldingRegisters:
		return "ReadHoldingRegisters" + suffix
	case FuncCodeWriteSingleRegister:
		return "WriteSingleRegister" + suffix
	case FuncCodeWriteMultipleRegisters:
		return "WriteMultipleRegisters" + suffix
	case FuncCodeReadWriteMultipleRegisters:
		return "ReadWriteMultipleRegisters" + suffix
	case FuncCodeMaskWriteRegister:
		return "MaskWriteRegister" + suffix
	case FuncCodeReadFIFOQueue:
		return "ReadFIFOQueue" + suffix
	default:
		return "Unknown" + suffix
	}
}

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

func (e ExceptionCode) Error() string {
	switch e {
	case 1:
		return "ILLEGAL FUNCTION"
	case 2:
		return "ILLEGAL DATA ADDRESS"
	case 3:
		return "ILLEGAL DATA VALUE"
	case 4:
		return "SERVER DEVICE FAILURE"
	case 5:
		return "ACKNOWLEDGE"
	case 6:
		return "SERVER DEVICE BUSY"
	case 8:
		return "MEMORY PARITY ERROR"
	case 0x0a:
		return "GATEWAY PATH UNAVAILABLE"
	case 0x0B:
		return "GATEWAY TARGET DEVICE FAILED TO RESPOND"
	}
	return fmt.Sprintf("unknown exception code %x", int(e))
}
