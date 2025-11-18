package cudo

import (
	"os"
	"os/user"
	"time"
)

// MockFileInfo is a mock implementation of os.FileInfo
type MockFileInfo struct {
	sys interface{}
}

func (m *MockFileInfo) Name() string       { return "" }
func (m *MockFileInfo) Size() int64        { return 0 }
func (m *MockFileInfo) Mode() os.FileMode  { return 0 }
func (m *MockFileInfo) ModTime() time.Time { return time.Time{} }
func (m *MockFileInfo) IsDir() bool        { return false }
func (m *MockFileInfo) Sys() interface{}   { return m.sys }

// MockOSFunc is a custom mock implementation of the OSFunc interface
type MockOSFunc struct {
	envVars     map[string]string
	statResult  os.FileInfo
	statErr     error
	lookupUser  *user.User
	lookupErr   error
	egid        int
	lookupGroup *user.Group
}

func (m *MockOSFunc) Getenv(key string) string {
	return m.envVars[key]
}

func (m *MockOSFunc) Stat(name string) (os.FileInfo, error) {
	return m.statResult, m.statErr
}

func (m *MockOSFunc) LookupID(uid string) (*user.User, error) {
	return m.lookupUser, m.lookupErr
}

func (m *MockOSFunc) LookupUser(name string) (*user.User, error) {
	return m.lookupUser, m.lookupErr
}

func (m *MockOSFunc) Getegid() int {
	return m.egid
}

func (m *MockOSFunc) LookupGroupID(id string) (*user.Group, error) {
	return m.lookupGroup, nil
}
