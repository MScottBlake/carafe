package brew

import "github.com/macadmins/carafe/exec"

func Upgrade(c exec.CarafeConfig, item string) error {
	args := []string{"upgrade", item}
	_, err := c.RunBrew(args)
	if err != nil {
		return err
	}

	err = Cleanup(c, item)
	return err
}

func EnsureMinimumVersion(c exec.CarafeConfig, item, version string) error {
	meetsMinimum, err := VersionMeetsOrExceedsMinimum(c, item, version)
	if err != nil {
		return err
	}

	if !meetsMinimum {
		return Upgrade(c, item)
	}

	return nil
}
