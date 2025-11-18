package cudo

import (
	"fmt"
	"os/user"
	"testing"

	"github.com/macadmins/carafe/shell/testshell"

	"github.com/stretchr/testify/assert"
)

func TestGetLinuxUserHome(t *testing.T) {
	tests := []struct {
		name          string
		sudoUser      string
		lookupUser    *user.User
		lookupUserErr error
		expectedHome  string
		expectedErr   error
	}{
		{
			name:         "Valid SUDO_USER",
			sudoUser:     "testuser",
			lookupUser:   &user.User{Username: "testuser", HomeDir: "/home/testuser"},
			expectedHome: "/home/testuser",
			expectedErr:  nil,
		},
		{
			name:         "SUDO_USER not set",
			sudoUser:     "",
			expectedHome: "",
			expectedErr:  fmt.Errorf("SUDO_USER environment variable is not set"),
		},
		{
			name:          "LookupUser error",
			sudoUser:      "testuser",
			lookupUserErr: fmt.Errorf("user not found"),
			expectedHome:  "",
			expectedErr:   fmt.Errorf("user not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockOSFunc := &MockOSFunc{
				envVars: map[string]string{
					"SUDO_USER": tt.sudoUser,
				},
				lookupUser: tt.lookupUser,
				lookupErr:  tt.lookupUserErr,
			}

			c := CUSudo{OSFunc: mockOSFunc}
			homeDir, err := c.getLinuxUserHome()

			assert.Equal(t, tt.expectedHome, homeDir)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetWindowsUserHome(t *testing.T) {
	tests := []struct {
		name         string
		expectedHome string
		expectedErr  error
	}{
		{
			name:         "Not Implemented",
			expectedHome: "not implemented",
			expectedErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CUSudo{}
			homeDir, err := c.getWindowsUserHome()

			assert.Equal(t, tt.expectedHome, homeDir)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestGetDarwinUserHome(t *testing.T) {
	tests := []struct {
		name         string
		currentUser  string
		runOutput    string
		runErr       error
		expectedHome string
		expectedErr  error
	}{
		{
			name:         "Valid user home",
			currentUser:  "testuser",
			runOutput:    "NFSHomeDirectory: /Users/testuser",
			expectedHome: "/Users/testuser",
			expectedErr:  nil,
		},
		{
			name:         "User not found",
			currentUser:  "testuser",
			runOutput:    "",
			expectedHome: "",
			expectedErr:  fmt.Errorf("could not find home directory for user testuser"),
		},
		{
			name:         "Run error",
			currentUser:  "testuser",
			runErr:       fmt.Errorf("command not found"),
			expectedHome: "",
			expectedErr:  fmt.Errorf("command not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExecutor := testshell.OutputExecutor(tt.runOutput)
			c := CUSudo{Executor: mockExecutor, CurrentUser: tt.currentUser, Platform: darwin}
			homeDir, err := c.getDarwinUserHome()

			if tt.expectedErr != nil {
				assert.Error(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedHome, homeDir)
			}
		})
	}
}

func TestSetUserHome(t *testing.T) {
	tests := []struct {
		name         string
		currentUser  string
		runOutput    string
		runErr       error
		expectedHome string
		expectedErr  error
	}{
		{
			name:         "Valid user home",
			currentUser:  "testuser",
			runOutput:    "NFSHomeDirectory: /Users/testuser",
			expectedHome: "/Users/testuser",
			expectedErr:  nil,
		},
		{
			name:         "User not found",
			currentUser:  "testuser",
			runOutput:    "",
			expectedHome: "",
			expectedErr:  fmt.Errorf("could not find home directory for user testuser"),
		},
		{
			name:         "Run error",
			currentUser:  "testuser",
			runErr:       fmt.Errorf("command not found"),
			expectedHome: "",
			expectedErr:  fmt.Errorf("command not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExecutor := testshell.OutputExecutor(tt.runOutput)
			c := CUSudo{Executor: mockExecutor, CurrentUser: tt.currentUser, Platform: darwin}
			err := c.SetUserHome()

			if tt.expectedErr != nil {
				assert.Error(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedHome, c.UserHome)
			}
		})
	}
}
