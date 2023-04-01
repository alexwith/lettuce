package connection

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"

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

		redisCommand := string(commandArgs[0])
		redisCommandArgs := commandArgs[1:]

		commandHandler := command.GetCommand(redisCommand)
		if commandHandler == nil {
			var unknownCommandArgs []string
			for i := 0; i < len(redisCommandArgs); i++ {
				unknownCommandArg := fmt.Sprintf("'%s'", string(redisCommandArgs[i]))
				unknownCommandArgs = append(unknownCommandArgs, unknownCommandArg)
			}

			error := fmt.Sprintf("ERR unknown command '%s', with args beginning with: %s", redisCommand, strings.Join(unknownCommandArgs, " "))
			respProtocol.WriteError(error)
			continue
		}

		commandContext := &command.CommandContext{
			Args: redisCommandArgs,
		}
		commandHandler.(func(*protocol.RESPProtocol, *command.CommandContext))(respProtocol, commandContext)
	}
}
