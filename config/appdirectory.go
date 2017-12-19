package config

import (
	"github.com/CloudHub360/ch360.go/fs"
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
	err := fs.CreateDirectoryIfNotExists(appDirectory.getPath(), DirRWPermissions)
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
