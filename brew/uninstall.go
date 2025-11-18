package brew

import "github.com/macadmins/carafe/exec"

func Uninstall(c exec.CarafeConfig, item string) error {
	args := []string{"uninstall", item}
	_, err := c.RunBrew(args)
	return err
}
