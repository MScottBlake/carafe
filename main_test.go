package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateFormulaArg_AllowsValid(t *testing.T) {
	valid := []string{
		"openssl",
		"openssl@1.1",
		"libxml++",
		"qt-5",
		"foo_bar.baz",
		"gcc@12",
		"a",
		strings.Repeat("a", 128),
	}
	for _, v := range valid {
		assert.NoError(t, validateFormulaArg(v))
		assert.True(t, formulaRe.MatchString(v))
	}
}

func TestValidateFormulaArg_RejectsInvalid(t *testing.T) {
	invalid := []string{
		"",
		"homebrew/cask/vlc",
		"foo bar",
		"naïve",
		"foo:bar",
		"abc!",
		"a\nb",
		strings.Repeat("a", 129),
	}
	for _, v := range invalid {
		assert.Error(t, validateFormulaArg(v))
		assert.False(t, formulaRe.MatchString(v))
	}
}

func TestFormulaRe_DoesNotPartiallyMatch(t *testing.T) {
	assert.False(t, formulaRe.MatchString("valid@1.0!")) // '!' not allowed anywhere
	assert.False(t, formulaRe.MatchString(" valid"))     // leading space
	assert.False(t, formulaRe.MatchString("valid "))     // trailing space
}
