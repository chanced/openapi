package openapi

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver"
	"github.com/chanced/uri"
	"github.com/go-playground/validator/v10"
)

type LoadOpts struct {
	// DefaultSchemaDialect is used when a document does not define a schema
	// dialect
	DefaultSchemaDialect *uri.URI
}

func mergeLoadOpts(opts []LoadOpts) LoadOpts {
	var l LoadOpts
	for _, o := range opts {
		if o.DefaultSchemaDialect != nil {
			l.DefaultSchemaDialect = o.DefaultSchemaDialect
		}
	}
	return l
}

// Load loads an OpenAPI document from a URI and validate it with the provided
// validator.
//
// Loading the raw data for OpenAPI Documents and externally referenced
// referenced JSON Schema components is done through the anonymous function fn.
// It is passed the URI of the resource and if known, the expected
// Kind. fn should return the Kind for the resource and the raw data if
// successful.
//
// Resources that can be referenced are:
//   - OpenAPI Document (KindDocument)
//   - JSON Schema (KindSchema)
//   - Components (KindComponents)
//   - Callbacks (KindCallbacks)
//   - Example (KindExample)
//   - Header (KindHeader)
//   - Link (KindLink)
//   - Parameter (KindParameter)
//   - PathItem (KindPathItem)
//   - Operation (KindOperation)
//   - Reference (KindReference)
//   - RequestBody (KindRequestBody)
//   - Response (KindResponse)
//   - SecurityScheme (KindSecurityScheme)
//
// fn will invoke fn with a URI containing a fragment; it will only ever
// be called to resolve to the root document data. This is why Kind must be
// returned from fn, as there may not be enough context to infer the shape of
// the data.
//
// Knowing the shape of root document prevents scenarios where we resolve
// "example.json#/foo/bar" and then later encounter a $ref to
// "example.json#/foo". Without knowing the shape of "example.json", we would
// have to extract out "example.json#/foo/bar" from the raw json/yaml, and then
// reparse "#/foo" when we hit the second $ref. As a result, there would then
// exist two references to the same object within the graph.
//
// Finally, being able to parse the root resource is necessary for anchors (i.e.
// $anchor, $dynamicAnchor, $recursiveAnchor) above referenced external
// resources. For example, if we have a reference to "example.json#/foo/bar"
// which has an anchor "#baz", that is located at the root of "example.json", it
// would not be found if example.json were not parsed entirely.
func Load(ctx context.Context, documentURI string, validator Validator, fn func(ctx context.Context, uri uri.URI, kind Kind) (Kind, []byte, error), opts ...LoadOpts) (*Document, error) {
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
		return nil, NewError(fmt.Errorf("documentURI may not contain a fragment: received \"%s\"", docURI), *docURI)
	}
	// du := *docURI
	// src, d, err := fn(ctx, du, KindDocument)
	// if err != nil {
	// 	return nil, NewError(fmt.Errorf("failed to load OpenAPI Document: %w", err), *docURI)
	// }
	// if src != KindDocument {
	// 	return nil, NewError(fmt.Errorf("documentURI must be an OpenAPI document: received \"%s\"", docURI), *docURI)
	// }

	l := newLoader(validator, fn, mergeLoadOpts(opts))
	n, err := l.load(ctx, *docURI, KindDocument)
	if err != nil {
		return nil, err
	}
	if doc, ok := n.(*Document); ok {
		return doc, nil
	} else {
		// this should never happen
		panic("node returned from load was not a *Document. This is a bug. Please report it to https://github.com/chanced/openapi/issues/new")
	}
}

func newLoader(v Validator, fn func(context.Context, uri.URI, Kind) (Kind, []byte, error), opts LoadOpts) *loader {
	nodes := make(map[uri.URI]nodectx)
	return &loader{
		validator: v,
		fn:        fn,
		nodes:     nodes,
		opts:      opts,
	}
}

type loader struct {
	opts      LoadOpts
	fn        func(context.Context, uri.URI, Kind) (Kind, []byte, error)
	validator Validator
	doc       *Document
	nodes     map[uri.URI]nodectx
	refs      []refctx
}

