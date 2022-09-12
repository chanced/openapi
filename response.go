package openapi

// ResponseMap is a container for the expected responses of an operation. The
// container maps a HTTP response code to the expected response.
//
// The documentation is not necessarily expected to cover all possible HTTP
// response codes because they may not be known in advance. However,
// documentation is expected to cover a successful operation response and any
// known errors.
//
// The default MAY be used as a default response object for all HTTP codes that
// are not covered individually by the ResponseMap Object.
//
// The ResponseMap Object MUST contain at least one response code, and if only one
// response code is provided it SHOULD be the response for a successful
// operation call.
type ResponseMap = ComponentMap[*Response]

// Response describes a single response from an API Operation, including
// design-time, static links to operations based on the response.
type Response struct {
	// A description of the response. CommonMark syntax MAY be used for rich
	// text representation.
	//
	// *required*
	Description Text `json:"description,omitempty"`
	// Maps a header name to its definition. RFC7230 states header names are
	// case insensitive. If a response header is defined with the name
	// "Content-Type", it SHALL be ignored.
	Headers HeaderMap `json:"headers,omitempty"`
	// A map containing descriptions of potential response payloads. The key is
	// a media type or media type range and the value describes it. For
	// responses that match multiple keys, only the most specific key is
	// applicable. e.g. text/plain overrides text/*
	Content ContentMap `json:"content,omitempty"`
	// A map of operations links that can be followed from the response. The key
	// of the map is a short name for the link, following the naming constraints
	// of the names for Component Objects.
	Links      LinkMap `json:"links,omitempty"`
	Extensions `json:"-"`

	Location *Location `json:"-"`
}

// MarshalJSON marshals r into JSON
func (r Response) MarshalJSON() ([]byte, error) {
	type response Response
	return marshalExtendedJSON(response(r))
}

// UnmarshalJSON unmarshals json into r
func (r *Response) UnmarshalJSON(data []byte) error {
	type response Response
	var v response
	err := unmarshalExtendedJSON(data, &v)
	*r = Response(v)
	return err
}

func (*Response) Kind() Kind      { return KindResponse }
func (*Response) mapKind() Kind   { return KindResponseMap }
func (*Response) sliceKind() Kind { return KindUndefined }

// setLocation implements node
func (r *Response) setLocation(loc Location) error {
	if r == nil {
		return nil
	}
	r.Location = &loc
	if err := r.Headers.setLocation(loc.Append("headers")); err != nil {
		return err
	}
	if err := r.Content.setLocation(loc.Append("content")); err != nil {
		return err
	}
	if err := r.Links.setLocation(loc.Append("links")); err != nil {
		return err
	}
	return nil
}

var _ node = (*Response)(nil)
