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
- [ ] Null elements in Arrays
- [x] Telnet commands
- [x] Pipelining

### Commands
I will only be implementing the most important commands, as I will not have time to implement the 450+ redis commands that exist. 
- [x] PING 
- [ ] KEYS 
- [ ] SET
- [ ] GET
- [ ] APPEND
- [ ] DEL
- [ ] FLUSHALL
- [ ] TTL
- [ ] EXPIRE
- [ ] EXISTS
- [ ] STRLEN
- [x] INCR

## License
Lettuce is free and open source software. The software is released under the terms of
the [GPL-3.0 license]("https://github.com/alexwith/lettuce/blob/main/LICENSE").