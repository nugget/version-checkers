package zonedirector

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExpandedVersion(t *testing.T) {
	var testTable = []struct {
		input       string // input
		expected    string // expected result
		description string // what we are testing
	}{
		{"10.20.30.40.50", "10.20.30.40 build 50", "valid version string"},
		{"10.20.30", "10.20.30", "Unrecognized version string"},
	}

	for _, tt := range testTable {
		actual := ExpandedVersion(tt.input)
		assert.Equal(t, actual, tt.expected, "Unexpected result on "+tt.description)
	}
}
