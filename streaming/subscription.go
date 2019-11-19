package streaming

import (
	"context"
	"io"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/golang/protobuf/proto"
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
	log := log.Ctx(s.context).With().Str("handler", s.handler.Name()).Logger()

	close(s.running)

	for m := range s.msgChannel {
		messageLogger := log.With().Str("message_uuid", m.UUID).Logger()
		topicLogger := messageLogger.With().Str("topic", s.topic).Logger()
		topicLogger.Info().Msg("Received message...")

		protoMsg := s.handler.NewMsg()
		payload := m.Payload

		if err := proto.Unmarshal(payload, protoMsg); err != nil {
			topicLogger.Error().Err(err).Msg("Could not deserialize message payload")
			m.Nack()
		}

		c := messageLogger.WithContext(m.Context())
		if err := s.handler.Handle(c, protoMsg); err != nil {
			topicLogger.Error().Err(err).Str("uuid", m.UUID).Msg("Error handling message.")
			m.Nack()
		}

		m.Ack()
	}
}

func (s *Subscription) Close() error {
	if c, ok := s.handler.(io.Closer); ok {
		log.Ctx(s.context).Info().Str("handler", s.handler.Name()).Msg("Closing the handler...")
		return c.Close()
	}
	return nil
}
