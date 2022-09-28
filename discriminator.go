package openapi

import (
	"encoding/json"

	"github.com/chanced/transcode"
	"gopkg.in/yaml.v3"
)

// Discriminator can be used to aid in serialization, deserialization, and
// validation of request bodies or response payloads which may be one of a
// number of different schemas. The discriminator is a specific object in a
// schema which is used to inform the consumer of the document of an alternative
// schema based on the value associated with it.
type Discriminator struct {
	Extensions `json:"-"`
	Location   `json:"-"`

	// The name of the property in the payload that will hold the discriminator
	// value.
	//
	// *required
	PropertyName Text `json:"propertyName"`
	// An object to hold mappings between payload values and schema names or
	// references.
	Mapping *Map[Text] `json:"mapping,omitempty"`
}

func (d *Discriminator) Clone() *Discriminator {
	if d == nil {
		return nil
	}
	var m *Map[Text]
	if d.Mapping != nil {
		m := Map[Text]{
			Items: make([]KeyValue[Text], len(d.Mapping.Items)),
		}
		copy(m.Items, d.Mapping.Items)
	}
	return &Discriminator{
		Extensions: d.Extensions,
		Location: Location{
			absolute: *d.Location.absolute.Clone(),
			relative: d.Location.relative,
		},
		PropertyName: d.PropertyName.Clone(),
		Mapping:      m,
	}
}

func (d *Discriminator) setLocation(loc Location) error {
	if d == nil {
		return nil
	}
	d.Location = loc
	return nil
}

// MarshalJSON marshals d into JSON
func (d Discriminator) MarshalJSON() ([]byte, error) {
	type discriminator Discriminator

	return marshalExtendedJSON(discriminator(d))
}

// UnmarshalJSON unmarshals json into d
func (d *Discriminator) UnmarshalJSON(data []byte) error {
	type discriminator Discriminator

	v := discriminator{}
	if err := unmarshalExtendedJSON(data, &v); err != nil {
		return err
	}
	*d = Discriminator(v)
	return nil
}

func (d Discriminator) MarshalYAML() (interface{}, error) {
	j, err := d.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return transcode.YAMLFromJSON(j)
}

// UnmarshalYAML implements yaml.Unmarshaler
func (d *Discriminator) UnmarshalYAML(value *yaml.Node) error {
	v, err := yaml.Marshal(value)
	if err != nil {
		return err
	}
	j, err := transcode.JSONFromYAML(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(j, d)
}

func (d *Discriminator) Anchors() (*Anchors, error) { return nil, nil }

func (*Discriminator) Kind() Kind { return KindDiscriminator }

func (*Discriminator) Refs() []Ref { return nil }

// func (d *Discriminator) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if err := ptr.Validate(); err != nil {
// 		return nil, err
// 	}
// 	return d.resolveNodeByPointer(ptr)
// }

// func (d *Discriminator) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if ptr.IsRoot() {
// 		return d, nil
// 	}
// 	tok, _ := ptr.NextToken()
// 	return nil, newErrNotResolvable(d.Location.AbsoluteLocation(), tok)
// }

func (d *Discriminator) Nodes() []Node {
	if d == nil {
		return nil
	}
	return downcastNodes(d.nodes())
}
func (d *Discriminator) nodes() []node { return nil }

func (d *Discriminator) isNil() bool   { return d == nil }
func (*Discriminator) mapKind() Kind   { return KindUndefined }
func (*Discriminator) sliceKind() Kind { return KindUndefined }

var _ node = (*Discriminator)(nil)
