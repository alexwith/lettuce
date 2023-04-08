package command

import (
	"time"

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
		key := context.StringArg(0)
		value := context.Args[1]

		storage.Set(key, value)

		protocol.WriteSimpleString("OK")
	})

	RegisterCommand("GET", func(protocol *protocol.RESPProtocol, context *CommandContext) {
		key := context.StringArg(0)

		value, present := storage.Get(key)
		if !present {
			protocol.WriteNullBulkString()
			return
		}

		protocol.WriteBulkString(string(value))
	})

	RegisterCommand("INCR", func(protocol *protocol.RESPProtocol, context *CommandContext) {
		key := context.StringArg(0)
		value, err := storage.Increment(key)
		if err != nil {
			protocol.WriteError("ERR value is not an integer or out of range")
			return
		}

		protocol.WriteInteger(value)
	})

	RegisterCommand("DEL", func(protocol *protocol.RESPProtocol, context *CommandContext) {
		var amount int
		for _, key := range context.Args {
			success := storage.Delete(string(key))
			if !success {
				continue
			}

			amount++
		}

		protocol.WriteInteger(amount)
	})

	RegisterCommand("EXISTS", func(protocol *protocol.RESPProtocol, context *CommandContext) {
		var amount int
		for _, key := range context.Args {
			_, present := storage.Get(string(key))
			if !present {
				continue
			}

			amount++
		}

		protocol.WriteInteger(amount)
	})

	RegisterCommand("STRLEN", func(protocol *protocol.RESPProtocol, context *CommandContext) {
		key := context.StringArg(0)

		var length int

		value, present := storage.Get(key)
		if !present {
			length = 0
		} else {
			length = len(value)
		}

		protocol.WriteInteger(length)
	})

	RegisterCommand("EXPIRE", func(protocol *protocol.RESPProtocol, context *CommandContext) {
		key := context.StringArg(0)
		seconds := context.IntegerArg(1)

		nx := context.HasOption("NX")
		xx := context.HasOption("XX")
		gt := context.HasOption("GT")
		lt := context.HasOption("LT")

		success := storage.Expire(key, seconds, nx, xx, gt, lt)

		status := 0
		if success {
			status = 1
		}

		protocol.WriteInteger(status)
	})

	RegisterCommand("PERSIST", func(protocol *protocol.RESPProtocol, context *CommandContext) {
		key := context.StringArg(0)

		success := storage.Persist(key)

		status := 0
		if success {
			status = 1
		}

		protocol.WriteInteger(status)
	})

	RegisterCommand("TTL", func(protocol *protocol.RESPProtocol, context *CommandContext) {
		key := context.StringArg(0)

		_, keyPresent := storage.Get(key)
		if !keyPresent {
			protocol.WriteInteger(-2)
			return
		}

		timeout, timeoutPresent := storage.GetTimeout(key)
		if !timeoutPresent {
			protocol.WriteInteger(-1)
			return
		}

		remainingTime := int(timeout-time.Now().UnixMilli()) / 1000

		protocol.WriteInteger(remainingTime)
	})
}
