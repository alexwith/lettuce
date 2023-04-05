package command

import (
	"github.com/alexwith/lettuce/protocol"
)

type CommandContext struct {
	Args [][]byte
}

var commands map[string]any = make(map[string]any)

func GetCommand(command string) any {
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
		protocol.WriteBulkString(response)
	})
}
