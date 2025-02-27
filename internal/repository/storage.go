package repository

import (
	"fmt"
	"io"
	"os"
)

// ReadFile reads file data into memory
func ReadFile(file io.Reader) ([]byte, error) {
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// SaveImageToStorage saves the image to the storage folder
func SaveImageToStorage(filename string, data []byte) (string, error) {
	path := fmt.Sprintf("storage/%s", filename)
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		return "", err
	}
	return path, nil
}
