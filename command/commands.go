package command

import (
	"fmt"

	"github.com/alexwith/lettuce/protocol"
)

type CommandContext struct {
	Args [][]byte
}

var commands map[string]interface{} = make(map[string]interface{})

func GetCommand(command string) interface{} {
	return commands[command]
}

func RegisterCommand(command string, handler func(protocol *protocol.RESPProtocol, context *CommandContext)) {
	commands[command] = handler
}

func RegisterCommands() {
	RegisterCommand("PING", func(protocol *protocol.RESPProtocol, context *CommandContext) {
		if len(context.Args) <= 0 {
			protocol.WriteSimpleString("PONG")
			return
		}

		response := string(context.Args[0])
		fmt.Println(response)

		yes := []int{1, 2, 3, 4, 5}
		protocol.WriteArray(yes)
		//protocol.WriteBulkString(response)
	})
}
