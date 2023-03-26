package protocol

import (
	"net"

	"github.com/alexwith/lettuce/buffer"
)

type RESPProtocol struct {
	Connection net.Conn
	Reader     *buffer.BufferReader
}

type DataType byte

const (
	SimpleStringType DataType = '+'
	ErrorType        DataType = '-'
	IntegerType      DataType = ':'
	BulkStringType   DataType = '$'
	ArrayType        DataType = '*'
)
