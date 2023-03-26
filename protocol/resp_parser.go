package protocol

import (
	"errors"
	"fmt"
	"strconv"

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
	dataType, err := reader.Handle.ReadByte()
	if err != nil {
		fmt.Println("Failed to get data type:", err.Error())
	}

	return DataType(dataType)
}

func ParseDataType(reader *buffer.BufferReader) (interface{}, error) {
	dataType := GetDataType(reader)
	switch dataType {
	case SimpleStringType:
		return ParseSimpleString(reader), nil
	case IntegerType:
		return ParseInteger(reader)
	case BulkStringType:
		return ParseBulkString(reader)
	case ArrayType:
		return ParseArray(reader)
	default:
		return nil, errors.New("Failed to parse the data type")
	}
}

func ParseArray(reader *buffer.BufferReader) ([]interface{}, error) {
	size := reader.ReadInt()

	var array []interface{}
	for i := 0; i < size; i++ {
		dataType, err := ParseDataType(reader)
		if err != nil {
			return array, err
		}

		array = append(array, dataType)
	}

	return array, nil
}

func ParseSimpleString(reader *buffer.BufferReader) string {
	return reader.ReadLine()
}

func ParseInteger(reader *buffer.BufferReader) (int, error) {
	return strconv.Atoi(reader.ReadLine())
}

func ParseBulkString(reader *buffer.BufferReader) (string, error) {
	length := reader.ReadInt()

	if length > 512*1024*1024 {
		return "", errors.New("A Bulk String cannot be longer than 512MB")
	}

	var bulkString []byte
	for i := 0; i < length; i++ {
		value, _ := reader.Handle.ReadByte()
		bulkString = append(bulkString, value)
	}

	for i := 0; i < 2; i++ {
		reader.Handle.ReadByte()
	}

	return string(bulkString), nil
}
