package brew

import "github.com/macadmins/carafe/exec"

func Install(c exec.CarafeConfig, item string) error {
	args := []string{"install", item}
	_, err := c.RunBrew(args)
	return err
}
