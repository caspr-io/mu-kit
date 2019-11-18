package kit

import (
	"testing"

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
