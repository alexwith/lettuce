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

		commandArgs, err := protocol.ParseArray(reader)
		if err != nil {
			fmt.Println("Failed to parse the incoming redis command:", err.Error())
		}

		fmt.Println("command args", commandArgs)

		response := []byte("+OK\r\n")
		connection.Write(response)
	}
}
