package river

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/rs/zerolog/log"
)

func NewTestRiver() (*SubSystem, error) {
	logger := log.Logger.With().Str("component", "µ-kit Streaming").Logger()

	watermillLogger := NewZerologLogger(&logger)

	pub, sub := goChannelPubSub(watermillLogger)

	return NewSubSystem(logger, watermillLogger, sub, pub)
}

func goChannelPubSub(logger watermill.LoggerAdapter) (message.Publisher, message.Subscriber) {
	pubSub := gochannel.NewGoChannel(gochannel.Config{
		OutputChannelBuffer: int64(1000),
	}, logger)

	return pubSub, pubSub
}
