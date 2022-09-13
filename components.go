package openapi

// ComponentMap is a pseudo map consisting of Components with type T.

func newComponent[T node](ref *Reference, obj T) Component[T] {
	return Component[T]{
		Reference: ref,
		Object:    obj,
	}
}

// Components holds a set of reusable objects for different aspects of the OAS.
// All objects defined within the components object will have no effect on the
// API unless they are explicitly referenced from properties outside the
// components object.
type Components struct {
	// OpenAPI extensions
	Extensions `json:"-"`
	Location   `json:"-"`

	Schemas         *SchemaMap         `json:"schemas,omitempty"`
	Responses       *ResponseMap       `json:"responses,omitempty"`
	Parameters      *ParameterMap      `json:"parameters,omitempty"`
	Examples        *ExampleMap        `json:"examples,omitempty"`
	RequestBodies   *RequestBodyMap    `json:"requestBodies,omitempty"`
	Headers         *HeaderMap         `json:"headers,omitempty"`
	SecuritySchemes *SecuritySchemeMap `json:"securitySchemes,omitempty"`
	Links           *LinkMap           `json:"links,omitempty"`
	Callbacks       *CallbacksMap      `json:"callbacks,omitempty"`
	PathItems       *PathItemMap       `json:"pathItems,omitempty"`
}

func (c *Components) Anchors() (*Anchors, error) {
	if c == nil {
		return nil, nil
	}
	var err error
	var anchors *Anchors
	if anchors, err = anchors.merge(c.Schemas.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(c.Responses.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(c.Parameters.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(c.Examples.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(c.RequestBodies.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(c.Headers.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(c.SecuritySchemes.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(c.Links.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(c.Callbacks.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(c.PathItems.Anchors()); err != nil {
		return nil, err
	}

	return anchors, nil
}

func (c *Components) setLocation(loc Location) error {
	if c == nil {
		return nil
	}
	c.Location = loc
	var err error
	if err = c.Schemas.setLocation(loc.Append("schemas")); err != nil {
		return err
	}
	if err = c.Responses.setLocation(loc.Append("responses")); err != nil {
		return err
	}
	if err = c.Parameters.setLocation(loc.Append("parameters")); err != nil {
		return err
	}
	if err = c.Examples.setLocation(loc.Append("examples")); err != nil {
		return err
	}
	if err = c.RequestBodies.setLocation(loc.Append("requestBodies")); err != nil {
		return err
	}
	if err = c.Headers.setLocation(loc.Append("headers")); err != nil {
		return err
	}
	if err = c.SecuritySchemes.setLocation(loc.Append("securitySchemes")); err != nil {
		return err
	}
	if err = c.Links.setLocation(loc.Append("links")); err != nil {
		return err
	}
	if err = c.Callbacks.setLocation(loc.Append("callbacks")); err != nil {
		return err
	}
	if err = c.PathItems.setLocation(loc.Append("pathItems")); err != nil {
		return err
	}

	return nil
}

// MarshalJSON marshals JSON
func (c Components) MarshalJSON() ([]byte, error) {
	type components Components
	return marshalExtendedJSON(components(c))
}

// UnmarshalJSON unmarshals JSON
func (c *Components) UnmarshalJSON(data []byte) error {
	type components Components
	var v components
	if err := unmarshalExtendedJSON(data, &v); err != nil {
		return err
	}
	*c = Components(v)
	return nil
}
