package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

//go:generate mockery -name "FileWriter"
type FileWriter interface {
	WriteFile(filepath string, data []byte) error
}

type DirectoryReaderWriter interface {
	Stat(name string) (os.FileInfo, error)
	MkDir(name string, perm os.FileMode) error
	CreateDirectoryIfNotExists(dir string) error
	JoinPath(path1 string, path2 string) string
}

type FileDirectoryReaderWriter interface {
	DirectoryReaderWriter
	FileWriter
}

type FileSystem struct{}

var userReadWritePermissions os.FileMode = 0600

func (fs *FileSystem) WriteFile(filepath string, data []byte) error {
	return ioutil.WriteFile(filepath, data, userReadWritePermissions)
}

func (fs *FileSystem) ReadFile(filepath string) ([]byte, error) {
	return ioutil.ReadFile(filepath)
}

func (fs *FileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (fs *FileSystem) MkDir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

func (fs *FileSystem) CreateDirectoryIfNotExists(dir string) error {
	exists, err := fs.DirectoryExists(dir)
	if err != nil {
		return err
	}

	if !exists {
		err := fs.MkDir(dir, userReadWritePermissions)
		return err
	}
	return nil
}

func (fs *FileSystem) DirectoryExists(dir string) (bool, error) {
	_, err := fs.Stat(dir)
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

func (fs *FileSystem) JoinPath(path1 string, path2 string) string {
	return filepath.Join(path1, path2)
}
