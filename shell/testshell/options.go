package testshell

import (
	"io"
	"regexp"

	"github.com/macadmins/carafe/shell"
)

// ExecutorOption modifies how a mock Executor works.
type ExecutorOption func(*capturingExecutor)

// AlwaysError causes all calls to Run to fail.
func AlwaysError(err error) ExecutorOption {
	return func(e *capturingExecutor) {
		e.err = err
	}
}

// AlwaysErrorOnCommandIndex makes the <index>th command run return the specified error. Index starts at 1.
func AlwaysErrorOnCommandIndex(err error, index int) ExecutorOption {
	return func(e *capturingExecutor) {
		e.err = err
		e.errIdx = index
	}
}

// WithStdout will make each call to Run/Start to write the next string to the command's stdout.
// If not enough strings are provided to account for all the invocations, Run will fail on the next one.
func WithStdout(stdout ...string) ExecutorOption {
	return func(e *capturingExecutor) {
		e.stdout = stdout
	}
}

func WithStderr(stderr ...string) ExecutorOption {
	return func(e *capturingExecutor) {
		e.stderr = stderr
	}
}

// WithMappedStdouts has the executor returns the stdout value for any commands which match
// the regexp key (and return nil error).
func WithMappedStdouts(stdouts map[*regexp.Regexp]string) ExecutorOption {
	return func(e *capturingExecutor) {
		e.mappedStdouts = stdouts
	}
}

// WithMappedStderrs has the executor returns the stdout value for any commands which match
// the regexp key (and return nil error).
func WithMappedStderrs(stderr map[*regexp.Regexp]string) ExecutorOption {
	return func(e *capturingExecutor) {
		e.mappedStderrs = stderr
	}
}

func WithPID(pid int) ExecutorOption {
	return func(e *capturingExecutor) {
		e.pid = pid
	}
}

// CaptureTo will put each command passed to Run/Start into the given slice.
func CaptureTo(captured *[]*shell.Cmd) ExecutorOption {
	return func(e *capturingExecutor) {
		e.captured = captured
	}
}

// CaptureInputs will puts each stdin for commands ran into the given slice.
func CaptureInputs(inputs *[]io.Reader) ExecutorOption {
	return func(e *capturingExecutor) {
		e.inputs = inputs
	}
}
