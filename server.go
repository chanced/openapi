package openapi

import (
	"encoding/json"

	"github.com/chanced/transcode"
	"gopkg.in/yaml.v3"
)

type (
	ServerSlice       = ObjSlice[*Server]
	ServerVariableMap = ObjMap[*ServerVariable]
)

// Server represention of a Server.
type Server struct {
	Location   `json:"-"`
	Extensions `json:"-"`

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
	Variables *ServerVariableMap `json:"variables,omitempty"`
}

func (s *Server) Nodes() []Node {
	if s == nil {
		return nil
	}
	return downcastNodes(s.nodes())
}

func (s *Server) nodes() []node {
	if s == nil {
		return nil
	}
	return appendEdges(nil, s.Variables)
}

func (s *Server) Refs() []Ref {
	if s == nil {
		return nil
	}
	var refs []Ref
	if s.Variables != nil {
		refs = append(refs, s.Variables.Refs()...)
	}
	return refs
}

// func (s *Server) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if err := ptr.Validate(); err != nil {
// 		return nil, err
// 	}
// 	return s.resolveNodeByPointer(ptr)
// }

// func (s *Server) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if ptr.IsRoot() {
// 		return s, nil
// 	}
// 	nxt, tok, _ := ptr.Next()
// 	switch tok {
// 	case "variables":
// 		if s.Variables == nil {
// 			return nil, newErrNotFound(s.AbsoluteLocation(), tok)
// 		}
// 		return s.Variables.resolveNodeByPointer(nxt)
// 	default:
// 		return nil, newErrNotResolvable(s.Location.AbsoluteLocation(), tok)
// 	}
// }

func (*Server) Kind() Kind      { return KindServer }
func (*Server) mapKind() Kind   { return KindUndefined }
func (*Server) sliceKind() Kind { return KindServerSlice }

func (*Server) Anchors() (*Anchors, error) { return nil, nil }

func (s *Server) setLocation(loc Location) error {
	if s == nil {
		return nil
	}
	s.Location = loc
	return s.Variables.setLocation(loc.AppendLocation("variables"))
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

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Marshaler interface
func (s Server) MarshalYAML() (interface{}, error) {
	j, err := s.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return transcode.YAMLFromJSON(j)
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Unmarshaler interface
func (s *Server) UnmarshalYAML(value *yaml.Node) error {
	j, err := transcode.YAMLFromJSON([]byte(value.Value))
	if err != nil {
		return err
	}
	return json.Unmarshal(j, s)
}

func (s *Server) isNil() bool { return s == nil }

var _ node = (*Server)(nil)
