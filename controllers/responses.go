package controllers

import (
	"net/http"

	"github.com/go-chi/render"
)

// ErrResponse ...
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

// Render ...
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// ErrBadRequest ...
func ErrBadRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "bad request",
		ErrorText:      err.Error(),
	}
}

// ErrUnprocessableEntity ...
func ErrUnprocessableEntity(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "unprocessable entity",
		ErrorText:      err.Error(),
	}
}

// ErrNotFound ...
var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "resource not found"}

// ErrWrongPassword ...
var ErrWrongPassword = &ErrResponse{HTTPStatusCode: 401, StatusText: "wrong password"}

// ErrInternalServer ...
func ErrInternalServer(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 500,
		StatusText:     "internal server error",
		ErrorText:      err.Error(),
	}
}
