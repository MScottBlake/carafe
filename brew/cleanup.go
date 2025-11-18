package brew

import (
	"github.com/macadmins/carafe/exec"
)

func Cleanup(c exec.CarafeConfig, item string) error {
	args := []string{"cleanup", item}
	_, err := c.RunBrew(args)
	return err
}
