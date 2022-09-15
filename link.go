package openapi

import (
	"encoding/json"

	"github.com/chanced/jsonpointer"
	"github.com/chanced/jsonx"
)

// LinkMap is a Map of either LinkMap or References to LinkMap
type (
	LinkMap = ComponentMap[*Link]
)

// Link represents a possible design-time link for a response. The presence of a
// link does not guarantee the caller's ability to successfully invoke it,
// rather it provides a known relationship and traversal mechanism between
// responses and other operations.
//
// Unlike dynamic links (i.e. links provided in the response payload), the OAS
// linking mechanism does not require link information in the runtime response.
//
// For computing links, and providing instructions to execute them, a runtime
// expression is used for accessing values in an operation and using them as
// parameters while invoking the linked operation.
type Link struct {
	Extensions `json:"-"`
	Location   `json:"-"`

	// A relative or absolute URI reference to an OAS operation. This field is
	// mutually exclusive of the operationId field, and MUST point to an
	// Operation Object. Relative operationRef values MAY be used to locate an
	// existing Operation Object in the OpenAPI definition. See the rules for
	// resolving Relative References.
	OperationRef *OperationRef `json:"operationRef,omitempty"`
	// The name of an existing, resolvable OAS operation, as defined with a
	// unique operationId. This field is mutually exclusive of the operationRef
	// field.
	OperationID Text `json:"operationId,omitempty"`
	// A map representing parameters to pass to an operation as specified with
	// operationId or identified via operationRef. The key is the parameter name
	// to be used, whereas the value can be a constant or an expression to be
	// evaluated and passed to the linked operation. The parameter name can be
	// qualified using the parameter location [{in}.]{name} for operations that
	// use the same parameter name in different locations (e.g. path.id).
	Parameters OrderedJSONObj `json:"parameters,omitempty"`
	// A literal value or {expression} to use as a request body when calling the
	// target operation.
	RequestBody jsonx.RawMessage `json:"requestBody,omitempty"`
	// A description of the link. CommonMark syntax MAY be used for rich text
	// representation.
	Description Text `json:"description,omitempty"`
}

func (l *Link) Refs() []Ref {
	if l == nil {
		return nil
	}
	var refs []Ref
	if l.OperationRef != nil {
		refs = append(refs, l.OperationRef)
	}
	return refs
}

// func (*Link) Walk(v Visitor) error {
// 	panic("Link.Walk() not implemented")
// }

func (l *Link) Anchors() (*Anchors, error) { return nil, nil }

func (l *Link) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
	if err := ptr.Validate(); err != nil {
		return nil, err
	}
	return l.resolveNodeByPointer(ptr)
}

func (l *Link) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
	if ptr.IsRoot() {
		return l, nil
	}
	tok, _ := ptr.NextToken()
	return nil, newErrNotResolvable(l.Location.AbsoluteLocation(), tok)
}

func (*Link) mapKind() Kind   { return KindLinkMap }
func (*Link) sliceKind() Kind { return KindUndefined }

// MarshalJSON marshals JSON
func (l Link) MarshalJSON() ([]byte, error) {
	type link Link
	return marshalExtendedJSON(link(l))
}

// UnmarshalJSON unmarshals JSON
func (l *Link) UnmarshalJSON(data []byte) error {
	type link Link
	var lv link
	if err := unmarshalExtendedJSON(data, &lv); err != nil {
		return err
	}
	*l = Link(lv)
	return nil
}

// DecodeRequestBody decodes l.RequestBody into dst
//
// dst should be a pointer to a concrete type
func (l *Link) DecodeRequestBody(dst interface{}) error {
	return json.Unmarshal(l.RequestBody, dst)
}

func (*Link) Kind() Kind { return KindLink }

func (l *Link) setLocation(loc Location) error {
	if l == nil {
		return nil
	}
	l.Location = loc
	if l.OperationRef != nil {
		l.OperationRef.Location = loc.Append("operationRef")
	}
	return nil
}

func (l *Link) isNil() bool { return l == nil }

// LinkParameters is a map representing parameters to pass to an operation as
// specified with operationId or identified via operationRef. The key is the
// parameter name to be used, whereas the value can be a constant or an
// expression to be evaluated and passed to the linked operation. The parameter
// name can be qualified using the parameter location [{in}.]{name} for
// operations that use the same parameter name in different locations (e.g.
// path.id).

var (
	_ node   = (*Link)(nil)
	_ Walker = (*Link)(nil)
)
