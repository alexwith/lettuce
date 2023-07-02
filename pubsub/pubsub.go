package pubsub

import (
	"github.com/alexwith/lettuce/protocol"
	glob "github.com/ganbarodigital/go_glob"
	"golang.org/x/exp/slices"
)

type PubSub struct {
	Subscribers map[string][]*protocol.Connection
	Channels    map[*protocol.Connection][]string
}

var pubsub = createPubSub()

func createPubSub() *PubSub {
	pubsub := &PubSub{}
	pubsub.Subscribers = make(map[string][]*protocol.Connection)
	pubsub.Channels = make(map[*protocol.Connection][]string)

	return pubsub
}

func GetPubSub() *PubSub {
	return pubsub
}

func (pubsub *PubSub) Subscribe(connection *protocol.Connection, channel string) {
	pubsub.Subscribers[channel] = append(pubsub.Subscribers[channel], connection)

	if _, present := pubsub.Channels[connection]; !present {
		connection.CloseListeners = append(connection.CloseListeners, func() {
			for _, channel := range pubsub.Channels[connection] {
				pubsub.Unsubscribe(connection, channel)
			}
		})
	}

	pubsub.Channels[connection] = append(pubsub.Channels[connection], channel)
}

func (pubsub *PubSub) Unsubscribe(connection *protocol.Connection, channel string) {
	protocolIndex := -1
	subscribers := pubsub.Subscribers[channel]
	for i := range subscribers {
		if subscribers[i] == connection {
			protocolIndex = i
			break
		}
	}

	if protocolIndex == -1 {
		return
	}

	pubsub.Subscribers[channel] = slices.Delete(pubsub.Subscribers[channel], protocolIndex, protocolIndex+1)

	channelIndex := -1
	channels := pubsub.Channels[connection]
	for i := range channels {
		if channels[i] == channel {
			channelIndex = i
			break
		}
	}

	pubsub.Channels[connection] = slices.Delete(pubsub.Channels[connection], channelIndex, channelIndex+1)
}

func (pubsub *PubSub) PSubscribe(connection *protocol.Connection, pattern string) {
	glob := glob.NewGlob(pattern)
	for channel := range pubsub.Subscribers {
		matches, err := glob.Match(channel)
		if err != nil || !matches {
			continue
		}

		pubsub.Subscribe(connection, channel)
		break
	}
}

func (pubsub *PubSub) Publish(connection *protocol.Connection, channel string, message string) int {
	subscribers := pubsub.Subscribers[channel]
	for _, connection := range subscribers {
		connection.WriteBulkString(message)
	}

	return len(subscribers)
}
