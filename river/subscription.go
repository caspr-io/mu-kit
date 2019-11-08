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

func (s *Subscription) Run() error {
	log := s.logger.With().Str("topic", s.topic).Logger()

	close(s.running)

	for m := range s.msgChannel {
		log.Info().Str("uuid", m.UUID).Msg("Received message...")

		protoMsg := s.handler.NewMsg()
		payload := m.Payload

		if err := proto.Unmarshal(payload, protoMsg); err != nil {
			return err
		}

		if err := s.handler.Handle(m.Context(), protoMsg); err != nil {
			return err
		}
	}

	return nil
}
