package protocol

import (
	"github.com/alexwith/lettuce/buffer"
)

type DataType byte

const (
	SimpleStringType DataType = '+'
	ErrorType        DataType = '-'
	IntegerType      DataType = ':'
	BulkStringType   DataType = '$'
	ArrayType        DataType = '*'
)

func GetDataType(reader *buffer.BufferReader) DataType {
	dataType, _ := reader.Handle.ReadByte()
	return DataType(dataType)
}

func ParseDataType(reader *buffer.BufferReader) interface{} {
	dataType := GetDataType(reader)
	switch dataType {
	case SimpleStringType:
		return ""
	case ErrorType:
		return ""
	case IntegerType:
		return 10
	case BulkStringType:
		return ParseBulkString(reader)
	case ArrayType:
		return ParseArray(reader)
	default:
		return ""
	}
}

func ParseArray(reader *buffer.BufferReader) []interface{} {
	size := reader.ReadInt()

	var array []interface{}
	for i := 0; i < size; i++ {
		array = append(array, ParseDataType(reader))
	}

	return array
}

func ParseBulkString(reader *buffer.BufferReader) string {
	length := reader.ReadInt()

	if length > 512*1024*1024 {
		return ""
	}

	var bulkString []byte
	for i := 0; i < length; i++ {
		value, _ := reader.Handle.ReadByte()
		bulkString = append(bulkString, value)
	}

	for i := 0; i < 2; i++ {
		reader.Handle.ReadByte()
	}

	return string(bulkString)
}
