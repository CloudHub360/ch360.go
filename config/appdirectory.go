package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type AppDirectory struct {
	homeDirectoryProvider DirectoryPathGetter
}

//go:generate mockery -name "ConfigurationWriter"
type ConfigurationWriter interface {
	WriteConfiguration(configuration *Configuration) error
}

var userReadWritePermissions os.FileMode = 0600

func NewAppDirectory(homeDirProvider DirectoryPathGetter) *AppDirectory {
	return &AppDirectory{
		homeDirectoryProvider: homeDirProvider,
	}
}

func (appDirectory *AppDirectory) WriteConfiguration(configuration *Configuration) error {
	contents, err := configuration.Serialise()
	if err != nil {
		return err
	}
	_, err = appDirectory.write(contents)
	return err
}

func (appDirectory *AppDirectory) ReadConfiguration() (*Configuration, error) {
	contents, err := appDirectory.read()
	if err != nil {
		return nil, err
	}

	configuration, err := DeserialiseConfiguration(contents)
	return configuration, err
}

func (appDirectory *AppDirectory) write(data []byte) (int, error) {
	appDirectory.createIfNotExists()

	fullFilePath := filepath.Join(appDirectory.getPath(), "config.yaml")
	err := ioutil.WriteFile(fullFilePath, data, userReadWritePermissions)
	return 0, err
}

func (appDirectory *AppDirectory) read() ([]byte, error) {
	fullFilePath := filepath.Join(appDirectory.getPath(), "config.yaml")
	return ioutil.ReadFile(fullFilePath)
}

func (appDirectory *AppDirectory) getPath() string {
	return filepath.Join(appDirectory.homeDirectoryProvider.GetPath(), ".ch360")
}

func (appDirectory *AppDirectory) createIfNotExists() error {
	dir := appDirectory.getPath()
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
