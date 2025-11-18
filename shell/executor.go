package shell

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

// Executor runs commands. Instead of running commands directly through this interface, you should probably pass this
// to the Run function for more flexibility and a variadic interface.
type Executor interface {
	// Run runs the command and waits for it to finish. See exec.Cmd#Run.
	Run(cmd *Cmd) error
	// Start runs the command and continues. See exec.Cmd#Start.
	Start(cmd *Cmd) error
}

type executor struct {
	config *executorConfig
}

type syncedWriter struct {
	w  io.Writer
	mu sync.Mutex
}

func (s *syncedWriter) Write(p []byte) (n int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.w.Write(p)
}

// NewExecutor returns a new Executor capable of running commands.
func NewExecutor(opts ...ExecutorOption) Executor {
	config := &executorConfig{
		// Ideally the caller would always provide a logger if needed and we could have a full decoupling of
		// the shell package from the log package. But, this would be a breaking change for a lot of CLIs.
		// So, we'll keep the discardLogger as default for now.
		debugLogger: &discardLogger{},
	}

	for _, opt := range opts {
		opt(config)
	}
	return &executor{config: config}
}

func (e *executor) Run(cmd *Cmd) error {
	e.config.debugLogger.Debug("Running %v", cmd)
	if e.config.verbose {
		if cmd.Stdout != os.Stdout {
			cmd.Stdout = multiWriter(cmd.Stdout, os.Stdout)
		}
		if cmd.Stderr != os.Stderr {
			cmd.Stderr = multiWriter(cmd.Stderr, os.Stderr)
		}
	}
	err := cmd.Native().Run()
	exitCode := cmd.ProcessState.ExitCode()
	e.config.debugLogger.Debug("Command exited with code %d: %v", exitCode, cmd)
	return err
}

func (e *executor) Start(cmd *Cmd) error {
	e.config.debugLogger.Debug("Running %v in background", cmd)
	err := cmd.Native().Start()
	return err
}

// Start will run the command with the given executor and options in the background.
func Start(e Executor, cmd *Cmd, opts ...StartOption) error {
	config := &startConfig{
		forceSilent: false,
	}

	for _, opt := range opts {
		opt(config)
	}

	if config.forceSilent {
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard
	}

	return e.Start(cmd)
}

func runCancellable(ctx context.Context, executor Executor, cmd *Cmd, cancel func()) error {
	errChan := make(chan error)
	if cancel != nil {
		defer cancel()
	}
	go func() {
		errChan <- executor.Run(cmd)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errChan:
		return err
	}
}

// Run will run the command with the given executor and options. This is to be preferred to calling Executor.Run.
func Run(e Executor, cmd *Cmd, opts ...ExecOption) error {
	config := &runConfig{
		// error output should include both stdout and stderr
		// Run will update stdout/err write to this combined buffer *in addition*
		// to their original values. This preserves the intention of people using Stdout()
		// in order to analyze stdout independently, while also giving a combined output if it errors.
		// this technique will probably fail if the output is too large
		combinedOutputBuffer: new(bytes.Buffer),
		logError:             true,
	}

	for _, opt := range opts {
		opt(config)
	}

	var logger Logger = &discardLogger{}
	if defaultExecutor, ok := e.(*executor); ok {
		logger = defaultExecutor.config.debugLogger
	}

	if config.jsonOutput != nil {
		config.jsonOutputBuffer = new(bytes.Buffer)
		config.stdout = config.jsonOutputBuffer
	}

	configureCmd(cmd, config)

	if config.stdin != nil {
		writeToStdin(logger, cmd, config.stdin)
	}

	var err error
	if config.ctx != nil {
		err = runCancellable(config.ctx, e, cmd, config.cancel)
	} else {
		err = e.Run(cmd)
	}

	if config.jsonOutput != nil {
		out, readErr := io.ReadAll(config.jsonOutputBuffer)
		if readErr != nil {
			return readErr
		}
		if jErr := json.Unmarshal(out, &config.jsonOutput); jErr != nil {
			return jErr
		}
	}

	if !config.logError || err == nil || errors.Cause(err) == context.Canceled {
		return err
	}

	logs := errorLogs(cmd, config.combinedOutputBuffer)
	if config.ctx != nil {
		return err
	}

	logger.Info(logs)
	return err
}

func configureCmd(cmd *Cmd, config *runConfig) {
	if config.env != nil {
		newEnvs := JoinEquals(config.env)
		cmd.Env = append(cmd.Env, newEnvs...)
	}

	if config.stdout != nil {
		cmd.Stdout = config.stdout
	}

	if config.stderr != nil {
		cmd.Stderr = config.stderr
	}

	if config.combinedOutputBuffer != nil {
		// when the command is executed, there will be two goroutines for capturing the output of stdout and stderr
		// since we are writing to a shared buffer, we need to sync these writes
		w := &syncedWriter{
			w: config.combinedOutputBuffer,
		}
		cmd.Stdout = multiWriter(cmd.Stdout, w)
		cmd.Stderr = multiWriter(cmd.Stderr, w)
	}

	if config.timeout != 0 {
		baseCtx := context.Background()
		if config.ctx != nil {
			baseCtx = config.ctx
		}
		config.ctx, config.cancel = context.WithTimeout(baseCtx, config.timeout)
	}

	if config.dir != "" {
		cmd.Dir = config.dir
	}
}

func writeToStdin(logger Logger, cmd *Cmd, r io.Reader) {
	stdin, err := cmd.Native().StdinPipe()
	if err != nil {
		// if the command already has stdin set up, assume it will be handled correctly
		return
	}

	// Write to the command's stdin pipe in a go routine so as to not block indefinitely when the stdin is very large.
	// This is necessary because linux pipes have a limited capacity and writes will block (or fail) once this
	// capacity is reached; so we need to allow the consuming command to read while the write goes through.
	go func() {
		bufferedWriter := bufio.NewWriter(stdin)
		// Since the most common failure mode here is when the command exits and closes stdin, we just log the failure
		// instead of bubbling it up through an error chan.
		if _, err := bufferedWriter.ReadFrom(r); err != nil {
			logStdinWriteFail(logger, err)
			return
		}

		if err := bufferedWriter.Flush(); err != nil {
			logStdinWriteFail(logger, err)
			return
		}
		AttemptClose(stdin)
	}()
}

func AttemptClose(closers ...io.Closer) {
	for _, c := range closers {
		if c == nil {
			continue
		}
		_ = c.Close()
	}
}

func logStdinWriteFail(logger Logger, err error) {
	logger.Debug("failed to write stdin (the likely cause is command failure): %s", err.Error())
}

func multiWriter(ws ...io.Writer) io.Writer {
	nonNil := make([]io.Writer, 0, len(ws))
	for _, w := range ws {
		if w == nil {
			continue
		}

		nonNil = append(nonNil, w)
	}

	return io.MultiWriter(nonNil...)
}

func errorLogs(cmd *Cmd, capturedOutput io.Reader) string {
	sb := new(strings.Builder)
	fmt.Fprintln(sb, "❌ Running", fmt.Sprintf("%q", strings.Join(cmd.Args, " ")), "FAILED")
	cmdName := filepath.Base(cmd.Path)
	fmt.Fprintln(sb, cmdName, cmdName, capturedOutput)
	return sb.String()
}

// JoinEquals converts a map into a slice of "key=value" pairs.
func JoinEquals(m map[string]string) []string {
	result := make([]string, 0, len(m))
	for key, value := range m {
		result = append(result, fmt.Sprintf("%s=%s", key, value))
	}

	return result
}
