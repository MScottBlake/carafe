package shell

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecutor_Run(t *testing.T) {
	cmd := NewCommand("echo", "Hello world")
	stdout := new(bytes.Buffer)
	cmd.Stdout = stdout
	err := NewExecutor().Run(cmd)
	require.NoError(t, err)
	assert.True(t, cmd.ProcessState.Exited())
	assert.Equal(t, 0, cmd.ProcessState.ExitCode())
	assert.Equal(t, "Hello world\n", stdout.String())
}

func TestExecutor_Run_Error(t *testing.T) {
	cmd := NewCommand("false")
	err := NewExecutor().Run(cmd)
	assert.Error(t, err)
	assert.True(t, cmd.ProcessState.Exited())
	assert.Equal(t, 1, cmd.ProcessState.ExitCode())
}

func TestExecutor_Start(t *testing.T) {
	cmd := NewCommand("echo", "hello world")
	err := NewExecutor().Start(cmd)
	assert.NoError(t, err)
}

func TestStart(t *testing.T) {
	cmd := NewCommand("true")
	err := Start(NewExecutor(), cmd)
	assert.NoError(t, err)
}

func TestStart_ForceSilent(t *testing.T) {
	cmd := NewCommand("echo", "hello world")
	cmd.Stdout = new(bytes.Buffer)
	err := Start(NewExecutor(), cmd, ForceSilent())
	assert.NoError(t, err)
	assert.Empty(t, cmd.Stdout)
}

func TestRun(t *testing.T) {
	// this command will actually run!
	cmd := NewCommand("echo", "Hello world")
	stdout := new(bytes.Buffer)
	cmd.Stdout = stdout
	err := Run(NewExecutor(), cmd)
	require.NoError(t, err)
	assert.True(t, cmd.ProcessState.Exited())
	assert.Equal(t, 0, cmd.ProcessState.ExitCode())
	assert.Equal(t, "Hello world\n", stdout.String())
	assert.Nil(t, cmd.Env)
}

func TestRun_Cwd(t *testing.T) {
	cmd := NewCommand("pwd")
	stdout := new(bytes.Buffer)
	cmd.Stdout = stdout
	dir, _ := os.MkdirTemp(os.TempDir(), "*")
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(dir)
	absPath, err := filepath.Abs(filepath.Dir(dir))
	require.NoError(t, err)
	absPath, err = filepath.EvalSymlinks(absPath)
	require.NoError(t, err)
	err = Run(NewExecutor(), cmd, Cwd(absPath))
	require.NoError(t, err)
	assert.True(t, cmd.ProcessState.Exited())
	assert.Equal(t, 0, cmd.ProcessState.ExitCode())
	assert.Equal(t, fmt.Sprintf("%s\n", absPath), stdout.String())
}

func TestRun_Output(t *testing.T) {
	cmd := NewCommand("echo", "Hello world")
	stdout := new(bytes.Buffer)
	err := Run(NewExecutor(), cmd, Output(stdout))
	require.NoError(t, err)
	assert.True(t, cmd.ProcessState.Exited())
	assert.Equal(t, 0, cmd.ProcessState.ExitCode())
	assert.Equal(t, "Hello world\n", stdout.String())
}

func TestRun_Input(t *testing.T) {
	cmd := NewCommand("tee")
	buff := new(bytes.Buffer)
	cmd.Stdout = buff
	err := Run(NewExecutor(), cmd, Stdin(strings.NewReader("hello world")))
	require.NoError(t, err)
	assert.Equal(t, "hello world", buff.String())
}

func TestRun_LargeInput(t *testing.T) {
	// Make a byte slice larger than the standard linux page size
	contentBytes := make([]byte, (1<<16)+2)
	contentBytes[0] = 1
	contentBytes[1<<16] = 2
	contentBytes[(1<<16)+1] = 3

	cmd := NewCommand("tee")
	buff := new(bytes.Buffer)
	err := Run(NewExecutor(), cmd, Output(buff), Stdin(bytes.NewBuffer(contentBytes)))
	require.NoError(t, err)
	outputBytes := buff.Bytes()
	assert.Equal(t, byte(1), outputBytes[0])
	assert.Equal(t, byte(2), outputBytes[1<<16])
	assert.Equal(t, byte(3), outputBytes[(1<<16)+1])
}

func TestRun_InputAlreadyHasInput(t *testing.T) {
	cmd := NewCommand("tee")
	cmd.Stdin = strings.NewReader("hello world")
	buff := new(bytes.Buffer)
	err := Run(NewExecutor(), cmd, Output(buff), Stdin(strings.NewReader("bye world")))
	require.NoError(t, err)
	assert.Equal(t, "hello world", buff.String())
}

func TestRun_ExtraEnvs(t *testing.T) {
	cmd := NewCommand("true")
	cmd.Env = []string{"KEY=value"}
	err := Run(
		NewExecutor(),
		cmd,
		ExtraEnvNoOverwrite(map[string]string{"PWD": "randomrandom", "LDAP_UNIQUE_FOR_TEST": "yak_yak"}),
	)
	require.NoError(t, err)
	assert.NotContains(t, cmd.Env, "PWD=randomrandom")
	assert.Contains(t, cmd.Env, "LDAP_UNIQUE_FOR_TEST=yak_yak")
	assert.Contains(t, cmd.Env, "KEY=value")
	// make sure there are additional env vars set, not just our's
	assert.Greater(t, len(cmd.Env), 1)
}

func TestRunCancellable(t *testing.T) {
	// this command will actually run!
	cmd := NewCommand("echo", "Hello world")
	stdout := new(bytes.Buffer)
	cmd.Stdout = stdout
	ctx := context.Background()
	err := Run(NewExecutor(), cmd, Context(ctx))
	require.NoError(t, err)
	assert.True(t, cmd.ProcessState.Exited())
	assert.Equal(t, 0, cmd.ProcessState.ExitCode())
	assert.Equal(t, "Hello world\n", stdout.String())
}

func TestRunCancellable_Cancel(t *testing.T) {
	cmd := NewCommand("echo", "Hello world")
	stdout := new(bytes.Buffer)
	cmd.Stdout = stdout
	ctx, cancelFn := context.WithCancel(context.Background())
	cancelFn()
	err := Run(NewExecutor(), cmd, Context(ctx))
	assert.Error(t, err)
}

func TestJSONOutput(t *testing.T) {
	// happy case
	cmd := NewCommand("echo", `{"a": "b"}`)
	testOut := make(map[string]string)
	err := Run(NewExecutor(), cmd, JSONOutput(&testOut))
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"a": "b"}, testOut)

	// error case
	cmd = NewCommand("echo", ":) this is not JSON! :(")
	testOut = make(map[string]string)
	err = Run(NewExecutor(), cmd, JSONOutput(&testOut))
	require.Error(t, err)
}
