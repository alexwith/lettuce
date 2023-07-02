package protocol

import (
	"fmt"
	"net"
	"reflect"

	"github.com/alexwith/lettuce/buffer"
)

type Connection struct {
	Handle         net.Conn
	Reader         *buffer.BufferReader
	CloseListeners []func()
}

type DataStringType byte

const (
	SimpleStringType DataStringType = '+'
	ErrorType        DataStringType = '-'
	IntegerType      DataStringType = ':'
	BulkStringType   DataStringType = '$'
	ArrayType        DataStringType = '*'
)

const (
	NullBulkString string = "$-1\r\n"
	NullArray      string = "*-1\r\n"
)

func (connection *Connection) Close() {
	for _, listener := range connection.CloseListeners {
		listener()
	}
}

func (connection *Connection) CreateSimpleString(value string) string {
	return fmt.Sprintf("%s%s\r\n", string(SimpleStringType), value)
}

func (connection *Connection) CreateError(value string) string {
	return fmt.Sprintf("%s%s\r\n", string(ErrorType), value)
}

func (connection *Connection) CreateInteger(value int) string {
	return fmt.Sprintf("%s%d\r\n", string(IntegerType), value)
}

func (connection *Connection) CreateBulkString(value string) string {
	return fmt.Sprintf("%s%d\r\n%s\r\n", string(BulkStringType), len(value), value)
}

func (connection *Connection) CreateArray(value []any) string {
	length := len(value)

	var respArray string
	for _, element := range value {
		respArray += connection.TryToCreateRESPType(element)
	}

	return fmt.Sprintf("%s%d\r\n%s", string(ArrayType), length, respArray)
}

func (connection *Connection) TryToCreateRESPType(value any) string {
	if value == nil {
		return NullBulkString
	}

	switch valueType := value.(type) {
	case int:
		return connection.CreateInteger(value.(int))
	case string:
		return connection.CreateBulkString(value.(string)) // TODO, check if its null, bulk or simple
	default:
		fmt.Printf("Invalid RESP type %v", reflect.TypeOf(valueType))
		return ""
	}
}
