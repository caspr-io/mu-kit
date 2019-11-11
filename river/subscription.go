package river

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog"
)

type Subscription struct {
	handler    MuMessageHandler
	msgChannel <-chan *message.Message
	topic      string
	logger     zerolog.Logger
	running    chan struct{}
}

func (s *Subscription) Run() {
	log := s.logger.With().Str("handler", s.handler.Name()).Str("topic", s.topic).Logger()

	close(s.running)

	for m := range s.msgChannel {
		messageLogger := log.With().Str("message_uuid", m.UUID).Logger()
		messageLogger.Info().Msg("Received message...")

		protoMsg := s.handler.NewMsg()
		payload := m.Payload

		if err := proto.Unmarshal(payload, protoMsg); err != nil {
			messageLogger.Error().Err(err).Msg("Could not deserialize message payload")
			m.Nack()
		}

		c := &MessageContext{m.Context(), &messageLogger}
		if err := s.handler.Handle(c, protoMsg); err != nil {
			log.Error().Err(err).Str("uuid", m.UUID).Msg("Error handling message.")
			m.Nack()
		}

		m.Ack()
	}
}
