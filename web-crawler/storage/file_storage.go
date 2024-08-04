package storage

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type LocalFileStore struct {
	basePath string
}

func NewFileStore(path string) *LocalFileStore {
	return &LocalFileStore{basePath: path}
}

func (f *LocalFileStore) Save(url, content string) (string, error) {
	// Create a filename based on URL and current timestamp
	timestamp := time.Now().Unix()
	safeURL := strings.ReplaceAll(url, "/", "_")
	filename := fmt.Sprintf("%d_%s.html", timestamp, safeURL)

	// Create full path for the file
	filepath := filepath.Join(f.basePath, filename)

	// Ensure the directory exists
	err := os.MkdirAll(f.basePath, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create directory %s: %v", f.basePath, err)
	}

	// Write the content to the file
	err = os.WriteFile(filepath, []byte(content), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write file %s: %v", filepath, err)
	}

	log.Printf("Content saved to %s", filepath)
	return filepath, nil
}
