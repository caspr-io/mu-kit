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
	Publish(ctx context.Context, msgs ...proto.Message) error
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

	return r.SubscribeTopic(mh, topic)
}

func (r *MuRouter) SubscribeTopic(mh MessageHandler, topic string) error {
	log.Ctx(r.context).Info().Str("handler", mh.Name()).Str("topic", topic).Msg("Subscribing to messages")

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
func (r *MuRouter) Publish(ctx context.Context, msgs ...proto.Message) error {
	for _, protoMsg := range msgs {
		topic := r.topicName(protoMsg)

		uuid := uuid.NewV4().String()

		log.Ctx(ctx).Debug().
			Str("pub-topic", topic).
			Str("pub-uuid", uuid).
			Msg("Publishing message")
		log.Ctx(ctx).Trace().Interface("payload", protoMsg).Send()

		payloadBytes, err := proto.Marshal(protoMsg)
		if err != nil {
			return err
		}

		msg := message.NewMessage(uuid, payloadBytes)

		if err := r.publisher.Publish(topic, msg); err != nil {
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
