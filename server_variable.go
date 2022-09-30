package openapi

import (
	"encoding/json"

	"github.com/chanced/transcode"
	"gopkg.in/yaml.v3"
)

// ServerVariable for server URL template substitution.
type ServerVariable struct {
	// An enumeration of string values to be used if the substitution options
	// are from a limited set. The array MUST NOT be empty.
	Enum Texts `json:"enum"`
	// The default value to use for substitution, which SHALL be sent if an
	// alternate value is not supplied. Note this behavior is different than the
	// Schema Object's treatment of default values, because in those cases
	// parameter values are optional. If the enum is defined, the value MUST
	// exist in the enum's values.
	//
	// 	*required*
	Default Text `json:"default"`
	// An optional description for the server variable. CommonMark syntax MAY be
	// used for rich text representation.
	Description Text `json:"description,omitempty"`

	Location   `json:"-"`
	Extensions `json:"-"`
}

func (sv *ServerVariable) Nodes() []Node {
	if sv == nil {
		return nil
	}
	return downcastNodes(sv.nodes())
}
func (sv *ServerVariable) nodes() []node { return nil }

func (*ServerVariable) Refs() []Ref { return nil }

// func (sv *ServerVariable) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if err := ptr.Validate(); err != nil {
// 		return nil, err
// 	}
// 	return sv.resolveNodeByPointer(ptr)
// }

// func (sv *ServerVariable) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if ptr.IsRoot() {
// 		return sv, nil
// 	}
// 	tok, _ := ptr.NextToken()
// 	return nil, newErrNotResolvable(sv.Location.AbsoluteLocation(), tok)
// }

func (*ServerVariable) Kind() Kind      { return KindServerVariable }
func (*ServerVariable) mapKind() Kind   { return KindServerVariableMap }
func (*ServerVariable) sliceKind() Kind { return KindUndefined }

func (sv *ServerVariable) setLocation(loc Location) error {
	if sv == nil {
		return nil
	}
	sv.Location = loc
	return nil
}

// MarshalJSON marshals JSON
func (sv ServerVariable) MarshalJSON() ([]byte, error) {
	type servervariable ServerVariable
	return marshalExtendedJSON(servervariable(sv))
}

// UnmarshalJSON unmarshals JSON
func (sv *ServerVariable) UnmarshalJSON(data []byte) error {
	type servervariable ServerVariable
	var v servervariable
	err := unmarshalExtendedJSON(data, &v)
	*sv = ServerVariable(v)
	return err
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Marshaler interface
func (sv ServerVariable) MarshalYAML() (interface{}, error) {
	j, err := sv.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var v interface{}
	err = json.Unmarshal(j, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Unmarshaler interface
func (sv *ServerVariable) UnmarshalYAML(value *yaml.Node) error {
	v, err := yaml.Marshal(value)
	if err != nil {
		return err
	}
	j, err := transcode.JSONFromYAML(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(j, sv)
}

func (sv *ServerVariable) Anchors() (*Anchors, error) {
	return nil, nil
}
func (sv *ServerVariable) isNil() bool { return sv == nil }

var _ node = (*ServerVariable)(nil)
