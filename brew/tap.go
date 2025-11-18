package brew

import "github.com/macadmins/carafe/exec"

func Tap(c exec.CarafeConfig, item string) error {
	args := []string{"tap", item}
	_, err := c.RunBrew(args)
	return err
}
