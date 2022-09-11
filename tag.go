package openapi

// Tag adds metadata that is used by the Operation Object.
//
// It is not mandatory to have a Tag Object per tag defined in the Operation
// Object instances.
type Tag struct {
	// The name of the tag.
	//
	// 	*required*
	Name string `json:"name"`
	//  A description for the tag.
	//
	// CommonMark syntax MAY be used for rich text representation.
	//
	// https://spec.commonmark.org/
	Description string `json:"description,omitempty"`
	// Additional external documentation for this tag.
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty" bson:"externalDocs,omitempty"`

	Extensions `json:"-"`
}

// MarshalJSON marshals t into JSON
func (t Tag) MarshalJSON() ([]byte, error) {
	type tag Tag

	return marshalExtendedJSON(tag(t))
}

// UnmarshalJSON unmarshals json into t
func (t *Tag) UnmarshalJSON(data []byte) error {
	type tag Tag

	v := tag{}
	err := unmarshalExtendedJSON(data, &v)
	*t = Tag(v)
	return err
}
