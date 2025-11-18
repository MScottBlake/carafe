package shell

import (
	"context"
	"io"
	"os"
	"strings"
	"time"
)

// DefaultEnvAllowlist is the default set of allowed environment variables when a command is executed with
// limited environment.
// nolint:gochecknoglobals
var DefaultEnvAllowlist = map[string]any{
	"PATH": struct{}{},
	"HOME": struct{}{},
	"USER": struct{}{},
}

type Logger interface {
	Info(format string, args ...any)
	Debug(format string, args ...any)
}

type discardLogger struct{}

func (d *discardLogger) Info(format string, args ...any)  {}
func (d *discardLogger) Debug(format string, args ...any) {}

type executorConfig struct {
	verbose     bool
	debugLogger Logger
}

// ExecutorOption controls how an Executor will run all commands executed with it.
type ExecutorOption func(*executorConfig)

// Verbose will result in the executor streaming all command output to its stdout/err, in addition
// to any buffering occurring.
func Verbose() ExecutorOption {
	return func(config *executorConfig) {
		config.verbose = true
	}
}

// WithLogger adds a debug logger to the executor.
// This doesn't change the command's stdout/stderr, but instead
// controls debugging statements from the executor itself.
// (eg, "running command X", error diagnostics, and error logs from the command itself)
func WithLogger(logger Logger) ExecutorOption {
	return func(config *executorConfig) {
		config.debugLogger = logger
	}
}

// NoLogging turns off diagnostic debugging from the executor.
// (eg, "running command X", error diagnostics, and error logs from the command itself)
func NoLogging() ExecutorOption {
	return func(config *executorConfig) {
		config.debugLogger = &discardLogger{}
	}
}

type runConfig struct {
	stdout               io.Writer
	stderr               io.Writer
	stdin                io.Reader
	combinedOutputBuffer io.ReadWriter
	ctx                  context.Context
	env                  map[string]string
	logError             bool
	timeout              time.Duration
	dir                  string

	jsonOutputBuffer io.ReadWriter
	jsonOutput       any

	cancel func()
}

// ExecOption controls how the Run function will run a command.
type ExecOption func(*runConfig)

// Output causes the command's stdout and stderr to go to the given writer.
func Output(w io.Writer) ExecOption {
	return func(config *runConfig) {
		config.stdout = w
		config.stderr = w
	}
}

func JSONOutput(out any) ExecOption {
	return func(config *runConfig) {
		config.jsonOutput = out
	}
}

// StreamOutput streams the output to the native stdout/stderr during command execution
// and skips error pretty-printing (which would just repeat the output)
func StreamOutput() ExecOption {
	return func(config *runConfig) {
		config.stdout = os.Stdout
		config.stderr = os.Stderr
		config.combinedOutputBuffer = nil
		config.logError = false
	}
}

// Stdout causes the command's stdout to go to the given writer.
func Stdout(w io.Writer) ExecOption {
	return func(config *runConfig) {
		config.stdout = w
	}
}

// Stderr causes the command's stderr to go to the given writer
func Stderr(w io.Writer) ExecOption {
	return func(config *runConfig) {
		config.stderr = w
	}
}

// Stdin causes the command's stdin to come from the given reader.
func Stdin(r io.Reader) ExecOption {
	return func(config *runConfig) {
		config.stdin = r
	}
}

// NoErrorReporting turns off pretty-printing the output of any failed commands.
func NoErrorReporting() ExecOption {
	return func(config *runConfig) {
		config.logError = false
		config.combinedOutputBuffer = nil
	}
}

func StreamStdErr() ExecOption {
	return func(config *runConfig) {
		config.stdout = nil
		config.stderr = os.Stderr
		config.combinedOutputBuffer = nil
		config.logError = false
	}
}

// Context runs the command in the background with a cancellable context.
func Context(ctx context.Context) ExecOption {
	return func(config *runConfig) {
		config.ctx = ctx
	}
}

// Timeout runs the command with a timeout.
func Timeout(timeout time.Duration) ExecOption {
	return func(config *runConfig) {
		config.timeout = timeout
	}
}

// Cwd sets the working directory of the command.
func Cwd(dir string) ExecOption {
	return func(config *runConfig) {
		config.dir = dir
	}
}

// ExtraEnvNoOverwrite adds the given environment variables to the process' environment.
// Any vars that already exist in the environment will be unchanged!
func ExtraEnvNoOverwrite(newEnvs map[string]string) ExecOption {
	return func(config *runConfig) {
		env := SplitEquals(os.Environ())
		// if the command's Env property is nil, it will use the environment inherited from the process
		// Otherwise, it will use the exact environment we give it.
		// Here, we want to *add* the above environment variables unless they are already present in the shell
		// (ie, shell envs take precedence over explicitly defined ones like in Docker Compose .env files)
		// So, we need to explicitly set the command's Env, but make sure the existing env vars are loaded too.
		for key, value := range newEnvs {
			if _, ok := env[key]; ok {
				continue
			}
			env[key] = value
		}

		config.env = env
	}
}

// SplitEquals converts strings in the form "key=value" to a map with the same key/values.
func SplitEquals(strs []string) map[string]string {
	result := make(map[string]string, len(strs))
	for _, str := range strs {
		split := strings.SplitN(str, "=", 2)
		// keys can be non-specified, which makes them empty
		if len(split) != 2 {
			continue
		}

		result[split[0]] = split[1]
	}

	return result
}

// ExtraEnvOverwrite adds the given environment variables to the process' environment.
// Any vars present already in the environment will be overwritten.
func ExtraEnvOverwrite(newEnvs map[string]string) ExecOption {
	return func(config *runConfig) {
		env := SplitEquals(os.Environ())
		for key, value := range newEnvs {
			env[key] = value
		}

		config.env = env
	}
}

func OverwriteEnvWithAllowlist(newEnvs map[string]string, allowlist map[string]any) ExecOption {
	return func(config *runConfig) {
		existingEnv := SplitEquals(os.Environ())
		cmdEnv := make(map[string]string, len(allowlist)+len(newEnvs))

		for key, value := range existingEnv {
			if _, ok := allowlist[key]; ok {
				cmdEnv[key] = value
			}
		}

		for key, value := range newEnvs {
			cmdEnv[key] = value
		}

		config.env = cmdEnv
	}
}

type startConfig struct {
	forceSilent bool
}

// StartOption controls how the Start function will run a command.
type StartOption func(*startConfig)

// ForceSilent will cause no logging (from neither log.Info nor log.Debug) to be emitted.
func ForceSilent() StartOption {
	return func(config *startConfig) {
		config.forceSilent = true
	}
}
