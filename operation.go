package openapi

import "github.com/chanced/jsonpointer"

// Operation describes a single API operation on a path.
type Operation struct {
	// Location contains information about the location of the node in the
	// document or referenced resource
	Location   `json:"-"`
	Extensions `json:"-"`

	// A list of tags for API documentation control. Tags can be used for
	// logical grouping of operations by resources or any other qualifier.
	Tags []Text `json:"tags,omitempty"`
	// A short summary of what the operation does.
	Summary Text `json:"summary,omitempty"`
	// A verbose explanation of the operation behavior. CommonMark syntax MAY be
	// used for rich text representation.
	Description Text `json:"description,omitempty"`
	// externalDocs	Additional external documentation.
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty"`
	// Unique string used to identify the operation. The id MUST be unique among
	// all operations described in the API. The operationId value is
	// case-sensitive. Tools and libraries MAY use the operationId to uniquely
	// identify an operation, therefore, it is RECOMMENDED to follow common
	// programming naming conventions.
	OperationID Text `json:"operationId,omitempty"`
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
	// Declares this operation to be deprecated. Consumers SHOULD refrain from
	// usage of the declared operation. Default value is false.
	Deprecated bool `json:"deprecated,omitempty"`
	// A declaration of which security mechanisms can be used for this
	// operation. The list of values includes alternative security requirement
	// objects that can be used. Only one of the security requirement objects
	// need to be satisfied to authorize a request. To make security optional,
	// an empty security requirement ({}) can be included in the array. This
	// definition overrides any declared top-level security. To remove a
	// top-level security declaration, an empty array can be used.
	Security *SecurityRequirements `json:"security,omitempty"`
	// An alternative server array to service this operation. If an alternative
	// server object is specified at the Path Item Object or Root level, it will
	// be overridden by this value.
	Servers *ServerSlice `json:"servers,omitempty"`
}

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

// ResolveNodeByPointers a Node by a json pointer
func (o *Operation) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
	err := ptr.Validate()
	if err != nil {
		return nil, err
	}
	return o.resolveNodeByPointer(ptr)
}

func (o *Operation) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
	if ptr.IsRoot() {
		return o, nil
	}
	nxt, tok, _ := ptr.Next()
	switch nxt {
	case "externalDocs":
		if o.ExternalDocs == nil {
			return nil, newErrNotFound(o.AbsoluteLocation(), tok)
		}
		return o.ExternalDocs.resolveNodeByPointer(nxt)
	case "parameters":
		if o.Parameters == nil {
			return nil, newErrNotFound(o.AbsoluteLocation(), tok)
		}
		return o.Parameters.resolveNodeByPointer(nxt)
	case "requestBody":
		if o.RequestBody == nil {
			return nil, newErrNotFound(o.AbsoluteLocation(), tok)
		}
		return o.RequestBody.resolveNodeByPointer(nxt)
	case "responses":
		if o.Responses == nil {
			return nil, newErrNotFound(o.AbsoluteLocation(), tok)
		}
		return o.Responses.resolveNodeByPointer(nxt)
	case "callbacks":
		if o.Callbacks == nil {
			return nil, newErrNotFound(o.AbsoluteLocation(), tok)
		}
		return o.Callbacks.resolveNodeByPointer(nxt)
	case "security":
		if o.Security == nil {
			return nil, newErrNotFound(o.AbsoluteLocation(), tok)
		}
		return o.Security.resolveNodeByPointer(nxt)
	case "servers":
		if o.Servers == nil {
			return nil, newErrNotFound(o.AbsoluteLocation(), tok)
		}
		return o.Servers.resolveNodeByPointer(nxt)
	default:
		return nil, newErrNotResolvable(o.Location.AbsoluteLocation(), tok)
	}
}

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

func (*Operation) Kind() Kind      { return KindOperation }
func (*Operation) mapKind() Kind   { return KindUndefined }
func (*Operation) sliceKind() Kind { return KindUndefined }

func (o *Operation) setLocation(loc Location) error {
	if o == nil {
		return nil
	}
	o.Location = loc

	if err := o.ExternalDocs.setLocation(loc.Append("externalDocs")); err != nil {
		return err
	}
	if err := o.Parameters.setLocation(loc.Append("parameters")); err != nil {
		return err
	}
	if err := o.RequestBody.setLocation(loc.Append("requestBody")); err != nil {
		return err
	}
	if err := o.Responses.setLocation(loc.Append("responses")); err != nil {
		return err
	}
	if err := o.Callbacks.setLocation(loc.Append("callbacks")); err != nil {
		return err
	}
	if err := o.Security.setLocation(loc.Append("security")); err != nil {
		return err
	}
	if err := o.Servers.setLocation(loc.Append("servers")); err != nil {
		return err
	}
	return nil
}

var _ node = (*Operation)(nil)
