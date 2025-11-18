package cudo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetElevationPrefix(t *testing.T) {
	tests := []struct {
		platform    string
		currentUser string
		expected    []string
	}{
		{
			platform:    "windows",
			currentUser: "testuser",
			expected:    []string{"runas", "/user:testuser"},
		},
		{
			platform:    "darwin",
			currentUser: "testuser",
			expected:    []string{"sudo", "-H", "-u", "testuser"},
		},
		{
			platform:    "linux",
			currentUser: "testuser",
			expected:    []string{"sudo", "-H", "-u", "testuser"},
		},
	}

	for _, test := range tests {
		c := CUSudo{
			Platform:    test.platform,
			CurrentUser: test.currentUser,
		}

		result, err := c.GetElevationPrefix()
		assert.NoError(t, err)
		assert.Equal(t, test.expected, result)
	}
}

func TestBuildCmd(t *testing.T) {
	tests := []struct {
		platform    string
		currentUser string
		args        []string
		expected    []string
	}{
		{
			platform:    "windows",
			currentUser: "testuser",
			args:        []string{"echo", "hello"},
			expected:    []string{"runas", "/user:testuser", "echo", "hello"},
		},
		{
			platform:    "darwin",
			currentUser: "testuser",
			args:        []string{"echo", "hello"},
			expected:    []string{"sudo", "-H", "-u", "testuser", "echo", "hello"},
		},
		{
			platform:    "linux",
			currentUser: "testuser",
			args:        []string{"echo", "hello"},
			expected:    []string{"sudo", "-H", "-u", "testuser", "echo", "hello"},
		},
	}

	for _, test := range tests {
		c := CUSudo{
			Platform:    test.platform,
			CurrentUser: test.currentUser,
		}

		cmd, err := c.BuildCmd(test.args)
		assert.NoError(t, err)
		assert.Equal(t, test.expected, cmd.Args)
	}
}
