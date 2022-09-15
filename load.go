package openapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/chanced/jsonpointer"
	"github.com/chanced/uri"
)

type resolver struct {
	primaryResource []byte
	resources       map[string][]byte
	fn              func(context.Context, *uri.URI) (Kind, []byte, error)
	uri             *uri.URI
	document        *Document
	// id rather have map[string]*Callback and so on, but that causes a circular dependency
	// and I have no idea transcodefmt
	components map[Kind]map[string]node
	nodes      map[string]node
}

// kind identifies source data's shape (e.g. JSON Schema or a second OpenAPI
// Document).
//
// Load will never ask for resolution of a URI with fragment. This only pertains
// to the data at the absolute, non fragmented URI.
//
// This is incredibly useful when loading data from a source as the result of a
// a reference (i.e. $ref, $dynamicRef) with a fragment (e.g.
// example.json#/$defs/foo). In this scenario, we know what $defs/foo is
// expected to be, but we may not know what example.json is yet, if ever.
//
// Knowing what the root document is prevents scenarios where we resolve
// "example.json#/foo/bar" and then later encounter a $ref to
// "example.json#/foo". Without knowing the shape of "example.json" defined, we
// would have to extract out "example.json#/foo/bar" from the raw json/yaml, and
// then reparse "#/foo" when we hit the second $ref. As a result, there would
// then exist two references to the same object within the graph.
//
// Finally, this is necessary for anchors (i.e. $anchor, $dynamicAnchor,
// $recursiveAnchor) above referenced external resources. For example, if we
// have a reference to "example.json#/foo/bar" which has an anchor "#baz", that
// is located at the root of "example.json", it will not be found and an error
// will be returned upon parsing.
func Load(ctx context.Context, documentURI string, compiler SchemaCompiler, fn func(context.Context, *uri.URI) (Kind, []byte, error)) (*Document, error) {
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
		return nil, fmt.Errorf("documentURI may not contain a fragment: received \"%s\"", docURI)
	}

	src, data, err := fn(ctx, docURI)

	_ = src

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

func (r *resolver) resolve(ctx context.Context, kind Kind, u *uri.URI, dst node) (node, error) {
	if u.Fragment != "" {
		nu := *u
		nu.Fragment = ""
		if n, ok := r.nodes[nu.String()]; ok {
			if u.Fragment == "" {
				return n, nil
			}
		}
	}
	panic("not done")
}

func (r *resolver) findTopMostNode(u *uri.URI) (node, error) {
	if u.Fragment != "" {
		// see if its a jsonpointer first
		ptr, err := jsonpointer.Parse(u.Fragment)
		_ = ptr
		if err != nil {
			// not a jsonpointer, so it might be an anchor
			// check to see if the top-level schema exists first
			c := *u
			c.Fragment = ""
			c.RawFragment = ""

			n, ok := r.nodes[c.String()]
			if ok {
				// if it isn't a schema, then we can't resolve it
				if n.Kind() != KindSchema {
					return nil, fmt.Errorf("error: the reference URI fragment contains an anchor %q but the top-level Node is not a Schema", u.Fragment)
				}
			}
			return nil, fmt.Errorf("failed to find node with id \"%s\"", u.Fragment)
		}
	}
	panic("not done")
}
