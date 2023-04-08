package main

import (
	"github.com/alexwith/lettuce/command"
	"github.com/alexwith/lettuce/connection"
	"github.com/alexwith/lettuce/storage"
)

func Setup(host string, port int16, setup func()) {
	go storage.RegisterExpireTask()

	command.RegisterCommands()
	setup()

	connection.Listen(host, port)
}
