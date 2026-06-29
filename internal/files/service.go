package files

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

const MAX_FILE_SIZE_BYTES = 5 * 1024 * 1024 * 1024 // 5GB

type FileService struct {
	repo    FileRepository
	storage FileStorage
}

func NewFileService(repo FileRepository, storage FileStorage) *FileService {
	return &FileService{
		repo:    repo,
		storage: storage,
	}
}

func (fs *FileService) Store(ctx context.Context, f *FileStore, r io.Reader) (*File, error) {
	if f.ContentLength > MAX_FILE_SIZE_BYTES {
		return nil, ErrFileTooLarge
	}

	filename, err := sanitizeFilename(f.Filename)
	if err != nil {
		return nil, err
	}

	file := NewFile(filename, f.ContentType)
	limited := io.LimitReader(r, MAX_FILE_SIZE_BYTES-1)

	if err := fs.repo.Insert(ctx, file); err != nil {
		return nil, fmt.Errorf("failed to store metadata: %w", err)
	}

	writenBytes, err := fs.storage.Write(ctx, file.ID, limited)
	if err != nil {
		file.Status = FileStatusFailed
		if _, err := fs.repo.Update(ctx, file); err != nil {
			return nil, fmt.Errorf("failed to update metadata: %w", err)
		}
		return nil, fmt.Errorf("error while storing file: %w", err)
	}

	file.Size = writenBytes
	file.Status = FileStatusReady

	file, err = fs.repo.Update(ctx, file)
	if err != nil {
		return nil, fmt.Errorf("failed to store metadata: %w", err)
	}

	return file, nil
}

func (fs *FileService) List(ctx context.Context) ([]*File, error) {
	return fs.repo.List(ctx)
}

func (fs *FileService) Get(ctx context.Context, id string) (*File, error) {
	file, err := fs.repo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find file: %w", err)
	}
	if file == nil {
		return nil, ErrFileNotFound
	}
	return file, nil
}

func (fs *FileService) Stream(ctx context.Context, id string) (io.ReadCloser, error) {
	file, err := fs.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if file.Status != FileStatusReady {
		return nil, ErrFileNotReady
	}
	r, err := fs.storage.Read(ctx, file.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to stream file %v: %w", id, err)
	}
	return r, nil
}

func sanitizeFilename(name string) (string, error) {
	if len(name) > 64 {
		return "", ErrInvalidFilename
	}
	name = filepath.Base(name)
	name = strings.ReplaceAll(name, "\x00", "")
	if name == "" {
		return "", ErrInvalidFilename
	}
	return name, nil
}
