package connection

import (
	"net"

	"github.com/alexwith/lettuce/protocol"
)

type LettuceConnection struct {
	Connection net.Conn
	Protocol   protocol.RESPProtocol
}
