package id

import (
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestKitIDStartsWithPrefix(t *testing.T) {
	id := New("pref-")
	assert.Assert(t, strings.HasPrefix(id, "pref-"))
}

func TestKitIDIs30CharactersLong(t *testing.T) {
	id := New("pref-")
	assert.Equal(t, len(id), 30)
}
