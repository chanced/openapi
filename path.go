package openapi

import (
	"encoding/json"
	"strings"

	"github.com/chanced/openapi/yamlutil"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v2"
)

// Path can either be a Path or a Reference
type Path interface {
	Node
	ResolvePath(func(ref string) (*PathObj, error)) (*PathObj, error)
}

// PathValue is relative path to an individual endpoint. The path is appended
// (no relative URL resolution) to the expanded URL from the Server Object's url
// field in order to construct the full URL. PathValue templating is allowed. When
// matching URLs, concrete (non-templated) paths would be matched before their
// templated counterparts. Templated paths with the same hierarchy but different
// templated names MUST NOT exist as they are identical. In case of ambiguous
// matching, it's up to the tooling to decide which one to use.
type PathValue string

func (pv PathValue) String() string {
	str := string(pv)
	if len(pv) == 0 {
		return "/"
	}
	if pv[0] != '/' {
		return "/" + str
	}
	return str
}

// // Params returns all params in the path
// func (pv PathValue) Params() []string {
// }

// MarshalJSON Marshals PathEntry to JSON
func (pv PathValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(pv.String())
}

// MarshalYAML Marshals PathEntry to YAML
func (pv PathValue) MarshalYAML() ([]byte, error) {
	return yaml.Marshal(pv.String())
}

// Paths holds the relative paths to the individual endpoints and their
// operations. The path is appended to the URL from the Server Object in order
// to construct the full URL. The Paths MAY be empty, due to Access Control List
// (ACL) constraints.
type Paths struct {
	Items      map[PathValue]*PathObj `json:"-"`
	Extensions `json:"-"`
}

// MarshalJSON marshals JSON
func (p Paths) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{}, len(p.Items)+len(p.Extensions))
	for k, v := range p.Items {
		m[k.String()] = v
	}
	for k, v := range p.Extensions {
		m[k] = v
	}
	return json.Marshal(m)
}

// UnmarshalJSON unmarshals JSON data into p
func (p *Paths) UnmarshalJSON(data []byte) error {
	*p = Paths{
		Items:      map[PathValue]*PathObj{},
		Extensions: Extensions{},
	}
	var err error
	gjson.ParseBytes(data).ForEach(func(key, value gjson.Result) bool {
		if strings.HasPrefix(key.String(), "x-") {
			p.SetEncodedExtension(key.String(), []byte(value.Raw))
		} else {
			var v PathObj
			err = json.Unmarshal([]byte(value.Raw), &v)
			p.Items[PathValue(key.String())] = &v
		}
		return err == nil
	})
	return err
}

// PathItems is a map of Paths that can either be a Path or a Reference
type PathItems map[string]Path

// Kind returns KindPathItems
func (PathItems) Kind() Kind {
	return KindPathItems
}

// UnmarshalJSON unmarshals JSON data into pi
func (pi *PathItems) UnmarshalJSON(data []byte) error {
	var rd map[string]json.RawMessage
	err := json.Unmarshal(data, &rd)
	if err != nil {
		return err
	}
	res := PathItems{}
	for k, d := range rd {
		if isRefJSON(data) {
			var v Reference
			if err = json.Unmarshal(d, &v); err != nil {
				return err
			}
			res[k] = &v
		} else {
			var v PathObj
			if err = json.Unmarshal(d, &v); err != nil {
				return err
			}
			res[k] = &v
		}
	}
	*pi = res
	return nil
}

// UnmarshalYAML unmarshals YAML data into pi
func (pi *PathItems) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return yamlutil.Unmarshal(unmarshal, pi)
}

// MarshalYAML marshals pi into YAML
func (pi PathItems) MarshalYAML() (interface{}, error) {
	b, err := json.Marshal(pi)
	if err != nil {
		return nil, err
	}
	var v interface{}
	err = json.Unmarshal(b, &v)
	return v, err
}

