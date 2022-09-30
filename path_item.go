package openapi

import (
	"encoding/json"

	"github.com/chanced/transcode"
	"gopkg.in/yaml.v3"
)

// PathItem describes the operations available on a single path. A PathItem Item MAY
// be empty, due to ACL constraints. The path itself is still exposed to the
// documentation viewer but they will not know which operations and parameters
// are available.
type PathItem struct {
	Location   `json:"-"`
	Extensions `json:"-"`

	// An optional, string summary, intended to apply to all operations in this path.
	Summary Text `json:"summary,omitempty"`

	// An optional, string description, intended to apply to all operations in
	// this path. CommonMark syntax MAY be used for rich text representation.
	Description Text `json:"description,omitempty"`

	// A list of parameters that are applicable for all the operations described
	// under this path. These parameters can be overridden at the operation
	// level, but cannot be removed there. The list MUST NOT include duplicated
	// parameters. A unique parameter is defined by a combination of a name and
	// location. The list can use the Reference Object to link to parameters
	// that are defined at the OpenAPI Object's components/parameters.
	Parameters *ParameterSlice `json:"parameters,omitempty"`

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
}

func (pi *PathItem) Nodes() []Node {
	if pi == nil {
		return nil
	}
	return downcastNodes(pi.nodes())
}

func (pi *PathItem) nodes() []node {
	if pi == nil {
		return nil
	}
	var edges []node
	edges = appendEdges(edges, pi.Servers)
	edges = appendEdges(edges, pi.Parameters)
	edges = appendEdges(edges, pi.Get)
	edges = appendEdges(edges, pi.Put)
	edges = appendEdges(edges, pi.Post)
	edges = appendEdges(edges, pi.Delete)
	edges = appendEdges(edges, pi.Options)
	edges = appendEdges(edges, pi.Head)
	edges = appendEdges(edges, pi.Patch)
	edges = appendEdges(edges, pi.Trace)
	return edges
}
func (pi *PathItem) ref() Ref { return nil }

func (pi *PathItem) Refs() []Ref {
	if pi == nil {
		return nil
	}
	var refs []Ref
	refs = append(refs, pi.Servers.Refs()...)
	refs = append(refs, pi.Parameters.Refs()...)
	refs = append(refs, pi.Get.Refs()...)
	refs = append(refs, pi.Put.Refs()...)
	refs = append(refs, pi.Post.Refs()...)
	refs = append(refs, pi.Delete.Refs()...)
	refs = append(refs, pi.Options.Refs()...)
	refs = append(refs, pi.Head.Refs()...)
	refs = append(refs, pi.Patch.Refs()...)
	refs = append(refs, pi.Trace.Refs()...)
	return refs
}

