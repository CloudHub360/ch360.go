package config

import (
	"fmt"
	"os"
	"os/user"
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

		user, _ := user.Current() //TODO: Check err
		fmt.Println(user.HomeDir)
		return home
	}
	return os.Getenv("HOME")
}
