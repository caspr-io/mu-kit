package util

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestDnsSubdomainName(t *testing.T) {
	assert.Assert(t, isDnsSubdomainName("a"))
	assert.Assert(t, isDnsSubdomainName("9"))
	assert.Assert(t, isDnsSubdomainName("aa"))
	assert.Assert(t, isDnsSubdomainName("a9"))
	assert.Assert(t, isDnsSubdomainName("a.a"))
	assert.Assert(t, isDnsSubdomainName("a-a"))
	assert.Assert(t, isDnsSubdomainName("0abcdefghi1abcdefghi2abcdefghi3abcdefghi4abcdefghi5abcdefghi6abcdefghi7abcdefghi8abcdefghi9abcdefghi0abcdefghi1abcdefghi2abcdefghi3abcdefghi4abcdefghi5abcdefghi6abcdefghi7abcdefghi8abcdefghi9abcdefghi0abcdefghi1abcdefghi2abcdefghi3abcdefghi4abcdefghi5ab"))
	assert.Assert(t, !isDnsSubdomainName("0abcdefghi1abcdefghi2abcdefghi3abcdefghi4abcdefghi5abcdefghi6abcdefghi7abcdefghi8abcdefghi9abcdefghi0abcdefghi1abcdefghi2abcdefghi3abcdefghi4abcdefghi5abcdefghi6abcdefghi7abcdefghi8abcdefghi9abcdefghi0abcdefghi1abcdefghi2abcdefghi3abcdefghi4abcdefghi5abc"))
	assert.Assert(t, !isDnsSubdomainName("A"))
	assert.Assert(t, !isDnsSubdomainName("a/a"))
}
func TestDnsLabelName(t *testing.T) {
	assert.Assert(t, isDnsLabelName("a"))
	assert.Assert(t, isDnsLabelName("9"))
	assert.Assert(t, isDnsLabelName("aa"))
	assert.Assert(t, isDnsLabelName("a9"))
	assert.Assert(t, !isDnsLabelName("a.a"))
	assert.Assert(t, isDnsLabelName("a-a"))
	assert.Assert(t, isDnsLabelName("0abcdefghi1abcdefghi2abcdefghi3abcdefghi4abcdefghi5abcdefghi6ab"))
	assert.Assert(t, !isDnsLabelName("0abcdefghi1abcdefghi2abcdefghi3abcdefghi4abcdefghi5abcdefghi6abc"))
	assert.Assert(t, !isDnsLabelName("A"))
	assert.Assert(t, !isDnsLabelName("a/a"))
}
