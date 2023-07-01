package pubsub

import (
	"github.com/alexwith/lettuce/protocol"
)

type PubSub struct {
	Connections map[string][]*protocol.RESPProtocol
}

var pubsub = createPubSub()

func createPubSub() *PubSub {
	pubsub := &PubSub{}
	pubsub.Connections = make(map[string][]*protocol.RESPProtocol)

	return pubsub
}

func GetPubSub() *PubSub {
	return pubsub
}

func (pubsub *PubSub) Subscribe(protocol *protocol.RESPProtocol, channel string) {
	pubsub.Connections[channel] = append(pubsub.Connections[channel], protocol)
}

func (pubsub *PubSub) Publish(protocol *protocol.RESPProtocol, channel string, message string) int {
	connections := pubsub.Connections[channel]
	for _, protocol := range connections {
		protocol.WriteBulkString(message)
	}

	return len(connections)
}
