package http

import (
	"github.com/go-chi/render"
	"net/http"
)

var _ render.Renderer = (*ErrorResponse)(nil)
var _ render.Renderer = (*IndexFoundResponse)(nil)

// ErrorResponse is a general response for errors
type ErrorResponse struct {
	statusCode int
	Message    string `json:"message,omitempty"`
}

func (e ErrorResponse) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.statusCode)
	return nil
}

var ErrIndexNotFound = ErrorResponse{
	statusCode: http.StatusNotFound,
	Message:    "index not found",
}

var ErrInvalidValue = ErrorResponse{
	statusCode: http.StatusBadRequest,
	Message:    "invalid value passed",
}

var ErrUnknown = ErrorResponse{
	statusCode: http.StatusInternalServerError,
	Message:    "unexpected error happened",
}

// IndexFoundResponse is a response for successful index search
type IndexFoundResponse struct {
	Index int   `json:"index"`
	Value int64 `json:"value"`
}

func (i IndexFoundResponse) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, http.StatusOK)
	return nil
}
