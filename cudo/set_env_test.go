package cudo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetPathEnv(t *testing.T) {
	tests := []struct {
		name     string
		platform string
		expected map[string]string
	}{
		{
			name:     "Windows platform",
			platform: windows,
			expected: map[string]string{"PATH": "C:\\Windows\\System32;C:\\Windows"},
		},
		{
			name:     "macOS platform",
			platform: "darwin",
			expected: map[string]string{"PATH": "/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin"},
		},
		{
			name:     "Linux platform",
			platform: "linux",
			expected: map[string]string{"PATH": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"},
		},
		{
			name:     "Unknown platform",
			platform: "unknown",
			expected: map[string]string{"PATH": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CUSudo{Platform: tt.platform}
			result := c.SetPathEnv()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSetPWD(t *testing.T) {
	c := CUSudo{UserHome: "/home/user"}

	expected := map[string]string{"PWD": "/home/user", "CWD": "/home/user"}
	result := c.SetPWD("/home/user")

	assert.Equal(t, expected, result)
}

func TestSetEnvOpts(t *testing.T) {
	c := CUSudo{}

	opts := []map[string]string{
		{"KEY1": "value1"},
		{"KEY2": "value2"},
		{"KEY3": "value3"},
	}

	result := c.SetEnvOpts(opts...)

	// Ensure the result is not empty
	assert.NotNil(t, result)
}
