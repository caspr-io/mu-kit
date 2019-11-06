package river

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

func NewTestRiver() (*SubSystem, error) {
	logger := NewZerologLogger()
	pub, sub := goChannelPubSub(logger)

	return NewSubSystem(logger, sub, pub)
}

func goChannelPubSub(logger watermill.LoggerAdapter) (message.Publisher, message.Subscriber) {
	pubSub := gochannel.NewGoChannel(gochannel.Config{
		OutputChannelBuffer: int64(1000),
	}, logger)

	return pubSub, pubSub
}
