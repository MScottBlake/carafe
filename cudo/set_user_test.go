package cudo

import (
	"fmt"
	"os"
	"os/user"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetConsoleUser(t *testing.T) {
	tests := []struct {
		platform   string
		envVars    map[string]string
		statResult os.FileInfo
		statErr    error
		lookupUser *user.User
		lookupErr  error
		expected   string
		expectErr  bool
	}{
		{
			platform: "windows",
			envVars:  map[string]string{"USERNAME": "testuser"},
			expected: "testuser",
		},
		{
			platform:  "windows",
			envVars:   map[string]string{"USERNAME": ""},
			expectErr: true,
		},
		{
			platform: "darwin",
			statResult: &MockFileInfo{
				sys: &syscall.Stat_t{Uid: 501},
			},
			lookupUser: &user.User{Username: "testuser"},
			expected:   "testuser",
		},
		{
			platform:  "darwin",
			statErr:   fmt.Errorf("stat error"),
			expectErr: true,
		},
		{
			platform: "linux",
			envVars:  map[string]string{"SUDO_USER": "testuser"},
			expected: "testuser",
		},
		{
			platform: "linux",
			envVars:  map[string]string{"USER": "testuser"},
			expected: "testuser",
		},
		{
			platform:  "linux",
			envVars:   map[string]string{"SUDO_USER": "", "USER": ""},
			expectErr: true,
		},
	}

	for _, test := range tests {
		mockOSFunc := &MockOSFunc{
			envVars:    test.envVars,
			statResult: test.statResult,
			statErr:    test.statErr,
			lookupUser: test.lookupUser,
			lookupErr:  test.lookupErr,
		}

		c := &CUSudo{
			Platform: test.platform,
			OSFunc:   mockOSFunc,
		}

		// Run the test
		err := c.SetConsoleUser()
		if test.expectErr {
			assert.Error(t, err, "expected error for platform %s, got nil", test.platform)
		} else {
			assert.NoError(t, err, "unexpected error for platform %s: %v", test.platform, err)
			assert.Equal(t, test.expected, c.CurrentUser, "unexpected current user for platform %s", test.platform)
		}
	}
}

func TestCheckForRoot(t *testing.T) {

	tests := []struct {
		gid       int
		expectErr bool
		platform  string
	}{
		{
			gid:      0,
			platform: "linux",
		},
		{
			gid:       1000,
			expectErr: true,
			platform:  "linux",
		},
	}

	for _, test := range tests {
		c := CUSudo{
			Platform: test.platform,
			OSFunc:   OSFunc(&MockOSFunc{egid: test.gid}),
		}

		err := c.CheckForRoot()
		if test.expectErr {
			assert.Error(t, err, "expected error for GID %d, got nil", test.gid)
		} else {
			assert.NoError(t, err, "unexpected error for GID %d: %v", test.gid, err)
		}
	}
}

func TestCheckUserNotRoot(t *testing.T) {
	tests := []struct {
		user     string
		expected error
	}{
		{"root", fmt.Errorf("this program must be run when a regular user is the console user, not root")},
		{
			"_mbsetupuser",
			fmt.Errorf("this program must be run when a regular user is the console user, not _mbsetupuser"),
		},
		{"testuser", nil},
	}

	for _, tt := range tests {
		err := checkUserNotRoot(tt.user)
		assert.Equal(t, tt.expected, err)
	}
}
