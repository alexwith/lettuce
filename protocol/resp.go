package protocol

import (
	"fmt"
	"net"
	"reflect"

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

func (protocol *RESPProtocol) CreateSimpleString(value string) string {
	return fmt.Sprintf("%s%s\r\n", string(SimpleStringType), value)
}

func (protocol *RESPProtocol) CreateError(value string) string {
	return fmt.Sprintf("%s%s\r\n", string(ErrorType), value)
}

func (protocol *RESPProtocol) CreateInteger(value int) string {
	return fmt.Sprintf("%s%d\r\n", string(IntegerType), value)
}

func (protocol *RESPProtocol) CreateBulkString(value string) string {
	return fmt.Sprintf("%s%d\r\n%s\r\n", string(BulkStringType), len(value), value)
}

func (protocol *RESPProtocol) CreateArray(value interface{}) string {
	array := value.([]int) // Fix casting to []interface{}
	length := len(array)

	var respArray string
	for _, element := range array {
		respArray += protocol.TryToCreateRESPType(element)
	}

	return fmt.Sprintf("%s%d\r\n%s", string(ArrayType), length, respArray)
}

func (protocol *RESPProtocol) TryToCreateRESPType(value interface{}) string {
	switch valueType := value.(type) {
	case int:
		return protocol.CreateInteger(value.(int))
	case string:
		return protocol.CreateBulkString(value.(string)) // TODO, check if its null, bulk or simple
	default:
		fmt.Printf("Invalid RESP type %v", reflect.TypeOf(valueType))
		return ""
	}
}
