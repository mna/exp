package redpubsub

import "github.com/PuerkitoBio/exp/juggler/msg"

// TODO : to move elsewhere...
type PubSuber interface { // or something...
	Publish(msg.Pub)
	Subscribe(msg.Sub)
	Unsubscribe(msg.Unsb)
}
