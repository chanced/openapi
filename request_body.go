package openapi

import (
	"encoding/json"

	"github.com/chanced/openapi/yamlutil"
)

// RequestBodies is a map of RequestBody
type RequestBodies map[string]RequestBody

// Kind returns KindRequestBodies
func (RequestBodies) Kind() Kind {
	return KindRequestBodies
}

// UnmarshalJSON unmarshals JSON
func (rbs *RequestBodies) UnmarshalJSON(data []byte) error {
	var dm map[string]json.RawMessage
	if err := json.Unmarshal(data, &dm); err != nil {
		return err
	}
	rv := make(map[string]RequestBody, len(dm))
	for k, d := range dm {
		var v RequestBody
		if err := unmarshalRequestBody(d, &v); err != nil {
			return err
		}
		rv[k] = v
	}
	*rbs = rv
	return nil
}

func (rbs *RequestBodies) Get(key string) (RequestBody, bool) {
	if rbs == nil || *rbs == nil {
		return nil, false
	}
	v, ok := (*rbs)[key]
	return v, ok
}

func (rbs *RequestBodies) Set(key string, val RequestBody) {
	if *rbs == nil {
		*rbs = RequestBodies{
			key: val,
		}
		return
	}
	(*rbs)[key] = val
}

func (rbs RequestBodies) Nodes() Nodes {
	if len(rbs) == 0 {
		return nil
	}
	n := make(Nodes, len(rbs))
	for k, v := range rbs {
		n.maybeAdd(k, v, KindRequestBody)
	}
	return n
}

func (rbs *RequestBodies) Len() int {
	if rbs == nil || *rbs == nil {
		return 0
	}
	return len(*rbs)
}

// RequestBody can either be a RequestBody or a Reference
type RequestBody interface {
	ResolveRequestBody(func(ref string) (*RequestBodyObj, error)) (*RequestBodyObj, error)
	Node
}

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

func (rb *RequestBodyObj) Nodes() Nodes {
	return makeNodes(nodes{{"content", rb.Content, KindContent}})
}

type requestbody RequestBodyObj

// Kind returns KindRequestBody
func (*RequestBodyObj) Kind() Kind {
	return KindRequestBody
}

// ResolveRequestBody resolves RequestBodyObj by returning itself. resolve is  not called.
func (rb *RequestBodyObj) ResolveRequestBody(func(ref string) (*RequestBodyObj, error)) (*RequestBodyObj, error) {
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

func (rrb *ResolvedRequestBody) Nodes() Nodes {
	return makeNodes(nodes{{"content", rrb.Content, KindResolvedContent}})
}

// Kind returns KindResolvedRequestBody
func (*ResolvedRequestBody) Kind() Kind {
	return KindResolvedRequestBody
}

// ResolvedRequestBodies is a map of *ResolvedRequestBody
type ResolvedRequestBodies map[string]*ResolvedRequestBody

func (rrbs *ResolvedRequestBodies) Get(key string) (*ResolvedRequestBody, bool) {
	if rrbs == nil || *rrbs == nil {
		return nil, false
	}
	v, ok := (*rrbs)[key]
	return v, ok
}

func (rrbs *ResolvedRequestBodies) Set(key string, val *ResolvedRequestBody) {
	if *rrbs == nil {
		*rrbs = ResolvedRequestBodies{
			key: val,
		}
		return
	}
	(*rrbs)[key] = val
}

func (rrbs ResolvedRequestBodies) Nodes() Nodes {
	if len(rrbs) == 0 {
		return nil
	}
	n := make(Nodes, len(rrbs))
	for k, v := range rrbs {
		n.maybeAdd(k, v, KindResolvedRequestBody)
	}
	return n
}

func (rrbs *ResolvedRequestBodies) Len() int {
	if rrbs == nil || *rrbs == nil {
		return 0
	}
	return len(*rrbs)
}

// Kind returns KindResolvedRequestBodies
func (ResolvedRequestBodies) Kind() Kind {
	return KindResolvedRequestBodies
}

var (
	_ Node = (*RequestBodyObj)(nil)
	_ Node = (RequestBodies)(nil)
	_ Node = (*ResolvedRequestBody)(nil)
	_ Node = (ResolvedRequestBodies)(nil)
)
