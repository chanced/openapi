package openapi

import "github.com/chanced/openapi/yamlutil"

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

// Kind returns KindExternalDocs
func (*ExternalDocs) Kind() Kind {
	return KindExternalDocs
}

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

// MarshalYAML marshals YAML
func (ed ExternalDocs) MarshalYAML() (interface{}, error) {
	return yamlutil.Marshal(ed)
}

// UnmarshalYAML unmarshals YAML
func (ed *ExternalDocs) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return yamlutil.Unmarshal(unmarshal, ed)
}
