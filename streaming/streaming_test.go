package streaming

import (
	"math/rand"
	"os"
	"strconv"
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

var baseClientId string
var lastClientId int64

func init() {
	kit.Init()
	zerolog.TimestampFieldName = "t"
	zerolog.LevelFieldName = "l"
	zerolog.MessageFieldName = "m"

	rand.Seed(time.Now().UnixNano())
	baseClientId = "streaming_test-" + strconv.Itoa(os.Getpid()) + "-"
	lastClientId = 0
}

func getNextClientId() string {
	return baseClientId + strconv.FormatInt(atomic.AddInt64(&lastClientId, 1), 10)
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

}
