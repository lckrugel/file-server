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
		apiErr := badRequest("missing filename")
		apiErr.Write(w)
		return
	}

	contentType := req.Header.Get("Content-Type")
	if contentType == "" {
		apiErr := badRequest("unkown file type")
		apiErr.Write(w)
		return
	}

	contentLength := req.Header.Get("Content-Length")
	if contentLength == "" {
		apiErr := badRequest("unkown file size")
		apiErr.Write(w)
		return
	}
	fileSize, err := strconv.ParseInt(contentLength, 10, 64)
	if err != nil {
		apiErr := badRequest("invalid file size")
		apiErr.Write(w)
		return
	}

	file, err := fh.fileService.Store(req.Context(), &files.FileStore{
		Filename:      filename,
		ContentType:   contentType,
		ContentLength: fileSize,
	}, req.Body)
	if err != nil {
		var apiErr *apiError
		switch {
		case errors.Is(err, files.ErrInvalidFilename):
			apiErr = badRequest("invalid filename")
		case errors.Is(err, files.ErrFileTooLarge):
			apiErr = badRequest("file too large")
		default:
			log.Printf("unexpected error while storing file: %v", err)
			apiErr = internalError("failed to store file", err)
		}
		apiErr.Write(w)
		return
	}

	resp := created("file uploaded successfully", toFileResponse(file))
	resp.Write(w)
}

func (fh *FileHandler) Download(w http.ResponseWriter, req *http.Request) {
	fileId := req.PathValue("fileId")
	if fileId == "" {
		apiErr := badRequest("invalid file id")
		apiErr.Write(w)
		return
	}

	fileMeta, err := fh.fileService.Get(req.Context(), fileId)
	if err != nil {
		var apiErr *apiError
		switch {
		case errors.Is(err, files.ErrFileNotFound):
			apiErr = notFound("file not found")

		default:
			log.Printf("failed to get file metadata: %v", err)
			apiErr = internalError("failed to get file", err)
		}
		apiErr.Write(w)
		return
	}

	fileReader, err := fh.fileService.Stream(req.Context(), fileMeta.ID)
	if err != nil {
		var apiErr *apiError
		switch {
		case errors.Is(err, files.ErrFileNotFound):
			apiErr = notFound("file not found")

		case errors.Is(err, files.ErrFileNotReady):
			apiErr = notFound("file not found")

		default:
			log.Printf("failed to stream file: %v", err)
			apiErr = internalError("failed to get file", err)
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
		apiErr := internalError("file uploaded successfully", err)
		apiErr.Write(w)
		return
	}

	resp := ok("", toFileResponses(files))
	resp.Write(w)
}
