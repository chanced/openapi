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
	Name string `json:"name,omitempty"`
	// The URI of the namespace definition. This MUST be in the form of an
	// absolute URI.
	Namespace string `json:"namespace,omitempty"`
	// The prefix to be used for the name.
	Prefix string `json:"prefix,omitempty"`
	// Declares whether the property definition translates to an attribute
	// instead of an element. Default value is false.
	Attribute bool `json:"attribute,omitempty"`
	// MAY be used only for an array definition. Signifies whether the array is
	// wrapped (for example, <books><book/><book/></books>) or unwrapped
	// (<book/><book/>). Default value is false. The definition takes effect
	// only when defined alongside type being array (outside the items).
	Wrapped    bool `json:"wrapped,omitempty"`
	Extensions `json:"-"`
}

// Kind returns KindXML
func (*XML) Kind() Kind {
	return KindXML
}

func (x XML) Nodes() Nodes {
	return nil
}

var _ Node = (*XML)(nil)
