package openapi

import (
	"encoding/json"

	"github.com/chanced/transcode"
	"github.com/chanced/uri"
	"gopkg.in/yaml.v3"
)

// License information for the exposed API.
type License struct {
	Location `json:"-"`

	// The license name used for the API.
	//
	// 	*required*
	Name Text `json:"name"`

	// An SPDX license expression for the API. The identifier field is mutually
	// exclusive of the url field.
	Identifier Text `json:"identifier,omitempty"`
	// A URL to the license used for the API. This MUST be in the form of a URL.
	// The url field is mutually exclusive of the identifier field.
	URL *uri.URI `json:"url,omitempty"`
}

func (*License) Anchors() (*Anchors, error) { return nil, nil }

// Kind returns KindLicense
func (*License) Kind() Kind      { return KindLicense }
func (*License) sliceKind() Kind { return KindUndefined }
func (*License) mapKind() Kind   { return KindUndefined }

func (*License) nodes() []node        { return nil }
func (l *License) isNil() bool        { return l == nil }
func (l *License) location() Location { return l.Location }

func (*License) Refs() []Ref { return nil }

//	func (l *License) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
//		if err := ptr.Validate(); err != nil {
//			return nil, err
//		}
//		return l.resolveNodeByPointer(ptr)
//	}
//
//	func (l *License) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
//		if ptr.IsRoot() {
//			return l, nil
//		}
//		tok, _ := ptr.NextToken()
//		return nil, newErrNotResolvable(l.absolute, tok)
//	}
//
// MarshalJSON marshals JSON
func (l License) MarshalJSON() ([]byte, error) {
	type license License
	return json.Marshal(license(l))
}

// UnmarshalJSON unmarshals JSON
func (l *License) UnmarshalJSON(data []byte) error {
	*l = License{}
	type license License
	var a license
	err := json.Unmarshal(data, &a)
	if err != nil {
		return err
	}
	*l = License(a)
	return nil
}

// MarshalYAML implements yaml.Marshaler
func (l License) MarshalYAML() (interface{}, error) {
	j, err := l.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return transcode.YAMLFromJSON(j)
}

// UnmarshalYAML implements yaml.Unmarshaler
func (l *License) UnmarshalYAML(value *yaml.Node) error {
	j, err := transcode.YAMLFromJSON([]byte(value.Value))
	if err != nil {
		return err
	}
	return json.Unmarshal(j, l)
}

func (l *License) setLocation(loc Location) error {
	if l == nil {
		return nil
	}
	l.Location = loc
	return nil
}

var _ node = (*License)(nil)
