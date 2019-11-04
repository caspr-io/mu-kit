package messaging

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	mumill "github.com/caspr-io/mu-kit/watermill"
	"github.com/kelseyhightower/envconfig"
	"github.com/nats-io/stan.go"
	"github.com/rs/zerolog"
)

type System struct {
	subscriber message.Subscriber
	publisher  message.Publisher
	router     *mumill.MuRouter
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

	watermillLogger := mumill.NewZerologLogger()

	logger.Info().Msg("Building NATS Subscriber...")

	subscriber, err := NewNatsSubscriber(config, watermillLogger)
	if err != nil {
		return nil, err
	}

	logger.Info().Msg("Building NATS Publisher...")

	publisher, err := NewNatsPublisher(config, watermillLogger)
	if err != nil {
		return nil, err
	}

	return NewSystem(watermillLogger, config, subscriber, publisher)
}

func NewSystem(logger watermill.LoggerAdapter, config *SystemConfig, subscriber message.Subscriber, publisher message.Publisher) (*System, error) {
	router, err := mumill.NewRouter(publisher, subscriber, context.Background(), logger)
	if err != nil {
		return nil, err
	}

	return &System{
		subscriber: subscriber,
		publisher:  publisher,
		router:     router,
	}, nil
}

func NewNatsPublisher(config *SystemConfig, logger watermill.LoggerAdapter) (message.Publisher, error) {
	publisher, err := nats.NewStreamingPublisher(nats.StreamingPublisherConfig{
		ClusterID: config.NatsClusterID,
		ClientID:  config.NatsClientID,
		StanOptions: []stan.Option{
			stan.NatsURL(config.NatsURL),
		},
	}, logger)
	if err != nil {
		return nil, err
	}

	return publisher, nil
}

func NewNatsSubscriber(config *SystemConfig, logger watermill.LoggerAdapter) (message.Subscriber, error) {
	subscriber, err := nats.NewStreamingSubscriber(nats.StreamingSubscriberConfig{
		ClusterID:  config.NatsClusterID,
		ClientID:   config.NatsClientID,
		QueueGroup: config.NatsQueueGroup,
		StanOptions: []stan.Option{
			stan.NatsURL(config.NatsURL),
		},
		Unmarshaler: nats.GobMarshaler{},
	}, logger)
	if err != nil {
		return nil, err
	}

	return subscriber, nil
}
