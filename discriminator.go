package openapi

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
