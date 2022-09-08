package openapi

// ExternalDocs allows referencing an external resource for extended
// documentation.
type ExternalDocs struct {
	// The URL for the target documentation. This MUST be in the form of a URL.
	//
	// 	*required*
	URL string `json:"url"`
	// A description of the target documentation. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Extensions  `json:"-"`
}
type externaldocs ExternalDocs

// MarshalJSON marshals JSON
func (ed ExternalDocs) MarshalJSON() ([]byte, error) {
	return marshalExtendedJSON(externaldocs(ed))
}

// UnmarshalJSON unmarshals JSON
func (ed *ExternalDocs) UnmarshalJSON(data []byte) error {
	var v externaldocs
	err := unmarshalExtendedJSON(data, &v)
	*ed = ExternalDocs(v)
	return err
}
