package cudo

import "github.com/macadmins/carafe/shell"

// SetPathEnv sets the PATH environment variable to the default for each platform.
func (c *CUSudo) SetPathEnv() map[string]string {
	// Set the PATH environment variable to the default for each platform
	switch c.Platform {
	case windows:
		return map[string]string{"PATH": "C:\\Windows\\System32;C:\\Windows"}
	case darwin: // macOS
		return map[string]string{"PATH": "/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin"}
	default: // Linux and other Unix-like systems
		return map[string]string{"PATH": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"}
	}
}

func (c *CUSudo) SetPWD(path string) map[string]string {
	return map[string]string{"PWD": path, "CWD": path}
}

func (c *CUSudo) SetEnvOpts(opts ...map[string]string) shell.ExecOption {
	// merge all the maps into one
	envs := make(map[string]string)
	for _, opt := range opts {
		for k, v := range opt {
			envs[k] = v
		}
	}

	return shell.ExtraEnvOverwrite(envs)
}
