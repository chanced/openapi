package openapi

import (
	"encoding/json"

	"github.com/chanced/openapi/yamlutil"
)

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

type tag Tag

// MarshalJSON marshals t into JSON
func (t Tag) MarshalJSON() ([]byte, error) {
	return marshalExtendedJSON(tag(t))
}

func (t *Tag) Nodes() map[string][]Node {
	if t.ExternalDocs != nil {
		return map[string][]Node{
			"externalDocs": {t.ExternalDocs},
		}
	}
	return nil
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

func (t *Tag) Kind() Kind {
	return KindTag
}

var _ Node = (*Tag)(nil)
