package errors

import (
	"errors"
	"fmt"
)

// Error provides information about a business error.
type Error struct {
	Code    uint32
	Message string
	Path    string
}

func (e *Error) Error() string {
	return fmt.Sprintf("message: %s, path: %s, code: %d", e.Message, e.Path, e.Code)
}

func (e *Error) extendPath(path string) error {
	e.Path = fmt.Sprintf("%s: %s", path, e.Path)
	return e
}

// NewInternal builds a new internal server error.
func NewInternal(path string, err error) error {
	return &Error{
		Path:    path,
		Message: err.Error(),
		Code:    InternalErrCode,
	}
}

// New builds a new Error.
func New(path, message string, code uint32) error {
	return &Error{
		Path:    path,
		Message: message,
		Code:    code,
	}
}

func ExtendPath(path string, err error) error {
	var e *Error
	if errors.As(err, &e) {
		return e.extendPath(path)
	}

	return fmt.Errorf("%s: %w", path, err)
}

func Is(err error, code uint32) bool {
	var e *Error
	if errors.As(err, &e) {
		if e.Code == code {
			return true
		}
	}

	return false
}
