package kit

import (
	"testing"

	"github.com/caspr-io/mu-kit/test"
	"gotest.tools/v3/assert"
)

func TestShouldFailWhenVariablesMissing(t *testing.T) {
	assert.Error(t, readConfig("TEST_SERVICE"), "required key TEST_SERVICE_GRPC_PORT missing value")
}

func TestShouldUpperUnderscoreConfigPrefix(t *testing.T) {
	assert.Error(t, readConfig("test-service"), "required key TEST_SERVICE_GRPC_PORT missing value")
}

func TestShouldSucceedWhenAllVariablesSet(t *testing.T) {
	defer test.PatchEnvFromFile("testenv")()
	assert.NilError(t, readConfig("TEST_SERVICE"))
}
