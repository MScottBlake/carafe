package brew

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/macadmins/carafe/cudo"
	"github.com/macadmins/carafe/exec"
	"github.com/macadmins/carafe/shell/testshell"
)

func TestUpgrade(t *testing.T) {
	tests := []struct {
		name        string
		item        string
		runError    error
		expectError bool
	}{
		{
			name:        "success",
			item:        "package",
			runError:    nil,
			expectError: false,
		},
		{
			name:        "failure",
			item:        "package",
			runError:    fmt.Errorf("run error"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := exec.CarafeConfig{
				Arch: "amd64",
				CUSudo: &cudo.CUSudo{
					CurrentUser: "test",
					Platform:    "darwin",
				},
			}

			if tt.runError != nil {
				c.CUSudo.Executor = testshell.NewExecutor(testshell.AlwaysError(tt.runError))
			} else {
				c.CUSudo.Executor = testshell.OutputExecutor("Hello, World!", "Hello, World!")
			}
			err := Upgrade(c, tt.item)
			if tt.expectError {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func TestEnsureMinimumVersion(t *testing.T) {
	tests := []struct {
		name         string
		item         string
		version      string
		meetsMinimum bool
		versionErr   error
		upgradeErr   error
		expectedErr  error
		outputs      []string
	}{
		{
			name:         "meets minimum version",
			item:         "htop",
			version:      "2.0.0",
			meetsMinimum: true,
			versionErr:   nil,
			upgradeErr:   nil,
			expectedErr:  nil,
			outputs:      []string{TestInfoInstalledOutput, TestInfoInstalledOutput},
		},
		{
			name:         "does not meet minimum version and upgrade succeeds",
			item:         "htop",
			version:      "4.0.0",
			meetsMinimum: false,
			versionErr:   nil,
			upgradeErr:   nil,
			expectedErr:  nil,
			outputs:      []string{TestInfoInstalledOutput, TestInfoInstalledOutput, TestInfoInstalledOutput},
		},
		{
			name:         "does not meet minimum version and upgrade fails",
			item:         "htop",
			version:      "4.0.0",
			meetsMinimum: false,
			versionErr:   nil,
			upgradeErr:   assert.AnError,
			expectedErr:  assert.AnError,
			outputs:      []string{TestInfoInstalledOutput, TestInfoInstalledOutput, TestInfoInstalledOutput},
		},
		{
			name:         "version check fails",
			item:         "htop",
			version:      "2.0.0",
			meetsMinimum: false,
			versionErr:   assert.AnError,
			upgradeErr:   nil,
			expectedErr:  assert.AnError,
			outputs:      []string{TestInfoInstalledOutput, TestInfoInstalledOutput},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := exec.CarafeConfig{
				Arch: "amd64",
				CUSudo: &cudo.CUSudo{
					CurrentUser: "test",
					Platform:    "darwin",
				},
			}

			if tt.expectedErr != nil {
				c.CUSudo.Executor = testshell.NewExecutor(testshell.AlwaysError(tt.expectedErr))
			} else {
				c.CUSudo.Executor = testshell.OutputExecutor(tt.outputs...)
			}

			// Call the function under test
			err := EnsureMinimumVersion(c, tt.item, tt.version)

			// Assert the error
			if tt.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
