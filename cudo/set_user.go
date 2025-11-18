package cudo

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
)

func (c *CUSudo) SetConsoleUser() error {
	var currentUser string
	switch c.Platform {
	case windows:
		realUser := c.OSFunc.Getenv("USERNAME")
		if realUser == "" {
			return fmt.Errorf("failed to get real user on Windows")
		}
		currentUser = realUser
	case darwin:
		info, err := c.OSFunc.Stat("/dev/console")
		if err != nil {
			return fmt.Errorf("failed to stat /dev/console: %v", err)
		}
		stat, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("failed to get Stat_t for /dev/console")
		}
		uid := strconv.Itoa(int(stat.Uid))
		theUser, err := c.OSFunc.LookupID(uid)
		if err != nil {
			return fmt.Errorf("failed to lookup user by UID: %v", err)
		}

		currentUser = theUser.Username

	default: // Linux and other Unix-like systems
		realUser := c.OSFunc.Getenv("SUDO_USER")
		if realUser == "" {
			realUser = c.OSFunc.Getenv("USER")
		}
		if realUser == "" {
			return fmt.Errorf("failed to get real user on Unix-like system")
		}

		currentUser = realUser
	}

	err := checkUserNotRoot(currentUser)
	if err != nil {
		return err
	}

	c.CurrentUser = currentUser

	return nil
}

func checkUserNotRoot(user string) error {
	if user == "root" || user == "_mbsetupuser" {
		return fmt.Errorf("this program must be run when a regular user is the console user, not %s", user)
	}
	return nil
}

func isAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	return err == nil
}

func (c *CUSudo) CheckForRoot() error {

	if c.Platform == windows {
		if isAdmin() {
			return nil
		} else {
			return fmt.Errorf("this program should be run as Administrator")
		}
	} else {
		if c.OSFunc.Getegid() != 0 {
			return fmt.Errorf("this program must be run as root")
		}
	}
	return nil
}
