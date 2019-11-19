package streaming

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog"
	"github.com/satori/uuid"
)

type Publisher interface {
	Publish(msgs ...proto.Message) error
}

type MuRouter struct {
	logger        zerolog.Logger
	publisher     message.Publisher
	subscriber    message.Subscriber
	context       context.Context
	topicName     func(interface{}) string
	subscriptions []*Subscription
}

func NewRouter(
	context context.Context,
	publisher message.Publisher,
	subscriber message.Subscriber,
	logger zerolog.Logger) (*MuRouter, error) {
	return &MuRouter{logger, publisher, subscriber, context, DefaultTopicName, []*Subscription{}}, nil
}

// Subscribe subscribes a MuMessageHandler to its specific topic and will call the Handle
// function for each incoming deserialized message.
func (r *MuRouter) Subscribe(mh MessageHandler) error {
	m := mh.NewMsg()
	topic := r.topicName(m)
	r.logger.Info().Str("handler", mh.Name()).Str("topic", topic).Msg("Subscribe to messages")

	subscription, err := r.subscriber.Subscribe(r.context, topic)
	if err != nil {
		return err
	}

	subscriptionRunning := make(chan struct{})

	s := &Subscription{handler: mh, msgChannel: subscription, topic: topic, logger: r.logger, running: subscriptionRunning}
	r.subscriptions = append(r.subscriptions, s)
	go s.Run()

	<-subscriptionRunning

	return nil
}

// Publish publishes one or more messages on their respective topics.
func (r *MuRouter) Publish(msgs ...proto.Message) error {
	for _, msg := range msgs {
		topic := r.topicName(msg)

		payload, err := proto.Marshal(msg)
		if err != nil {
			return err
		}

		watermillMessage := message.NewMessage(uuid.NewV4().String(), payload)
		if err := r.publisher.Publish(topic, watermillMessage); err != nil {
			return err
		}
	}

	return nil
}

// Start starts the MuRouter in the background using a go channel
func (r *MuRouter) Start() {
}

func (r *MuRouter) Close() error {
	for _, s := range r.subscriptions {
		s.Close()
	}

	pubErr := r.publisher.Close()
	subErr := r.subscriber.Close()

	if pubErr != nil {
		return pubErr
	}

	if subErr != nil {
		return subErr
	}

	return nil
}
