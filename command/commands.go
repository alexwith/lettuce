package command

import (
	"github.com/alexwith/lettuce/protocol"
	"github.com/alexwith/lettuce/storage"
)

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

	RegisterCommand("SET", func(protocol *protocol.RESPProtocol, context *CommandContext) {
		key := string(context.Args[0])
		value := context.Args[1]

		storage.Set(key, value)

		protocol.WriteSimpleString("OK")
	})

	RegisterCommand("GET", func(protocol *protocol.RESPProtocol, context *CommandContext) {
		key := string(context.Args[0])

		value := storage.Get(key)

		protocol.WriteBulkString(string(value))
	})

	RegisterCommand("INCR", func(protocol *protocol.RESPProtocol, context *CommandContext) {
		key := string(context.Args[0])
		value, err := storage.Increment(key)
		if err != nil {
			protocol.WriteError("ERR value is not an integer or out of range")
			return
		}

		protocol.WriteInteger(value)
	})
}
