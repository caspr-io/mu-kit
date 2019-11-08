package river

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gogo/protobuf/proto"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type SubSystem struct {
	subscriber message.Subscriber
	publisher  message.Publisher
	router     *MuRouter
	logger     zerolog.Logger
}

type SubSystemConfig struct {
	NatsURL        string `split_words:"true" required:"true"`
	NatsClusterID  string `split_words:"true" required:"true"`
	NatsClientID   string `split_words:"true" required:"true"`
	NatsQueueGroup string `split_words:"true" required:"true"`
}

func New(config *SubSystemConfig) (*SubSystem, error) {
	logger := log.Logger.With().Str("component", "µ-kit Streaming").Logger()

	logger.Info().Interface("config", config).Msg("Initializing µ-Kit Messaging subsystem...")

	watermillLogger := NewZerologLogger(&logger)

	stanConn, err := connectToStan(config)
	if err != nil {
		return nil, err
	}

	logger.Info().Msg("Building NATS Subscriber...")

	subscriber, err := newNatsSubscriber(config, stanConn, watermillLogger)
	if err != nil {
		return nil, err
	}

	logger.Info().Msg("Building NATS Publisher...")

	publisher, err := newNatsPublisher(config, stanConn, watermillLogger)
	if err != nil {
		return nil, err
	}

	return NewSubSystem(logger, watermillLogger, subscriber, publisher)
}

func NewSubSystem(logger zerolog.Logger, watermillLogger watermill.LoggerAdapter, subscriber message.Subscriber, publisher message.Publisher) (*SubSystem, error) {
	router, err := NewRouter(context.Background(), publisher, subscriber, logger)
	if err != nil {
		return nil, err
	}

	return &SubSystem{
		logger:     logger,
		subscriber: subscriber,
		publisher:  publisher,
		router:     router,
	}, nil
}

func (s *SubSystem) Publish(msg proto.Message) error {
	return s.router.Publish(msg)
}

func (s *SubSystem) Subscribe(m MuMessageHandler) error {
	return s.router.Subscribe(m)
}

func (s *SubSystem) Run() {
	s.router.Start()
}

func (s *SubSystem) Close() error {
	s.router.Close()
	return nil
}
