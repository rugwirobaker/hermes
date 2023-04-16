package hermes

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrAlreadyExists = errors.New("already exists")
)

type ErrInvalid struct {
	// Message of the error
	Message string `json:"message"`
	//
	err error `json:"-"`
}

func (e ErrInvalid) Error() string {
	return fmt.Sprintf("%s: %s", e.Message, e.err)
}

func NewErrInvalid(text string) error {
	return &ErrInvalid{Message: text}
}

func NewErrInvalidWithErr(text string, err error) error {
	return &ErrInvalid{Message: text, err: err}
}
