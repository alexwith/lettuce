package command

import (
	"fmt"

	"github.com/alexwith/lettuce/protocol"
)

var commands map[string]interface{} = make(map[string]interface{})

func GetCommand(command string) interface{} {
	return commands[command]
}

func HandleCommand(args []interface{}, protocol *protocol.RESPProtocol) {
	response := "PONG"
	if len(args) > 0 {
		response = fmt.Sprint(args[0])
	}

	protocol.WriteSimpleString(response)
}

func RegisterCommands() {
	commands["PING"] = HandleCommand
}
