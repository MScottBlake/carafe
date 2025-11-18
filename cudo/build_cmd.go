package cudo

import (
	"errors"

	"github.com/macadmins/carafe/shell"
)

func (c *CUSudo) BuildCmd(args []string) (*shell.Cmd, error) {
	elevationPrefix, err := c.GetElevationPrefix()
	if err != nil {
		return nil, err
	}

	args = append(elevationPrefix, args...)

	cmd := shell.NewCommand(args[0], args[1:]...)

	return cmd, nil
}

func (c *CUSudo) GetElevationPrefix() ([]string, error) {
	if c.Platform == "" {
		return nil, errors.New("platform is not set")
	}
	if c.Platform == windows {
		return []string{"runas", "/user:" + c.CurrentUser}, nil
	}

	return []string{"sudo", "-H", "-u", c.CurrentUser}, nil
}
