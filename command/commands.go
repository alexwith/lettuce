package command

import "github.com/alexwith/lettuce/protocol"

var commands map[string]interface{} = make(map[string]interface{})

func GetCommand(command string) interface{} {
	return commands[command]
}

func HandleCommand(args []interface{}, protocol *protocol.RESPProtocol) {
	protocol.WriteSimpleString("PONG")
}

func RegisterCommands() {
	commands["PING"] = HandleCommand
}
