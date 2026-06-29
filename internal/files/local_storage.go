package files

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	basePath string
}

// Validates that LocalStorage implements FileStorage
var _ FileStorage = (*LocalStorage)(nil)

func NewLocalStorage(path string) *LocalStorage {
	return &LocalStorage{
		basePath: path,
	}
}

func (local *LocalStorage) Write(ctx context.Context, id string, r io.Reader) (int64, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	path := filepath.Join(local.basePath, id)
	dest, err := os.Create(path)
	if err != nil {
		return 0, fmt.Errorf("failed to write file to %v: %w", path, err)
	}
	defer dest.Close()
	writen, err := io.Copy(dest, r)
	if err != nil {
		return 0, fmt.Errorf("failed to copy file content: %w", err)
	}
	return writen, nil
}

func (local *LocalStorage) Read(ctx context.Context, id string) (io.ReadCloser, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	path := filepath.Join(local.basePath, id)
	r, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file at %v: %w", path, err)
	}
	return r, nil
}

func (local *LocalStorage) Delete(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	path := filepath.Join(local.basePath, id)
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to delete file %v: %w", id, err)
	}
	return nil
}
