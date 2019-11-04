package river

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gogo/protobuf/proto"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
)

type System struct {
	subscriber message.Subscriber
	publisher  message.Publisher
	router     *MuRouter
}

type SystemConfig struct {
	NatsURL        string `split_words:"true" required:"true"`
	NatsClusterID  string `split_words:"true" required:"true"`
	NatsClientID   string `split_words:"true" required:"true"`
	NatsQueueGroup string `split_words:"true" required:"true"`
}

func New(logger *zerolog.Logger) (*System, error) {
	config := &SystemConfig{}
	if err := envconfig.Process("MESSAGING", config); err != nil {
		return nil, err
	}

	logger.Info().Interface("config", config).Msg("Initializing messaging system...")

	watermillLogger := NewZerologLogger()

	logger.Info().Msg("Building NATS Subscriber...")

	subscriber, err := newNatsSubscriber(config, watermillLogger)
	if err != nil {
		return nil, err
	}

	logger.Info().Msg("Building NATS Publisher...")

	publisher, err := newNatsPublisher(config, watermillLogger)
	if err != nil {
		return nil, err
	}

	return NewSystem(watermillLogger, config, subscriber, publisher)
}

func NewSystem(logger watermill.LoggerAdapter, config *SystemConfig, subscriber message.Subscriber, publisher message.Publisher) (*System, error) {
	router, err := NewRouter(context.Background(), publisher, subscriber, logger)
	if err != nil {
		return nil, err
	}

	return &System{
		subscriber: subscriber,
		publisher:  publisher,
		router:     router,
	}, nil
}

func (s *System) Publish(msg proto.Message) error {
	return s.router.Publish(msg)
}

func (s *System) Subscribe(m MuMessageHandler) {

}
