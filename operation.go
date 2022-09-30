package openapi

import (
	"encoding/json"

	"github.com/chanced/transcode"
	"gopkg.in/yaml.v3"
)

// OperationItem is an *Operation and the HTTP Method it is associated with
//
// The primary purpose of this type is for use with a Visitor
// type OperationItem struct {
// 	Operation *Operation
// 	Method    Method
// }

// func (oi *OperationItem) Walk(v Visitor) error {
// 	if v == nil {
// 		return nil
// 	}
// 	if oi == nil {
// 		return nil
// 	}
// 	if oi.Operation == nil {
// 		return nil
// 	}
// 	var err error
// 	v, err = v.VisitOperationItem(oi)
// 	if v == nil {
// 		return err
// 	}
// 	if err != nil {
// 		return err
// 	}
// 	return oi.Operation.Walk(v)
// }

// Operation describes a single API operation on a path.
type Operation struct {
	// Location contains information about the location of the node in the
	// document or referenced resource
	Location   `json:"-"`
	Extensions `json:"-"`

	// Unique string used to identify the operation. The id MUST be unique among
	// all operations described in the API. The operationId value is
	// case-sensitive. Tools and libraries MAY use the operationId to uniquely
	// identify an operation, therefore, it is RECOMMENDED to follow common
	// programming naming conventions.
	OperationID Text `json:"operationId,omitempty"`

	// Declares this operation to be deprecated. Consumers SHOULD refrain from
	// usage of the declared operation. Default value is false.
	Deprecated bool `json:"deprecated,omitempty"`

	// A short summary of what the operation does.
	Summary Text `json:"summary,omitempty"`

	// A verbose explanation of the operation behavior. CommonMark syntax MAY be
	// used for rich text representation.
	Description Text `json:"description,omitempty"`

	// A list of tags for API documentation control. Tags can be used for
	// logical grouping of operations by resources or any other qualifier.
	Tags Texts `json:"tags,omitempty"`

	// A list of parameters that are applicable for this operation. If a
	// parameter is already defined at the Path Item, the new definition will
	// override it but can never remove it. The list MUST NOT include duplicated
	// parameters. A unique parameter is defined by a combination of a name and
	// location. The list can use the Reference Object to link to parameters
	// that are defined at the OpenAPI Object's components/parameters.
	Parameters *ParameterSlice `json:"parameters,omitempty"`

	// The request body applicable for this operation. The requestBody is fully
	// supported in HTTP methods where the HTTP 1.1 specification RFC7231 has
	// explicitly defined semantics for request bodies. In other cases where the
	// HTTP spec is vague (such as GET, HEAD and DELETE), requestBody is
	// permitted but does not have well-defined semantics and SHOULD be avoided
	// if possible.
	RequestBody *Component[*RequestBody] `json:"requestBody,omitempty"`

	// The list of possible responses as they are returned from executing this
	// operation.
	Responses *ResponseMap `json:"responses,omitempty"`

	// A map of possible out-of band callbacks related to the parent operation.
	// The key is a unique identifier for the Callback Object. Each value in the
	// map is a Callback Object that describes a request that may be initiated
	// by the API provider and the expected responses.
	Callbacks *CallbacksMap `json:"callbacks,omitempty"`

	// A declaration of which security mechanisms can be used for this
	// operation. The list of values includes alternative security requirement
	// objects that can be used. Only one of the security requirement objects
	// need to be satisfied to authorize a request. To make security optional,
	// an empty security requirement ({}) can be included in the array. This
	// definition overrides any declared top-level security. To remove a
	// top-level security declaration, an empty array can be used.
	Security *SecurityRequirementMap `json:"security,omitempty"`

	// An alternative server array to service this operation. If an alternative
	// server object is specified at the Path Item Object or Root level, it will
	// be overridden by this value.
	Servers *ServerSlice `json:"servers,omitempty"`

	// externalDocs	Additional external documentation.
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty"`
}

func (o *Operation) Nodes() []Node {
	if o == nil {
		return nil
	}
	return downcastNodes(o.nodes())
}

func (o *Operation) nodes() []node {
	if o == nil {
		return nil
	}
	return appendEdges(nil,
		o.Parameters,
		o.RequestBody,
		o.Responses,
		o.Callbacks,
		o.Security,
		o.Servers,
		o.ExternalDocs,
	)
}

func (*Operation) ref() Ref { return nil }

func (o *Operation) Refs() []Ref {
	if o == nil {
		return nil
	}
	var refs []Ref
	refs = append(refs, o.ExternalDocs.Refs()...)
	refs = append(refs, o.Parameters.Refs()...)
	refs = append(refs, o.RequestBody.Refs()...)
	refs = append(refs, o.Responses.Refs()...)
	refs = append(refs, o.Callbacks.Refs()...)
	refs = append(refs, o.Security.Refs()...)
	refs = append(refs, o.Servers.Refs()...)
	return refs
}

func (o *Operation) isNil() bool { return o == nil }

