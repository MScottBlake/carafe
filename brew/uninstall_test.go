package brew

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/macadmins/carafe/cudo"
	"github.com/macadmins/carafe/exec"
	"github.com/macadmins/carafe/shell/testshell"
)

func TestUninstall(t *testing.T) {
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
				Arch: "arm64",
				CUSudo: &cudo.CUSudo{
					CurrentUser: "testuser",
					Platform:    "darwin",
					UserHome:    "/Users/testuser",
				},
			}

			if tt.runError != nil {
				c.CUSudo.Executor = testshell.NewExecutor(testshell.AlwaysError(tt.runError))
			} else {
				c.CUSudo.Executor = testshell.OutputExecutor("Hello, World!")
			}
			err := Uninstall(c, tt.item)
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
