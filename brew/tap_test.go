package brew

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/macadmins/carafe/cudo"
	"github.com/macadmins/carafe/exec"
	"github.com/macadmins/carafe/shell/testshell"
)

func TestTap(t *testing.T) {
	tests := []struct {
		name        string
		item        string
		runError    error
		expectError bool
	}{
		{
			name:        "success",
			item:        "tap",
			runError:    nil,
			expectError: false,
		},
		{
			name:        "failure",
			item:        "tap",
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
			err := Tap(c, tt.item)
			if tt.expectError {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
		})
	}
}
