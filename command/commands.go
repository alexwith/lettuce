package command

import (
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
		response := "PONG"
		if len(context.Args) > 0 {
			response = string(context.Args[0])
		}

		protocol.WriteSimpleString(response)
	})
}
