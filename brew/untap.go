package brew

import "github.com/macadmins/carafe/exec"

func Untap(c exec.CarafeConfig, item string) error {
	args := []string{"untap", item}
	_, err := c.RunBrew(args)
	return err
}
