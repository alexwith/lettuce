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

		switch dataType {
		case protocol.ArrayType:
			HandleRespCommand(respProtocol)
		default:
			HandleRawCommand(respProtocol)
		}
	}
}

func HandleRespCommand(respProtocol *protocol.RESPProtocol) {
	commandArgs, err := respProtocol.ParseArray()
	if err != nil {
		fmt.Println("Failed to parse the incoming redis command:", err.Error())
	}

	redisCommand := string(commandArgs[0])
	redisCommandArgs := commandArgs[1:]

	HandleCommand(respProtocol, redisCommand, redisCommandArgs)
}

func HandleRawCommand(respProtocol *protocol.RESPProtocol) {
	reader := respProtocol.Reader
	reader.Handle.UnreadByte() // Undo the check we did for data type

	for !reader.IsEmpty() {
		line := reader.ReadLine()

		commandArgs := strings.Fields(line)

		redisCommand := commandArgs[0]
		redisCommandArgs := make([][]byte, len(commandArgs)-1)

		for i := 0; i < len(commandArgs)-1; i++ {
			redisCommandArgs[i] = []byte(commandArgs[i+1])
		}

		HandleCommand(respProtocol, redisCommand, redisCommandArgs)
	}
}

func HandleCommand(respProtocol *protocol.RESPProtocol, redisCommand string, redisCommandArgs [][]byte) {
	commandData := command.GetCommand(redisCommand)
	if commandData == nil {
		var unknownCommandArgs []string
		for i := 0; i < len(redisCommandArgs); i++ {
			unknownCommandArg := fmt.Sprintf("'%s'", string(redisCommandArgs[i]))
			unknownCommandArgs = append(unknownCommandArgs, unknownCommandArg)
		}

		error := fmt.Sprintf("ERR unknown command '%s', with args beginning with: %s", redisCommand, strings.Join(unknownCommandArgs, " "))
		respProtocol.WriteError(error)
		return
	}

	requiredArgumentsSize := commandData.ArgumentsSize
	if requiredArgumentsSize != -1 && len(redisCommandArgs) < requiredArgumentsSize {
		respProtocol.WriteError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", redisCommand))
		return
	}

	var stringifiedArgs map[string]int = make(map[string]int)
	for index, arg := range redisCommandArgs {
		stringifiedArgs[strings.ToUpper(string(arg))] = index
	}

	commandContext := &command.CommandContext{
		Args:            redisCommandArgs,
		StringifiedArgs: stringifiedArgs,
	}

	commandData.Handler(respProtocol, commandContext)
}
