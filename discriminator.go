package openapi

import (
	"github.com/chanced/openapi/yamlutil"
)

// Discriminator can be used to aid in serialization, deserialization, and
// validation of request bodies or response payloads which may be one of a
// number of different schemas. The discriminator is a specific object in a
// schema which is used to inform the consumer of the document of an alternative
// schema based on the value associated with it.
type Discriminator struct {
	// The name of the property in the payload that will hold the discriminator
	// value.
	//
	// *required
	PropertyName string `json:"propertyName"`
	// An object to hold mappings between payload values and schema names or
	// references.
	Mapping map[string]string `json:"mapping,omitempty"`

	Extensions `json:"-"`
}

type discriminator Discriminator

// MarshalJSON marshals d into JSON
func (d Discriminator) MarshalJSON() ([]byte, error) {
	return marshalExtendedJSON(discriminator(d))
}

// UnmarshalJSON unmarshals json into d
func (d *Discriminator) UnmarshalJSON(data []byte) error {
	v := discriminator{}
	if err := unmarshalExtendedJSON(data, &v); err != nil {
		return err
	}
	*d = Discriminator(v)
	return nil
}

// MarshalYAML first marshals and unmarshals into JSON and then marshals into
// YAML
func (d *Discriminator) MarshalYAML() (interface{}, error) {
	return yamlutil.Marshal(d)
}

// UnmarshalYAML unmarshals yaml into s
func (d *Discriminator) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return yamlutil.Unmarshal(unmarshal, d)
}
