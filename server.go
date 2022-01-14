package openapi

import (
	"github.com/chanced/openapi/yamlutil"
)

type ServerVariables map[string]*ServerVariable

func (ServerVariables) Kind() Kind {
	return KindServerVariables
}

func (svs *ServerVariables) Len() int {
	if svs == nil || *svs == nil {
		return 0
	}
	return len(*svs)
}

func (svs *ServerVariables) Get(key string) (*ServerVariable, bool) {
	if svs == nil || *svs == nil {
		return nil, false
	}
	v, ok := (*svs)[key]
	return v, ok
}

func (svs *ServerVariables) Set(key string, val *ServerVariable) {
	if *svs == nil {
		*svs = ServerVariables{
			key: val,
		}
		return
	}
	(*svs)[key] = val
}

func (svs ServerVariables) Nodes() Nodes {
	if len(svs) == 0 {
		return nil
	}
	n := make(Nodes, len(svs))
	for k, v := range svs {
		n.maybeAdd(k, v, KindServerVariable)
	}
	if len(n) == 0 {
		return nil
	}
	return n
}

type Servers []*Server

func (ss Servers) Nodes() Nodes {
	if ss.Len() == 0 {
		return nil
	}
	nodes := make(Nodes, ss.Len())
	for i, s := range ss {
		nodes.maybeAdd(i, s, KindServer)
	}
	return nodes
}

func (Servers) Kind() Kind {
	return KindServers
}

func (ss *Servers) Get(idx int) (*Server, bool) {
	if *ss == nil {
		return nil, false
	}
	if idx < 0 || idx >= len(*ss) {
		return nil, false
	}
	return (*ss)[idx], true
}

func (ss *Servers) Append(val *Server) {
	if *ss == nil {
		*ss = Servers{val}
		return
	}
	(*ss) = append(*ss, val)
}

func (ss *Servers) Remove(s *Server) {
	if *ss == nil {
		return
	}
	for k, v := range *ss {
		if v == s {
			ss.RemoveIndex(k)
			return
		}
	}
}

func (ss *Servers) RemoveIndex(i int) {
	if *ss == nil {
		return // nothing to do
	}
	if i < 0 || i >= len(*ss) {
		return
	}
	copy((*ss)[i:], (*ss)[i+1:])
	(*ss)[len(*ss)-1] = nil
	(*ss) = (*ss)[:ss.Len()-1]
}

// Len returns the length of ss
func (ss *Servers) Len() int {
	if ss == nil || *ss == nil {
		return 0
	}
	return len(*ss)
}

type server Server

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
	Variables  ServerVariables `json:"variables,omitempty"`
	Extensions `json:"-"`
}

func (s *Server) Nodes() Nodes {
	return makeNodes(nodes{{"variables", s.Variables, KindServerVariables}})
}

// Kind returns KindServer
func (*Server) Kind() Kind {
	return KindServer
}

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
func (s Server) MarshalYAML() (interface{}, error) {
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

func (*ServerVariable) Kind() Kind { return KindServerVariable }

func (*ServerVariable) Nodes() Nodes { return nil }

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
func (sv ServerVariable) MarshalYAML() (interface{}, error) {
	return yamlutil.Marshal(sv)
}

// UnmarshalYAML unmarshals YAML data into s
func (sv *ServerVariable) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return yamlutil.Unmarshal(unmarshal, sv)
}

var (
	_ Node = (*Server)(nil)
	_ Node = (Servers)(nil)
)
