package openapi

import (
	"encoding/json"

	"github.com/chanced/openapi/yamlutil"
)

// RequestBodyKind distinguishes a RequestBodyObj as either a RequestBody or
// Reference
type RequestBodyKind int

const (
	// RequestBodyKindObj = RequestBodyObj
	RequestBodyKindObj RequestBodyKind = iota
	// RequestBodyKindRef = Reference
	RequestBodyKindRef
)

// RequestBodyObj describes a single request body.
type RequestBodyObj struct {
	// A brief description of the request body. This could contain examples of
	// use. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`
	// The content of the request body. The key is a media type or media type range and the value describes it. For requests that match multiple keys, only the most specific key is applicable. e.g. text/plain overrides text
	//
	// *required*
	Content Content `json:"content,omitempty"`
	// Determines if the request body is required in the request. Defaults to false.
	Required bool `json:"required,omitempty"`

	Extensions `json:"-"`
}

type requestbody RequestBodyObj

// RequestBodyKind returns RequestBodyKindRequestBody
func (rb *RequestBodyObj) RequestBodyKind() RequestBodyKind { return RequestBodyKindObj }

// ResolveRequestBody resolves RequestBodyObj by returning itself. resolve is  not called.
func (rb *RequestBodyObj) ResolveRequestBody(RequestBodyResolverFunc) (*RequestBodyObj, error) {
	return rb, nil
}

// MarshalJSON marshals h into JSON
func (rb RequestBodyObj) MarshalJSON() ([]byte, error) {
	return marshalExtendedJSON(requestbody(rb))
}

// UnmarshalJSON unmarshals json into rb
func (rb *RequestBodyObj) UnmarshalJSON(data []byte) error {
	var v requestbody
	err := unmarshalExtendedJSON(data, &v)
	*rb = RequestBodyObj(v)
	return err
}

// UnmarshalYAML unmarshals YAML data into rb
func (rb *RequestBodyObj) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return yamlutil.Unmarshal(unmarshal, rb)
}

// MarshalYAML marshals rb into YAML
func (rb RequestBodyObj) MarshalYAML() (interface{}, error) {
	b, err := json.Marshal(rb)
	if err != nil {
		return nil, err
	}
	var v interface{}
	err = json.Unmarshal(b, &v)
	return v, err
}

// RequestBody can either be a RequestBody or a Reference
type RequestBody interface {
	ResolveRequestBody(RequestBodyResolverFunc) (*RequestBodyObj, error)
	RequestBodyKind() RequestBodyKind
}

func unmarshalRequestBody(data []byte, rb *RequestBody) error {
	if isRefJSON(data) {
		v, err := unmarshalReferenceJSON(data)
		*rb = v
		return err
	}
	var v RequestBodyObj
	err := json.Unmarshal(data, &v)
	*rb = &v
	return err
}

// RequestBodies is a map of RequestBody
type RequestBodies map[string]RequestBody

// UnmarshalJSON unmarshals JSON
func (rb RequestBodies) UnmarshalJSON(data []byte) error {
	var dm map[string]json.RawMessage
	if err := json.Unmarshal(data, &dm); err != nil {
		return err
	}
	for k, d := range dm {
		var v RequestBody
		if err := unmarshalRequestBody(d, &v); err != nil {
			return err
		}
		rb[k] = v
	}

	return nil
}

// ResolvedRequestBody describes a single request body.
type ResolvedRequestBody struct {

	// TODO: reference

	// A brief description of the request body. This could contain examples of
	// use. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`
	// The content of the request body. The key is a media type or media type range and the value describes it. For requests that match multiple keys, only the most specific key is applicable. e.g. text/plain overrides text
	//
	// *required*
	Content Content `json:"content,omitempty"`
	// Determines if the request body is required in the request. Defaults to false.
	Required bool `json:"required,omitempty"`

	Extensions `json:"-"`
}

// RequestBodies is a map of RequestBody
type ResolvedRequestBodies map[string]*ResolvedRequestBody
