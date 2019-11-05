package river

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/nats-io/stan.go"
)

func newNatsPublisher(config *SystemConfig, logger watermill.LoggerAdapter) (message.Publisher, error) {
	publisher, err := nats.NewStreamingPublisher(nats.StreamingPublisherConfig{
		ClusterID: config.NatsClusterID,
		ClientID:  config.NatsClientID,
		StanOptions: []stan.Option{
			stan.NatsURL(config.NatsURL),
		},
		Marshaler: nats.GobMarshaler{},
	}, logger)
	if err != nil {
		return nil, err
	}

	return publisher, nil
}

func newNatsSubscriber(config *SystemConfig, logger watermill.LoggerAdapter) (message.Subscriber, error) {
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
