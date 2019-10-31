package watermill

import (
	"context"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/golang/protobuf/proto"
	"github.com/satori/uuid"
)

type MuRouter struct {
	r          *message.Router
	publisher  message.Publisher
	subscriber message.Subscriber
	context    context.Context
	topicName  func(interface{}) string
}

type MuMessageHandler interface {
	Name() string
	NewMsg() proto.Message
	Handle(ctx context.Context, m proto.Message) error
}

func NewRouter(
	publisher message.Publisher,
	subscriber message.Subscriber,
	context context.Context,
	logger watermill.LoggerAdapter) (*MuRouter, error) {

	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return nil, err
	}

	router.AddPlugin(plugin.SignalsHandler)
	router.AddMiddleware(
		middleware.CorrelationID,
		middleware.Retry{
			MaxRetries:      3,
			InitialInterval: time.Millisecond * 100,
			Logger:          logger,
		}.Middleware,
		middleware.Recoverer,
	)

	// mup := &MuPublisher{publisher: publisher, topicName: DefaultTopicName}
	// mus := &MuSubscriber{subscriber: subscriber, context: context, topicName: DefaultTopicName}
	return &MuRouter{router, publisher, subscriber, context, DefaultTopicName}, nil
}

func (r *MuRouter) Subscribe(mh MuMessageHandler) error {
	m := mh.NewMsg()
	topic := r.topicName(m)

	r.r.AddNoPublisherHandler(mh.Name(), topic, r.subscriber, func(msg *message.Message) error {
		protoMsg := mh.NewMsg()

		if err := proto.Unmarshal(msg.Payload, protoMsg); err != nil {
			return err
		}

		if err := mh.Handle(msg.Context(), protoMsg); err != nil {
			return err
		}

		return nil
	})

	return nil
}

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
