# 🥬 Lettuce

Lettuce is a redis server written in Go. It implements the [RESP protocol](https://redis.io/docs/reference/protocol-spec/) and should therefore work with any redis client. 

It is worth noting that this is just a fun project for learning purposes, but I do attempt to make it as fast and functional as possible.

## Getting Started
If you want to use the pre-existing Lettuce server, navigate to it's location and execute: `go run .`
This will start the lettuce server on port `6380`.

You can also create your own custom redis server. Execute `go get github.com/alexwith/lettuce` in your own project, then initiate the server. The following example will also register a custom command:
```go
package main

import (
  "github.com/alexwith/lettuce/command"
  "github.com/alexwith/lettuce/lettuce"
  "github.com/alexwith/lettuce/protocol"
)

const HOST string = "127.0.0.1"
const PORT int16 = 6380

func main()  {
  lettuce.Setup(HOST, PORT, func() {
    registerCommands()
  })
}

func registerCommands()  {
  command.RegisterCommand("CUSTOMPING", func(protocol *protocol.RESPProtocol, context *command.CommandContext)  {
    if len(context.Args) <= 0 {
      protocol.WriteSimpleString("CUSTOMPONG")
      return
    }
  
    response := context.StringArg(0)
    protocol.WriteBulkString(response)
  })
}
```

## Features
### RESP (REdis Serialization Protocol)
- [x] Simple Strings
- [x] Errors
- [x] Integers
- [x] Bulk Strings
- [x] Arrays
- [x] Null Array and Bulk Strings
- [x] Telnet commands
- [x] Pipelining

### Commands
I will only be implementing the most important commands, as I will not have time to implement the 450+ redis commands that exist. 
- [x] PING 
- [x] KEYS 
- [x] SET
- [x] GET
- [x] APPEND
- [x] DEL
- [x] TTL
- [x] EXPIRE
- [x] PERSIST
- [x] EXISTS
- [x] STRLEN
- [x] INCR



The Lettuce server has been built based on documentation from the following sources:
- https://redis.io/docs/reference/protocol-spec/
- https://redis.io/commands
- https://redis.io/docs/management/persistence/

## License
Lettuce is free and open source software. The software is released under the terms of
the [GPL-3.0 license]("https://github.com/alexwith/lettuce/blob/main/LICENSE").