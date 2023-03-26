package connection

import (
	"bufio"
	"fmt"
	"net"

	"github.com/alexwith/lettuce/buffer"
	"github.com/alexwith/lettuce/protocol"
)

func HandleConnection(connection net.Conn) {
	defer connection.Close()

	reader := &buffer.BufferReader{
		Handle: bufio.NewReader(connection),
	}

	for {
		dataType := protocol.GetDataType(reader)
		if dataType != protocol.ArrayType {
			continue
		}

		commandArgs := protocol.ParseArray(reader)

		fmt.Println("command args", commandArgs)

		response := []byte("+OK\r\n")
		connection.Write(response)
	}
}
