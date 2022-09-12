package openapi

import (
	"encoding/json"
	"strings"

	"github.com/chanced/jsonpointer"
	"github.com/tidwall/gjson"
)

type PathItemEntry struct {
	Key      string
	PathItem *PathItem
}

type PathItemObjs = ObjMap[*PathItem]

// PathItemMap is a map of Paths that can either be a Path or a Reference
type PathItemMap = ComponentMap[*PathItem]

// Paths holds the relative paths to the individual endpoints and their
// operations. The path is appended to the URL from the Server Object in order
// to construct the full URL. The Paths MAY be empty, due to Access Control List
// (ACL) constraints.
type Paths struct {
	Location   `json:"-"`
	Extensions `json:"-"`

	// Items are the Path
	Items PathItemObjs `json:"-"`
}

func (p *Paths) Resolve(ptr jsonpointer.Pointer) (Node, error) {
	if err := ptr.Validate(); err != nil {
		return nil, err
	}
	return p.resolve(ptr)
}

func (p *Paths) resolve(ptr jsonpointer.Pointer) (Node, error) {
	if ptr.IsRoot() {
		return p, nil
	}
	nxt, tok, _ := ptr.Next()
	v := p.Items.Get(Text(tok))
	if v == nil {
		return nil, newErrNotFound(p.Location.AbsoluteLocation(), tok)
	}
	return v.resolve(nxt)
}

func (*Paths) Kind() Kind      { return KindPaths }
func (*Paths) mapKind() Kind   { return KindUndefined }
func (*Paths) sliceKind() Kind { return KindUndefined }

func (p *Paths) setLocation(loc Location) error {
	if p == nil {
		return nil
	}
	p.Location = loc
	return p.Items.setLocation(loc)
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
			p.Items.Set(Text(key.String()), &v)
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
	Servers *ServerSlice `json:"servers,omitempty"`
	// A list of parameters that are applicable for all the operations described
	// under this path. These parameters can be overridden at the operation
	// level, but cannot be removed there. The list MUST NOT include duplicated
	// parameters. A unique parameter is defined by a combination of a name and
	// location. The list can use the Reference Object to link to parameters
	// that are defined at the OpenAPI Object's components/parameters.
	Parameters *ParameterSlice `json:"parameters,omitempty"`
	Location   `json:"-"`
	Extensions `json:"-"`
}

func (pi *PathItem) Resolve(ptr jsonpointer.Pointer) (Node, error) {
	if err := ptr.Validate(); err != nil {
		return nil, err
	}

	return pi.resolve(ptr)
}

func (pi *PathItem) resolve(ptr jsonpointer.Pointer) (Node, error) {
	if ptr.IsRoot() {
		return pi, nil
	}
	nxt, tok, _ := ptr.Next()
	var node node
	switch tok {
	case "get":
		node = pi.Get
	case "put":
		node = pi.Put
	case "post":
		node = pi.Post
	case "delete":
		node = pi.Delete
	case "options":
		node = pi.Options
	case "head":
		node = pi.Head
	case "patch":
		node = pi.Patch
	case "trace":
		node = pi.Trace
	case "servers":
		node = pi.Servers
	case "parameters":
		node = pi.Parameters
	default:
		return nil, newErrNotResolvable(pi.Location.AbsoluteLocation(), tok)
	}
	if nxt.IsRoot() {
		return node, nil
	}

	if node == nil {
		return nil, newErrNotFound(pi.Location.AbsoluteLocation(), tok)
	}
	return node.resolve(nxt)
}

func (*PathItem) mapKind() Kind { return KindPathItemMap }

func (*PathItem) sliceKind() Kind { return KindUndefined }

func (p *PathItem) setLocation(loc Location) error {
	if p == nil {
		return nil
	}
	p.Location = loc
	if err := p.Delete.setLocation(loc.Append("delete")); err != nil {
		return err
	}
	if err := p.Get.setLocation(loc.Append("get")); err != nil {
		return err
	}
	if err := p.Head.setLocation(loc.Append("head")); err != nil {
		return err
	}
	if err := p.Options.setLocation(loc.Append("options")); err != nil {
		return err
	}
	if err := p.Patch.setLocation(loc.Append("patch")); err != nil {
		return err
	}
	if err := p.Post.setLocation(loc.Append("post")); err != nil {
		return err
	}
	if err := p.Put.setLocation(loc.Append("put")); err != nil {
		return err
	}
	if err := p.Trace.setLocation(loc.Append("trace")); err != nil {
		return err
	}
	if err := p.Parameters.setLocation(loc.Append("parameters")); err != nil {
		return err
	}
	if err := p.Servers.setLocation(loc.Append("servers")); err != nil {
		return err
	}

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

var (
	_ node = (*PathItem)(nil)
	_ node = (*Paths)(nil)
)
