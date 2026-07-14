package handlers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/lckrugel/file-server/internal/files"
	"github.com/lckrugel/file-server/internal/http_api"
)

type fileResponse struct {
	ID          string    `json:"id"`
	Filename    string    `json:"filename"`
	ContentType string    `json:"content_type"`
	Size        int64     `json:"size"`
	UploadedAt  time.Time `json:"uploaded_at"`
}

func toFileResponse(f *files.File) fileResponse {
	return fileResponse{
		ID:          f.ID,
		Filename:    f.Filename,
		ContentType: f.ContentType,
		Size:        f.Size,
		UploadedAt:  f.UploadedAt,
	}
}

func toFileResponses(fs []*files.File) []fileResponse {
	response := make([]fileResponse, len(fs))
	for i, f := range fs {
		response[i] = toFileResponse(f)
	}
	return response
}

type FileHandler struct {
	fileService *files.FileService
}

func NewFileHandler(fs *files.FileService) *FileHandler {
	return &FileHandler{
		fileService: fs,
	}
}

func (fh *FileHandler) Upload(w http.ResponseWriter, req *http.Request) {
	filename := req.Header.Get("X-Filename")
	if filename == "" {
		apiErr := http_api.BadRequest("missing filename")
		apiErr.Write(w)
		return
	}

	contentType := req.Header.Get("Content-Type")
	if contentType == "" {
		apiErr := http_api.BadRequest("unkown file type")
		apiErr.Write(w)
		return
	}

	contentLength := req.Header.Get("Content-Length")
	if contentLength == "" {
		apiErr := http_api.BadRequest("unkown file size")
		apiErr.Write(w)
		return
	}
	fileSize, err := strconv.ParseInt(contentLength, 10, 64)
	if err != nil {
		apiErr := http_api.BadRequest("invalid file size")
		apiErr.Write(w)
		return
	}

	file, err := fh.fileService.Store(req.Context(), &files.FileStore{
		Filename:      filename,
		ContentType:   contentType,
		ContentLength: fileSize,
	}, req.Body)
	if err != nil {
		var apiErr *http_api.APIError
		switch {
		case errors.Is(err, files.ErrInvalidFilename):
			apiErr = http_api.BadRequest("invalid filename")
		case errors.Is(err, files.ErrFileTooLarge):
			apiErr = http_api.BadRequest("file too large")
		default:
			log.Printf("unexpected error while storing file: %v", err)
			apiErr = http_api.InternalError("failed to store file", err)
		}
		apiErr.Write(w)
		return
	}

	resp := http_api.Created("file uploaded successfully", toFileResponse(file))
	resp.Write(w)
}

func (fh *FileHandler) Download(w http.ResponseWriter, req *http.Request) {
	fileId := req.PathValue("fileId")
	if fileId == "" {
		apiErr := http_api.BadRequest("invalid file id")
		apiErr.Write(w)
		return
	}

	fileMeta, err := fh.fileService.Get(req.Context(), fileId)
	if err != nil {
		var apiErr *http_api.APIError
		switch {
		case errors.Is(err, files.ErrFileNotFound):
			apiErr = http_api.NotFound("file not found")

		default:
			log.Printf("failed to get file metadata: %v", err)
			apiErr = http_api.InternalError("failed to get file", err)
		}
		apiErr.Write(w)
		return
	}

	fileReader, err := fh.fileService.Stream(req.Context(), fileMeta.ID)
	if err != nil {
		var apiErr *http_api.APIError
		switch {
		case errors.Is(err, files.ErrFileNotFound):
			apiErr = http_api.NotFound("file not found")

		case errors.Is(err, files.ErrFileNotReady):
			apiErr = http_api.NotFound("file not found")

		default:
			log.Printf("failed to stream file: %v", err)
			apiErr = http_api.InternalError("failed to get file", err)
		}
		apiErr.Write(w)
		return
	}
	defer fileReader.Close()

	w.Header().Set("Content-Type", fileMeta.ContentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileMeta.Filename))
	w.Header().Set("Content-Length", strconv.FormatInt(fileMeta.Size, 10))

	io.Copy(w, fileReader)
}

func (fh *FileHandler) List(w http.ResponseWriter, req *http.Request) {
	files, err := fh.fileService.List(req.Context())
	if err != nil {
		log.Printf("error while listing files: %v", err)
		apiErr := http_api.InternalError("file uploaded successfully", err)
		apiErr.Write(w)
		return
	}

	resp := http_api.Ok("", toFileResponses(files))
	resp.Write(w)
}
