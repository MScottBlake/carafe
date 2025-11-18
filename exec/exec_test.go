package exec

import (
	"fmt"
	"testing"

	"github.com/macadmins/carafe/cudo"
	"github.com/macadmins/carafe/shell/testshell"

	"github.com/stretchr/testify/assert"
)

func TestGetBrewPath(t *testing.T) {
	tests := []struct {
		name             string
		arch             string
		expectedBrewPath string
	}{
		{
			name:             "Apple Silicon",
			arch:             "arm64",
			expectedBrewPath: "/opt/homebrew/bin/brew",
		},
		{
			name:             "Intel",
			arch:             "amd64",
			expectedBrewPath: "/usr/local/bin/brew",
		},
		{
			name:             "Default",
			arch:             "unknown",
			expectedBrewPath: "/usr/local/bin/brew",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CarafeConfig{
				Arch: tt.arch,
			}

			brewPath := c.GetBrewPath()
			assert.Equal(t, tt.expectedBrewPath, brewPath)
		})
	}
}

func TestRunBrew(t *testing.T) {
	executor := testshell.OutputExecutor("Hello, World!")
	c := CarafeConfig{
		Arch: "amd64",
		CUSudo: &cudo.CUSudo{
			Executor:    executor,
			CurrentUser: "test",
			Platform:    "darwin",
		},
	}

	out, err := c.RunBrew([]string{"hello", "world"})
	assert.Equal(t, "Hello, World!", out)

	assert.NoError(t, err)
}

func TestRunBrewError(t *testing.T) {
	executor := testshell.NewExecutor(testshell.AlwaysError(fmt.Errorf("error")))
	c := CarafeConfig{
		Arch: "amd64",
		CUSudo: &cudo.CUSudo{
			Executor:    executor,
			CurrentUser: "test",
			Platform:    "darwin",
		},
	}

	out, err := c.RunBrew([]string{"hello", "world"})
	assert.Equal(t, "", out)
	assert.Error(t, err)
}

func TestRunBrewWithOutput(t *testing.T) {
	executor := testshell.OutputExecutor("Hello, World!")
	c := CarafeConfig{
		Arch: "amd64",
		CUSudo: &cudo.CUSudo{
			Executor:    executor,
			CurrentUser: "test",
			Platform:    "darwin",
		},
	}

	out, err := c.RunBrewWithOutput([]string{"hello", "world"})
	assert.Equal(t, "Hello, World!", out)

	assert.NoError(t, err)
}

func TestRunBrewWithOutputError(t *testing.T) {
	executor := testshell.NewExecutor(testshell.AlwaysError(fmt.Errorf("error")))
	c := CarafeConfig{
		Arch: "amd64",
		CUSudo: &cudo.CUSudo{
			Executor:    executor,
			CurrentUser: "test",
			Platform:    "darwin",
		},
	}

	out, err := c.RunBrewWithOutput([]string{"hello", "world"})
	assert.Equal(t, "", out)
	assert.Error(t, err)
}
