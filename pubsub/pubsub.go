package pubsub

import (
	"github.com/alexwith/lettuce/protocol"
)

type PubSub struct {
	Connections []*protocol.RESPProtocol
}

var pubsub = &PubSub{}

func GetPubSub() *PubSub {
	return pubsub
}

func (pubsub *PubSub) Subscribe(protocol *protocol.RESPProtocol, channel string) {
	pubsub.Connections = append(pubsub.Connections, protocol)
}

func (pubsub *PubSub) Publish(protocol *protocol.RESPProtocol, channel string, message string) int {
	for _, protocol := range pubsub.Connections {
		protocol.WriteBulkString(message)
	}

	return len(pubsub.Connections)
}
