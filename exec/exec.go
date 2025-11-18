package exec

import (
	"bytes"
	"io"
	"os"
	"runtime"

	"github.com/macadmins/carafe/cudo"
	"github.com/macadmins/carafe/shell"
	"github.com/pkg/errors"
)

type CarafeConfig struct {
	Arch   string
	CUSudo *cudo.CUSudo
}

func NewConfig() (CarafeConfig, error) {
	config := CarafeConfig{
		Arch:   runtime.GOARCH,
		CUSudo: cudo.NewCUSudo(cudo.WithExecutor(shell.NewExecutor())),
	}

	err := config.CUSudo.SetConsoleUser()
	if err != nil {
		return CarafeConfig{}, errors.Wrap(err, "failed to set console user")
	}

	err = config.CUSudo.SetUserHome()
	if err != nil {
		return CarafeConfig{}, errors.Wrap(err, "failed to set user home")
	}

	config.CUSudo.CWD = config.CUSudo.UserHome

	return config, nil
}

func (c CarafeConfig) GetBrewPath() string {
	// return the path to the brew executable for the architecture
	if c.Arch == "arm64" {
		return "/opt/homebrew/bin/brew"
	}
	return "/usr/local/bin/brew"
}

func (c CarafeConfig) RunBrew(args []string) (string, error) {
	args = append([]string{c.GetBrewPath()}, args...)
	out, err := c.CUSudo.Run(args)
	if err != nil {
		return out, err
	}

	return out, nil
}

func (c CarafeConfig) RunBrewWithOutput(args []string) (string, error) {
	var outputBuffer bytes.Buffer
	multiWriter := io.MultiWriter(os.Stdout, &outputBuffer)
	args = append([]string{c.GetBrewPath()}, args...)

	envSlice := []map[string]string{
		c.CUSudo.SetPathEnv(),
	}

	if c.CUSudo.CWD != "" {
		envSlice = append(envSlice, c.CUSudo.SetPWD(c.CUSudo.CWD))
	}

	envOpt := c.CUSudo.SetEnvOpts(envSlice...)
	opts := append([]shell.ExecOption{envOpt}, shell.Output(multiWriter))
	if c.CUSudo.CWD != "" {
		opts = append(opts, shell.Cwd(c.CUSudo.CWD))
	}

	err := c.CUSudo.RunWithOpts(args, opts...)
	if err != nil {
		return "", err
	}

	return outputBuffer.String(), nil
}
