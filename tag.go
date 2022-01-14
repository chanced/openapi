package openapi

import (
	"encoding/json"

	"github.com/chanced/openapi/yamlutil"
)

type Tags []*Tag

func (ts Tags) Kind() Kind { return KindTags }

func (ts Tags) Nodes() Nodes {
	if ts.Len() == 0 {
		return nil
	}
	n := make(Nodes, len(ts))
	for i, s := range ts {
		n.maybeAdd(i, s, KindSchema)
	}
	if len(n) == 0 {
		return nil
	}
	return n
}

func (ts *Tags) Get(idx int) (*Tag, bool) {
	if *ts == nil {
		return nil, false
	}
	if idx < 0 || idx >= len(*ts) {
		return nil, false
	}
	return (*ts)[idx], true
}

func (ts *Tags) Append(val *Tag) {
	if *ts == nil {
		*ts = Tags{val}
		return
	}
	(*ts) = append(*ts, val)
}

func (ts *Tags) Remove(s *Tag) {
	if *ts == nil {
		return
	}
	for k, v := range *ts {
		if v == s {
			ts.RemoveIndex(k)
			return
		}
	}
}

func (ts *Tags) RemoveIndex(i int) {
	if *ts == nil {
		return // nothing to do
	}
	if i < 0 || i >= len(*ts) {
		return
	}
	copy((*ts)[i:], (*ts)[i+1:])
	(*ts)[len(*ts)-1] = nil
	(*ts) = (*ts)[:ts.Len()-1]
}

// Len returns the length of s
func (ts *Tags) Len() int {
	if ts == nil || *ts == nil {
		return 0
	}
	return len(*ts)
}

// Tag adds metadata that is used by the Operation Object.
//
// It is not mandatory to have a Tag Object per tag defined in the Operation
// Object instances.
type Tag struct {
	// The name of the tag.
	//
	// 	*required*
	Name string `json:"name" yaml:"name"`
	//  A description for the tag.
	//
	// CommonMark syntax MAY be used for rich text representation.
	//
	// https://spec.commonmark.org/
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Additional external documentation for this tag.
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty" bson:"externalDocs,omitempty"`

	Extensions `json:"-"`
}

func (t *Tag) Nodes() Nodes {
	return makeNodes(nodes{
		"externalDocs": {t.ExternalDocs, KindExternalDocs},
	})
}

type tag Tag

// MarshalJSON marshals t into JSON
func (t Tag) MarshalJSON() ([]byte, error) {
	return marshalExtendedJSON(tag(t))
}

// UnmarshalJSON unmarshals json into t
func (t *Tag) UnmarshalJSON(data []byte) error {
	v := tag{}
	err := unmarshalExtendedJSON(data, &v)
	*t = Tag(v)
	return err
}

// MarshalYAML first marshals and unmarshals into JSON and then marshals into
// YAML
func (t Tag) MarshalYAML() (interface{}, error) {
	b, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	var v interface{}
	err = json.Unmarshal(b, &v)
	return v, err
}

// UnmarshalYAML unmarshals yaml into t
func (t *Tag) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return yamlutil.Unmarshal(unmarshal, t)
}

// Kind returns KindTag
func (*Tag) Kind() Kind {
	return KindTag
}

var _ Node = (*Tag)(nil)
