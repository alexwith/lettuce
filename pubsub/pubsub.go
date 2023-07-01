package pubsub

import (
	"github.com/alexwith/lettuce/protocol"
	glob "github.com/ganbarodigital/go_glob"
	"golang.org/x/exp/slices"
)

type PubSub struct {
	Subscribers map[string][]*protocol.RESPProtocol
	Channels    map[*protocol.RESPProtocol][]string
}

var pubsub = createPubSub()

func createPubSub() *PubSub {
	pubsub := &PubSub{}
	pubsub.Subscribers = make(map[string][]*protocol.RESPProtocol)
	pubsub.Channels = make(map[*protocol.RESPProtocol][]string)

	return pubsub
}

func GetPubSub() *PubSub {
	return pubsub
}

func (pubsub *PubSub) Subscribe(protocol *protocol.RESPProtocol, channel string) {
	pubsub.Subscribers[channel] = append(pubsub.Subscribers[channel], protocol)

	if _, present := pubsub.Channels[protocol]; !present {
		protocol.CloseListeners = append(protocol.CloseListeners, func() {
			for _, channel := range pubsub.Channels[protocol] {
				pubsub.Unsubscribe(protocol, channel)
			}
		})
	}

	pubsub.Channels[protocol] = append(pubsub.Channels[protocol], channel)
}

func (pubsub *PubSub) Unsubscribe(protocol *protocol.RESPProtocol, channel string) {
	protocolIndex := -1
	subscribers := pubsub.Subscribers[channel]
	for i := range subscribers {
		if subscribers[i] == protocol {
			protocolIndex = i
			break
		}
	}

	if protocolIndex == -1 {
		return
	}

	pubsub.Subscribers[channel] = slices.Delete(pubsub.Subscribers[channel], protocolIndex, protocolIndex+1)

	channelIndex := -1
	channels := pubsub.Channels[protocol]
	for i := range channels {
		if channels[i] == channel {
			channelIndex = i
			break
		}
	}

	pubsub.Channels[protocol] = slices.Delete(pubsub.Channels[protocol], channelIndex, channelIndex+1)
}

func (pubsub *PubSub) PSubscribe(protocol *protocol.RESPProtocol, pattern string) {
	glob := glob.NewGlob(pattern)
	for channel := range pubsub.Subscribers {
		matches, err := glob.Match(channel)
		if err != nil || !matches {
			continue
		}

		pubsub.Subscribe(protocol, channel)
		break
	}
}

func (pubsub *PubSub) Publish(protocol *protocol.RESPProtocol, channel string, message string) int {
	subscribers := pubsub.Subscribers[channel]
	for _, protocol := range subscribers {
		protocol.WriteBulkString(message)
	}

	return len(subscribers)
}
