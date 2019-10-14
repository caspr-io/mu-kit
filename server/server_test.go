package server

import (
	"os"
	"testing"

	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
	"gotest.tools/v3/env"
)

func TestFailServerNewWithoutEnvSet(t *testing.T) {
	assert.Assert(t, is.Len(os.Getenv("MUKIT_GRPC_PORT"), 0))

	_, err := New()
	if err == nil {
		t.Error("Should fail with missing Environment")
	}
}

func TestRunServerWithEnv(t *testing.T) {
	defer env.Patch(t, "MUKIT_GRPC_PORT", "8080")()
	_, err := New()
	assert.NilError(t, err)
}
