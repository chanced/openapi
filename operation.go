package openapi

import (
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// Operation describes a single API operation on a path.
type Operation struct {
	// A list of tags for API documentation control. Tags can be used for
	// logical grouping of operations by resources or any other qualifier.
	Tags []string `json:"tags,omitempty"`
	// A short summary of what the operation does.
	Summary string `json:"summary,omitempty"`
	// A verbose explanation of the operation behavior. CommonMark syntax MAY be
	// used for rich text representation.
	Description string `json:"description,omitempty"`
	// externalDocs	Additional external documentation.
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty"`
	// Unique string used to identify the operation. The id MUST be unique among
	// all operations described in the API. The operationId value is
	// case-sensitive. Tools and libraries MAY use the operationId to uniquely
	// identify an operation, therefore, it is RECOMMENDED to follow common
	// programming naming conventions.
	OperationID string `json:"operationId,omitempty"`
	// A list of parameters that are applicable for this operation. If a
	// parameter is already defined at the Path Item, the new definition will
	// override it but can never remove it. The list MUST NOT include duplicated
	// parameters. A unique parameter is defined by a combination of a name and
	// location. The list can use the Reference Object to link to parameters
	// that are defined at the OpenAPI Object's components/parameters.
	Parameters *ParameterList `json:"parameters,omitempty"`

	// The request body applicable for this operation. The requestBody is fully
	// supported in HTTP methods where the HTTP 1.1 specification RFC7231 has
	// explicitly defined semantics for request bodies. In other cases where the
	// HTTP spec is vague (such as GET, HEAD and DELETE), requestBody is
	// permitted but does not have well-defined semantics and SHOULD be avoided
	// if possible.
	RequestBody RequestBody `json:"requestBody,omitempty"`
	// The list of possible responses as they are returned from executing this
	// operation.
	Responses Responses `json:"responses,omitempty"`

	// A map of possible out-of band callbacks related to the parent operation.
	// The key is a unique identifier for the Callback Object. Each value in the
	// map is a Callback Object that describes a request that may be initiated
	// by the API provider and the expected responses.
	Callbacks Callbacks `json:"callbacks,omitempty"`
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
	Security SecurityRequirements `json:"security,omitempty"`
	// An alternative server array to service this operation. If an alternative
	// server object is specified at the Path Item Object or Root level, it will
	// be overridden by this value.
	Servers    []*Server `json:"servers,omitempty"`
	Extensions `json:"-"`
}
type operation struct {
	Tags         []string             `json:"tags,omitempty"`
	Summary      string               `json:"summary,omitempty"`
	Description  string               `json:"description,omitempty"`
	ExternalDocs *ExternalDocs        `json:"externalDocs,omitempty"`
	OperationID  string               `json:"operationId,omitempty"`
	Parameters   *ParameterList       `json:"parameters,omitempty"`
	RequestBody  RequestBody          `json:"-"`
	Responses    Responses            `json:"responses,omitempty"`
	Callbacks    Callbacks            `json:"callbacks,omitempty"`
	Deprecated   bool                 `json:"deprecated,omitempty"`
	Security     SecurityRequirements `json:"security,omitempty"`
	Servers      []*Server            `json:"servers,omitempty"`
	Extensions   `json:"-"`
}

// MarshalJSON marshals JSON
func (o Operation) MarshalJSON() ([]byte, error) {
	b, err := marshalExtendedJSON(operation(o))
	if err != nil {
		return b, err
	}
	if o.RequestBody != nil {
		b, err = sjson.SetBytes(b, "requestBody", o.RequestBody)
	}
	return b, err
}

// UnmarshalJSON unmarshals JSON
func (o *Operation) UnmarshalJSON(data []byte) error {
	var v operation
	if err := unmarshalExtendedJSON(data, &v); err != nil {
		return err
	}
	r := gjson.GetBytes(data, "requestBody")
	if len(r.Raw) > 0 {
		var rb RequestBody
		if err := unmarshalRequestBody([]byte(r.Raw), &rb); err != nil {
			return err
		}
		v.RequestBody = rb
	}
	*o = Operation(v)
	return nil
}
