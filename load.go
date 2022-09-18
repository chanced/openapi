package openapi

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver"
	"github.com/chanced/jsonpointer"
	"github.com/chanced/uri"
)

type LoadOpts struct {
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
	l := newLoader(validator, fn, mergeLoadOpts(opts))
	n, err := l.load(ctx, *docURI, KindDocument, nil, nil)
	if err != nil {
		return nil, err
	}
	return n.(*Document), nil
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

func (l *loader) load(ctx context.Context, u uri.URI, ek Kind, v *semver.Version, d *uri.URI) (Node, error) {
	if n, ok := l.nodes[u]; ok {
		return n.node, nil
	}
	k, data, err := l.loadData(ctx, u, ek)
	if err != nil {
		return nil, err
	}
	switch k {
	case KindDocument:
		return l.loadDocument(ctx, data, u)
	case KindSchema:
		return l.loadSchema(ctx, data, u, *v)
	case KindCallbacks, KindExample, KindHeader, KindPathItem, KindOperation,
		KindRequestBody, KindResponse, KindLink, KindSecurityScheme:
		return l.loadNode(ctx, data, *v, *d)
	default:
		return nil, NewError(fmt.Errorf("loading %s as an external resource is not currently supported", k), u)
	}
}

func (l *loader) loadData(ctx context.Context, u uri.URI, ek Kind) (Kind, []byte, error) {
	k, d, err := l.fn(ctx, u, ek)
	if err != nil {
		return k, d, err
	}

	if k == KindUndefined && ek != KindUndefined {
		k = ek
	}
	if ek != KindUndefined && k != ek {
		return k, nil, NewError(fmt.Errorf("expected %s, but received %s", ek, k), u)
	}
	return k, d, nil
}

func (l *loader) loadDocument(ctx context.Context, data []byte, u uri.URI) (*Document, error) {
	var err error

	vs, ok := TryGetOpenAPIVersion(data)
	var v *semver.Version
	if ok {
		v, err = semver.NewVersion(vs)
		if err != nil {
			return nil, NewError(fmt.Errorf("failed to parse OpenAPI version: %w", err), u)
		}
	}
	if v == nil {
		return nil, NewError(fmt.Errorf("failed to determine OpenAPI version; ensure that the OpenAPI document has an openapi field"), u)
	}

	sd, err := l.getJSONSchemaDialect(data, v)
	if err != nil {
		return nil, NewError(fmt.Errorf("failed to determine OpenAPI schema dialect: %w", err), u)
	}
	if sd == nil {
		return nil, NewError(fmt.Errorf("failed to determine OpenAPI schema dialect"), u)
	}

	if err = l.validator.Validate(data, u, KindDocument, *v, *sd); err != nil {
		return nil, NewValidationError(err, KindDocument, u)
	}

	var doc Document
	loc, err := NewLocation(u)
	if err != nil {
		return nil, NewError(err, u)
	}
	if err := doc.UnmarshalJSON(data); err != nil {
		return nil, NewError(fmt.Errorf("failed to unmarshal OpenAPI Document: %w", err), u)
	}
	if err = doc.setLocation(loc); err != nil {
		return nil, NewError(err, u)
	}
	l.nodes[u] = nodectx{
		node:       &doc,
		openapi:    doc.OpenAPI,
		jsonschema: *sd,
	}

	if err = l.traverse(&doc, doc.edges(), *v, *sd); err != nil {
		return nil, err
	}
	// we only traverse the references after the top-level document is fully
	// materialized.
	if l.doc == nil {
		l.doc = &doc
	} else {
		return &doc, nil
	}

	var r refctx
	var nodes []nodectx
	for len(l.refs) > 0 {
		for len(l.refs) > 0 {
			r, l.refs = l.refs[0], l.refs[1:]
			n, err := l.resolveRef(ctx, r)
			nodes = append(nodes, n)
			if err != nil {
				return nil, err
			}
		}
		for _, n := range nodes {
			if err = l.traverse(n.root, n.edges(), n.openapi, n.jsonschema); err != nil {
				return nil, err
			}
		}
		nodes = nil
	}
	if err = l.validator.ValidateDocument(&doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

func (l *loader) resolveRef(ctx context.Context, r refctx) (nodectx, error) {
	u := r.ref.URI()
	if u == nil {
		return nodectx{}, NewValidationError(fmt.Errorf("error: ref URI cannot be empty"), r.Kind(), r.AbsolutePath())
	}

	// check to see if the node has already been loaded
	if pn, ok := l.nodes[*u]; ok {
		// if so, return it
		return pn, nil
	}

	// check to see if the node has fragment if so, we need to load the data and
	// then traverse to the fragment.
	if u.Fragment != "" || u.RawFragment != "" {
		var a string
		// The fragment should be a jsonpointer unless it the referenced node is
		// schema. If so, it can be either a json pointer or an anchor
		//
		// checking to see if its a json pointer first
		p, err := jsonpointer.Parse(u.Fragment)
		_ = p
		if err != nil {
			// if its not a json pointer, then it could be an anchor if the node kind is a schema
			if r.Kind() == KindSchema {
				a = u.Fragment
			} else {
				return nodectx{}, NewValidationError(fmt.Errorf("ref URI fragment must be a json pointer: %w", err), r.Kind(), r.AbsolutePath())
			}
		}
		// checking to see if this is a ref relative to the current resource
		if u.Host == "" && u.Path == "" {
			if a != "" {
				as, err := r.Anchors()
				if err != nil {
					return nodectx{}, NewValidationError(fmt.Errorf("failed to get anchors: %w", err), r.Kind(), r.AbsolutePath())
				}
				if as == nil {
					return nodectx{}, NewError(fmt.Errorf("failed to retrieve anchors"), r.AbsolutePath())
				}
			}
			r.ref.re
		}

		// now if the resource is
	}
	uc := *u
	uc.Fragment = ""
	uc.RawFragment = ""

	if rn, ok := l.nodes[uc]; ok {
		_ = rn
		// this node has already been loaded
		// if a != "" {
		// 	n, err := rn.ResolveNodeByPointer()
		// 	_ = n
		// 	if err != nil {
		// 		return nodectx{}, NewValidationError(fmt.Errorf("failed to resolve ref URI fragment: %w", err), r.Kind(), r.AbsolutePath())
		// 	}

		// }

		// checking to see if the loader has encountered the referenced root resource
		var nk Kind
		for _, r := range l.refs {
			if *r.ref.URI() == uc {
				nk = r.Kind()
				break
			}
		}
		if nk != KindUndefined {
			// the loader has encountered the referenced root resource so we can
			// load it first
			k, d, err := l.loadData(ctx, uc, nk)
			_ = d
			if err != nil {
				return nodectx{}, err
			}

			switch k {
			case KindDocument:
				// doc, err := l.load(ctx, )

			case KindSchema:

			}
		}

	}
	panic("not done")
}

func (l *loader) getDocumentSchemaDialect(doc *Document) (*uri.URI, error) {
	if doc.JSONSchemaDialect != nil {
		return doc.JSONSchemaDialect, nil
	}
	if l.opts.DefaultSchemaDialect != nil {
		return l.opts.DefaultSchemaDialect, nil
	}
	if VersionConstraints3_1.Check(&doc.OpenAPI) {
		return &JSONSchemaDialect202012, nil
	}
	if VersionConstraints3_0.Check(&doc.OpenAPI) {
		return &JSONSchemaDialect201909, nil
	}
	return nil, fmt.Errorf("failed to determine OpenAPI schema dialect")
}

func (l *loader) traverse(root Node, nodes []node, openapi semver.Version, jsonschema uri.URI) error {
	for _, n := range nodes {
		nc, err := newNodeCtx(n, &openapi, &jsonschema)
		if err != nil {
			return err
		}
		l.nodes[n.AbsolutePath()] = nc

		if !IsRef(n) {
			if err := l.traverse(root, n.edges(), nc.openapi, nc.jsonschema); err != nil {
				return err
			}
		} else {
			r := n.(ref)
			l.refs = append(l.refs, refctx{root: root, ref: r, openapi: nc.openapi, jsonschema: nc.jsonschema})

		}
	}

	return nil
}

func (l *loader) loadSchema(ctx context.Context, data []byte, u uri.URI, v semver.Version) (*Schema, error) {
	var s Schema
	if err := s.UnmarshalJSON(data); err != nil {
		return nil, NewError(fmt.Errorf("failed to unmarshal JSON Schema: %w", err), u)
	}
	return &s, nil
}

func (l *loader) loadNode(ctx context.Context, data []byte, v semver.Version, s uri.URI) (nodectx, error) {
	panic("not implemented")
}

func (l *loader) getJSONSchemaDialect(data []byte, v *semver.Version) (*uri.URI, error) {
	sds, ok := TryGetSchemaDialect(data)
	var sd *uri.URI
	var err error
	switch {
	case ok:
		sd, err = uri.Parse(sds)
		if err != nil {
			return nil, fmt.Errorf("failed to parse JSON Schema dialect: %w", err)
		}
	case l.opts.DefaultSchemaDialect != nil:
		sd = l.opts.DefaultSchemaDialect
	case checkVersion(VersionConstraints3_1, v):
		sd = &JSONSchemaDialect202012
	case checkVersion(VersionConstraints3_0, v):
		sd = &JSONSchemaDialect201909
	default:
		return nil, nil
	}
	return sd, nil
}

// func (l *loader) resolve(ctx context.Context, kind Kind, u uri.URI, dst node) (node, error) {
// 	if n, ok := l.nodes[u]; ok {
// 		return n, nil
// 	}

// 	src, d, err := l.fn(ctx, u)
// 	if err != nil {
// 		return nil, NewError(fmt.Errorf("failed to load %s: %w", u, err), u)
// 	}
// 	switch src {
// 	case KindDocument:
// 	}
// }

type nodectx struct {
	node
	openapi    semver.Version
	jsonschema uri.URI
	root       Node
}
type refctx struct {
	ref
	root       Node
	openapi    semver.Version
	jsonschema uri.URI
}

func newNodeCtx(n node, openapi *semver.Version, jsonschema *uri.URI) (nodectx, error) {
	switch t := n.(type) {
	case *Document:
		panic("not implemented")
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
