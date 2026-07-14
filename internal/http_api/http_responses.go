package http_api

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	status  int
	Message string    `json:"message,omitempty"`
	Error   *APIError `json:"error,omitempty"`
	Data    any       `json:"data,omitempty"`
}

func newAPIResponse(msg string, status int, err *APIError, data any) *APIResponse {
	return &APIResponse{
		Message: msg,
		status:  status,
		Error:   err,
		Data:    data,
	}
}

func Ok(msg string, data any) *APIResponse {
	return &APIResponse{
		Message: msg,
		status:  http.StatusOK,
		Data:    data,
	}
}

func Created(msg string, data any) *APIResponse {
	return &APIResponse{
		Message: msg,
		status:  http.StatusCreated,
		Data:    data,
	}
}

func (resp *APIResponse) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.status)
	json.NewEncoder(w).Encode(resp)
}

type APIError struct {
	Message string `json:"message,omitempty"`
	Status  int    `json:"status"`
	stack   string
}

func newAPIError(message string, status int, err error) *APIError {
	return &APIError{
		Message: message,
		Status:  status,
		stack:   err.Error(),
	}
}

func BadRequest(msg string) *APIError {
	return &APIError{
		Message: msg,
		Status:  http.StatusBadRequest,
	}
}

func NotFound(msg string) *APIError {
	return &APIError{
		Message: msg,
		Status:  http.StatusNotFound,
	}
}

func InternalError(msg string, err error) *APIError {
	return &APIError{
		Message: msg,
		Status:  http.StatusInternalServerError,
		stack:   err.Error(),
	}
}

func (aErr *APIError) Write(w http.ResponseWriter) {
	resp := newAPIResponse(aErr.Message, aErr.Status, aErr, nil)
	resp.Write(w)
}
