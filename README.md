# 🥬 Lettuce

Lettuce is a redis server written in Go. It implements the [RESP protocol](https://redis.io/docs/reference/protocol-spec/) and should therefore work with any redis client. 

It is worth noting that this is just a fun project for learning purposes, but I do attempt to make it as fast and functional as possible.

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
- [ ] KEYS 
- [ ] SET
- [x] GET
- [ ] APPEND
- [x] DEL
- [ ] FLUSHALL
- [ ] TTL
- [ ] EXPIRE
- [x] EXISTS
- [ ] STRLEN
- [x] INCR

The Lettuce server has been built based on documentation from the following sources:
- https://redis.io/docs/reference/protocol-spec/
- https://redis.io/commands
- https://redis.io/docs/management/persistence/

## License
Lettuce is free and open source software. The software is released under the terms of
the [GPL-3.0 license]("https://github.com/alexwith/lettuce/blob/main/LICENSE").