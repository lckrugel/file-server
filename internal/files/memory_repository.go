package files

import (
	"context"
	"sync"
)

type MemoryRepository struct {
	files map[string]*File
	mu    sync.Mutex
}

// Validates that MemoryRepository implementes FileRepository
var _ FileRepository = (*MemoryRepository)(nil)

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		files: make(map[string]*File),
	}
}

func (memo *MemoryRepository) Insert(ctx context.Context, f *File) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	memo.mu.Lock()
	defer memo.mu.Unlock()
	memo.files[f.ID] = f
	return nil
}

func (memo *MemoryRepository) GetById(ctx context.Context, id string) (*File, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	memo.mu.Lock()
	defer memo.mu.Unlock()
	for _, file := range memo.files {
		if file.ID == id {
			return file, nil
		}
	}
	return nil, nil
}

func (memo *MemoryRepository) List(ctx context.Context) ([]*File, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	memo.mu.Lock()
	defer memo.mu.Unlock()
	var files []*File
	for _, file := range memo.files {
		if file.Status == FileStatusReady {
			files = append(files, file)
		}
	}
	return files, nil
}

func (memo *MemoryRepository) Update(ctx context.Context, file *File) (*File, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	memo.mu.Lock()
	defer memo.mu.Unlock()
	memo.files[file.ID] = file
	return file, nil
}
