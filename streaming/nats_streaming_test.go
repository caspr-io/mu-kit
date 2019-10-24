package streaming

import (
	"math/rand"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/nats-io/stan.go"
	"github.com/rs/zerolog"
	"gotest.tools/v3/assert"

	"github.com/caspr-io/mu-kit/kit"
	mutest "github.com/caspr-io/mu-kit/test"
)

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

func getSubject(t *testing.T) string {
	return baseId + t.Name()
}

type StanConfig struct {
	NatsUrl       string `split_words:"true" default:"localhost:4222"`
	StanClusterId string `split_words:"true"`
}

func connect(t *testing.T, stanClientId string) stan.Conn {
	defer mutest.PatchEnvFromFile("../.env")()

	var config StanConfig
	err := envconfig.Process("", &config)
	assert.NilError(t, err)

	if config.StanClusterId == "" {
		t.Skip("STAN_CLUSTER_ID not set, skipping test")
	}

	sc, err := stan.Connect(config.StanClusterId, stanClientId, stan.NatsURL(config.NatsUrl))
	assert.NilError(t, err)

	return sc
}

func TestDoNotReceiveMessagePublishedBeforeSubscribed(t *testing.T) {
	sc := connect(t, getNextClientId())
	defer sc.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)

	subject := getSubject(t)
	message := "Hello World"

	publish(t, sc, subject, message)
	defer subscribeAndReceive(t, sc, subject, message, &wg).Unsubscribe()

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

	defer subscribeAndReceive(t, sc, subject, message, &wg).Unsubscribe()
	publish(t, sc, subject, message)

	receivedMessage := waitWithTimeout(&wg, time.Second*2)
	assert.Assert(t, receivedMessage)
}

func publish(t *testing.T, sc stan.Conn, subject string, message string) {
	err := sc.Publish(subject, []byte(message))
	assert.NilError(t, err)
}

func subscribeAndReceive(t *testing.T, sc stan.Conn, subject string, message string, wg *sync.WaitGroup) stan.Subscription {
	subs, err := sc.Subscribe(subject, func(m *stan.Msg) {
		if string(m.Data) == message {
			wg.Done()
		}
	})
	assert.NilError(t, err)
	return subs
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
