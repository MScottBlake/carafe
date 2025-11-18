package cudo

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

func (c *CUSudo) getUserHome() (string, error) {
	switch c.Platform {
	case darwin:
		return c.getDarwinUserHome()
	case windows:
		return c.getWindowsUserHome()
	default:
		return c.getLinuxUserHome()
	}
}

func (c *CUSudo) getLinuxUserHome() (string, error) {
	sudoUser := c.OSFunc.Getenv("SUDO_USER")
	if sudoUser == "" {
		return "", fmt.Errorf("SUDO_USER environment variable is not set")
	}

	usr, err := c.OSFunc.LookupUser(sudoUser)
	if err != nil {
		return "", err
	}

	return usr.HomeDir, nil
}

func (c *CUSudo) getWindowsUserHome() (string, error) {
	return "not implemented", nil
}

func (c *CUSudo) getDarwinUserHome() (string, error) {
	args := []string{"/usr/bin/dscl", ".", "-read", fmt.Sprintf("/Users/%s", c.CurrentUser), "NFSHomeDirectory"}
	output, err := c.Run(args)
	if err != nil {
		return "", err
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "NFSHomeDirectory:") {
			homeDir := strings.TrimSpace(strings.TrimPrefix(line, "NFSHomeDirectory:"))
			return homeDir, nil
		}
	}

	return "", fmt.Errorf("could not find home directory for user %s", c.CurrentUser)
}

func (c *CUSudo) SetUserHome() error {
	home, err := c.getUserHome()
	if err != nil {
		return errors.Wrap(err, "failed to get user home")
	}
	c.UserHome = home
	return nil
}
