package openapi

import (
	"github.com/chanced/jsonpointer"
	"github.com/chanced/uri"
)

type Location struct {
	Absolute *uri.URI
	Relative jsonpointer.Pointer
}