// PathObj describes the operations available on a single path. A PathObj Item MAY
// be empty, due to ACL constraints. The path itself is still exposed to the
// documentation viewer but they will not know which operations and parameters
// are available.
type PathObj struct {
	// Allows for a referenced definition of this path item. The referenced
	// structure MUST be in the form of a Path Item Object. In case a Path Item
	// Object field appears both in the defined object and the referenced
	// object, the behavior is undefined. See the rules for resolving Relative
	// References.
	Ref string `json:"$ref,omitempty"`
	// An optional, string summary, intended to apply to all operations in this path.
	Summary string `json:"summary,omitempty"`
	// An optional, string description, intended to apply to all operations in
	// this path. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`
	// A definition of a GET operation on this path.
	Get *Operation `json:"get,omitempty"`
	// A definition of a PUT operation on this path.
	Put *Operation `json:"put,omitempty"`
	// A definition of a POST operation on this path.
	Post *Operation `json:"post,omitempty"`
	// A definition of a DELETE operation on this path.
	Delete *Operation `json:"delete,omitempty"`
	// A definition of a OPTIONS operation on this path.
	Options *Operation `json:"options,omitempty"`
	// A definition of a HEAD operation on this path.
	Head *Operation `json:"head,omitempty"`
	// A definition of a PATCH operation on this path.
	Patch *Operation `json:"patch,omitempty"`
	// A definition of a TRACE operation on this path.
	Trace *Operation `json:"trace,omitempty"`
	// An alternative server array to service all operations in this path.
	Servers []*Server `json:"servers,omitempty"`
	// A list of parameters that are applicable for all the operations described
	// under this path. These parameters can be overridden at the operation
	// level, but cannot be removed there. The list MUST NOT include duplicated
	// parameters. A unique parameter is defined by a combination of a name and
	// location. The list can use the Reference Object to link to parameters
	// that are defined at the OpenAPI Object's components/parameters.
	Parameters *ParameterSet `json:"parameters,omitempty"`
	Extensions `json:"-"`
}

func (p *PathObj) Nodes() map[string]*NodeDetail {
	m := make(map[string]*NodeDetail)
	if p.Get != nil {
		m["get"] = &NodeDetail{
			Node: p.Get,
		}
	}
	if p.Put != nil {
		m["put"] = &NodeDetail{
			Node:       p.Put,
			TargetKind: KindOperation,
		}
	}
	if p.Post != nil {
		m["post"] = &NodeDetail{
			Node:       p.Post,
			TargetKind: KindOperation,
		}
	}
	if p.Delete != nil {
		m["delete"] = &NodeDetail{
			Node:       p.Delete,
			TargetKind: KindOperation,
		}
	}
	if p.Options != nil {
		m["options"] = &NodeDetail{
			Node:       p.Options,
			TargetKind: KindOperation,
		}
	}
	if p.Head != nil {
		m["head"] = &NodeDetail{
			Node:       p.Head,
			TargetKind: KindOperation,
		}
	}
	if p.Patch != nil {
		m["patch"] = &NodeDetail{
			Node:       p.Patch,
			TargetKind: KindOperation,
		}
	}
	if p.Trace != nil {
		m["trace"] = &NodeDetail{
			Node:       p.Trace,
			TargetKind: KindOperation,
		}
	}
	if p.Parameters != nil {
		m["parameters"] = &NodeDetail{
			Node:       p.Parameters,
			TargetKind: KindParameterSet,
		}
	}
	return m
}

// Kind returns KindPath
func (*PathObj) Kind() Kind {
	return KindPath
}

type pathobj PathObj

// ResolvePath resolves PathObj by returning itself. resolve is  not called.
func (p *PathObj) ResolvePath(func(ref string) (*PathObj, error)) (*PathObj, error) {
	return p, nil
}

// MarshalJSON marshals p into JSON
func (p PathObj) MarshalJSON() ([]byte, error) {
	return marshalExtendedJSON(pathobj(p))
}

// UnmarshalJSON unmarshals json into p
func (p *PathObj) UnmarshalJSON(data []byte) error {
	var v pathobj
	if err := unmarshalExtendedJSON(data, &v); err != nil {
		return err
	}
	*p = PathObj(v)
	return nil
}

