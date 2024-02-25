package fileutils

import "os"

// CreateTempFile creates a temporary file and returns its path
func CreateTempFile(pattern string) (string, error) {
	file, err := os.CreateTemp("/tmp", pattern)
	return file.Name(), err
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
