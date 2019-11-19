package streaming

import (
	"testing"

	"gotest.tools/v3/assert"
)

type TestStruct struct{}

func TestReturnNameOfPointer(t *testing.T) {
	subject := &TestStruct{}

	assert.Equal(t, DefaultTopicName(subject), "test.struct")
}

func TestReturnNameOfValue(t *testing.T) {
	subject := TestStruct{}

	assert.Equal(t, DefaultTopicName(subject), "test.struct")
}
