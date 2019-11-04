package test

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/caspr-io/mu-kit/river"
)

func NewTestRiver() (*river.System, error) {
	logger := river.NewZerologLogger()
	pub, sub := goChannelPubSub(logger)

	return river.NewSystem(logger, sub, pub)
}

func goChannelPubSub(logger watermill.LoggerAdapter) (message.Publisher, message.Subscriber) {
	pubSub := gochannel.NewGoChannel(gochannel.Config{
		OutputChannelBuffer: int64(1000),
	}, logger)

	return pubSub, pubSub
}
