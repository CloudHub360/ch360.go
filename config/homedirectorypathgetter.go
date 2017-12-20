package config

import (
	"os"
	"runtime"
)

type PathGetter interface {
	Path() string
}

type HomeDirectoryPathGetter struct{}

func (provider HomeDirectoryPathGetter) Path() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
