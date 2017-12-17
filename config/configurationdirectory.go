package config

import (
	"path/filepath"
)

type ConfigurationDirectory struct {
	homeDirectoryProvider DirectoryPathGetter
	fileSystem            FileDirectoryReaderWriter
}

func NewConfigurationDirectory(homeDirProvider DirectoryPathGetter, directoryReaderWriter FileDirectoryReaderWriter) *ConfigurationDirectory {
	return &ConfigurationDirectory{
		homeDirectoryProvider: homeDirProvider,
		fileSystem:            directoryReaderWriter,
	}
}

func (configDirectory *ConfigurationDirectory) WriteFile(filename string, data []byte) error {
	configDirectory.createIfNotExists()
	filepath := configDirectory.fileSystem.JoinPath(configDirectory.getPath(), filename)
	return configDirectory.fileSystem.WriteFile(filepath, data)
}

func (configDirectory *ConfigurationDirectory) getPath() string {
	return filepath.Join(configDirectory.homeDirectoryProvider.GetPath(), ".ch360")
}

func (configDirectory *ConfigurationDirectory) createIfNotExists() error {
	dir := configDirectory.getPath()
	return configDirectory.fileSystem.CreateDirectoryIfNotExists(dir)
}
