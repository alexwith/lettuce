package buffer

import (
	"bufio"
	"strconv"
)

type BufferReader struct {
	Handle *bufio.Reader
}

func (reader BufferReader) ReadInt() int {
	line := reader.ReadLine()
	value, _ := strconv.Atoi(line)

	return value
}

func (reader BufferReader) ReadLine() string {
	for {
		value, err := reader.Handle.ReadBytes('\n')
		if err != nil {
			continue
		}

		if value[len(value)-2] != '\r' {
			continue
		}

		return string(value[:len(value)-2])
	}
}
