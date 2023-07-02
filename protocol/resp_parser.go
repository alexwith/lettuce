package protocol

import (
	"errors"
	"strconv"
)

func (connection *Connection) GetDataType() (DataStringType, error) {
	dataType, err := connection.Reader.ReadByte()

	return DataStringType(dataType), err
}

func (connection *Connection) ParseDataType() ([]byte, error) {
	dataType, err := connection.GetDataType()
	if err != nil {
		return []byte{byte(dataType)}, err
	}

	switch dataType {
	case SimpleStringType:
		return []byte(connection.ParseSimpleString()), nil
	case IntegerType:
		value, err := connection.ParseInteger()
		return []byte(strconv.Itoa(value)), err
	case BulkStringType:
		value, err := connection.ParseBulkString()
		return []byte(value), err
	default:
		return nil, errors.New("Failed to parse the data type")
	}
}

func (connection *Connection) ParseArray() ([][]byte, error) {
	size := connection.Reader.ReadInt()

	var array [][]byte
	for i := 0; i < size; i++ {
		dataType, err := connection.ParseDataType()
		if err != nil {
			return array, err
		}

		array = append(array, dataType)
	}

	return array, nil
}

func (connection *Connection) ParseSimpleString() string {
	return connection.Reader.ReadLine()
}

func (connection *Connection) ParseInteger() (int, error) {
	return strconv.Atoi(connection.Reader.ReadLine())
}

func (connection *Connection) ParseBulkString() (string, error) {
	length := connection.Reader.ReadInt()

	if length > 512*1024*1024 {
		return "", errors.New("A Bulk String cannot be longer than 512MB")
	}

	var bulkString []byte
	for i := 0; i < length; i++ {
		value, _ := connection.Reader.ReadByte()
		bulkString = append(bulkString, value)
	}

	for i := 0; i < 2; i++ {
		connection.Reader.ReadByte()
	}

	return string(bulkString), nil
}
