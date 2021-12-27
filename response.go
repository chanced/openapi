package openapi

import (
	"encoding/json"

	"github.com/chanced/openapi/yamlutil"
)

// Response is either a Response or a Reference
type Response interface {
	Node
	ResolveResponse(func(ref string) (*ResponseObj, error)) (*ResponseObj, error)
}

// Responses is a container for the expected responses of an operation. The
// container maps a HTTP response code to the expected response.
//
// The documentation is not necessarily expected to cover all possible HTTP
// response codes because they may not be known in advance. However,
// documentation is expected to cover a successful operation response and any
// known errors.
//
// The default MAY be used as a default response object for all HTTP codes that
// are not covered individually by the Responses Object.
//
// The Responses Object MUST contain at least one response code, and if only one
// response code is provided it SHOULD be the response for a successful
// operation call.
type Responses map[string]Response

// Kind returns KindResponses
func (Responses) Kind() Kind {
	return KindResponses
}

// UnmarshalJSON unmarshals JSON data into r
func (r *Responses) UnmarshalJSON(data []byte) error {
	var m map[string]json.RawMessage
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}
	if *r == nil {
		*r = make(Responses, len(m))
	}
	rv := *r
	for k, j := range m {
		v, err := unmarshalResponse(j)
		if err != nil {
			return err
		}
		rv[k] = v
	}
	return nil
}

// UnmarshalYAML unmarshals YAML data into r
func (r *Responses) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return yamlutil.Unmarshal(unmarshal, r)
}

// MarshalYAML marshals r into YAML
func (r Responses) MarshalYAML() (interface{}, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	var v interface{}
	err = json.Unmarshal(b, &v)
	return v, err
}

// ResponseObj describes a single response from an API Operation, including
// design-time, static links to operations based on the response.
type ResponseObj struct {
	// A description of the response. CommonMark syntax MAY be used for rich
	// text representation.
	//
	// *required*
	Description string `json:"description,omitempty"`
	// Maps a header name to its definition. RFC7230 states header names are
	// case insensitive. If a response header is defined with the name
	// "Content-Type", it SHALL be ignored.
	Headers Headers `json:"headers,omitempty"`
	// A map containing descriptions of potential response payloads. The key is
	// a media type or media type range and the value describes it. For
	// responses that match multiple keys, only the most specific key is
	// applicable. e.g. text/plain overrides text/*
	Content Content `json:"content,omitempty"`
	// A map of operations links that can be followed from the response. The key
	// of the map is a short name for the link, following the naming constraints
	// of the names for Component Objects.
	Links      Links `json:"links,omitempty"`
	Extensions `json:"-"`
}

type response ResponseObj

// Kind returns KindResponse
func (*ResponseObj) Kind() Kind {
	return KindResponse
}

// ResolveResponse resolves ResponseObj by returning itself. resolve is  not called.
func (r *ResponseObj) ResolveResponse(func(ref string) (*ResponseObj, error)) (*ResponseObj, error) {
	return r, nil
}

// MarshalJSON marshals r into JSON
func (r ResponseObj) MarshalJSON() ([]byte, error) {
	return marshalExtendedJSON(response(r))
}

// UnmarshalJSON unmarshals json into r
func (r *ResponseObj) UnmarshalJSON(data []byte) error {
	var v response
	err := unmarshalExtendedJSON(data, &v)
	*r = ResponseObj(v)
	return err
}

// MarshalYAML first marshals and unmarshals into JSON and then marshals into
// YAML
func (r ResponseObj) MarshalYAML() (interface{}, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	var v interface{}
	err = json.Unmarshal(b, &v)
	return v, err
}

// UnmarshalYAML unmarshals yaml into s
func (r *ResponseObj) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return yamlutil.Unmarshal(unmarshal, r)
}

func unmarshalResponse(data []byte) (Response, error) {
	if isRefJSON(data) {
		return unmarshalReferenceJSON(data)
	}
	var v ResponseObj
	err := json.Unmarshal(data, &v)
	return &v, err
}

// ResolvedResponses is a container for the expected responses of an operation. The
// container maps a HTTP response code to the expected response.
//
// The documentation is not necessarily expected to cover all possible HTTP
// response codes because they may not be known in advance. However,
// documentation is expected to cover a successful operation response and any
// known errors.
//
// The default MAY be used as a default response object for all HTTP codes that
// are not covered individually by the Responses Object.
//
// The Responses Object MUST contain at least one response code, and if only one
// response code is provided it SHOULD be the response for a successful
// operation call.
type ResolvedResponses map[string]*ResolvedResponse

// Kind returns KindResolvedResponses
func (ResolvedResponses) Kind() Kind {
	return KindResolvedResponses
}

// ResolvedResponse describes a single response from an API Operation, including
// design-time, static links to operations based on the response.
type ResolvedResponse struct {
	// A description of the response. CommonMark syntax MAY be used for rich
	// text representation.
	//
	// *required*
	Description string `json:"description,omitempty"`
	// Maps a header name to its definition. RFC7230 states header names are
	// case insensitive. If a response header is defined with the name
	// "Content-Type", it SHALL be ignored.
	Headers ResolvedHeaders `json:"headers,omitempty"`
	// A map containing descriptions of potential response payloads. The key is
	// a media type or media type range and the value describes it. For
	// responses that match multiple keys, only the most specific key is
	// applicable. e.g. text/plain overrides text/*
	Content Content `json:"content,omitempty"`
	// A map of operations links that can be followed from the response. The key
	// of the map is a short name for the link, following the naming constraints
	// of the names for Component Objects.
	Links      Links `json:"links,omitempty"`
	Extensions `json:"-"`
}

// Kind returns KindResolvedResponse
func (*ResolvedResponse) Kind() Kind {
	return KindResolvedResponse
}

var (
	_ Node = (*ResolvedResponses)(nil)
	_ Node = (*ResolvedResponse)(nil)
	_ Node = (*ResponseObj)(nil)
	_ Node = (Responses)(nil)
)
