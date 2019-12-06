package kit

import (
	"testing"

	"github.com/caspr-io/mu-kit/test"
	"gotest.tools/v3/assert"
)

type EmbeddedConfig struct {
	*MuKitConfig
}

func TestShouldBeCastableToMuServerConfig(t *testing.T) {
	var c interface{} = new(EmbeddedConfig)
	_, ok := c.(MuServerConfig)
	assert.Check(t, ok)
}

func TestShouldFailWhenVariablesMissing(t *testing.T) {
	assert.Error(t, ReadConfig("TEST_SERVICE", &MuKitConfig{}), "required key TEST_SERVICE_GRPC_PORT missing value")
}

func TestShouldUpperUnderscoreConfigPrefix(t *testing.T) {
	assert.Error(t, ReadConfig("test-service", &MuKitConfig{}), "required key TEST_SERVICE_GRPC_PORT missing value")
}

func TestShouldSucceedWhenAllVariablesSet(t *testing.T) {
	defer test.PatchEnvFromFile("testenv")()
	assert.NilError(t, ReadConfig("TEST_SERVICE", &MuKitConfig{}))
}
