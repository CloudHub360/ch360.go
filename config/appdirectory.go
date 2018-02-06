package config

import (
	"errors"
	"fmt"
	"github.com/CloudHub360/ch360.go/fs"
	"github.com/mitchellh/go-homedir"
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

//go:generate mockery -name "ConfigurationReader"
type ConfigurationReader interface {
	ReadConfiguration() (*Configuration, error)
}

const FileRWPermissions os.FileMode = 0600
const DirRWPermissions os.FileMode = 0700

func NewAppDirectoryInDir(dir string) *AppDirectory {
	return &AppDirectory{
		homeDirectory: dir,
	}
}

func NewAppDirectory() (*AppDirectory, error) {
	dir, err := homedir.Dir()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not determine home directory. Details: %v", err))
	}

	return NewAppDirectoryInDir(dir), nil
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

	return ioutil.WriteFile(appDirectory.configFilePath(), data, FileRWPermissions)
}

func (appDirectory *AppDirectory) read() ([]byte, error) {
	return ioutil.ReadFile(appDirectory.configFilePath())
}

func (appDirectory *AppDirectory) getPath() string {
	return filepath.Join(appDirectory.homeDirectory, ".surf")
}

func (appDirectory *AppDirectory) configFilePath() string {
	return filepath.Join(appDirectory.getPath(), "config.yaml")
}
