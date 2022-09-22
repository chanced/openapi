package openapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/chanced/transcode"
	"github.com/chanced/uri"
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
	opts        LoadOpts
	fn          func(context.Context, uri.URI, Kind) (Kind, []byte, error)
	validator   Validator
	doc         *Document
	nodes       map[string]nodectx
	dynamicRefs []refctx
	refs        []refctx
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
		return l.loadNode(ctx, k, data, *openapi, *dialect)
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

	dc := nodectx{
		node:       &doc,
		openapi:    *doc.OpenAPI,
		jsonschema: *sd,
	}
	dc.root = &dc
	anchors, err := doc.Anchors()
	if err != nil {
		return nil, NewError(fmt.Errorf("failed to get anchors: %w", err), u)
	}
	dc.anchors = anchors

	l.nodes[u.String()] = dc
	if err = l.init(&dc, &dc, doc.nodes(), *v, *sd); err != nil {
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
			// r, l.refs = l.refs[len(l.refs)-1], l.refs[:len(l.refs)-1]
			r, l.refs = l.refs[0], l.refs[1:]
			n, err := l.resolveRef(ctx, r)
			if err != nil {
				return nil, err
			}
			if n != nil {
				nodes = append(nodes, *n)
				r.resolved = n
			}

			r.root.resolvedRefs = append(r.root.resolvedRefs, r)
		}
		for _, n := range nodes {
			if err = l.init(&dc, n.root, n.nodes(), n.openapi, n.jsonschema); err != nil {
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

func (l *loader) resolveRef(ctx context.Context, r refctx) (*nodectx, error) {
	u := r.URI()

	if u == nil {
		return nil, NewValidationError(fmt.Errorf("openapi: ref URI cannot be empty"), r.Kind(), r.AbsoluteLocation())
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
	rooturi := *u
	au := r.AbsoluteLocation()
	if u.Host == "" {
		rur := au.ResolveReference(u)
		rooturi = *rur
	}
	rooturi.Fragment = ""
	rooturi.RawFragment = ""

	if _, ok := l.nodes[rooturi.String()]; ok {
		switch r.RefType() {
		case RefTypeSchemaDynamicRef:
			r.root.dynamicRefs = append(r.root.dynamicRefs, r)
		case RefTypeSchemaRecursiveRef:
			r.root.recursiveRefs = append(r.root.recursiveRefs, r)
		}
		// then we should have the node in stock
		uc := r.URI()
		if u.Host == "" {
			uc = au.ResolveReference(u)
		}
		if n, ok := l.nodes[uc.String()]; ok {
			return &n, r.resolve(n.node)
		} else if u.Fragment == "" || strings.HasPrefix(u.Fragment, "/") {
			// something went sideways
			return nil, NewError(fmt.Errorf("openapi: ref URI not found: %s", u), r.AbsoluteLocation())
		}
	} else {
		// we need to load the root resource first we need to load the resource
		// we need to check to see if there is a reference pointing to the root first
		// so we know what the expected type is
		rus := rooturi.String()
		for _, x := range l.refs {
			if x.URI().String() == rus {
				// found it. we load that one first.
				if _, err := l.load(ctx, rooturi, x.RefKind(), nil, nil); err != nil {
					return nil, err
				}
				break
			}
		}

		// now check to see if we've found it.
		if _, ok := l.nodes[rooturi.String()]; !ok {
			if _, err := l.load(ctx, rooturi, KindUndefined, &r.openapi, &r.jsonschema); err != nil {
				return nil, err
			}
		}
		// checking to make sure the root node is loaded
		_, ok := l.nodes[rooturi.String()]
		if !ok {
			// otherwise we return an error
			return nil, NewError(fmt.Errorf("openapi: ref URI not found: %s", u), r.AbsoluteLocation())
		}
	}

	switch r.RefType() {
	case RefTypeSchemaDynamicRef:
		r.root.dynamicRefs = append(r.root.dynamicRefs, r)
	case RefTypeSchemaRecursiveRef:
		r.root.recursiveRefs = append(r.root.recursiveRefs, r)
	}

	// resetting ur
	rooturi = *u
	if u.Host == "" {
		rur := au.ResolveReference(u)
		rooturi = *rur
		rooturi.Fragment = u.Fragment
	}
	// we check to see if the node is in stock
	us := rooturi.String()
	if n, ok := l.nodes[us]; ok {
		return &n, r.resolve(n.node)
	}
	if u.Fragment == "" {
		return nil, NewError(fmt.Errorf("openapi: ref URI not found: %s", u), r.AbsoluteLocation())
	}

	// otherwise we may be dealing with an anchor
	a := Text(u.Fragment)

	rn, ok := l.nodes[rooturi.String()]
	if !ok {
		return nil, NewError(fmt.Errorf("openapi: ref URI not found: %s", u), r.AbsoluteLocation())
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
		return nil, NewError(fmt.Errorf("openapi: ref URI not found: %s", u), r.AbsoluteLocation())
	}

	as, err := rn.Anchors()
	if err != nil {
		return nil, NewError(fmt.Errorf("openapi: failed to resolve anchors: %w", err), r.AbsoluteLocation())
	}

	an, ok := as.Standard[a]
	if !ok {
		return nil, NewError(fmt.Errorf("openapi: ref URI not found: %s", u), r.AbsoluteLocation())
	}

	x, ok := l.nodes[an.AbsoluteLocation().String()]
	if !ok {
		return nil, NewError(fmt.Errorf("openapi: ref URI not found: %s", u), r.AbsoluteLocation())
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
		r.root.dynamicRefs = append(r.root.dynamicRefs, r)
	case RefTypeSchemaRecursiveRef:
		r.root.recursiveRefs = append(r.root.recursiveRefs, r)
	}
	// check to see if this node has already been loaded
	if n, ok := l.nodes[u.String()]; ok {
		// resolve it and move along
		if err := r.resolve(n.node); err != nil {
			return nil, err
		}
		return &n, nil
	} else if strings.HasPrefix(u.Fragment, "/") || r.ref.RefKind() != KindSchema {
		// otherwise something went awry
		return nil, NewError(fmt.Errorf("openapi: ref URI not found: %s", u), r.AbsoluteLocation())
	}

	// we are dealing with an anchor

	switch r.RefType() {
	case RefTypeSchemaDynamicRef:
		a, ok := r.root.anchors.Dynamic[Text(r.URI().Fragment)]
		if !ok {
			return nil, NewError(fmt.Errorf("openapi: ref URI not found: %s", u), r.AbsoluteLocation())
		}
		err := r.resolve(a.In)
		if err != nil {
			return nil, fmt.Errorf("openapi: failed to resolve node for anchor \"#%s\": %w", r.URI().Fragment, err)
		}
		return &nodectx{
			node:       a.In,
			openapi:    r.openapi,
			jsonschema: r.jsonschema,
			root:       r.root,
			anchors:    r.root.anchors,
		}, nil
	case RefTypeSchemaRecursiveRef:
		a := r.root.anchors.Recursive
		if a == nil {
			return nil, NewError(fmt.Errorf("openapi: node does not have a $recursiveAnchor but $recursiveRef was found: %s", u), r.root.AbsoluteLocation())
		}

		return &nodectx{
			node:       a.In,
			openapi:    r.openapi,
			jsonschema: r.jsonschema,
			root:       r.root,
			anchors:    r.root.anchors,
		}, nil
	case RefTypeSchema:
		a, ok := r.root.anchors.Standard[Text(r.URI().Fragment)]
		if !ok {
			return nil, NewError(fmt.Errorf("openapi:  not found: %s", u), r.AbsoluteLocation())
		}
		err := r.resolve(a.In)
		if err != nil {
			return nil, fmt.Errorf("openapi: failed to resolve node for anchor \"#%s\": %w", r.URI().Fragment, err)
		}
		return &nodectx{
			node:       a.In,
			openapi:    r.openapi,
			jsonschema: r.jsonschema,
			root:       r.root,
			anchors:    r.root.anchors,
		}, nil
	default:
		return nil, NewError(fmt.Errorf("openapi: anchors are not supported for %s references: #%s", r.RefKind(), u.Fragment), r.AbsoluteLocation())
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
	// if VersionConstraints3_0.Check(doc.OpenAPI) {
	// 	return &JSONSchemaDialect201909, nil
	// }
	return nil, fmt.Errorf("failed to determine OpenAPI schema dialect")
}

func (l *loader) init(node *nodectx, root *nodectx, nodes []node, openapi semver.Version, jsonschema uri.URI) error {
	for _, n := range nodes {
		nc, err := newNodeCtx(n, root, &openapi, &jsonschema)
		if err != nil {
			return err
		}

		l.nodes[n.AbsoluteLocation().String()] = nc

		if IsRef(n) {
			r := n.(ref)
			if !r.IsResolved() {
				l.refs = append(l.refs, refctx{root: root, in: node, ref: r, openapi: nc.openapi, jsonschema: nc.jsonschema})
			}
			return nil
		}
		if err := l.init(&nc, root, n.nodes(), nc.openapi, nc.jsonschema); err != nil {
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
	loc, err := NewLocation(u)
	if err != nil {
		return nil, err
	}
	s.setLocation(loc)
	nc := nodectx{node: &s, openapi: v, jsonschema: u}
	nc.root = &nc
	a, err := s.Anchors()
	if err != nil {
		return nil, fmt.Errorf("failed to load anchors: %w", err)
	}
	nc.anchors = a

	if s.ID != nil && s.ID.String() != u.String() {
		loc, err := NewLocation(*s.ID)
		if err != nil {
			return nil, NewError(fmt.Errorf("failed to parse schema ID: %w", err), u)
		}
		s.setLocation(loc)
		l.nodes[loc.String()] = nc
	} else {
		l.nodes[u.String()] = nc
	}
	return &s, nil
}

func (l *loader) loadNode(ctx context.Context, k Kind, data []byte, v semver.Version, s uri.URI) (nodectx, error) {
	panic("not impl")
}

// func (l *loader) resolveDynamicRefs(n *nodectx) error {
// 	var r refctx
// 	var sr refctx
// 	da := n.anchors.Dynamic
// 	for _, r = range n.resolvedRefs {
// 		for _, sr = range r.resolved.resolvedRefs {
// 			if sr.RefType() == RefTypeSchemaDynamicRef {
// 				a, ok := da[Text(sr.URI().Fragment)]
// 				if ok {
// 					x, ok := sr.in.node.(*Schema)
// 					if !ok {
// 						return fmt.Errorf("openapi: expected schema but got %T", sr.in.node)
// 					}
// 					x.DynamicRef.Resolved = a.In
// 					err := sr.resolve(a.In)
// 					if err != nil {
// 						return fmt.Errorf("openapi: failed to resolve node for anchor \"#%s\": %w", sr.URI().Fragment, err)
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return nil
// }

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
	// case checkVersion(VersionConstraints3_0, v):
	// 	sd = &JSONSchemaDialect201909
	default:
		return nil, nil
	}
	return sd, nil
}

type nodectx struct {
	node
	openapi       semver.Version
	jsonschema    uri.URI
	root          *nodectx
	anchors       *Anchors
	recursiveRefs []refctx
	dynamicRefs   []refctx
	resolvedRefs  []refctx
}
type refctx struct {
	ref
	in         *nodectx
	resolved   *nodectx
	root       *nodectx
	openapi    semver.Version
	jsonschema uri.URI
}

func newNodeCtx(n node, root *nodectx, openapi *semver.Version, jsonschema *uri.URI) (nodectx, error) {
	switch t := n.(type) {
	case *Document:

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
		root:       root,
	}, nil
}

func checkVersion(c semver.Constraints, v *semver.Version) bool {
	if v == nil {
		return false
	}
	return c.Check(v)
}
