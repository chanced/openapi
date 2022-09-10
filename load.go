package openapi

import (
	"context"
	"fmt"

	"github.com/chanced/uri"
)

type resolver struct {
	primaryResource []byte
	resources       map[string][]byte
	fn              func(context.Context, *uri.URI) ([]byte, error)
	uri             *uri.URI
	document        *Document
	schemas         map[string]*Schema
	responses       map[string]*Response
	parameters      map[string]*Parameter
	Examples        map[string]*Example
	requestBodies   map[string]*RequestBody
	headers         map[string]*Header
	securitySchemes map[string]*SecurityScheme
	Links           map[string]*Link
	Callbacks       map[string]*Callbacks
	PathItems       map[string]*PathItem
}

// NewLoader returns a new Loader where documentURI is the URI of root OpenAPI document
// to load.
//
// NewLoader panics if fn is nil
func Load(ctx context.Context, documentURI string, fn func(context.Context, *uri.URI) ([]byte, error)) (*Document, error) {
	if fn == nil {
		panic("fn cannot be nil")
	}
	if documentURI == "" {
		return nil, fmt.Errorf("documentURI cannot be empty")
	}

	docURI, err := uri.Parse(documentURI)
	if err != nil {
		return nil, fmt.Errorf("failed to parse documentURI: %w", err)
	}

	if docURI.Fragment != "" {
		return nil, fmt.Errorf("documentURI must not contain a fragment: received \"%s\"", docURI)
	}

	data, err := fn(ctx, docURI)
	if err != nil {
		return nil, fmt.Errorf("failed to load OpenAPI Document: %w", err)
	}

	var doc Document
	if err := doc.UnmarshalJSON(data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal OpenAPI Document: %w", err)
	}
	loader := resolver{
		resources: map[string][]byte{
			docURI.String(): data,
		},
		fn:  fn,
		uri: docURI,
	}
	return loader, nil
}
