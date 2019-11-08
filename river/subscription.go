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
		log.Info().Str("uuid", m.UUID).Msg("Received message...")

		protoMsg := s.handler.NewMsg()
		payload := m.Payload

		if err := proto.Unmarshal(payload, protoMsg); err != nil {
			log.Error().Err(err).Str("uuid", m.UUID).Msg("Could not deserialize message payload")
			m.Nack()
		}

		if err := s.handler.Handle(m.Context(), protoMsg); err != nil {
			log.Error().Err(err).Str("uuid", m.UUID).Msg("Error handling message.")
			m.Nack()
		}

		m.Ack()
	}
}
