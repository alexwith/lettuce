package main

import (
	"github.com/alexwith/lettuce/lettuce"
)

const HOST string = "127.0.0.1"
const PORT int16 = 6380

func main() {
	lettuce.Setup(HOST, PORT, func() {})
}
