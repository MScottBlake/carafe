package brew

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-version"

	"github.com/macadmins/carafe/exec"
)

type HomebrewFormula struct {
	Name      string      `json:"name"`
	Installed []Installed `json:"installed"`
}

type Installed struct {
	Version string `json:"version"`
}

func AllInfo(c exec.CarafeConfig) error {
	args := []string{"info", "--json", "--installed"}
	_, err := c.RunBrewWithOutput(args)
	return err
}

func Info(c exec.CarafeConfig, item string) error {
	args := []string{"info", "--json", item}
	_, err := c.RunBrewWithOutput(args)
	return err
}

func infoOutput(c exec.CarafeConfig, item string) (string, error) {
	args := []string{"info", "--json", item}
	out, err := c.RunBrew(args)
	return out, err
}

func installed(output string) (bool, error) {
	var info []HomebrewFormula
	err := json.Unmarshal([]byte(output), &info)
	if err != nil {
		return false, err
	}

	if len(info) == 0 {
		return false, fmt.Errorf("empty JSON array")
	}

	if len(info[0].Installed) == 0 {
		return false, nil
	}
	return true, nil
}

func IsInstalled(c exec.CarafeConfig, item string) (bool, error) {
	out, err := infoOutput(c, item)
	if err != nil {
		return false, err
	}

	return installed(out)
}

func getVersion(output string) (string, error) {
	var info []HomebrewFormula
	err := json.Unmarshal([]byte(output), &info)
	if err != nil {
		return "", err
	}

	if len(info) == 0 || len(info[0].Installed) == 0 {
		return "", nil
	}
	return info[0].Installed[0].Version, nil
}

func InstalledVersion(c exec.CarafeConfig, item string) (string, error) {
	out, err := infoOutput(c, item)
	if err != nil {
		return "", err
	}

	return getVersion(out)
}

func VersionMeetsOrExceedsMinimum(c exec.CarafeConfig, item, minimumVersion string) (bool, error) {
	out, err := infoOutput(c, item)
	if err != nil { // couldn't get the state, return true to be safe
		return true, err
	}

	isInstalled, err := installed(out)
	if err != nil {
		return true, err
	}

	if !isInstalled {
		return true, nil // not installed, so it meets the minimum
	}

	installedVersion, err := getVersion(out)
	if err != nil {
		return true, err
	}

	if installedVersion == "" {
		return true, nil
	}

	parsedInstalledVersion, err := version.NewVersion(installedVersion)
	if err != nil {
		return true, err
	}

	parsedMinimumVersion, err := version.NewVersion(minimumVersion)
	if err != nil {
		return true, err
	}

	return parsedInstalledVersion.GreaterThanOrEqual(parsedMinimumVersion), nil
}
