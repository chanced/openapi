package openapi

// XML is a metadata object that allows for more fine-tuned XML model
// definitions.
//
// When using arrays, XML element names are not inferred (for singular/plural
// forms) and the name property SHOULD be used to add that information. See
// examples for expected behavior.
type XML struct {
	// Replaces the name of the element/attribute used for the described schema
	// property. When defined within items, it will affect the name of the
	// individual XML elements within the list. When defined alongside type
	// being array (outside the items), it will affect the wrapping element and
	// only if wrapped is true. If wrapped is false, it will be ignored.
	Name Text `json:"name,omitempty"`
	// The URI of the namespace definition. This MUST be in the form of an
	// absolute URI.
	Namespace Text `json:"namespace,omitempty"`
	// The prefix to be used for the name.
	Prefix Text `json:"prefix,omitempty"`
	// Declares whether the property definition translates to an attribute
	// instead of an element. Default value is false.
	Attribute bool `json:"attribute,omitempty"`
	// MAY be used only for an array definition. Signifies whether the array is
	// wrapped (for example, <books><book/><book/></books>) or unwrapped
	// (<book/><book/>). Default value is false. The definition takes effect
	// only when defined alongside type being array (outside the items).
	Wrapped    bool `json:"wrapped,omitempty"`
	Extensions `json:"-"`
	Location   *Location `json:"-"`
}

func (*XML) Kind() Kind      { return KindXML }
func (*XML) mapKind() Kind   { return KindUndefined }
func (*XML) sliceKind() Kind { return KindUndefined }

// MarshalJSON implements node
func (x XML) MarshalJSON() ([]byte, error) {
	type xml XML
	return marshalExtendedJSON(xml(x))
}

// UnmarshalJSON implements node
func (x *XML) UnmarshalJSON(data []byte) error {
	type xml XML
	var v xml
	err := unmarshalExtendedJSON(data, &v)
	if err != nil {
		return err
	}
	*x = XML(v)
	return nil
}

func (xml *XML) setLocation(loc Location) error {
	if xml == nil {
		return nil
	}
	xml.Location = &loc
	return nil
}

var _ node = (*XML)(nil)
