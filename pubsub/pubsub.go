package pubsub

type PubSub struct {
}

var pubsub = &PubSub{}

func GetPubSub() *PubSub {
	return pubsub
}

func (pubsub *PubSub) Subscribe(channel string) {

}
