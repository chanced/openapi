package openapi

import (
	"encoding/json"

	"github.com/chanced/transcode"
	"gopkg.in/yaml.v3"
)

// RequestBodyMap is a map of RequestBody
type RequestBodyMap = ComponentMap[*RequestBody]

// RequestBody describes a single request body.
type RequestBody struct {
	Location   `json:"-"`
	Extensions `json:"-"`

	// A brief description of the request body. This could contain examples of
	// use. CommonMark syntax MAY be used for rich text representation.
	Description Text `json:"description,omitempty"`
	// The content of the request body. The key is a media type or media type range and the value describes it. For requests that match multiple keys, only the most specific key is applicable. e.g. text/plain overrides text
	//
	// *required*
	Content *ContentMap `json:"content,omitempty"`
	// Determines if the request body is required in the request. Defaults to false.
	Required bool `json:"required,omitempty"`
}

func (rb *RequestBody) Nodes() []Node {
	if rb == nil {
		return nil
	}
	return downcastNodes(rb.nodes())
}

func (rb *RequestBody) nodes() []node {
	if rb == nil {
		return nil
	}
	return appendEdges(nil, rb.Content)
}

func (rb *RequestBody) Refs() []Ref {
	if rb == nil {
		return nil
	}
	return rb.Content.Refs()
}

func (rb *RequestBody) isNil() bool { return rb == nil }

func (rb *RequestBody) Anchors() (*Anchors, error) {
	if rb == nil {
		return nil, nil
	}
	return rb.Content.Anchors()
}

// func (rb *RequestBody) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if err := ptr.Validate(); err != nil {
// 		return nil, err
// 	}
// 	return rb.resolveNodeByPointer(ptr)
// }

// func (rb *RequestBody) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if ptr.IsRoot() {
// 		return rb, nil
// 	}
// 	nxt, tok, _ := ptr.Next()
// 	switch tok {
// 	case "content":
// 		if rb.Content == nil {
// 			return nil, newErrNotFound(rb.AbsoluteLocation(), tok)
// 		}
// 		return rb.Content.resolveNodeByPointer(nxt)
// 	default:
// 		return nil, newErrNotResolvable(rb.Location.AbsoluteLocation(), tok)
// 	}
// }

func (*RequestBody) Kind() Kind      { return KindRequestBody }
func (*RequestBody) mapKind() Kind   { return KindRequestBodyMap }
func (*RequestBody) sliceKind() Kind { return KindUndefined }

// MarshalJSON marshals h into JSON
func (rb RequestBody) MarshalJSON() ([]byte, error) {
	type requestbody RequestBody

	return marshalExtendedJSON(requestbody(rb))
}

// UnmarshalJSON unmarshals json into rb
func (rb *RequestBody) UnmarshalJSON(data []byte) error {
	type requestbody RequestBody

	var v requestbody
	err := unmarshalExtendedJSON(data, &v)
	*rb = RequestBody(v)
	return err
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Marshaler interface
func (rb RequestBody) MarshalYAML() (interface{}, error) {
	j, err := rb.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return transcode.YAMLFromJSON(j)
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Unmarshaler interface
func (rb *RequestBody) UnmarshalYAML(value *yaml.Node) error {
	v, err := yaml.Marshal(value)
	if err != nil {
		return err
	}
	j, err := transcode.JSONFromYAML(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(j, rb)
}

func (rb *RequestBody) setLocation(loc Location) error {
	if rb == nil {
		return nil
	}
	rb.Location = loc
	if err := rb.Content.setLocation(loc.AppendLocation("content")); err != nil {
		return err
	}
	return nil
}
func (*RequestBody) refable() {}

var _ node = (*RequestBody)(nil)
