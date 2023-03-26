package main

import (
	"github.com/alexwith/lettuce/connection"
)

const HOST string = "127.0.0.1"
const PORT int16 = 6380

func main() {
	connection.Listen(HOST, PORT)
}
