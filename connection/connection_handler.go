package connection

import (
	"bufio"
	"fmt"
	"io"
	"net"

	"github.com/alexwith/lettuce/buffer"
	"github.com/alexwith/lettuce/command"
	"github.com/alexwith/lettuce/protocol"
)

func HandleConnection(connection net.Conn) {
	defer connection.Close()

	respProtocol := &protocol.RESPProtocol{
		Connection: connection,
		Reader: &buffer.BufferReader{
			Handle: bufio.NewReader(connection),
		},
	}

	for {
		dataType, err := respProtocol.GetDataType()
		if err != nil && err == io.EOF {
			break
		}

		if dataType != protocol.ArrayType {
			continue
		}

		commandArgs, err := respProtocol.ParseArray()
		if err != nil {
			fmt.Println("Failed to parse the incoming redis command:", err.Error())
		}

		redisCommand := fmt.Sprint(commandArgs[0])

		commandHandler := command.GetCommand(redisCommand)
		if commandHandler == nil {
			respProtocol.WriteSimpleString("HELLO")
			continue
		}

		commandHandler.(func([]interface{}, *protocol.RESPProtocol))(commandArgs[1:], respProtocol)
	}
}
