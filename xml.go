package openapi

import "github.com/chanced/jsonpointer"

// XML is a metadata object that allows for more fine-tuned XML model
// definitions.
//
// When using arrays, XML element names are not inferred (for singular/plural
// forms) and the name property SHOULD be used to add that information. See
// examples for expected behavior.
type XML struct {
	Extensions `json:"-"`
	Location   `json:"-"`

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
	Wrapped bool `json:"wrapped,omitempty"`
}

func (*XML) Anchors() (*Anchors, error) { return nil, nil }

func (x *XML) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
	if err := ptr.Validate(); err != nil {
		return nil, err
	}
	return x.resolveNodeByPointer(ptr)
}

func (x *XML) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
	if ptr.IsRoot() {
		return x, nil
	}
	tok, _ := ptr.NextToken()
	return nil, newErrNotResolvable(x.Location.AbsolutePath(), tok)
}

func (*XML) Kind() Kind      { return KindXML }
func (*XML) mapKind() Kind   { return KindUndefined }
func (*XML) sliceKind() Kind { return KindUndefined }

func (x XML) MarshalJSON() ([]byte, error) {
	type xml XML
	return marshalExtendedJSON(xml(x))
}

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
	xml.Location = loc
	return nil
}

func (xml *XML) isNil() bool { return xml == nil }
func (xml *XML) Refs() []Ref { return nil }
func (xml *XML) Edges() []Node {
	if xml == nil {
		return nil
	}
	return downcastNodes(xml.edges())
}
func (xml *XML) edges() []node { return nil }

var _ node = (*XML)(nil) // _ Walker = (*XML)(nil)
