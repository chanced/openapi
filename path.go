package openapi

import (
	"encoding/json"
	"strings"

	"github.com/chanced/openapi/yamlutil"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v2"
)

// PathKind indicates whether the PathObj is a Path or a Reference
type PathKind uint8

const (
	// PathKindObj = PathObj
	PathKindObj PathKind = iota
	// PathKindRef = Reference
	PathKindRef
)

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

// Params returns all params in the path
func (pv PathValue) Params() []string {
	panic("not impl")
}

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
	Parameters *ParameterList `json:"parameters,omitempty"`
	Extensions `json:"-"`
}

type path PathObj

// PathKind returns PathKindPath
func (p *PathObj) PathKind() PathKind { return PathKindObj }

// MarshalJSON marshals p into JSON
func (p PathObj) MarshalJSON() ([]byte, error) {
	return marshalExtendedJSON(path(p))
}

// UnmarshalJSON unmarshals json into p
func (p *PathObj) UnmarshalJSON(data []byte) error {
	var v path
	if err := unmarshalExtendedJSON(data, &v); err != nil {
		return err
	}
	*p = PathObj(v)
	return nil
}

// MarshalYAML first marshals and unmarshals into JSON and then marshals into
// YAML
func (p *PathObj) MarshalYAML() (interface{}, error) {
	return yamlutil.Marshal(p)
}

// UnmarshalYAML unmarshals yaml into s
func (p *PathObj) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return yamlutil.Unmarshal(unmarshal, p)
}

// Path can either be a Path or a Reference
type Path interface {
	PathKind() PathKind
}

// PathItems is a map of Paths that can either be a Path or a Reference
type PathItems map[string]Path

// UnmarshalJSON unmarshals JSON data into rp
func (rp *PathItems) UnmarshalJSON(data []byte) error {
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
	*rp = res
	return nil

}

// UnmarshalYAML unmarshals YAML data into rp
func (rp *PathItems) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return yamlutil.Unmarshal(unmarshal, rp)
}

// MarshalYAML marshals rp into YAML
func (rp PathItems) MarshalYAML() (interface{}, error) {
	b, err := json.Marshal(rp)
	if err != nil {
		return nil, err
	}
	var v interface{}
	err = json.Unmarshal(b, &v)
	return v, err
}

func unmarshalPathJSON(data []byte) (Path, error) {
	if isRefJSON(data) {
		return unmarshalReferenceJSON(data)
	}
	var p path
	err := json.Unmarshal(data, &p)
	v := PathObj(p)
	return &v, err
}
