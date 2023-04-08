package main

import (
	"github.com/alexwith/lettuce/command"
	"github.com/alexwith/lettuce/protocol"
)

const HOST string = "127.0.0.1"
const PORT int16 = 6380

func main() {
	Setup(HOST, PORT, func() {
		registerCommands()
	})
}

func registerCommands() {
	command.RegisterCommand("CUSTOMPING", func(protocol *protocol.RESPProtocol, context *command.CommandContext) {
		if len(context.Args) <= 0 {
			protocol.WriteSimpleString("CUSTOMPONG")
			return
		}

		response := context.StringArg(0)
		protocol.WriteBulkString(response)
	})
}
