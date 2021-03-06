package streaming

import (
	"context"

	"google.golang.org/protobuf/proto"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/caspr-io/mu-kit/util"

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
	log.Ctx(r.context).Trace().
		Str("handler", mh.Name()).
		Str("topic", topic).
		Msg("Subscribing to messages")

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

		if err := r.PublishTopic(ctx, topic, protoMsg); err != nil {
			return err
		}
	}

	return nil
}

// PublishTopic publishes a protobuf message to a topic
func (r *MuRouter) PublishTopic(ctx context.Context, topic string, protoMsg proto.Message) error {
	uuid := uuid.NewV4().String()

	log.Ctx(ctx).Trace().
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

	return nil
}

// Start starts the MuRouter in the background using a go channel
func (r *MuRouter) Start() {
}

func (r *MuRouter) Close() error {
	errorCollector := new(util.ErrorCollector)

	log.Ctx(r.context).Debug().Msg("Closing µ-Kit Streaming system...")

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
