package cudo

import (
	"testing"

	"github.com/macadmins/carafe/shell"

	"github.com/macadmins/carafe/shell/testshell"

	"github.com/stretchr/testify/assert"
)

func TestWithPlatform(t *testing.T) {
	c := &CUSudo{}
	opt := WithPlatform("linux")
	opt(c)
	assert.Equal(t, "linux", c.Platform)
}

func TestWithExecutor(t *testing.T) {
	c := &CUSudo{}
	mockExec := testshell.NewExecutor()
	opt := WithExecutor(mockExec)
	opt(c)
	assert.Equal(t, mockExec, c.Executor)
}

func TestWithOSFunc(t *testing.T) {
	c := &CUSudo{}
	mockOS := &MockOSFunc{}
	opt := WithOSFunc(mockOS)
	opt(c)
	assert.Equal(t, mockOS, c.OSFunc)
}

func TestWithCWD(t *testing.T) {
	c := &CUSudo{}
	opt := WithCWD("/home/user")
	opt(c)
	assert.Equal(t, "/home/user", c.CWD)
}

func TestNewCUSudo(t *testing.T) {
	mockOS := &MockOSFunc{}

	tests := []struct {
		name     string
		opts     []CUSudoOption
		expected *CUSudo
	}{
		{
			name: "default options",
			opts: nil,
			expected: &CUSudo{
				OSFunc:   &StdlibOSFunc{},
				Executor: shell.NewExecutor(shell.Verbose()),
				Platform: "darwin",
			},
		},
		{
			name: "custom platform",
			opts: []CUSudoOption{WithPlatform("linux")},
			expected: &CUSudo{
				OSFunc:   &StdlibOSFunc{},
				Executor: shell.NewExecutor(shell.Verbose()),
				Platform: "linux",
			},
		},
		{
			name: "custom OSFunc",
			opts: []CUSudoOption{WithOSFunc(mockOS)},
			expected: &CUSudo{
				OSFunc:   mockOS,
				Executor: shell.NewExecutor(shell.Verbose()),
				Platform: "darwin",
			},
		},
		{
			name: "custom Executor",
			opts: []CUSudoOption{WithExecutor(testshell.NewExecutor())},
			expected: &CUSudo{
				OSFunc:   &StdlibOSFunc{},
				Executor: testshell.NewExecutor(),
				Platform: "darwin",
			},
		},
		{
			name: "all custom options",
			opts: []CUSudoOption{
				WithPlatform("linux"),
				WithOSFunc(mockOS),
				WithExecutor(testshell.NewExecutor()),
			},
			expected: &CUSudo{
				OSFunc:   mockOS,
				Executor: testshell.NewExecutor(),
				Platform: "linux",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCUSudo(tt.opts...)
			assert.Equal(t, tt.expected.Platform, c.Platform)
			assert.IsType(t, tt.expected.OSFunc, c.OSFunc)
			assert.IsType(t, tt.expected.Executor, c.Executor)
		})
	}
}
