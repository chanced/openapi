package openapi

import "github.com/chanced/jsonpointer"

type (
	ServerSlice       = ComponentSlice[*Server]
	ServerVariableMap = ObjMap[*ServerVariable]
)

// Server represention of a Server.
type Server struct {
	// A URL to the target host. This URL supports Server Variables and MAY be
	// relative, to indicate that the host location is relative to the location
	// where the OpenAPI document is being served. Variable substitutions will
	// be made when a variable is named in {brackets}.
	URL Text `json:"url"`
	// Description of the host designated by the URL. CommonMark syntax MAY be
	// used for rich text representation.
	Description Text `json:"description,omitempty"`
	// A map between a variable name and its value. The value is used for
	// substitution in the server's URL template.
	Variables  *ServerVariableMap `json:"variables,omitempty"`
	Location   `json:"-"`
	Extensions `json:"-"`
}

func (s *Server) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
	if err := ptr.Validate(); err != nil {
		return nil, err
	}
	return s.resolveNodeByPointer(ptr)
}

func (s *Server) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
	if ptr.IsRoot() {
		return s, nil
	}
	nxt, tok, _ := ptr.Next()
	switch tok {
	case "variables":
		if s.Variables == nil {
			return nil, newErrNotFound(s.AbsoluteLocation(), tok)
		}
		return s.Variables.resolveNodeByPointer(nxt)
	default:
		return nil, newErrNotResolvable(s.Location.AbsoluteLocation(), tok)
	}
}

func (*Server) Kind() Kind      { return KindServer }
func (*Server) mapKind() Kind   { return KindUndefined }
func (*Server) sliceKind() Kind { return KindServerSlice }

func (*Server) Anchors() (*Anchors, error) { return nil, nil }

func (s *Server) setLocation(loc Location) error {
	if s == nil {
		return nil
	}
	s.Location = loc
	return s.Variables.setLocation(loc.Append("variables"))
}

// MarshalJSON marshals JSON
func (s Server) MarshalJSON() ([]byte, error) {
	type server Server
	return marshalExtendedJSON(server(s))
}

// UnmarshalJSON unmarshals JSON
func (s *Server) UnmarshalJSON(data []byte) error {
	type server Server
	var v server
	err := unmarshalExtendedJSON(data, &v)
	*s = Server(v)
	return err
}

// ServerVariable for server URL template substitution.
type ServerVariable struct {
	// An enumeration of string values to be used if the substitution options
	// are from a limited set. The array MUST NOT be empty.
	Enum []string `json:"enum"`
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

func (sv *ServerVariable) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
	if err := ptr.Validate(); err != nil {
		return nil, err
	}
	return sv.resolveNodeByPointer(ptr)
}

func (sv *ServerVariable) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
	if ptr.IsRoot() {
		return sv, nil
	}
	tok, _ := ptr.NextToken()
	return nil, newErrNotResolvable(sv.Location.AbsoluteLocation(), tok)
}

func (*ServerVariable) Kind() Kind      { return KindServerVariable }
func (*ServerVariable) mapKind() Kind   { return KindServerVariableMap }
func (*ServerVariable) sliceKind() Kind { return KindUndefined }

func (sv *ServerVariable) setLocation(loc Location) error {
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

func (sv *ServerVariable) Anchors() (*Anchors, error) {
	return nil, nil
}

var (
	_ node = (*Server)(nil)
	_ node = (*ServerVariable)(nil)
)
