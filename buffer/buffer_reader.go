package buffer

import (
	"bufio"
	"strconv"
)

type BufferReader struct {
	Handle *bufio.Reader
}

var CRLF = []byte{'\r', '\n'}

func (reader BufferReader) IsEmpty() bool {
	return reader.Handle.Buffered() <= 0
}

func (reader BufferReader) ReadByte() (byte, error) {
	return reader.Handle.ReadByte()
}

func (reader BufferReader) ReadInt() int {
	line := reader.ReadLine()
	value, _ := strconv.Atoi(line)

	return value
}

func (reader BufferReader) ReadLine() string {
	for {
		value, err := reader.Handle.ReadBytes(CRLF[1])
		if err != nil {
			continue
		}

		if value[len(value)-2] != CRLF[0] {
			continue
		}

		return string(value[:len(value)-2])
	}
}
