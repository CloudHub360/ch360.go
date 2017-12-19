package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type AppDirectory struct {
	homeDirectory string
}

//go:generate mockery -name "ConfigurationWriter"
type ConfigurationWriter interface {
	WriteConfiguration(configuration *Configuration) error
}

const FileRWPermissions os.FileMode = 0600
const DirRWPermissions os.FileMode = 0700

func NewAppDirectory(homeDirectory string) *AppDirectory {
	return &AppDirectory{
		homeDirectory: homeDirectory,
	}
}

func (appDirectory *AppDirectory) WriteConfiguration(configuration *Configuration) error {
	contents, err := configuration.Serialise()
	if err != nil {
		return err
	}
	return appDirectory.write(contents)
}

func (appDirectory *AppDirectory) ReadConfiguration() (*Configuration, error) {
	contents, err := appDirectory.read()
	if err != nil {
		return nil, err
	}

	configuration, err := DeserialiseConfiguration(contents)
	return configuration, err
}

func (appDirectory *AppDirectory) write(data []byte) error {
	err := createDirectoryIfNotExists(appDirectory.getPath())
	if err != nil {
		return err
	}

	fullFilePath := filepath.Join(appDirectory.getPath(), "config.yaml")
	err = ioutil.WriteFile(fullFilePath, data, FileRWPermissions)
	return err
}

func (appDirectory *AppDirectory) read() ([]byte, error) {
	fullFilePath := filepath.Join(appDirectory.getPath(), "config.yaml")
	return ioutil.ReadFile(fullFilePath)
}

func (appDirectory *AppDirectory) getPath() string {
	return filepath.Join(appDirectory.homeDirectory, ".ch360")
}

func createDirectoryIfNotExists(dir string) error {
	exists, err := directoryExists(dir)
	if err != nil {
		return err
	}

	if !exists {
		err := os.Mkdir(dir, DirRWPermissions)
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
