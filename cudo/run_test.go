package cudo

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/macadmins/carafe/shell"

	"github.com/macadmins/carafe/shell/testshell"
)

func TestRunWithOpts(t *testing.T) {
	c := CUSudo{
		Platform:    "darwin",
		CurrentUser: "testuser",
		Executor:    testshell.OutputExecutor("Hello, World!"),
	}

	args := []string{"echo", "hello"}
	opts := []shell.ExecOption{
		shell.Stdout(os.Stdout),
		shell.Stdout(os.Stderr),
	}

	err := c.RunWithOpts(args, opts...)
	assert.NoError(t, err)
}

func TestRunWithOutput(t *testing.T) {
	c := CUSudo{
		Platform:    "darwin",
		CurrentUser: "testuser",
		Executor:    testshell.OutputExecutor("Hello, World!"),
	}
	out, err := c.RunWithOutput([]string{"echo", "hello"})
	assert.NoError(t, err)
	assert.Equal(t, "Hello, World!", out)
}

func TestRun(t *testing.T) {
	c := CUSudo{
		Platform:    "darwin",
		CurrentUser: "testuser",
		Executor:    testshell.OutputExecutor("Hello, World!"),
	}

	out, err := c.Run([]string{"echo", "hello"})
	assert.NoError(t, err)
	assert.Equal(t, "Hello, World!", out)
}

func TestLimitedBufferWriter_UnderLimit(t *testing.T) {
	var buf bytes.Buffer
	lw := newLimitedBufferWriter(&buf, 10)

	n, err := lw.Write([]byte("hello"))
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, "hello", buf.String())
	assert.False(t, lw.Exceeded())
}

func TestLimitedBufferWriter_AtLimit(t *testing.T) {
	var buf bytes.Buffer
	lw := newLimitedBufferWriter(&buf, 5)

	n, err := lw.Write([]byte("hello"))
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, "hello", buf.String())
	assert.False(t, lw.Exceeded())
}

func TestLimitedBufferWriter_OverLimitSingleWrite(t *testing.T) {
	var buf bytes.Buffer
	lw := newLimitedBufferWriter(&buf, 5)

	n, err := lw.Write([]byte("helloworld"))
	assert.NoError(t, err)
	assert.Equal(t, 10, n)
	assert.Equal(t, "hello", buf.String())
	assert.True(t, lw.Exceeded())
}

func TestLimitedBufferWriter_OverLimitMultipleWrites(t *testing.T) {
	var buf bytes.Buffer
	lw := newLimitedBufferWriter(&buf, 5)

	n1, err1 := lw.Write([]byte("he"))
	n2, err2 := lw.Write([]byte("llo"))
	n3, err3 := lw.Write([]byte("world"))

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NoError(t, err3)

	assert.Equal(t, 2, n1)
	assert.Equal(t, 3, n2)
	assert.Equal(t, 5, n3)

	assert.Equal(t, "hello", buf.String())
	assert.True(t, lw.Exceeded())
}

func TestLimitedBufferWriter_ZeroLimit(t *testing.T) {
	var buf bytes.Buffer
	lw := newLimitedBufferWriter(&buf, 0)

	n, err := lw.Write([]byte("data"))
	assert.NoError(t, err)
	assert.Equal(t, 4, n)
	assert.Equal(t, "", buf.String())
	assert.True(t, lw.Exceeded())
}

func TestLimitedBufferWriter_WithMultiWriter_OtherSinkGetsAllData(t *testing.T) {
	var capped bytes.Buffer
	var pass bytes.Buffer

	lw := newLimitedBufferWriter(&capped, 5)
	mw := io.MultiWriter(lw, &pass)

	payload := strings.Repeat("x", 12)
	n, err := mw.Write([]byte(payload))
	assert.NoError(t, err)
	assert.Equal(t, len(payload), n)

	assert.Equal(t, strings.Repeat("x", 5), capped.String())
	assert.True(t, lw.Exceeded())
	assert.Equal(t, payload, pass.String())
}
