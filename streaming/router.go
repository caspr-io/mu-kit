package streaming

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/caspr-io/mu-kit/util"
	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog/log"
	"github.com/satori/uuid"
)

type Publisher interface {
	Publish(msgs ...proto.Message) error
}

type MuRouter struct {
	publisher     message.Publisher
	subscriber    message.Subscriber
	context       context.Context
	topicName     func(interface{}) string
	subscriptions []*Subscription
}

func NewRouter(
	context context.Context,
	publisher message.Publisher,
	subscriber message.Subscriber) (*MuRouter, error) {
	return &MuRouter{publisher, subscriber, context, DefaultTopicName, []*Subscription{}}, nil
}

// Subscribe subscribes a MuMessageHandler to its specific topic and will call the Handle
// function for each incoming deserialized message.
func (r *MuRouter) Subscribe(mh MessageHandler) error {
	m := mh.NewMsg()
	topic := r.topicName(m)
	log.Ctx(r.context).Info().Str("handler", mh.Name()).Str("topic", topic).Msg("Subscribe to messages")

	subscription, err := r.subscriber.Subscribe(r.context, topic)
	if err != nil {
		return err
	}

	subscriptionRunning := make(chan struct{})

	s := &Subscription{handler: mh, msgChannel: subscription, topic: topic, context: r.context, running: subscriptionRunning}
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
	errorCollector := new(util.ErrorCollector)
	log.Ctx(r.context).Info().Msg("Closing Router...")

	for _, s := range r.subscriptions {
		errorCollector.Collect(s.Close())
	}

	errorCollector.Collect(r.publisher.Close())
	errorCollector.Collect(r.subscriber.Close())

	if errorCollector.HasErrors() {
		return errorCollector
	}
	return nil
}
