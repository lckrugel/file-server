package files

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/google/uuid"
)

type FileRepository interface {
	Insert(ctx context.Context, f *File) error
	GetById(ctx context.Context, id string) (*File, error)
	List(ctx context.Context) ([]*File, error)
	Update(ctx context.Context, file *File) (*File, error)
}

type FileStorage interface {
	Write(ctx context.Context, id string, r io.Reader) (int64, error)
	Read(ctx context.Context, id string) (io.ReadCloser, error)
	Delete(ctx context.Context, id string) error
}

// Domain errors
var (
	ErrFileNotFound    = errors.New("file not found")
	ErrFileNotReady    = errors.New("file not ready")
	ErrFileTooLarge    = errors.New("file too large")
	ErrInvalidFilename = errors.New("invalid filename")
)

type FileStatus string

const (
	FileStatusPending FileStatus = "pending"
	FileStatusReady   FileStatus = "ready"
	FileStatusFailed  FileStatus = "failed"
)

type File struct {
	ID          string
	Filename    string
	ContentType string
	Size        int64
	Path        string
	Status      FileStatus
	UploadedAt  time.Time
}

func NewFile(filename, contentType string) *File {
	return &File{
		ID:          uuid.New().String(),
		Filename:    filename,
		ContentType: contentType,
		Status:      FileStatusPending,
		UploadedAt:  time.Now(),
	}
}

type FileStore struct {
	Filename      string
	ContentType   string
	ContentLength int64
}
