package cudo

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/macadmins/carafe/shell"
)

const defaultMaxOutputBytes int64 = 4 << 20 // 4 MiB

// limitedBufferWriter writes to an in-memory buffer up to a limit.
// After the limit is reached, it discards further writes while reporting success
// to callers (so MultiWriter doesn't fail), and records that the limit was exceeded.
type limitedBufferWriter struct {
	buf      *bytes.Buffer
	limit    int64
	written  int64
	exceeded bool
}

func newLimitedBufferWriter(buf *bytes.Buffer, limit int64) *limitedBufferWriter {
	return &limitedBufferWriter{buf: buf, limit: limit}
}

func (lw *limitedBufferWriter) Write(p []byte) (int, error) {
	remaining := lw.limit - lw.written
	if remaining <= 0 {
		lw.exceeded = true
		return len(p), nil
	}

	if int64(len(p)) > remaining {
		n, _ := lw.buf.Write(p[:remaining])
		lw.written += int64(n)
		lw.exceeded = true
		return len(p), nil
	}

	n, err := lw.buf.Write(p)
	lw.written += int64(n)
	return n, err
}

func (lw *limitedBufferWriter) Exceeded() bool { return lw.exceeded }

// RunWithOpts runs the command with the provided options
func (c *CUSudo) RunWithOpts(args []string, opts ...shell.ExecOption) error {
	cmd, err := c.BuildCmd(args)
	if err != nil {
		return err
	}

	err = shell.Run(c.Executor, cmd, opts...)
	if err != nil {
		return err
	}

	return nil
}

// Run runs the command, with a default HOME and CWD environment variables set and returns the output
func (c *CUSudo) Run(args []string) (string, error) {
	var outputBuffer bytes.Buffer
	var errBuffer bytes.Buffer

	stdout := newLimitedBufferWriter(&outputBuffer, defaultMaxOutputBytes)
	stderr := newLimitedBufferWriter(&errBuffer, defaultMaxOutputBytes)

	envSlice := []map[string]string{
		c.SetPathEnv(),
	}

	if c.CWD != "" {
		envSlice = append(envSlice, c.SetPWD(c.CWD))
	}

	envOpt := c.SetEnvOpts(envSlice...)
	opts := append([]shell.ExecOption{envOpt}, shell.Stdout(stdout), shell.Stderr(stderr))
	if c.CWD != "" {
		opts = append(opts, shell.Cwd(c.CWD))
	}

	err := c.RunWithOpts(args, opts...)
	truncated := stdout.Exceeded() || stderr.Exceeded()
	if truncated {
		return outputBuffer.String(), fmt.Errorf(
			"command output exceeded limit; stderr: %s",
			strings.TrimSpace(errBuffer.String()),
		)
	}
	if err != nil {
		return errBuffer.String(), fmt.Errorf(
			"command failed: %v; stderr: %s",
			err,
			strings.TrimSpace(errBuffer.String()),
		)
	}

	return outputBuffer.String(), nil
}

func (c *CUSudo) RunWithOutput(args []string) (string, error) {
	var outputBuffer bytes.Buffer

	lw := newLimitedBufferWriter(&outputBuffer, defaultMaxOutputBytes)
	multiWriter := io.MultiWriter(os.Stdout, lw)

	err := c.RunWithOpts(args, shell.Output(multiWriter))
	if lw.Exceeded() {
		return outputBuffer.String(), fmt.Errorf("command output exceeded limit (combined output printed to stdout)")
	}
	if err != nil {
		return outputBuffer.String(), err
	}

	return outputBuffer.String(), nil
}
