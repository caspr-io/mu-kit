package watermill

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/rs/zerolog"
	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"

	"github.com/caspr-io/mu-kit/kit"
)

func Not(c is.Comparison) is.Comparison {
	return func() is.Result {
		r := c()
		if r.Success() {
			return is.ResultFailure("Failed due to success")
		} else {
			return is.ResultSuccess
		}
	}
}

func TestShouldLogDebug(t *testing.T) {
	b, l := InitTest()
	l.Debug("testing message", nil)
	assert.Assert(t, is.Contains(b.String(), "\"l\":\"debug\""))
	assert.Assert(t, Not(is.Contains(b.String(), "\"foo\":\"bar\"")))
	assert.Assert(t, Not(is.Contains(b.String(), "\"trace\":\"true\"")))
	assert.Assert(t, is.Contains(b.String(), "\"m\":\"testing message\""))
}

func TestShouldLogDebugWithExtraFields(t *testing.T) {
	b, l := InitTest()
	fields := watermill.LogFields{"foo": "bar"}
	l.Debug("testing message", fields)
	assert.Assert(t, is.Contains(b.String(), "\"l\":\"debug\""))
	assert.Assert(t, is.Contains(b.String(), "\"foo\":\"bar\""))
	assert.Assert(t, is.Contains(b.String(), "\"m\":\"testing message\""))
}

func TestShouldLogWithFields(t *testing.T) {
	b, l := InitTest()
	fields := watermill.LogFields{"foo": "bar"}
	zl := l.With(fields)
	zl.Debug("testing message", nil)
	assert.Assert(t, is.Contains(b.String(), "\"l\":\"debug\""))
	assert.Assert(t, is.Contains(b.String(), "\"foo\":\"bar\""))
	assert.Assert(t, is.Contains(b.String(), "\"m\":\"testing message\""))

}

func TestShouldLogTrace(t *testing.T) {
	b, l := InitTest()
	l.Trace("testing message", nil)
	assert.Assert(t, is.Contains(b.String(), "\"l\":\"debug\""))
	assert.Assert(t, is.Contains(b.String(), "\"trace\":\"true\""))
	assert.Assert(t, is.Contains(b.String(), "\"m\":\"testing message\""))
}

func TestShouldLogInfo(t *testing.T) {
	b, l := InitTest()
	l.Info("testing message", nil)
	assert.Assert(t, is.Contains(b.String(), "\"l\":\"info\""))
	assert.Assert(t, Not(is.Contains(b.String(), "\"trace\":\"true\"")))
	assert.Assert(t, is.Contains(b.String(), "\"m\":\"testing message\""))

}
func TestShouldLogError(t *testing.T) {
	b, l := InitTest()
	l.Error("testing message", fmt.Errorf("testing error"), nil)
	assert.Assert(t, is.Contains(b.String(), "\"l\":\"error\""))
	assert.Assert(t, is.Contains(b.String(), "\"error\":\"testing error\""))
	assert.Assert(t, Not(is.Contains(b.String(), "\"trace\":\"true\"")))
	assert.Assert(t, is.Contains(b.String(), "\"m\":\"testing message\""))
}

func TestShouldNotLogTraceIfDisabled(t *testing.T) {
	b, l := InitTest()
	l.traceEnabled = false
	l.Trace("testing message", nil)
	assert.Assert(t, is.Len(b.String(), 0))
}

func InitTest() (*strings.Builder, *ZeroLogger) {
	kit.Init()
	b := strings.Builder{}
	l := zerolog.New(&b)
	zl := ZeroLogger{true, &l}
	return &b, &zl
}
