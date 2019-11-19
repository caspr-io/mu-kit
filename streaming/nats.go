package streaming

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/nats-io/stan.go"
)

func connectToStan(config *Config) (*stan.Conn, error) {
	conn, err := stan.Connect(config.NatsClusterID, config.NatsClientID, stan.NatsURL(config.NatsURL))

	if err != nil {
		return nil, err
	}

	return &conn, nil
}

func newNatsPublisher(config *Config, conn *stan.Conn, logger watermill.LoggerAdapter) (message.Publisher, error) {
	publisher, err := nats.NewStreamingPublisherWithStanConn(*conn, nats.StreamingPublisherPublishConfig{
		Marshaler: nats.GobMarshaler{},
	}, logger)
	if err != nil {
		return nil, err
	}

	return publisher, nil
}

func newNatsSubscriber(config *Config, conn *stan.Conn, logger watermill.LoggerAdapter) (message.Subscriber, error) {
	subscriber, err := nats.NewStreamingSubscriberWithStanConn(*conn, nats.StreamingSubscriberSubscriptionConfig{
		QueueGroup:  config.NatsQueueGroup,
		Unmarshaler: nats.GobMarshaler{},
	}, logger)
	if err != nil {
		return nil, err
	}

	return subscriber, nil
}
