package cudo

import "fmt"

func (c *CUSudo) SetGroup() error {

	if c.CurrentUser == "" {
		return fmt.Errorf("current user is not set")
	}

	usr, err := c.OSFunc.LookupUser(c.CurrentUser)
	if err != nil {
		return fmt.Errorf("failed to lookup user: %v", err)
	}

	group, err := c.OSFunc.LookupGroupID(usr.Gid)
	if err != nil {
		return fmt.Errorf("failed to lookup group: %v", err)
	}

	if group.Name == "" {
		return fmt.Errorf("group name is empty")
	}

	c.CurrentGroup = group.Name
	return nil
}