func (l *loader) load(ctx context.Context, docURI uri.URI, kind Kind) (Node, error) {
	var err error
	vers, ok := TryGetOpenAPIVersion(data)
	var v *semver.Version
	if ok {
		v, err = semver.NewVersion(vers)
		if err != nil {
			return nil, NewError(fmt.Errorf("failed to parse OpenAPI version: %w", err), *docURI)
		}
	}

	ds, ok := TryGetSchemaDialect(data)
	var dialect *uri.URI
	switch {
	case ok:
		dialect, err = uri.Parse(ds)
		if err != nil {
			return nil, NewError(fmt.Errorf("failed to parse JSON Schema dialect: %w", err), *docURI)
		}
	case checkVersion(SemanticVersion3_0, v):
		dialect = &JSONSchemaDialect202012
	case checkVersion(SemanticVersion3_1, v):
		dialect = &JSONSchemaDialect201909
	default:
		return nil, NewError(fmt.Errorf("failed to determine JSON Schema dialect"), *docURI)
	}

	if err = validator.Validate(data, KindDocument, *v, *dialect); err != nil {
		return nil, NewError(err, *docURI)
	}

	var doc Document
	loc, err := NewLocation(*docURI)
	if err != nil {
		return nil, NewError(err, *docURI)
	}
	if err := doc.UnmarshalJSON(data); err != nil {
		return nil, NewError(fmt.Errorf("failed to unmarshal OpenAPI Document: %w", err), *docURI)
	}
	if err = doc.setLocation(loc); err != nil {
		return nil, NewError(err, *docURI)
	}
}

func (l *loader) traverse(nodes []node, openapi semver.Version, jsonschema *uri.URI) error {
	for _, n := range nodes {
		nc, err := newNodeCtx(n, openapi, jsonschema)
		if err != nil {
			return err
		}
		l.nodes[n.AbsolutePath()] = nc
		if !IsRef(n) {
			if err := l.traverse(n.edges(), &nc.openapi, &nc.jsonschema); err != nil {
				return err
			}
		} else {
			r := n.(ref)
			l.refs = append(l.refs, refctx{ref: r, openapi: nc.openapi, jsonschema: nc.jsonschema})
		}
	}
	return nil
}

func (l *loader) resolve(ctx context.Context, kind Kind, u uri.URI, dst node) (node, error) {
	if n, ok := l.nodes[u]; ok {
		return n, nil
	}

	src, d, err := l.fn(ctx, u)
	if err != nil {
		return nil, NewError(fmt.Errorf("failed to load %s: %w", u, err), u)
	}
	switch src {
	case KindDocument:
	}
}

type nodectx struct {
	node
	openapi    semver.Version
	jsonschema uri.URI
}
type refctx struct {
	ref
	openapi    semver.Version
	jsonschema uri.URI
}

func newNodeCtx(n node, openapi *semver.Version, jsonschema *uri.URI) (nodectx, error) {
	switch t := n.(type) {
	case *Document:
		switch {
		case t.JSONSchemaDialect != nil:
			jsonschema = t.JSONSchemaDialect
		case SemanticVersion3_1.Check(&t.OpenAPI):
			jsonschema = &JSONSchemaDialect202012
		case SemanticVersion3_1.Check(&t.OpenAPI):
			jsonschema = &JSONSchemaDialect201909
		}
		if jsonschema == nil {
			return nodectx{}, fmt.Errorf("failed to determine jsonschema dialect for Document", t.AbsolutePath())
		}
		if len(t.OpenAPI.String()) == 0 {
			return nodectx{}, fmt.Errorf("OpenAPI version must be defined")
		}
	case *Schema:
		if t.Schema != nil {
			jsonschema = t.Schema
		}
	}
	if jsonschema == nil {
		return nodectx{}, fmt.Errorf("failed to determine JSON Schema dialect")
	}
	if openapi == nil {
		return nodectx{}, fmt.Errorf("failed to determine OpenAPI version")
	}
	return nodectx{
		node:       n,
		jsonschema: *jsonschema,
		openapi:    *openapi,
	}, nil
}

func checkVersion(c semver.Constraints, v *semver.Version) bool {
	if v == nil {
		return false
	}
	return c.Check(v)
}
