package test

import (
	"os"
	"testing"

	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
	"gotest.tools/v3/fs"
)

func TestShouldSetEnvironmentVariableFromFile(t *testing.T) {
	assert.Assert(t, !is.Contains(os.Environ(), "MUTEST=test.1")().Success())
	file := fs.NewFile(t, "env_test", fs.WithContent("MUTEST=test.1\n"))
	defer file.Remove()
	defer PatchEnvFromFile(file.Path())()

	assert.Assert(t, is.Contains(os.Environ(), "MUTEST=test.1"))
}
