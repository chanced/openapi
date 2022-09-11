package openapi

// ExternalDocs allows referencing an external resource for extended
// documentation.
type ExternalDocs struct {
	// The URL for the target documentation. This MUST be in the form of a URL.
	//
	// 	*required*
	URL Text `json:"url"`
	// A description of the target documentation. CommonMark syntax MAY be used for rich text representation.
	Description Text `json:"description,omitempty"`

	Location   *Location `json:"-"`
	Extensions `json:"-"`
}

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

func (ed *ExternalDocs) setLocation(loc Location) error {
	if ed == nil {
		return nil
	}
	ed.Location = &loc
	return nil
}
