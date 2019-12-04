package common

import (
	"os"
	"path/filepath"
)

// GetProjectPath returns the root path of project.
func GetProjectPath() string {
	wd, _ := os.Getwd()
	return wd
}

func GetWebSocketClient(f string) string {
	return filepath.Join(GetProjectPath(), "svc", "ws", "view", f)
}
