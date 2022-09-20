package openapi

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver"
	"github.com/chanced/transcode"
	"github.com/chanced/uri"
	"github.com/sanity-io/litter"
	"github.com/tidwall/gjson"
)

// TryGetSchemaDialect attempts to extract the schema dialect from raw JSON
// data.
//
// TryGetSchemaDialect will check the following fields in order:
//   - $schema
//   - jsonSchemaDialect
func TryGetSchemaDialect(data []byte) (string, bool) {
	id := gjson.GetBytes(data, "$schema")
	if id.Exists() {
		return id.String(), true
	}
	id = gjson.GetBytes(data, "jsonSchemaDialect")
	if id.Exists() {
		return id.String(), true
	}
	return "", false
}

// TryGetOpenAPIVersion attempts to extract the OpenAPI version from raw JSON
// data and parse it as a semver.Version.
func TryGetOpenAPIVersion(data []byte) (string, bool) {
	v := gjson.GetBytes(data, "openapi")
	if v.Exists() {
		return v.String(), true
	}
	return "", false
}

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
	nodes := make(map[string]nodectx)
	return &loader{
		validator: v,
		fn:        fn,
		nodes:     nodes,
		opts:      opts,
	}
}

type loader struct {
	opts          LoadOpts
	fn            func(context.Context, uri.URI, Kind) (Kind, []byte, error)
	validator     Validator
	doc           *Document
	nodes         map[string]nodectx
	refs          []refctx
	recursiveRefs []refctx
	dynamicRefs   []refctx
}

func (l *loader) load(ctx context.Context, location uri.URI, ek Kind, openapi *semver.Version, dialect *uri.URI) (Node, error) {
	if n, ok := l.nodes[location.String()]; ok {
		return n.node, nil
	}
	k, data, err := l.loadData(ctx, location, ek)
	if err != nil {
		return nil, err
	}
	switch k {
	case KindDocument:
		return l.loadDocument(ctx, data, location)
	case KindSchema:
		return l.loadSchema(ctx, data, location, *openapi)
	case KindCallbacks, KindExample, KindHeader, KindPathItem, KindOperation,
		KindRequestBody, KindResponse, KindLink, KindSecurityScheme:
		return l.loadNode(ctx, data, *openapi, *dialect)
	default:
		return nil, NewError(fmt.Errorf("loading %s as an external resource is not currently supported", k), location)
	}
}

