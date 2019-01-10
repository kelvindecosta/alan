package utils

import (
	"os"
	"path/filepath"
)

// CreateFileIfNotExist creates a file if it doesn't exist along with the directories in its path
func CreateFileIfNotExist(filename string) (*(os.File), error) {
	directory := filepath.Dir(filename)

	// Create necessary directories
	err := os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return os.OpenFile(filename, os.O_CREATE|os.O_RDONLY, 0666)
}
