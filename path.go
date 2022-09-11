package openapi

import (
	"encoding/json"
	"strings"

	"github.com/tidwall/gjson"
)

// PathItemMap is a map of Paths that can either be a Path or a Reference
type PathItemMap = ComponentMap[*PathItem]

// Paths holds the relative paths to the individual endpoints and their
// operations. The path is appended to the URL from the Server Object in order
// to construct the full URL. The Paths MAY be empty, due to Access Control List
// (ACL) constraints.
type Paths struct {
	// Items are the Path
	Items      PathItemMap `json:"-"`
	Extensions `json:"-"`
}

// MarshalJSON marshals JSON
func (p Paths) MarshalJSON() ([]byte, error) {
	j, err := p.Items.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return p.marshalExtensionsInto(j)
}

// UnmarshalJSON unmarshals JSON data into p
func (p *Paths) UnmarshalJSON(data []byte) error {
	*p = Paths{
		Extensions: Extensions{},
	}
	var err error
	gjson.ParseBytes(data).ForEach(func(key, value gjson.Result) bool {
		if strings.HasPrefix(key.String(), "x-") {
			p.SetRawExtension(key.String(), []byte(value.Raw))
		} else {
			var v PathItem
			err = json.Unmarshal([]byte(value.Raw), &v)
			p.Items = append(p.Items, ComponentEntry[*PathItem]{
				Key:       key.String(),
				Component: Component[*PathItem]{Object: &v},
			})
		}
		return err == nil
	})
	return err
}

// PathItem describes the operations available on a single path. A PathItem Item MAY
// be empty, due to ACL constraints. The path itself is still exposed to the
// documentation viewer but they will not know which operations and parameters
// are available.
type PathItem struct {
	// An optional, string summary, intended to apply to all operations in this path.
	Summary Text `json:"summary,omitempty"`
	// An optional, string description, intended to apply to all operations in
	// this path. CommonMark syntax MAY be used for rich text representation.
	Description Text `json:"description,omitempty"`
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
	Parameters ParameterSet `json:"parameters,omitempty"`
	Location   *Location     `json:"-"`
	Extensions `json:"-"`
}

// setLocation implements node
func (p *PathItem) setLocation(loc Location) error {
	p.Location = &loc
	p.
	return nil
}

// MarshalJSON marshals p into JSON
func (p PathItem) MarshalJSON() ([]byte, error) {
	type path PathItem
	return marshalExtendedJSON(path(p))
}

// UnmarshalJSON unmarshals json into p
func (p *PathItem) UnmarshalJSON(data []byte) error {
	type path PathItem

	var v path
	if err := unmarshalExtendedJSON(data, &v); err != nil {
		return err
	}
	*p = PathItem(v)
	return nil
}

func (*PathItem) Kind() Kind { return KindPathItem }

var _ node = (*PathItem)(nil)
