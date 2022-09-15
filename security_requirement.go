package openapi

import (
	"encoding/json"
	"fmt"

	"github.com/chanced/jsonpointer"
	"github.com/chanced/jsonx"
)

// TODO: make SecurityRequirement an ordered slice.

// SecurityRequirements is a list of SecurityRequirement
type SecurityRequirements = ObjMap[*SecurityRequirement]

type SecurityRequirementItem struct {
	Location
	Key   Text
	Value []Text
}

func (sri *SecurityRequirementItem) Refs() []Ref { return nil }
func (sri *SecurityRequirementItem) isNil() bool { return sri == nil }

func (sri *SecurityRequirementItem) Anchors() (*Anchors, error) { return nil, nil }

func (sri *SecurityRequirementItem) setLocation(loc Location) error {
	sri.Location = loc
	return nil
}

func (*SecurityRequirementItem) Kind() Kind      { return KindSecurityRequirementItem }
func (*SecurityRequirementItem) mapKind() Kind   { return KindSecurityRequirement }
func (*SecurityRequirementItem) sliceKind() Kind { return KindUndefined }

func (sri *SecurityRequirementItem) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
	if err := ptr.Validate(); err != nil {
		return nil, err
	}
	return sri.resolveNodeByPointer(ptr)
}

func (sri *SecurityRequirementItem) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
	if ptr.IsRoot() {
		return sri, nil
	}
	tok, _ := ptr.NextToken()
	return nil, newErrNotResolvable(sri.AbsoluteLocation(), tok)
}

func (sri SecurityRequirementItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(sri.Value)
}

func (sri *SecurityRequirementItem) UnmarshalJSON(data []byte) error {
	*sri = SecurityRequirementItem{}
	if len(data) == 0 {
		return nil
	}
	t := jsonx.TypeOf(data)
	switch t {
	case jsonx.TypeString:
		var v []Text
		err := json.Unmarshal(data, &v)
		if err != nil {
			return err
		}
		sri.Value = v
		return nil
	default:
		var v map[Text][]Text
		err := json.Unmarshal(data, &v)
		if err != nil {
			return err
		}
		if len(v) > 1 {
			return fmt.Errorf("can not unmarshal more than a single key/value pair into a Scope")
		}
		for k, v := range v {
			sri.Key = k
			sri.Value = v
			break
		}
		return nil
	}
}

// SecurityRequirement lists the required security schemes to execute this
// operation. The name used for each property MUST correspond to a security
// scheme declared in the Security Schemes under the Components Object.
//
// Security Requirement Objects that contain multiple schemes require that all
// schemes MUST be satisfied for a request to be authorized. This enables
// support for scenarios where multiple query parameters or HTTP headers are
// required to convey security information.
//
// When a list of Security Requirement Objects is defined on the OpenAPI Object
// or Operation Object, only one of the Security Requirement Objects in the list
// needs to be satisfied to authorize the request.
//
// Each name MUST correspond to a security scheme which is declared in the
// Security Schemes under the Components Object. If the security scheme is of
// type "oauth2" or "openIdConnect", then the value is a list of scope names
// required for the execution, and the list MAY be empty if authorization does
// not require a specified scope. For other security scheme types, the array MAY
// contain a list of role names which are required for the execution, but are
// not otherwise defined or exchanged in-band.
type SecurityRequirement = ObjMap[*SecurityRequirementItem]

var (
	_ node   = (*SecuritySchemeMap)(nil)
	_ Walker = (*SecuritySchemeMap)(nil)
	_ node   = (*SecurityRequirements)(nil)
	_ Walker = (*SecurityRequirements)(nil)
	_ node   = (*SecurityRequirement)(nil)
	_ Walker = (*SecurityRequirement)(nil)
	_ node   = (*SecurityRequirementItem)(nil)
	_ Walker = (*SecurityRequirementItem)(nil)
)
