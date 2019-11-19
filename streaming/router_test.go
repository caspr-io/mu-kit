package streaming

import (
	"context"
	"sync"
	"testing"

	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog/log"
	"gotest.tools/v3/assert"
)

type TestMessageHandler struct {
	wg       *sync.WaitGroup
	messages []TestMessage
}

func (th *TestMessageHandler) Name() string          { return "TestMessageHandler" }
func (th *TestMessageHandler) NewMsg() proto.Message { return &TestMessage{} }
func (th *TestMessageHandler) Handle(ctx *MessageContext, m proto.Message) error {
	defer th.wg.Done()
	th.messages = append(th.messages, *m.(*TestMessage))

	return nil
}

func TestShouldReceiveOnePublishedMessage(t *testing.T) {
	router := makeRouter(t)

	var wg sync.WaitGroup

	tmh := TestMessageHandler{wg: &wg}

	wg.Add(1)

	if err := router.Subscribe(&tmh); err != nil {
		t.Error(err)
	}

	router.Start()
	defer router.Close()

	if err := router.Publish(&TestMessage{Contents: "test message"}); err != nil {
		t.Error(err)
	}

	wg.Wait()
	assert.Assert(t, len(tmh.messages) == 1)
}

func TestShouldReceiveAllPublishedMessage(t *testing.T) {
	router := makeRouter(t)

	var wg sync.WaitGroup

	tmh := TestMessageHandler{wg: &wg}

	wg.Add(3)

	if err := router.Subscribe(&tmh); err != nil {
		t.Error(err)
	}

	router.Start()
	defer router.Close()

	if err := router.Publish(&TestMessage{Contents: "test message"}, &TestMessage{Contents: "2nd message"}, &TestMessage{Contents: "third message"}); err != nil {
		t.Error(err)
	}

	wg.Wait()
	assert.Assert(t, len(tmh.messages) == 3)
}

func makeRouter(t *testing.T) *MuRouter {
	logger := NewZerologLogger(&log.Logger)
	pubSub := gochannel.NewGoChannel(gochannel.Config{}, logger)

	router, err := NewRouter(context.Background(), pubSub, pubSub, log.Logger)
	if err != nil {
		t.Error(err)
	}

	return router
}
