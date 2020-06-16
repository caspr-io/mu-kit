package streaming

import (
	"context"
	"io"

	"google.golang.org/protobuf/proto"

	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/rs/zerolog/log"
)

var _ io.Closer = (*Subscription)(nil) // compile-time check for io.Closer assignability

type Subscription struct {
	handler    MessageHandler
	msgChannel <-chan *message.Message
	topic      string
	context    context.Context
	running    chan struct{}
}

func (s *Subscription) Run() {
	close(s.running)

	for m := range s.msgChannel {
		logger := log.Ctx(s.context).With().Str("uuid", m.UUID).Logger()
		logger.Trace().
			Str("handler", s.handler.Name()).
			Str("topic", s.topic).
			Msg("Received message")

		protoMsg := s.handler.NewMsg()

		err := proto.Unmarshal(m.Payload, protoMsg)
		if err != nil {
			logger.Error().Err(err).Msg("Could not deserialize message payload")
			m.Nack()
		}

		logger.Trace().Interface("payload", protoMsg).Send()

		err = s.handler.Handle(logger.WithContext(m.Context()), protoMsg)
		if err != nil {
			logger.Error().Err(err).Msg("Could not handle message, nacking it")
			m.Nack()
		}

		logger.Trace().Msg("Handled message, acking it")
		m.Ack()
	}
}

func (s *Subscription) Close() error {
	if c, ok := s.handler.(io.Closer); ok {
		log.Ctx(s.context).Trace().Str("handler", s.handler.Name()).Msg("Closing the handler...")
		return c.Close()
	}

	return nil
}
