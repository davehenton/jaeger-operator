package util

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDnsName(t *testing.T) {
	var tests = []struct {
		in  string
		out string
	}{
		{"simplest", "simplest"},
		{"instance.with.dots-collector-headless", "instance-with-dots-collector-headless"},
		{"TestQueryDottedServiceName.With.Dots", "testquerydottedservicename-with-dots"},
		{"Service🦄", "service-z"},
		{"📈Stock-Tracker", "a-stock-tracker"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.out, DNSName(tt.in))

		matched, err := regexp.MatchString(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`, tt.out)
		assert.NoError(t, err)
		assert.True(t, matched, "%v is not a valid name", tt.out)
	}
}
