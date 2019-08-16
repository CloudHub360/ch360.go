package commands

import "github.com/mattn/go-zglob"

// GlobMany searches multiple patterns for files.
func GlobMany(filePatterns []string) ([]string, error) {
	var files []string
	for _, filePattern := range filePatterns {
		globResult, err := Glob(filePattern)

		if err != nil {
			return nil, err
		}
		files = append(files, globResult...)
	}
	return files, nil
}

// Glob searches provided the provided file pattern for files.
func Glob(filePattern string) ([]string, error) {
	return zglob.Glob(filePattern)
}
