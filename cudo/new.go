package cudo

import (
	"os"
	"os/user"

	"github.com/macadmins/carafe/shell"
)

type CUSudo struct {
	Platform     string
	Executor     shell.Executor
	CurrentUser  string
	CurrentGroup string
	OSFunc       OSFunc
	UserHome     string
	CWD          string
}

type OSFunc interface {
	Stat(string) (os.FileInfo, error)
	LookupID(string) (*user.User, error)
	LookupUser(string) (*user.User, error)
	Getegid() int
	LookupGroupID(string) (*user.Group, error)
	Getenv(string) string
}

type StdlibOSFunc struct{}

func (s *StdlibOSFunc) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (s *StdlibOSFunc) LookupID(uid string) (*user.User, error) {
	return user.LookupId(uid)
}

func (s *StdlibOSFunc) LookupUser(name string) (*user.User, error) {
	return user.Lookup(name)
}

func (s *StdlibOSFunc) Getegid() int {
	return os.Getegid()
}

func (s *StdlibOSFunc) LookupGroupID(id string) (*user.Group, error) {
	return user.LookupGroupId(id)
}

func (s *StdlibOSFunc) Getenv(key string) string {
	return os.Getenv(key)
}

const windows = "windows"
const darwin = "darwin"

type CUSudoOption func(*CUSudo)

func WithPlatform(platform string) CUSudoOption {
	return func(c *CUSudo) {
		c.Platform = platform
	}
}

func WithCWD(cwd string) CUSudoOption {
	return func(c *CUSudo) {
		c.CWD = cwd
	}
}

func WithExecutor(exec shell.Executor) CUSudoOption {
	return func(c *CUSudo) {
		c.Executor = exec
	}
}

func WithOSFunc(osf OSFunc) CUSudoOption {
	return func(c *CUSudo) {
		c.OSFunc = osf
	}
}

func NewCUSudo(opts ...CUSudoOption) *CUSudo {
	c := &CUSudo{
		OSFunc:   &StdlibOSFunc{},
		Executor: shell.NewExecutor(shell.Verbose()),
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.Platform == "" {
		c.Platform = darwin
	}

	return c
}