// MarshalYAML first marshals and unmarshals into JSON and then marshals into
// YAML
func (p PathObj) MarshalYAML() (interface{}, error) {
	return yamlutil.Marshal(p)
}

// UnmarshalYAML unmarshals yaml into s
func (p *PathObj) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return yamlutil.Unmarshal(unmarshal, p)
}

func unmarshalPathJSON(data []byte) (Path, error) {
	if isRefJSON(data) {
		return unmarshalReferenceJSON(data)
	}
	var p pathobj
	err := json.Unmarshal(data, &p)
	v := PathObj(p)
	return &v, err
}

// ResolvedPath is a Path Object which has beeen resolved. It describes the
// operations available on a single path. A PathObj Item MAY be empty, due to
// ACL constraints. The path itself is still exposed to the documentation viewer
// but they will not know which operations and parameters are available.
type ResolvedPath struct {
	// Allows for a referenced definition of this path item. The referenced
	// structure MUST be in the form of a Path Item Object. In case a Path Item
	// Object field appears both in the defined object and the referenced
	// object, the behavior is undefined. See the rules for resolving Relative
	// References.
	Ref string `json:"$ref,omitempty"`
	// An optional, string summary, intended to apply to all operations in this path.
	Summary string `json:"summary,omitempty"`
	// An optional, string description, intended to apply to all operations in
	// this path. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`
	// A definition of a GET operation on this path.
	Get *Operation `json:"get,omitempty"`
	// A definition of a PUT operation on this path.
	Put *Operation `json:"put,omitempty"`
	// A definition of a POST operation on this path.
	Post *Operation `json:"post,omitempty"`
	// A definition of a DELETE operation on this path.
	Delete *Operation `json:"delete,omitempty"`
	// A definition of a OPTIONS operation on this path.
	Options *Operation `json:"options,omitempty"`
	// A definition of a HEAD operation on this path.
	Head *Operation `json:"head,omitempty"`
	// A definition of a PATCH operation on this path.
	Patch *Operation `json:"patch,omitempty"`
	// A definition of a TRACE operation on this path.
	Trace *Operation `json:"trace,omitempty"`
	// An alternative server array to service all operations in this path.
	Servers []*Server `json:"servers,omitempty"`
	// A list of parameters that are applicable for all the operations described
	// under this path. These parameters can be overridden at the operation
	// level, but cannot be removed there. The list MUST NOT include duplicated
	// parameters. A unique parameter is defined by a combination of a name and
	// location. The list can use the Reference Object to link to parameters
	// that are defined at the OpenAPI Object's components/parameters.
	Parameters *ResolvedParameterSet `json:"parameters,omitempty"`
	Extensions `json:"-"`
}

// Kind returns KindResolvedPath
func (*ResolvedPath) Kind() Kind {
	return KindResolvedPath
}

// ResolvedPathItems is a map of resolved Path objects
type ResolvedPathItems map[string]*ResolvedPath

// Kind returns KindResolvedPathItems
func (ResolvedPathItems) Kind() Kind {
	return KindResolvedPathItems
}

// ResolvedPaths holds the relative paths to the individual endpoints and their
// operations. The path is appended to the URL from the Server Object in order
// to construct the full URL. The Paths MAY be empty, due to Access Control List
// (ACL) constraints.
type ResolvedPaths struct {
	Items      map[PathValue]*ResolvedPath `json:"-"`
	Extensions `json:"-"`
}

// Kind returns KindResolvedPaths
func (*ResolvedPaths) Kind() Kind {
	return KindResolvedPaths
}

// MarshalJSON marshals JSON
func (rp ResolvedPaths) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{}, len(rp.Items)+len(rp.Extensions))
	for k, v := range rp.Items {
		m[k.String()] = v
	}
	for k, v := range rp.Extensions {
		m[k] = v
	}
	return json.Marshal(m)
}

var (
	_ Node = (*PathObj)(nil)
	_ Node = (*PathItems)(nil)
	_ Node = (PathItems)(nil)
	_ Node = (*ResolvedPath)(nil)
	_ Node = (*ResolvedPaths)(nil)
	_ Node = (ResolvedPathItems)(nil)
)
