package openapi

// RequestBodyMap is a map of RequestBody
type RequestBodyMap = ComponentMap[*RequestBody]

// RequestBody describes a single request body.
type RequestBody struct {
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

func (RequestBody) Kind() Kind { return KindRequestBody }

// MarshalJSON marshals h into JSON
func (rb RequestBody) MarshalJSON() ([]byte, error) {
	type requestbody RequestBody

	return marshalExtendedJSON(requestbody(rb))
}

// UnmarshalJSON unmarshals json into rb
func (rb *RequestBody) UnmarshalJSON(data []byte) error {
	type requestbody RequestBody

	var v requestbody
	err := unmarshalExtendedJSON(data, &v)
	*rb = RequestBody(v)
	return err
}
