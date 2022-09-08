package openapi

import (
	"encoding/json"
	"strings"

	"github.com/tidwall/gjson"
)

// PathItems is a map of Paths that can either be a Path or a Reference
type PathItems Map[*Path]

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
// 	panic("not impl")
// }

// MarshalJSON Marshals PathEntry to JSON
func (pv PathValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(pv.String())
}

// Paths holds the relative paths to the individual endpoints and their
// operations. The path is appended to the URL from the Server Object in order
// to construct the full URL. The Paths MAY be empty, due to Access Control List
// (ACL) constraints.
type Paths struct {
	Items      map[PathValue]*Path `json:"-"`
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
		Items:      map[PathValue]*Path{},
		Extensions: Extensions{},
	}
	var err error
	gjson.ParseBytes(data).ForEach(func(key, value gjson.Result) bool {
		if strings.HasPrefix(key.String(), "x-") {
			p.SetEncodedExtension(key.String(), []byte(value.Raw))
		} else {
			var v Path
			err = json.Unmarshal([]byte(value.Raw), &v)
			p.Items[PathValue(key.String())] = &v
		}
		return err == nil
	})
	return err
}

// Path describes the operations available on a single path. A Path Item MAY
// be empty, due to ACL constraints. The path itself is still exposed to the
// documentation viewer but they will not know which operations and parameters
// are available.
type Path struct {
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

// MarshalJSON marshals p into JSON
func (p Path) MarshalJSON() ([]byte, error) {
	type path Path
	return marshalExtendedJSON(path(p))
}

// UnmarshalJSON unmarshals json into p
func (p *Path) UnmarshalJSON(data []byte) error {
	type path Path

	var v path
	if err := unmarshalExtendedJSON(data, &v); err != nil {
		return err
	}
	*p = Path(v)
	return nil
}

func (Path) Kind() Kind { return KindPath }
