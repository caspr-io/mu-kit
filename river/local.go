package river

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/golang/protobuf/proto"
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

type ChannelMessageHandler struct {
	Received chan (proto.Message)
	newMsg   func() proto.Message
}

func (tmh *ChannelMessageHandler) NewMsg() proto.Message { return tmh.newMsg() }
func (tmh *ChannelMessageHandler) Name() string          { return "ChannelMessageHandler" }
func (tmh *ChannelMessageHandler) Handle(ctx *MessageContext, m proto.Message) error {
	tmh.Received <- m
	return nil
}

func NewChannelHandler(f func() proto.Message) *ChannelMessageHandler {
	return &ChannelMessageHandler{
		Received: make(chan (proto.Message)),
		newMsg:   f,
	}
}
