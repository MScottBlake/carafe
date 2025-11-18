package shell

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCmd_Format(t *testing.T) {
	tests := []struct {
		name      string
		cmd       *Cmd
		expectedS string
		expectedV string
	}{
		{
			name:      "simple command, no args",
			cmd:       NewCommand("whoami"),
			expectedS: "whoami",
			expectedV: "`whoami`",
		},
		{
			name: "command with env and args",
			cmd: &Cmd{
				Path: "yak",
				Args: []string{"deploy", "-g", "--profile=yak", "--name", "Yak Yak"},
				Env:  []string{"YAK_USE_PLUGINS=y"},
			},
			expectedS: "YAK_USE_PLUGINS=y deploy -g --profile=yak --name 'Yak Yak'",
			expectedV: "`YAK_USE_PLUGINS=y deploy -g --profile=yak --name 'Yak Yak'`",
		},
		{
			name: "command with working directory",
			cmd: &Cmd{
				Path: "yak",
				Args: []string{"deploy", "-g"},
				Env:  []string{"YAK_USE_PLUGINS=y"},
				Dir:  "/tmp",
			},
			expectedS: "YAK_USE_PLUGINS=y deploy -g",
			expectedV: "`YAK_USE_PLUGINS=y deploy -g` in /tmp",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedS, fmt.Sprintf("%s", tc.cmd)) //nolint:gocritic
			assert.Equal(t, tc.expectedS, tc.cmd.String())
			assert.Equal(t, tc.expectedV, fmt.Sprintf("%v", tc.cmd)) //nolint:gocritic
		})
	}
}
