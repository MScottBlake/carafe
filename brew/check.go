package brew

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/macadmins/carafe/exec"
)

type CheckResult struct {
	Installed           bool   `json:"installed"`
	Version             string `json:"version"`
	MeetsMinimumVersion bool   `json:"meets_minimum_version"`
	Name                string `json:"name"`
}

func Check(c exec.CarafeConfig, item, minVersion string, munkiInstallsCheck, skipNotInstalled bool) (int, error) {
	pass := 0
	fail := 1

	if munkiInstallsCheck {
		pass = 1
		fail = 0
	}

	if minVersion == "" && skipNotInstalled {
		return fail, fmt.Errorf("--min-version must be set when --skip-not-installed is set")
	}

	result := CheckResult{
		Name: item,
	}

	// make sure brew is installed - check the path is present
	_, err := c.CUSudo.OSFunc.Stat(c.GetBrewPath())
	if err != nil {
		if os.IsNotExist(err) {
			return pass, nil // if brew is not installed we do not need to do anything
		}
	}

	output, err := infoOutput(c, item)
	if err != nil {
		return fail, err
	}

	// is it installed?
	isPresent, err := installed(output)
	if err != nil {
		return fail, err
	}

	result.Installed = isPresent
	if !isPresent {
		err = printResultJSON(result)
		if err != nil {
			return fail, err
		}
		if skipNotInstalled {
			return pass, nil
		}
		return fail, nil
	}

	// get the version
	version, err := getVersion(output)
	if err != nil {
		return pass, err
	}

	result.Version = version

	if minVersion != "" {
		// does it meet the minimum version?
		meetsMinimum, err := VersionMeetsOrExceedsMinimum(c, item, minVersion)
		if err != nil {
			return pass, err
		}

		result.MeetsMinimumVersion = meetsMinimum
		if !meetsMinimum {
			err = printResultJSON(result)
			if err != nil {
				return fail, err
			}
			return fail, nil
		}
	}
	return pass, printResultJSON(result)
}

func printResultJSON(result CheckResult) error {
	// Print the result as JSON
	b, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(b))
	return nil
}
