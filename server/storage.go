package server

import (
	"os"
	"path/filepath"
	"sync"
)

// Storage handles file storage operations
type Storage struct {
	basePath string
	mu       sync.Mutex
}

// NewStorage creates a new storage instance
func NewStorage(basePath string) (*Storage, error) {
	// Ensure base path exists
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, err
	}

	return &Storage{
		basePath: basePath,
	}, nil
}

// SaveFile saves data to a file
func (s *Storage) SaveFile(filename string, data []byte) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create the full path
	fullPath := filepath.Join(s.basePath, filename)

	// Ensure directory exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	// Write the file
	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return "", err
	}

	return fullPath, nil
}

// GetFilePath returns the full path for a filename
func (s *Storage) GetFilePath(filename string) string {
	return filepath.Join(s.basePath, filename)
}
