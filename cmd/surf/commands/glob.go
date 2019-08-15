package commands

import "github.com/mattn/go-zglob"

// GlobMany searches multiple patterns for files.
func GlobMany(filePatterns []string) []string {
	var files []string
	for _, filePattern := range filePatterns {
		globResult := Glob(filePattern)
		files = append(files, globResult...)
	}
	return files
}

// Glob searches provided the provided file pattern for files.
func Glob(filePattern string) []string {
	var files, _ = zglob.Glob(filePattern)
	return files
}
