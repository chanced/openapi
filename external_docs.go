package openapi

import (
	"encoding/json"

	"github.com/chanced/transcode"
	"github.com/chanced/uri"
	"gopkg.in/yaml.v3"
)

// ExternalDocs allows referencing an external resource for extended
// documentation.
type ExternalDocs struct {
	Location   `json:"-"`
	Extensions `json:"-"`

	// The URL for the target documentation. This MUST be in the form of a URL.
	//
	// 	*required*
	URL *uri.URI `json:"url"`

	// A description of the target documentation. CommonMark syntax MAY be used
	// for rich text representation.
	Description Text `json:"description,omitempty"`
}

func (ed *ExternalDocs) Nodes() []Node {
	if ed == nil {
		return nil
	}
	return downcastNodes(ed.nodes())
}
func (ed *ExternalDocs) nodes() []node { return nil }

func (*ExternalDocs) Refs() []Ref     { return nil }
func (*ExternalDocs) Kind() Kind      { return KindExternalDocs }
func (*ExternalDocs) mapKind() Kind   { return KindUndefined }
func (*ExternalDocs) sliceKind() Kind { return KindUndefined }

func (*ExternalDocs) Anchors() (*Anchors, error) { return nil, nil }

// // ResolveNodeByPointer resolves a Node by a jsonpointer. It validates the pointer and then
// // attempts to resolve the Node.
// //
// // # Errors
// //
// // - [ErrNotFound] indicates that the component was not found
// //
// // - [ErrNotResolvable] indicates that the pointer path can not resolve to a
// // Node
// //
// // - [jsonpointer.ErrMalformedEncoding] indicates that the pointer encoding
// // is malformed
// //
// // - [jsonpointer.ErrMalformedStart] indicates that the pointer is not empty
// // and does not start with a slash
// func (ed *ExternalDocs) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	err := ptr.Validate()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return ed.resolveNodeByPointer(ptr)
// }

// func (ed *ExternalDocs) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	tok, _ := ptr.NextToken()
// 	if !ptr.IsRoot() {
// 		return nil, newErrNotResolvable(ed.Location.AbsoluteLocation(), tok)
// 	}
// 	return ed, nil
// }

// MarshalJSON marshals JSON
func (ed ExternalDocs) MarshalJSON() ([]byte, error) {
	type externaldocs ExternalDocs

	return marshalExtendedJSON(externaldocs(ed))
}

// UnmarshalJSON unmarshals JSON
func (ed *ExternalDocs) UnmarshalJSON(data []byte) error {
	type externaldocs ExternalDocs

	var v externaldocs
	err := unmarshalExtendedJSON(data, &v)
	*ed = ExternalDocs(v)
	return err
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Marshaler interface
func (ed ExternalDocs) MarshalYAML() (interface{}, error) {
	j, err := ed.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return transcode.YAMLFromJSON(j)
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Unmarshaler interface
func (ed *ExternalDocs) UnmarshalYAML(value *yaml.Node) error {
	v, err := yaml.Marshal(value)
	if err != nil {
		return err
	}
	j, err := transcode.JSONFromYAML(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(j, ed)
}

func (ed *ExternalDocs) setLocation(loc Location) error {
	if ed == nil {
		return nil
	}
	ed.Location = loc
	return nil
}
func (ed *ExternalDocs) isNil() bool { return ed == nil }

var _ node = (*ExternalDocs)(nil)
