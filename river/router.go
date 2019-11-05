package river

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
	context context.Context,
	publisher message.Publisher,
	subscriber message.Subscriber,
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

	return &MuRouter{router, publisher, subscriber, context, DefaultTopicName}, nil
}

// Subscribe subscribes a MuMessageHandler to its specific topic and will call the Handle
// function for each incoming deserialized message.
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
	go r.r.Run(r.context)
	<-r.r.Running()
}

func (r *MuRouter) Close() error {
	return r.r.Close()
}
