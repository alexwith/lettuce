package command

import (
	"strings"
	"time"

	"github.com/alexwith/lettuce/protocol"
	"github.com/alexwith/lettuce/pubsub"
	"github.com/alexwith/lettuce/storage"
	glob "github.com/ganbarodigital/go_glob"
)

type CommandData struct {
	Command      string
	MinArguments int
	Handler      func(protocol *protocol.Connection, context *CommandContext)
}

var commands map[string]*CommandData = make(map[string]*CommandData)

func GetCommand(command string) *CommandData {
	return commands[strings.ToUpper(command)]
}

func RegisterCommand(command string, minArguments int, handler func(connection *protocol.Connection, context *CommandContext)) {
	commands[command] = &CommandData{
		Command:      command,
		MinArguments: minArguments,
		Handler:      handler,
	}
}

func RegisterCommands() {
	RegisterCommand("PING", -1, func(connection *protocol.Connection, context *CommandContext) {
		if len(context.Args) <= 0 {
			connection.WriteSimpleString("PONG")
			return
		}

		response := context.StringArg(0)
		connection.WriteBulkString(response)
	})

	RegisterCommand("SET", 2, func(connection *protocol.Connection, context *CommandContext) {
		key := context.StringArg(0)
		value := context.Args[1]

		expireTime, ex := context.ReadOptionAsInt("EX")
		expireTimeMs, px := context.ReadOptionAsInt("PX")
		timeout, exat := context.ReadOptionAsInt("EXAT")
		timeoutMs, pxat := context.ReadOptionAsInt("PXAT")
		nx := context.HasOption("NX")
		xx := context.HasOption("XX")
		keepttl := context.HasOption("KEEPTTL")
		get := context.HasOption("GET")

		oldValue, oldPresent := storage.Get(key)
		if nx && oldPresent {
			connection.WriteNullBulkString()
			return
		}

		if xx && !oldPresent {
			connection.WriteNullBulkString()
			return
		}

		storage.Set(key, value, !keepttl)

		if ex {
			storage.ExpireIn(key, int64(expireTime*1000))
		}

		if px {
			storage.ExpireIn(key, int64(expireTimeMs))
		}

		if exat {
			storage.Expire(key, int64(timeout*1000))
		}

		if pxat {
			storage.Expire(key, int64(timeoutMs))
		}

		if get {
			if oldPresent {
				connection.WriteBulkString(string(oldValue))
			} else {
				connection.WriteNullBulkString()
			}

			return
		}

		connection.WriteSimpleString("OK")
	})

	RegisterCommand("GET", 1, func(connection *protocol.Connection, context *CommandContext) {
		key := context.StringArg(0)

		value, present := storage.Get(key)
		if !present {
			connection.WriteNullBulkString()
			return
		}

		connection.WriteBulkString(string(value))
	})

	RegisterCommand("INCR", 2, func(connection *protocol.Connection, context *CommandContext) {
		key := context.StringArg(0)
		value, err := storage.Increment(key)
		if err != nil {
			connection.WriteError("ERR value is not an integer or out of range")
			return
		}

		connection.WriteInteger(value)
	})

	RegisterCommand("DEL", 1, func(connection *protocol.Connection, context *CommandContext) {
		var amount int
		for _, key := range context.Args {
			success := storage.Delete(string(key))
			if !success {
				continue
			}

			amount++
		}

		connection.WriteInteger(amount)
	})

	RegisterCommand("EXISTS", 1, func(connection *protocol.Connection, context *CommandContext) {
		var amount int
		for _, key := range context.Args {
			_, present := storage.Get(string(key))
			if !present {
				continue
			}

			amount++
		}

		connection.WriteInteger(amount)
	})

	RegisterCommand("STRLEN", 1, func(connection *protocol.Connection, context *CommandContext) {
		key := context.StringArg(0)

		var length int

		value, present := storage.Get(key)
		if !present {
			length = 0
		} else {
			length = len(value)
		}

		connection.WriteInteger(length)
	})

	RegisterCommand("EXPIRE", 2, func(connection *protocol.Connection, context *CommandContext) {
		key := context.StringArg(0)
		seconds := context.IntegerArg(1)

		failed := func() {
			connection.WriteInteger(0)
		}

		nx := context.HasOption("NX")
		xx := context.HasOption("XX")
		gt := context.HasOption("GT")
		lt := context.HasOption("LT")

		_, keyPresent := storage.Get(key)
		if !keyPresent {
			failed()
			return
		}

		currentTimout, timeoutPresent := storage.GetTimeout(key)

		if nx && timeoutPresent {
			failed()
			return
		}

		if xx && !timeoutPresent {
			failed()
			return
		}

		timeout := time.Now().UnixMilli() + int64(seconds*1000)
		if gt && timeout <= currentTimout {
			failed()
			return
		}

		if lt && timeout >= currentTimout {
			failed()
			return
		}

		storage.Expire(key, timeout)
		connection.WriteInteger(1)
	})

	RegisterCommand("PERSIST", 1, func(connection *protocol.Connection, context *CommandContext) {
		key := context.StringArg(0)

		success := storage.Persist(key)

		status := 0
		if success {
			status = 1
		}

		connection.WriteInteger(status)
	})

	RegisterCommand("TTL", 1, func(connection *protocol.Connection, context *CommandContext) {
		key := context.StringArg(0)

		_, keyPresent := storage.Get(key)
		if !keyPresent {
			connection.WriteInteger(-2)
			return
		}

		timeout, timeoutPresent := storage.GetTimeout(key)
		if !timeoutPresent {
			connection.WriteInteger(-1)
			return
		}

		remainingTime := int(timeout-time.Now().UnixMilli()) / 1000

		connection.WriteInteger(remainingTime)
	})

	RegisterCommand("APPEND", 2, func(connection *protocol.Connection, context *CommandContext) {
		key := context.StringArg(0)
		value := context.StringArg(1)

		length := storage.Append(key, value)

		connection.WriteInteger(length)
	})

	RegisterCommand("KEYS", 1, func(connection *protocol.Connection, context *CommandContext) {
		pattern := context.StringArg(0)

		glob := glob.NewGlob(pattern)

		var keys []any
		for key := range storage.Storage {
			matches, err := glob.Match(key)
			if !matches || err != nil {
				continue
			}

			keys = append(keys, key)
		}

		connection.WriteArray(keys)
	})

	RegisterCommand("SUBSCRIBE", 1, func(connection *protocol.Connection, context *CommandContext) {
		for i := 0; i < len(context.Args); i++ {
			channel := context.StringArg(i)
			pubsub.GetPubSub().Subscribe(connection, channel)
		}
	})

	RegisterCommand("UNSUBSCRIBE", 1, func(connection *protocol.Connection, context *CommandContext) {
		for i := 0; i < len(context.Args); i++ {
			channel := context.StringArg(i)
			pubsub.GetPubSub().Unsubscribe(connection, channel)
		}
	})

	RegisterCommand("PUBLISH", 2, func(connection *protocol.Connection, context *CommandContext) {
		channel := context.StringArg(0)
		message := context.StringArg(1)

		clients := pubsub.GetPubSub().Publish(connection, channel, message)
		connection.WriteInteger(clients)
	})

	RegisterCommand("PSUBSCRIBE", 1, func(connection *protocol.Connection, context *CommandContext) {
		pubsub := pubsub.GetPubSub()

		for i := 0; i < len(context.Args); i++ {
			pattern := context.StringArg(i)
			channel, err := pubsub.FindChannelByGlob(pattern)
			if err != nil {
				continue
			}

			pubsub.Subscribe(connection, channel)
		}
	})

	RegisterCommand("PUNSUBSCRIBE", 1, func(connection *protocol.Connection, context *CommandContext) {
		pubsub := pubsub.GetPubSub()

		for i := 0; i < len(context.Args); i++ {
			pattern := context.StringArg(i)
			channel, err := pubsub.FindChannelByGlob(pattern)
			if err != nil {
				continue
			}

			pubsub.Unsubscribe(connection, channel)
		}
	})
}
