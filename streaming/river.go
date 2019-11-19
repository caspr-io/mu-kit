package streaming

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gogo/protobuf/proto"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	NatsURL        string `split_words:"true" required:"true"`
	NatsClusterID  string `split_words:"true" required:"true"`
	NatsClientID   string `split_words:"true" required:"true"`
	NatsQueueGroup string `split_words:"true" required:"true"`
}

type River struct {
	subscriber message.Subscriber
	publisher  message.Publisher
	router     *MuRouter
	logger     zerolog.Logger
}

func NewRiver(config *Config) (*River, error) {
	logger := log.Logger.With().Str("component", "streaming").Logger()

	logger.Info().Interface("config", config).Msg("Initializing µ-Kit Streaming system...")

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

	return CreateRiver(logger, subscriber, publisher)
}

func CreateRiver(logger zerolog.Logger, subscriber message.Subscriber, publisher message.Publisher) (*River, error) {
	ctx := logger.WithContext(context.Background())
	router, err := NewRouter(ctx, publisher, subscriber)
	if err != nil {
		return nil, err
	}

	return &River{
		logger:     logger,
		subscriber: subscriber,
		publisher:  publisher,
		router:     router,
	}, nil
}

func (s *River) Publish(msg proto.Message) error {
	return s.router.Publish(msg)
}

func (s *River) Publisher() Publisher {
	return s.router
}

func (s *River) Subscribe(m MessageHandler) error {
	return s.router.Subscribe(m)
}

func (s *River) Run() {
	s.router.Start()
}

func (s *River) Close() error {
	s.router.Close()
	return nil
}
