package pubsub

import (
	"github.com/alexwith/lettuce/protocol"
	glob "github.com/ganbarodigital/go_glob"
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

func (pubsub *PubSub) PSubscribe(protocol *protocol.RESPProtocol, pattern string) {
	glob := glob.NewGlob(pattern)
	for channel := range pubsub.Connections {
		matches, err := glob.Match(channel)
		if err != nil || !matches {
			continue
		}

		pubsub.Subscribe(protocol, channel)
		break
	}
}

func (pubsub *PubSub) Publish(protocol *protocol.RESPProtocol, channel string, message string) int {
	connections := pubsub.Connections[channel]
	for _, protocol := range connections {
		protocol.WriteBulkString(message)
	}

	return len(connections)
}
