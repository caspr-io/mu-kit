package watermill

import (
	"context"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/kelseyhightower/envconfig"
	"github.com/nats-io/stan.go"
	"github.com/rs/zerolog"
	"gotest.tools/v3/assert"

	"github.com/caspr-io/mu-kit/kit"
	mutest "github.com/caspr-io/mu-kit/test"
)

func TestDoNotReceiveMessagePublishedBeforeSubscribed(t *testing.T) {
	sc := connect(t, getNextClientId())
	defer sc.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)

	subject := getSubject(t)
	message := "Hello World"

	publish(t, sc, subject, message)
	defer subscribeAndReceive(t, sc, subject, message, &wg).Close()

	receivedMessage := waitWithTimeout(&wg, time.Second*2)
	assert.Assert(t, !receivedMessage)
}

func TestDoReceiveMessagePublishedAfterSubscribed(t *testing.T) {
	sc := connect(t, getNextClientId())
	defer sc.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)

	subject := getSubject(t)
	message := "Hello World"

	defer subscribeAndReceive(t, sc, subject, message, &wg).Close()
	publish(t, sc, subject, message)

	receivedMessage := waitWithTimeout(&wg, time.Second*2)
	assert.Assert(t, receivedMessage)

}

func TestDoReceiveMessagePublishedBeforeSubscribedDurably(t *testing.T) {
	sc := connect(t, getNextClientId())
	defer sc.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)

	subject := getSubject(t)
	message := "Hello World"

	defer subscribeAndReceive(t, sc, subject, message, &wg, stan.DurableName(subject+"-durable")).Close()
	publish(t, sc, subject, message)

	receivedMessage := waitWithTimeout(&wg, time.Second*2)
	assert.Assert(t, receivedMessage)
}

var baseId string
var lastClientId int64

func init() {
	kit.Init()
	zerolog.TimestampFieldName = "t"
	zerolog.LevelFieldName = "l"
	zerolog.MessageFieldName = "m"

	rand.Seed(time.Now().UnixNano())
	baseId = "nats_streaming_test-" + strconv.Itoa(os.Getpid()) + "-"
	lastClientId = 0
}

func getNextClientId() string {
	return baseId + strconv.FormatInt(atomic.AddInt64(&lastClientId, 1), 10)
}

func connect(t *testing.T, stanClientId string) stan.Conn {
	defer mutest.PatchEnvFromFile("../../.env")()

	var config struct {
		NatsUrl       string `split_words:"true" default:"localhost:4222"`
		StanClusterId string `split_words:"true"`
	}
	err := envconfig.Process("", &config)
	assert.NilError(t, err)

	if config.StanClusterId == "" {
		t.Skip("STAN_CLUSTER_ID not set, skipping test")
	}

	sc, err := stan.Connect(config.StanClusterId, stanClientId, stan.NatsURL(config.NatsUrl))
	assert.NilError(t, err)

	return sc
}

func getSubject(t *testing.T) string {
	return baseId + t.Name()
}

func publish(t *testing.T, sc stan.Conn, subject string, payload string) {
	publisher, err := nats.NewStreamingPublisherWithStanConn(sc,
		nats.StreamingPublisherPublishConfig{
			Marshaler: nats.GobMarshaler{},
		},
		watermill.NewStdLogger(false, false),
	)
	assert.NilError(t, err)

	err = publisher.Publish(subject, message.NewMessage(watermill.NewUUID(), []byte(payload)))
	assert.NilError(t, err)
}

func subscribeAndReceive(t *testing.T, sc stan.Conn, subject string, expectedPayload string, wg *sync.WaitGroup, opts ...stan.SubscriptionOption) *nats.StreamingSubscriber {
	subscriber, err := nats.NewStreamingSubscriberWithStanConn(sc,
		nats.StreamingSubscriberSubscriptionConfig{
			// SubscribersCount and CloseTimeout need to be set explicitly.
			// `nats.NewStreamingSubscriber` will set default values, but `nats.NewStreamingSubscriberWithStanConn` won't.
			SubscribersCount: 1,
			CloseTimeout:     time.Second * 30,
			Unmarshaler:      nats.GobMarshaler{},
		},
		watermill.NewStdLogger(false, false),
	)
	assert.NilError(t, err)

	messages, err := subscriber.Subscribe(context.Background(), subject)
	assert.NilError(t, err)

	go func(messages <-chan *message.Message) {
		for msg := range messages {
			if string(msg.Payload) == expectedPayload {
				wg.Done()
			}
		}
	}(messages)

	return subscriber
}

func waitWithTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()

	select {
	case <-c:
		return true // completed normally
	case <-time.After(timeout):
		return false // timed out
	}
}
