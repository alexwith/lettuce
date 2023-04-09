package protocol

import (
	"errors"
	"strconv"
)

func (protocol *RESPProtocol) GetDataType() (DataStringType, error) {
	dataType, err := protocol.Reader.ReadByte()

	return DataStringType(dataType), err
}

func (protocol *RESPProtocol) ParseDataType() ([]byte, error) {
	dataType, err := protocol.GetDataType()
	if err != nil {
		return []byte{byte(dataType)}, err
	}

	switch dataType {
	case SimpleStringType:
		return []byte(protocol.ParseSimpleString()), nil
	case IntegerType:
		value, err := protocol.ParseInteger()
		return []byte(strconv.Itoa(value)), err
	case BulkStringType:
		value, err := protocol.ParseBulkString()
		return []byte(value), err
	default:
		return nil, errors.New("Failed to parse the data type")
	}
}

func (protocol *RESPProtocol) ParseArray() ([][]byte, error) {
	size := protocol.Reader.ReadInt()

	var array [][]byte
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
		value, _ := protocol.Reader.ReadByte()
		bulkString = append(bulkString, value)
	}

	for i := 0; i < 2; i++ {
		protocol.Reader.ReadByte()
	}

	return string(bulkString), nil
}