func (pi *PathItem) Anchors() (*Anchors, error) {
	if pi == nil {
		return nil, nil
	}
	var anchors *Anchors
	var err error
	if anchors, err = anchors.merge(pi.Get.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(pi.Put.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(pi.Post.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(pi.Delete.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(pi.Options.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(pi.Head.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(pi.Patch.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(pi.Trace.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(pi.Servers.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(pi.Parameters.Anchors()); err != nil {
		return nil, err
	}
	return anchors, err
}

func (*PathItem) mapKind() Kind { return KindPathItemMap }

func (*PathItem) sliceKind() Kind { return KindUndefined }

func (p *PathItem) setLocation(loc Location) error {
	if p == nil {
		return nil
	}
	p.Location = loc
	var err error
	if err = p.Delete.setLocation(loc.AppendLocation("delete")); err != nil {
		return err
	}
	if err = p.Get.setLocation(loc.AppendLocation("get")); err != nil {
		return err
	}
	if err = p.Head.setLocation(loc.AppendLocation("head")); err != nil {
		return err
	}
	if err = p.Options.setLocation(loc.AppendLocation("options")); err != nil {
		return err
	}
	if err = p.Patch.setLocation(loc.AppendLocation("patch")); err != nil {
		return err
	}
	if err = p.Post.setLocation(loc.AppendLocation("post")); err != nil {
		return err
	}
	if err = p.Put.setLocation(loc.AppendLocation("put")); err != nil {
		return err
	}
	if err = p.Trace.setLocation(loc.AppendLocation("trace")); err != nil {
		return err
	}
	if err = p.Parameters.setLocation(loc.AppendLocation("parameters")); err != nil {
		return err
	}
	if err = p.Servers.setLocation(loc.AppendLocation("servers")); err != nil {
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

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Marshaler interface
func (p PathItem) MarshalYAML() (interface{}, error) {
	j, err := p.MarshalJSON()
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
func (p *PathItem) UnmarshalYAML(value *yaml.Node) error {
	v, err := yaml.Marshal(value)
	if err != nil {
		return err
	}
	j, err := transcode.JSONFromYAML(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(j, p)
}

func (*PathItem) Kind() Kind { return KindPathItem }

func (pi *PathItem) isNil() bool { return pi == nil }

func (*PathItem) refable() {}

var _ node = (*PathItem)(nil)

// func (pi *PathItem) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if err := ptr.Validate(); err != nil {
// 		return nil, err
// 	}

// 	return pi.resolveNodeByPointer(ptr)
// }

// func (pi *PathItem) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if ptr.IsRoot() {
// 		return pi, nil
// 	}
// 	nxt, tok, _ := ptr.Next()
// 	switch tok {
// 	case "get":
// 		if pi.Get == nil {
// 			return nil, newErrNotFound(pi.Location.AbsoluteLocation(), tok)
// 		}
// 		return pi.resolveNodeByPointer(nxt)
// 	case "put":
// 		if pi.Put == nil {
// 			return nil, newErrNotFound(pi.Location.AbsoluteLocation(), tok)
// 		}
// 		return pi.Put.resolveNodeByPointer(nxt)
// 	case "post":
// 		if pi.Post == nil {
// 			return nil, newErrNotFound(pi.Location.AbsoluteLocation(), tok)
// 		}
// 		return pi.Post.resolveNodeByPointer(nxt)
// 	case "delete":
// 		if pi.Delete == nil {
// 			return nil, newErrNotFound(pi.AbsoluteLocation(), tok)
// 		}
// 		return pi.Delete.resolveNodeByPointer(nxt)
// 	case "options":
// 		if pi.Options == nil {
// 			return nil, newErrNotFound(pi.AbsoluteLocation(), tok)
// 		}
// 		return pi.Options.resolveNodeByPointer(nxt)
// 	case "head":
// 		if pi.Head == nil {
// 			return nil, newErrNotFound(pi.AbsoluteLocation(), tok)
// 		}
// 		return pi.Head.resolveNodeByPointer(nxt)
// 	case "patch":
// 		if pi.Patch == nil {
// 			return nil, newErrNotFound(pi.AbsoluteLocation(), tok)
// 		}
// 		return pi.Patch.resolveNodeByPointer(nxt)
// 	case "trace":
// 		if pi.Trace == nil {
// 			return nil, newErrNotFound(pi.AbsoluteLocation(), tok)
// 		}
// 		return pi.Trace.resolveNodeByPointer(nxt)
// 	case "servers":
// 		if pi.Servers == nil {
// 			return nil, newErrNotFound(pi.AbsoluteLocation(), tok)
// 		}
// 		return pi.Servers.resolveNodeByPointer(nxt)
// 	case "parameters":
// 		if pi.Parameters == nil {
// 			return nil, newErrNotFound(pi.AbsoluteLocation(), tok)
// 		}
// 		return pi.Parameters.resolveNodeByPointer(nxt)
// 	default:
// 		return nil, newErrNotResolvable(pi.Location.AbsoluteLocation(), tok)
// 	}
// }

// func (pi *PathItem) Walk(v Visitor) error {
// 	if v == nil {
// 		return nil
// 	}
// 	var err error
// 	v, err = v.Visit(pi)
// 	if err != nil {
// 		return err
// 	}
// 	if v == nil {
// 		return nil
// 	}
// 	v, err = v.VisitPathItem(pi)
// 	if err != nil {
// 		return err
// 	}
// 	if v == nil {
// 		return nil
// 	}

// 	if pi.Parameters != nil {
// 		if err = pi.Parameters.Walk(v); err != nil {
// 			return err
// 		}
// 	}

// 	if pi.Servers != nil {
// 		if err = pi.Servers.Walk(v); err != nil {
// 			return err
// 		}
// 	}

// 	var op OperationItem
// 	if pi.Get != nil {
// 		op = OperationItem{Operation: pi.Get, Method: MethodGet}
// 		if err = op.Walk(v); err != nil {
// 			return err
// 		}
// 	}
// 	if pi.Put != nil {
// 		op = OperationItem{Operation: pi.Put, Method: MethodPut}
// 		if err = op.Walk(v); err != nil {
// 			return err
// 		}
// 	}
// 	if pi.Post != nil {
// 		op = OperationItem{Operation: pi.Post, Method: MethodPost}
// 		if err = op.Walk(v); err != nil {
// 			return err
// 		}
// 	}
// 	if pi.Delete != nil {
// 		op = OperationItem{Operation: pi.Delete, Method: MethodDelete}
// 		if err = op.Walk(v); err != nil {
// 			return err
// 		}
// 	}
// 	if pi.Options != nil {
// 		op = OperationItem{Operation: pi.Options, Method: MethodOptions}
// 		if err = op.Walk(v); err != nil {
// 			return err
// 		}
// 	}
// 	if pi.Head != nil {
// 		op = OperationItem{Operation: pi.Head, Method: MethodHead}
// 		if err = op.Walk(v); err != nil {
// 			return err
// 		}
// 	}
// 	if pi.Patch != nil {
// 		op = OperationItem{Operation: pi.Patch, Method: MethodPatch}
// 		if err = op.Walk(v); err != nil {
// 			return err
// 		}
// 	}
// 	if pi.Trace != nil {
// 		if err = pi.Trace.Walk(v); err != nil {
// 			return err
// 		}
// 	}

//		return nil
//	}
