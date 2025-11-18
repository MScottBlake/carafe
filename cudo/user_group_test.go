package cudo

import (
	"os/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetGroup(t *testing.T) {
	tests := []struct {
		name      string
		group     *user.Group
		user      *user.User
		expected  string
		expectErr bool
	}{
		{
			name:     "Valid group",
			group:    &user.Group{Name: "testgroup", Gid: "1000"},
			user:     &user.User{Username: "testuser", Gid: "1000"},
			expected: "testgroup",
		},
		{
			name:      "Invalid group",
			group:     &user.Group{},
			user:      &user.User{Username: "testuser", Gid: "1000"},
			expectErr: true,
		},
		{
			name:      "Invalid user",
			group:     &user.Group{Name: "testgroup", Gid: "1000"},
			user:      &user.User{},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CUSudo{
				CurrentUser: tt.user.Username,
				OSFunc: &MockOSFunc{
					lookupUser:  tt.user,
					lookupGroup: tt.group,
				},
			}
			err := c.SetGroup()

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