func (o *Operation) Anchors() (*Anchors, error) {
	if o == nil {
		return nil, nil
	}
	var anchors *Anchors
	var err error
	if anchors, err = anchors.merge(o.ExternalDocs.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(o.Parameters.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(o.RequestBody.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(o.Responses.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(o.Callbacks.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(o.Security.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(o.Servers.Anchors()); err != nil {
		return nil, err
	}
	return anchors, nil
}

// // ResolveNodeByPointer resolves a Node by a json pointer
// func (o *Operation) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if err := ptr.Validate(); err != nil {
// 		return nil, err
// 	}
// 	return o.resolveNodeByPointer(ptr)
// }

// func (o *Operation) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if ptr.IsRoot() {
// 		return o, nil
// 	}
// 	nxt, tok, _ := ptr.Next()
// 	switch nxt {
// 	case "externalDocs":
// 		if o.ExternalDocs == nil {
// 			return nil, newErrNotFound(o.AbsoluteLocation(), tok)
// 		}
// 		return o.ExternalDocs.resolveNodeByPointer(nxt)
// 	case "parameters":
// 		if o.Parameters == nil {
// 			return nil, newErrNotFound(o.AbsoluteLocation(), tok)
// 		}
// 		return o.Parameters.resolveNodeByPointer(nxt)
// 	case "requestBody":
// 		if o.RequestBody == nil {
// 			return nil, newErrNotFound(o.AbsoluteLocation(), tok)
// 		}
// 		return o.RequestBody.resolveNodeByPointer(nxt)
// 	case "responses":
// 		if o.Responses == nil {
// 			return nil, newErrNotFound(o.AbsoluteLocation(), tok)
// 		}
// 		return o.Responses.resolveNodeByPointer(nxt)
// 	case "callbacks":
// 		if o.Callbacks == nil {
// 			return nil, newErrNotFound(o.AbsoluteLocation(), tok)
// 		}
// 		return o.Callbacks.resolveNodeByPointer(nxt)
// 	case "security":
// 		if o.Security == nil {
// 			return nil, newErrNotFound(o.AbsoluteLocation(), tok)
// 		}
// 		return o.Security.resolveNodeByPointer(nxt)
// 	case "servers":
// 		if o.Servers == nil {
// 			return nil, newErrNotFound(o.AbsoluteLocation(), tok)
// 		}
// 		return o.Servers.resolveNodeByPointer(nxt)
// 	default:
// 		return nil, newErrNotResolvable(o.Location.AbsoluteLocation(), tok)
// 	}
// }

// MarshalJSON marshals JSON
func (o Operation) MarshalJSON() ([]byte, error) {
	type operation Operation
	b, err := marshalExtendedJSON(operation(o))
	if err != nil {
		return b, err
	}
	return b, err
}

// UnmarshalJSON unmarshals JSON
func (o *Operation) UnmarshalJSON(data []byte) error {
	type operation Operation
	var v operation
	if err := unmarshalExtendedJSON(data, &v); err != nil {
		return err
	}
	*o = Operation(v)
	return nil
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Marshaler interface
func (o Operation) MarshalYAML() (interface{}, error) {
	j, err := o.MarshalJSON()
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
func (o *Operation) UnmarshalYAML(value *yaml.Node) error {
	v, err := yaml.Marshal(value)
	if err != nil {
		return err
	}
	j, err := transcode.JSONFromYAML(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(j, o)
}

func (*Operation) Kind() Kind      { return KindOperation }
func (*Operation) mapKind() Kind   { return KindUndefined }
func (*Operation) sliceKind() Kind { return KindUndefined }

func (o *Operation) setLocation(loc Location) error {
	if o == nil {
		return nil
	}
	o.Location = loc
	var err error
	if err = o.ExternalDocs.setLocation(loc.AppendLocation("externalDocs")); err != nil {
		return err
	}
	if err = o.Parameters.setLocation(loc.AppendLocation("parameters")); err != nil {
		return err
	}
	if err = o.RequestBody.setLocation(loc.AppendLocation("requestBody")); err != nil {
		return err
	}
	if err = o.Responses.setLocation(loc.AppendLocation("responses")); err != nil {
		return err
	}
	if err = o.Callbacks.setLocation(loc.AppendLocation("callbacks")); err != nil {
		return err
	}
	if err = o.Security.setLocation(loc.AppendLocation("security")); err != nil {
		return err
	}
	if err = o.Servers.setLocation(loc.AppendLocation("servers")); err != nil {
		return err
	}
	return nil
}

// func (o *Operation) Walk(v Visitor) error {
// 	if v == nil {
// 		return nil
// 	}
// 	if o == nil {
// 		return nil
// 	}

// 	var err error
// 	v, err = v.Visit(o)
// 	if err != nil {
// 		return err
// 	}
// 	if v == nil {
// 		return nil
// 	}
// 	v, err = v.VisitOperation(o)
// 	if err != nil {
// 		return err
// 	}
// 	if v == nil {
// 		return nil
// 	}

// 	if o.Parameters != nil {
// 		err = o.Parameters.Walk(v)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	if o.RequestBody != nil {
// 		err = o.RequestBody.Walk(v)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	if o.Responses != nil {
// 		err = o.Responses.Walk(v)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	if o.Callbacks != nil {
// 		err = o.Callbacks.Walk(v)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	if o.Security != nil {
// 		err = o.Security.Walk(v)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	if o.Servers != nil {
// 		err = o.Servers.Walk(v)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	if o.ExternalDocs != nil {
// 		err = o.ExternalDocs.Walk(v)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

var _ node = (*Operation)(nil)
