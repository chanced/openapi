package openapi

import "github.com/chanced/openapi/yamlutil"

// Server represention of a Server.
type Server struct {
	// A URL to the target host. This URL supports Server Variables and MAY be
	// relative, to indicate that the host location is relative to the location
	// where the OpenAPI document is being served. Variable substitutions will
	// be made when a variable is named in {brackets}.
	URL string `json:"url"`
	// Description of the host designated by the URL. CommonMark syntax MAY be
	// used for rich text representation.
	Description string `json:"description,omitempty"`
	// A map between a variable name and its value. The value is used for
	// substitution in the server's URL template.
	Variables  map[string]*ServerVariable `json:"variables,omitempty"`
	Extensions `json:"-"`
}
type server Server

// MarshalJSON marshals JSON
func (s Server) MarshalJSON() ([]byte, error) {
	return marshalExtendedJSON(server(s))
}

// UnmarshalJSON unmarshals JSON
func (s *Server) UnmarshalJSON(data []byte) error {
	var v server
	err := unmarshalExtendedJSON(data, &v)
	*s = Server(v)
	return err
}

// MarshalYAML marshals YAML
func (s *Server) MarshalYAML() (interface{}, error) {
	return yamlutil.Marshal(s)
}

// UnmarshalYAML unmarshals YAML data into s
func (s *Server) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return yamlutil.Unmarshal(unmarshal, s)
}

// ServerVariable for server URL template substitution.
type ServerVariable struct {
	// An enumeration of string values to be used if the substitution options
	// are from a limited set. The array MUST NOT be empty.
	Enum []string `json:"enum" yaml:"enum"`
	// The default value to use for substitution, which SHALL be sent if an
	// alternate value is not supplied. Note this behavior is different than the
	// Schema Object's treatment of default values, because in those cases
	// parameter values are optional. If the enum is defined, the value MUST
	// exist in the enum's values.
	//
	// 	*required*
	Default string `json:"default" yaml:"default"`
	// An optional description for the server variable. CommonMark syntax MAY be
	// used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Extensions  `json:"-"`
}

type servervariable ServerVariable

// MarshalJSON marshals JSON
func (sv ServerVariable) MarshalJSON() ([]byte, error) {
	return marshalExtendedJSON(servervariable(sv))
}

// UnmarshalJSON unmarshals JSON
func (sv *ServerVariable) UnmarshalJSON(data []byte) error {
	var v servervariable
	err := unmarshalExtendedJSON(data, &v)
	*sv = ServerVariable(v)
	return err
}

// MarshalYAML marshals YAML
func (sv *ServerVariable) MarshalYAML() (interface{}, error) {
	return yamlutil.Marshal(sv)
}

// UnmarshalYAML unmarshals YAML data into s
func (sv *ServerVariable) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return yamlutil.Unmarshal(unmarshal, sv)
}
