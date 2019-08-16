package fs

import (
	"io"
	"os"
)

func DirectoryOrFileExists(dir string) (bool, error) {
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

func CreateDirectoryIfNotExists(dir string, permissions os.FileMode) error {
	exists, err := DirectoryOrFileExists(dir)
	if err != nil {
		return err
	}

	if !exists {
		err := os.Mkdir(dir, permissions)
		return err
	}
	return nil
}

// OpenForWriting is a convenience function that wraps os.Create,
// but which returns os.Stdout if the provided filename is "-" or the empty string.
func OpenForWriting(filename string) (*os.File, error) {
	if filename == "-" || filename == "" {
		return os.Stdout, nil
	}

	return os.Create(filename)
}

// TryClose will attempt to cast all provided writers to io.Closer,
// and call Close on them if the cast succeeds.
func TryClose(writers ...io.Writer) {
	for _, writer := range writers {
		if closer, ok := writer.(io.Closer); ok {
			_ = closer.Close()
		}
	}
}
