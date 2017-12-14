package config

import (
	"os"
	"path/filepath"
	"runtime"
)

type configurationDirectory struct {
}

func (configDirectory *configurationDirectory) CreateIfNotExists() error {
	dir := configDirectory.GetPath()

	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			// directory does not exist
			os.Mkdir(dir, 0644) //TODO: Permissions?
			return nil
		} else {
			// other error
			return err
		}
	}

	return nil
}

func (dir *configurationDirectory) GetPath() string {
	return filepath.Join(userHomeDir(), ".ch360")
}

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
