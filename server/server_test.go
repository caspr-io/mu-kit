package server

import (
	"os"
	"testing"

	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

func TestFailServerNewWithoutEnvSet(t *testing.T) {
	os.Setenv("MUKIT_GRPC_PORT", "")
	assert.Assert(t, is.Len(os.Getenv("MUKIT_GRPC_PORT"), 0))

	_, err := New()
	if err == nil {
		t.Error("Should fail with missing Environment")
	}
}

func TestRunServerWithEnv(t *testing.T) {
	os.Setenv("MUKIT_GRPC_PORT", "8080")
	_, err := New()
	assert.NilError(t, err)
}
