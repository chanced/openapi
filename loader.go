package openapi

import (
	"fmt"

	"github.com/chanced/uri"
)

type Loader struct {
	primaryResource []byte
	resources       map[string][]byte
	fn              func(*uri.URI) ([]byte, error)
	uri             *uri.URI
}

// NewLoader returns a new Loader where documentURI is the URI of root OpenAPI document
// to load.
//
// NewLoader panics if fn is nil or documentURI is an empty string
func NewLoader(documentURI string, fn func(*uri.URI) ([]byte, error)) (Loader, error) {
	if fn == nil {
		panic("fn cannot be nil")
	}
	if documentURI == "" {
		panic("documentURI cannot be empty")
	}

	docURI, err := uri.Parse(documentURI)
	if err != nil {
		return Loader{}, fmt.Errorf("failed to parse documentURI: %w", err)
	}

	if docURI.Fragment != "" {
		return Loader{}, fmt.Errorf("documentURI must not contain a fragment: received \"%s\"", docURI)
	}

	doc, err := fn(docURI)
	if err != nil {
		return Loader{}, fmt.Errorf("failed to load OpenAPI Document: %w", err)
	}
	loader := Loader{
		resources: map[string][]byte{
			docURI.String(): doc,
		},
		fn:  fn,
		uri: docURI,
	}
	return loader, nil
}
