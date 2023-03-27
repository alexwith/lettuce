package protocol

import (
	"errors"
	"strconv"
)

func (protocol *RESPProtocol) GetDataType() (DataType, error) {
	dataType, err := protocol.Reader.Handle.ReadByte()

	return DataType(dataType), err
}

func (protocol *RESPProtocol) ParseDataType() (interface{}, error) {
	dataType, err := protocol.GetDataType()
	if err != nil {
		return dataType, err
	}

	switch dataType {
	case SimpleStringType:
		return protocol.ParseSimpleString(), nil
	case IntegerType:
		return protocol.ParseInteger()
	case BulkStringType:
		return protocol.ParseBulkString()
	case ArrayType:
		return protocol.ParseArray()
	default:
		return nil, errors.New("Failed to parse the data type")
	}
}

func (protocol *RESPProtocol) ParseArray() ([]interface{}, error) {
	size := protocol.Reader.ReadInt()

	var array []interface{}
	for i := 0; i < size; i++ {
		dataType, err := protocol.ParseDataType()
		if err != nil {
			return array, err
		}

		array = append(array, dataType)
	}

	return array, nil
}

func (protocol *RESPProtocol) ParseSimpleString() string {
	return protocol.Reader.ReadLine()
}

func (protocol *RESPProtocol) ParseInteger() (int, error) {
	return strconv.Atoi(protocol.Reader.ReadLine())
}

func (protocol *RESPProtocol) ParseBulkString() (string, error) {
	length := protocol.Reader.ReadInt()

	if length > 512*1024*1024 {
		return "", errors.New("A Bulk String cannot be longer than 512MB")
	}

	var bulkString []byte
	for i := 0; i < length; i++ {
		value, _ := protocol.Reader.Handle.ReadByte()
		bulkString = append(bulkString, value)
	}

	for i := 0; i < 2; i++ {
		protocol.Reader.Handle.ReadByte()
	}

	return string(bulkString), nil
}
