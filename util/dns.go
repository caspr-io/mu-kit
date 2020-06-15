package util

import (
	"regexp"
)

var DnsDomainNameRe *regexp.Regexp
var DnsLabelNameRe *regexp.Regexp

func init() {
	DnsDomainNameRe = regexp.MustCompile("^[a-z0-9]([a-z0-9]|[a-z0-9.-]*[a-z0-9])?$")
	DnsLabelNameRe = regexp.MustCompile("^[a-z0-9]([a-z0-9]|[a-z0-9-]*[a-z0-9])?$")
}

func isDnsSubdomainName(name string) bool {
	if len(name) > 253 {
		return false
	}

	return DnsDomainNameRe.MatchString(name)
}

func isDnsLabelName(name string) bool {
	if len(name) > 63 {
		return false
	}
	return DnsLabelNameRe.MatchString(name)
}
