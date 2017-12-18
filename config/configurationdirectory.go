package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type ConfigurationDirectory struct {
	homeDirectoryProvider DirectoryPathGetter
}

var userReadWritePermissions os.FileMode = 0600

func NewConfigurationDirectory(homeDirProvider DirectoryPathGetter) *ConfigurationDirectory {
	return &ConfigurationDirectory{
		homeDirectoryProvider: homeDirProvider,
	}
}

func (configDirectory *ConfigurationDirectory) Write(data []byte) (int, error) {
	configDirectory.createIfNotExists()

	fullFilePath := filepath.Join(configDirectory.getPath(), "config.yaml")
	err := ioutil.WriteFile(fullFilePath, data, userReadWritePermissions)
	return 0, err
}

func (configDirectory *ConfigurationDirectory) getPath() string {
	return filepath.Join(configDirectory.homeDirectoryProvider.GetPath(), ".ch360")
}

func (configDirectory *ConfigurationDirectory) createIfNotExists() error {
	dir := configDirectory.getPath()
	return createDirectoryIfNotExists(dir)
}

func createDirectoryIfNotExists(dir string) error {
	exists, err := directoryExists(dir)
	if err != nil {
		return err
	}

	if !exists {
		err := os.Mkdir(dir, userReadWritePermissions)
		return err
	}
	return nil
}

func directoryExists(dir string) (bool, error) {
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			// directory does not exist
			return false, nil
		} else {
			// other error
			return false, err
		}
	}

	return true, nil
}