func (l *loader) loadData(ctx context.Context, u uri.URI, ek Kind) (Kind, []byte, error) {
	k, d, err := l.fn(ctx, u, ek)
	if err != nil {
		return k, d, err
	}

	d, err = transcode.JSONFromYAML(d)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to transcode data: %w", err)
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
	l.nodes[u.String()] = nodectx{
		node:       &doc,
		openapi:    *doc.OpenAPI,
		jsonschema: *sd,
	}
	dc := &nodectx{
		node:       &doc,
		openapi:    *doc.OpenAPI,
		jsonschema: *sd,
	}
	dc.root = dc

	if err = l.traverse(dc, doc.nodes(), *v, *sd); err != nil {
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
			if n != nil {
				nodes = append(nodes, *n)
			}
			if err != nil {
				return nil, err
			}
		}
		for _, n := range nodes {
			if err = l.traverse(n.root, n.nodes(), n.openapi, n.jsonschema); err != nil {
				return nil, err
			}
		}
		nodes = nil
	}
	for len(l.recursiveRefs) > 0 {
		panic("not done")
	}
	for len(l.dynamicRefs) > 0 {
		panic("not done")
	}
	if err = l.validator.ValidateDocument(&doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

func (l *loader) resolveRef(ctx context.Context, r refctx) (*nodectx, error) {
	u := r.URI()

	if u == nil {
		return nil, NewValidationError(fmt.Errorf("error: ref URI cannot be empty"), r.Kind(), r.AbsoluteLocation())
	}

	if u.Host == "" && u.Path == "" {
		return l.resolveLocalRef(ctx, r)
	} else {
		return l.resolveRemoteRef(ctx, r)
	}
}

func (l *loader) resolveRemoteRef(ctx context.Context, r refctx) (*nodectx, error) {
	u := r.URI()

	// let's see if we've already loaded this node
	ur := *u
	ur.Fragment = ""
	ur.RawFragment = ""
	if _, ok := l.nodes[ur.String()]; ok {
		switch r.RefType() {
		case RefTypeSchemaDynamicRef:
			l.dynamicRefs = append(l.dynamicRefs, r)
			return nil, nil
		case RefTypeSchemaRecursiveRef:
			l.recursiveRefs = append(l.recursiveRefs, r)
			return nil, nil
		}
		// then we should have the node in stock
		if n, ok := l.nodes[u.String()]; ok {
			return &n, r.resolve(n.node)
		} else {
			// otherwise something went awry
			return nil, NewError(fmt.Errorf("error: ref URI not found: %s", u), r.AbsoluteLocation())
		}
	}
	// otherwise we need to load the root resource
	// first we need to load the resource

	// we need to check to see if there is a reference pointing to the root first
	// so we know what the expected type is
	urs := ur.String()
	for _, x := range l.refs {
		if x.URI().String() == urs {
			// found it. we load that one first.
			if _, err := l.load(ctx, ur, x.RefKind(), nil, nil); err != nil {
				return nil, err
			}
			break
		}
	}

	// now check to see if we've found it.
	if _, ok := l.nodes[ur.String()]; !ok {
		if _, err := l.load(ctx, ur, KindUndefined, &r.openapi, &r.jsonschema); err != nil {
			return nil, err
		}
	}
	// checking to make sure the root node is loaded
	_, ok := l.nodes[ur.String()]
	if !ok {
		// otherwise we return an error
		return nil, NewError(fmt.Errorf("error: ref URI not found: %s", u), r.AbsoluteLocation())
	}
	switch r.RefType() {
	// since we found it, we need to bail if this is a dynamicRef or a recursiveRef
	case RefTypeSchemaDynamicRef:
		l.dynamicRefs = append(l.dynamicRefs, r)
		return nil, nil
	case RefTypeSchemaRecursiveRef:
		l.recursiveRefs = append(l.recursiveRefs, r)
		return nil, nil
	}

	// we should have the node in stock
	if n, ok := l.nodes[u.String()]; ok {
		return &n, r.resolve(n.node)
	}

	// otherwise we may be dealing with an anchor
	a := Text(u.Fragment).TrimPrefix("#")

	rn, ok := l.nodes[ur.String()]
	if !ok {
		return nil, NewError(fmt.Errorf("error: ref URI not found: %s", u), r.AbsoluteLocation())
	}

	if a == "" {
		// we should have already resolved it?
		if err := r.resolve(rn.node); err != nil {
			return nil, err
		}
		return &rn, nil
	}

	if a.HasPrefix("/") {
		// we aren't dealing with an anchor
		return nil, NewError(fmt.Errorf("error: ref URI not found: %s", u), r.AbsoluteLocation())
	}

	as, err := rn.Anchors()
	if err != nil {
		return nil, NewError(fmt.Errorf("error: failed to resolve anchors: %w", err), r.AbsoluteLocation())
	}

	an, ok := as.Standard[a]
	if !ok {
		return nil, NewError(fmt.Errorf("error: ref URI not found: %s", u), r.AbsoluteLocation())
	}

	x, ok := l.nodes[an.AbsoluteLocation().String()]
	if !ok {
		return nil, NewError(fmt.Errorf("error: ref URI not found: %s", u), r.AbsoluteLocation())
	}
	if err := r.resolve(x.node); err != nil {
		return nil, err
	}
	return &x, nil
}

func (l *loader) resolveLocalRef(ctx context.Context, r refctx) (*nodectx, error) {
	u := r.AbsoluteLocation()
	u.Fragment = r.URI().Fragment
	u.RawFragment = r.URI().RawFragment

	// if this is a $recursiveRef or a $dynamicRef, we need to cycle
	// through all the refs first so we kick the can down the road.
	switch r.RefType() {
	case RefTypeSchemaDynamicRef:
		l.dynamicRefs = append(l.dynamicRefs, r)
		return nil, nil
	case RefTypeSchemaRecursiveRef:
		l.recursiveRefs = append(l.recursiveRefs, r)
		return nil, nil
	}
	// check to see if this node has already been loaded
	if n, ok := l.nodes[u.String()]; ok {
		// resolve it and move along
		if err := r.resolve(n.node); err != nil {
			return nil, err
		}
		return &n, nil
	} else {

		for k, v := range l.nodes {
			if k == u.String() {
				fmt.Println("found it:", k, v)
			}
		}
		// otherwise something went awry
		return nil, NewError(fmt.Errorf("error: ref URI not found: %s", u), r.AbsoluteLocation())
	}
}

func (l *loader) getDocumentSchemaDialect(doc *Document) (*uri.URI, error) {
	if doc.JSONSchemaDialect != nil {
		return doc.JSONSchemaDialect, nil
	}
	if l.opts.DefaultSchemaDialect != nil {
		return l.opts.DefaultSchemaDialect, nil
	}
	if VersionConstraints3_1.Check(doc.OpenAPI) {
		return &JSONSchemaDialect202012, nil
	}
	if VersionConstraints3_0.Check(doc.OpenAPI) {
		return &JSONSchemaDialect201909, nil
	}
	return nil, fmt.Errorf("failed to determine OpenAPI schema dialect")
}

func (l *loader) traverse(root *nodectx, nodes []node, openapi semver.Version, jsonschema uri.URI) error {
	for _, n := range nodes {
		nc, err := newNodeCtx(n, &openapi, &jsonschema)
		if err != nil {
			return err
		}
		if n.AbsoluteLocation().String() == "#" {
			fmt.Println("\nnode with empty absolute location:")
			litter.Dump(n)
		}
		l.nodes[n.AbsoluteLocation().String()] = nc

		if IsRef(n) {
			r := n.(ref)
			if !r.IsResolved() {
				l.refs = append(l.refs, refctx{root: root, ref: r, openapi: nc.openapi, jsonschema: nc.jsonschema})
			}
			return nil
		}
		if err := l.traverse(root, n.nodes(), nc.openapi, nc.jsonschema); err != nil {
			return err
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

type nodectx struct {
	node
	openapi    semver.Version
	jsonschema uri.URI
	root       *nodectx
}
type refctx struct {
	ref
	root       *nodectx
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
