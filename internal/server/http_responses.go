package server

import (
	"encoding/json"
	"net/http"
)

type apiResponse struct {
	status  int
	Message string    `json:"message,omitempty"`
	Error   *apiError `json:"error,omitempty"`
	Data    any       `json:"data,omitempty"`
}

func newAPIResponse(msg string, status int, err *apiError, data any) *apiResponse {
	return &apiResponse{
		Message: msg,
		status:  status,
		Error:   err,
		Data:    data,
	}
}

func ok(msg string, data any) *apiResponse {
	return &apiResponse{
		Message: msg,
		status:  http.StatusOK,
		Data:    data,
	}
}

func created(msg string, data any) *apiResponse {
	return &apiResponse{
		Message: msg,
		status:  http.StatusCreated,
		Data:    data,
	}
}

func (resp *apiResponse) write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.status)
	json.NewEncoder(w).Encode(resp)
}

type apiError struct {
	Message string `json:"message,omitempty"`
	Status  int    `json:"status"`
	stack   string
}

func newAPIError(message string, status int, err error) *apiError {
	return &apiError{
		Message: message,
		Status:  status,
		stack:   err.Error(),
	}
}

func badRequest(msg string) *apiError {
	return &apiError{
		Message: msg,
		Status:  http.StatusBadRequest,
	}
}

func notFound(msg string) *apiError {
	return &apiError{
		Message: msg,
		Status:  http.StatusNotFound,
	}
}

func unauthorized(msg string) *apiError {
	return &apiError{
		Message: msg,
		Status: http.StatusUnauthorized,
	}
}


func internalError(msg string, err error) *apiError {
	return &apiError{
		Message: msg,
		Status:  http.StatusInternalServerError,
		stack:   err.Error(),
	}
}


func (aErr *apiError) write(w http.ResponseWriter) {
	resp := newAPIResponse(aErr.Message, aErr.Status, aErr, nil)
	resp.write(w)
}
