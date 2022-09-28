package openapi

import (
	"encoding/json"

	"github.com/chanced/transcode"
	"gopkg.in/yaml.v3"
)

// ResponseMap is a container for the expected responses of an operation. The
// container maps a HTTP response code to the expected response.
//
// The documentation is not necessarily expected to cover all possible HTTP
// response codes because they may not be known in advance. However,
// documentation is expected to cover a successful operation response and any
// known errors.
//
// The default MAY be used as a default response object for all HTTP codes that
// are not covered individually by the ResponseMap Object.
//
// The ResponseMap Object MUST contain at least one response code, and if only one
// response code is provided it SHOULD be the response for a successful
// operation call.
type ResponseMap = ComponentMap[*Response]

// Response describes a single response from an API Operation, including
// design-time, static links to operations based on the response.
type Response struct {
	// A description of the response. CommonMark syntax MAY be used for rich
	// text representation.
	//
	// *required*
	Description Text `json:"description,omitempty"`
	// Maps a header name to its definition. RFC7230 states header names are
	// case insensitive. If a response header is defined with the name
	// "Content-Type", it SHALL be ignored.
	Headers *HeaderMap `json:"headers,omitempty"`
	// A map containing descriptions of potential response payloads. The key is
	// a media type or media type range and the value describes it. For
	// responses that match multiple keys, only the most specific key is
	// applicable. e.g. text/plain overrides text/*
	Content *ContentMap `json:"content,omitempty"`
	// A map of operations links that can be followed from the response. The key
	// of the map is a short name for the link, following the naming constraints
	// of the names for Component Objects.
	Links      *LinkMap `json:"links,omitempty"`
	Extensions `json:"-"`

	Location `json:"-"`
}

func (r *Response) Nodes() []Node {
	if r == nil {
		return nil
	}
	return downcastNodes(r.nodes())
}

func (r *Response) nodes() []node {
	if r == nil {
		return nil
	}
	return appendEdges(nil, r.Headers, r.Content, r.Links)
}

func (r *Response) Refs() []Ref {
	if r == nil {
		return nil
	}
	var refs []Ref
	refs = append(refs, r.Headers.Refs()...)
	refs = append(refs, r.Content.Refs()...)
	refs = append(refs, r.Links.Refs()...)
	return refs
}

func (r *Response) Anchors() (*Anchors, error) {
	if r == nil {
		return nil, nil
	}
	var anchors *Anchors
	var err error
	if anchors, err = anchors.merge(r.Headers.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(r.Content.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(r.Links.Anchors()); err != nil {
		return nil, err
	}
	return anchors, nil
}

// MarshalJSON marshals r into JSON
func (r Response) MarshalJSON() ([]byte, error) {
	type response Response
	return marshalExtendedJSON(response(r))
}

// UnmarshalJSON unmarshals json into r
func (r *Response) UnmarshalJSON(data []byte) error {
	type response Response
	var v response
	err := unmarshalExtendedJSON(data, &v)
	*r = Response(v)
	return err
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Marshaler interface
func (r Response) MarshalYAML() (interface{}, error) {
	j, err := r.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return transcode.YAMLFromJSON(j)
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Unmarshaler interface
func (r *Response) UnmarshalYAML(value *yaml.Node) error {
	v, err := yaml.Marshal(value)
	if err != nil {
		return err
	}
	j, err := transcode.JSONFromYAML(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(j, r)
}

func (*Response) Kind() Kind      { return KindResponse }
func (*Response) mapKind() Kind   { return KindResponseMap }
func (*Response) sliceKind() Kind { return KindUndefined }

func (r *Response) isNil() bool { return r == nil }

func (r *Response) setLocation(loc Location) error {
	if r == nil {
		return nil
	}
	r.Location = loc
	if err := r.Headers.setLocation(loc.AppendLocation("headers")); err != nil {
		return err
	}
	if err := r.Content.setLocation(loc.AppendLocation("content")); err != nil {
		return err
	}
	if err := r.Links.setLocation(loc.AppendLocation("links")); err != nil {
		return err
	}
	return nil
}

func (*Response) refable() {}

var _ node = (*Response)(nil)

// ResolveNodeByPointer resolves a Node by a jsonpointer. It validates the pointer and then
// attempts to resolve the Node.
//
// # Errors
//
// - [ErrNotFound] indicates that the component was not found
//
// - [ErrNotResolvable] indicates that the pointer path can not resolve to a
// Node
//
// - [jsonpointer.ErrMalformedEncoding] indicates that the pointer encoding
// is malformed
//
// - [jsonpointer.ErrMalformedStart] indicates that the pointer is not empty
// and does not start with a slash
// func (r *Response) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	err := ptr.Validate()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return r.resolveNodeByPointer(ptr)
// }

// func (r *Response) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if ptr.IsRoot() {
// 		return r, nil
// 	}
// 	nxt, tok, _ := ptr.Next()
// 	switch tok {
// 	case "headers":
// 		if r.Headers == nil {
// 			return nil, newErrNotFound(r.AbsoluteLocation(), tok)
// 		}
// 		return r.Headers.resolveNodeByPointer(nxt)
// 	case "content":
// 		if r.Content == nil {
// 			return nil, newErrNotFound(r.AbsoluteLocation(), tok)
// 		}
// 		return r.Content.resolveNodeByPointer(nxt)
// 	case "links":
// 		if r.Links == nil {
// 			return nil, newErrNotFound(r.AbsoluteLocation(), tok)
// 		}
// 		return r.Links.resolveNodeByPointer(nxt)
// 	default:
// 		return nil, newErrNotResolvable(r.Location.AbsoluteLocation(), tok)
// 	}
// }
