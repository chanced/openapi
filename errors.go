package openapi

import (
	"errors"
	"fmt"

	"github.com/chanced/jsonpointer"
	"github.com/chanced/uri"
)

var (
	ErrEmptyRef      = errors.New("openapi: empty $ref")
	ErrNotFound      = fmt.Errorf("openapi: component not found")
	ErrNotResolvable = errors.New("openapi: pointer path not resolvable")
)

type Error struct {
	Err error
	URI uri.URI
}

func NewError(err error, uri uri.URI) *Error {
	return &Error{
		Err: err,
		URI: uri,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Err, e.URI.String())
}

func (e *Error) Unwrap() error {
	return e.Err
}

func newErrNotFound(uri uri.URI, tok jsonpointer.Token) error {
	return NewError(fmt.Errorf("%w: %q", ErrNotFound, tok), uri)
}

func newErrNotResolvable(uri uri.URI, tok jsonpointer.Token) error {
	return NewError(fmt.Errorf("%w: %q", ErrNotResolvable, tok), uri)
}
