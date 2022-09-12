package openapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/chanced/uri"
)

type resolver struct {
	primaryResource []byte
	resources       map[string][]byte
	fn              func(context.Context, *uri.URI) ([]byte, error)
	uri             *uri.URI
	document        *Document
	// id rather have map[string]*Callback and so on, but that causes a circular dependency
	// and I have no idea transcodefmt
	components componentNodes
	nodes      map[string]node
}

func (r *resolver) resolve(ctx context.Context, kind Kind, u *uri.URI, dst node) (node, error) {
	if u.Fragment != "" {
		nu := *u
		nu.Fragment = ""
		if n, ok := r.nodes[nu.String()]; ok {
			if u.Fragment == "" {
				return n, nil
			}
			panic("not done")

		}

	}
	panic("not done")
}

func resolveComponent[T node](ctx context.Context, c *Component[T], r resolver) error {
	if c.Reference == nil || (any)(c.Object) != nil {
		return nil
	}
	if c.Reference.Ref == nil || c.Reference.Ref.String() == "" {
		return fmt.Errorf("%w: %s", ErrEmptyRef, c.Reference.Location)
	}
	t, err := r.resolve(ctx, c.Object.Kind(), c.Reference.Ref, c.Object)
	if err != nil {
		return err
	}
	c.Object = (any)(t).(T)
	return nil
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
	var h http.Request
	h.Context()
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
	_ = loader
	panic("not done")
	// return loader, nil
}
