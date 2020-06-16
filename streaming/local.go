package streaming

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

func NewTestRiver() (*River, error) {
	logger := log.Logger.With().Str("component", "streaming").Logger()

	watermillLogger := NewZerologLogger(&logger)

	pub, sub := goChannelPubSub(watermillLogger)

	return CreateRiver(logger, sub, pub)
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
func (tmh *ChannelMessageHandler) Handle(ctx context.Context, m proto.Message) error {
	log.Ctx(ctx).Trace().Interface("message", m).Send()
	tmh.Received <- m

	return nil
}

func NewChannelHandler(f func() proto.Message) *ChannelMessageHandler {
	return &ChannelMessageHandler{
		Received: make(chan (proto.Message)),
		newMsg:   f,
	}
}

func (tmh *ChannelMessageHandler) Close() error {
	close(tmh.Received)
	return nil
}
