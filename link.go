package openapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// ErrLinkParameterNotFound is returned if a
var ErrLinkParameterNotFound = errors.New("error: link parameter not found")

// Links is a Map of either Links or References to Links
type Links Map[*Link]

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
	// A relative or absolute URI reference to an OAS operation. This field is
	// mutually exclusive of the operationId field, and MUST point to an
	// Operation Object. Relative operationRef values MAY be used to locate an
	// existing Operation Object in the OpenAPI definition. See the rules for
	// resolving Relative References.
	OperationRef string `json:"operationRef,omitempty"`
	// The name of an existing, resolvable OAS operation, as defined with a
	// unique operationId. This field is mutually exclusive of the operationRef
	// field.
	OperationID string `json:"operationId,omitempty"`
	// A map representing parameters to pass to an operation as specified with
	// operationId or identified via operationRef. The key is the parameter name
	// to be used, whereas the value can be a constant or an expression to be
	// evaluated and passed to the linked operation. The parameter name can be
	// qualified using the parameter location [{in}.]{name} for operations that
	// use the same parameter name in different locations (e.g. path.id).
	Parameters LinkParameters `json:"parameters,omitempty"`
	// A literal value or {expression} to use as a request body when calling the
	// target operation.
	RequestBody json.RawMessage `json:"requestBody,omitempty"`
	// A description of the link. CommonMark syntax MAY be used for rich text
	// representation.
	Description string `json:"description,omitempty"`
	Extensions  `json:"-"`
}

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

func (Link) Kind() Kind { return KindLink }

// LinkParameters is a map representing parameters to pass to an operation as
// specified with operationId or identified via operationRef. The key is the
// parameter name to be used, whereas the value can be a constant or an
// expression to be evaluated and passed to the linked operation. The parameter
// name can be qualified using the parameter location [{in}.]{name} for
// operations that use the same parameter name in different locations (e.g.
// path.id).
type LinkParameters map[string]json.RawMessage

// Decode decodes all parameters into dst
func (lp LinkParameters) Decode(dst interface{}) error {
	data, err := json.Marshal(lp)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}

// DecodeParameter decodes a given parameter by name.
// It returns an error if the parameter can not be found.
func (lp LinkParameters) DecodeParameter(key string, dst interface{}) error {
	if v, ok := lp[key]; ok {
		if reflect.TypeOf(dst).Kind() == reflect.Ptr {
			return json.Unmarshal(v, dst)
		}
		return json.Unmarshal(v, &dst)
	}
	return fmt.Errorf("%w {%s}", ErrLinkParameterNotFound, key)
}

// Has returns true if key exists in lp
func (lp LinkParameters) Has(key string) bool {
	_, exists := lp[key]
	return exists
}

// Set concrete object to lp. To add JSON, use SetEncoded
func (lp *LinkParameters) Set(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	lp.SetEncoded(key, data)
	return nil
}

// SetEncoded sets the value of key to value.
//
// Value should be a json encoded byte slice
func (lp *LinkParameters) SetEncoded(key string, value []byte) {
	(*lp)[key] = json.RawMessage(value)
}
