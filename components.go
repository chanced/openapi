package openapi

import "github.com/chanced/openapi/yamlutil"

type components Components

// Components holds a set of reusable objects for different aspects of the OAS.
// All objects defined within the components object will have no effect on the
// API unless they are explicitly referenced from properties outside the
// components object.
type Components struct {
	// An object to hold reusable Schema Objects.
	Schemas Schemas `json:"schemas,omitempty"`
	// An object to hold reusable Response Objects.
	Responses Responses `json:"responses,omitempty"`
	// An object to hold reusable Parameter Objects.
	Parameters Parameters `json:"parameters,omitempty"`
	// An object to hold reusable Example Objects.
	Examples Examples `json:"examples,omitempty"`
	// An object to hold reusable Request Body Objects.
	RequestBodies RequestBodies `json:"requestBodies,omitempty"`
	// An object to hold reusable Header Objects.
	Headers Headers `json:"headers,omitempty"`
	// An object to hold reusable Security Scheme Objects.
	SecuritySchemes SecuritySchemes `json:"securitySchemes,omitempty"`
	// An object to hold reusable Link Objects.
	Links *Links `json:"links,omitempty"`
	// An object to hold reusable Callback Objects.
	Callbacks Callbacks `json:"callbacks,omitempty"`
	// An object to hold reusable Path Item Object.
	PathItems  PathItems `json:"pathItems,omitempty"`
	Extensions `json:"-"`
}

func (c Components) Nodes() Nodes {
	return makeNodes(nodes{
		{"schemas", c.Schemas, KindSchemas},
		{"responses", c.Responses, KindResponses},
		{"parameters", c.Parameters, KindParameters},
		{"examples", c.Examples, KindExamples},
		{"requestBodies", c.RequestBodies, KindRequestBodies},
		{"headers", c.Headers, KindHeaders},
		{"securitySchemes", c.SecuritySchemes, KindSecuritySchemes},
		{"links", c.Links, KindLinks},
		{"callbacks", c.Callbacks, KindCallbacks},
		{"pathItems", c.PathItems, KindPathItems},
	})
}

// Kind returns KindComponents
func (*Components) Kind() Kind {
	return KindComponents
}

// MarshalJSON marshals JSON
func (c Components) MarshalJSON() ([]byte, error) {
	return marshalExtendedJSON(components(c))
}

// UnmarshalJSON unmarshals JSON
func (c *Components) UnmarshalJSON(data []byte) error {
	var v components
	if err := unmarshalExtendedJSON(data, &v); err != nil {
		return err
	}
	*c = Components(v)
	return nil
}

// MarshalYAML marshals YAML
func (c Components) MarshalYAML() (interface{}, error) {
	return yamlutil.Marshal(c)
}

// UnmarshalYAML unmarshals YAML
func (c *Components) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return yamlutil.Unmarshal(unmarshal, c)
}

// ResolvedComponents holds a set of resolved reusable objects for different
// aspects of the OAS. All objects defined within the components object will
// have no effect on the API unless they are explicitly referenced from
// properties outside the components object.
type ResolvedComponents struct {
	// An object to hold reusable Schema Objects.
	Schemas ResolvedSchemas `json:"schemas,omitempty"`
	// An object to hold reusable Response Objects.
	Responses ResolvedResponses `json:"responses,omitempty"`
	// An object to hold reusable Parameter Objects.
	Parameters ResolvedParameters `json:"parameters,omitempty"`
	// An object to hold reusable Example Objects.
	Examples ResolvedExamples `json:"examples,omitempty"`
	// An object to hold reusable Request Body Objects.
	RequestBodies ResolvedRequestBodies `json:"requestBodies,omitempty"`
	// An object to hold reusable Header Objects.
	Headers ResolvedHeaders `json:"headers,omitempty"`
	// An object to hold reusable Security Scheme Objects.
	SecuritySchemes ResolvedSecuritySchemes `json:"securitySchemes,omitempty"`
	// An object to hold reusable Link Objects.
	Links ResolvedLinks `json:"links,omitempty"`
	// An object to hold reusable Callback Objects.
	Callbacks ResolvedCallbacks `json:"callbacks,omitempty"`
	// An object to hold reusable Path Item Object.
	PathItems  ResolvedPathItems `json:"pathItems,omitempty"`
	Extensions `json:"-"`
}

func (rc *ResolvedComponents) Nodes() Nodes {
	return makeNodes(nodes{
		{"schemas", rc.Schemas, KindResolvedSchemas},
		{"responses", rc.Responses, KindResolvedResponses},
		{"parameters", rc.Parameters, KindResolvedParameters},
		{"examples", rc.Examples, KindResolvedExamples},
		{"requestBodies", rc.RequestBodies, KindResolvedRequestBodies},
		{"headers", rc.Headers, KindResolvedHeaders},
		{"securitySchemes", rc.SecuritySchemes, KindResolvedSecuritySchemes},
		{"links", rc.Links, KindResolvedLinks},
		{"callbacks", rc.Callbacks, KindResolvedCallbacks},
		{"pathItems", rc.PathItems, KindResolvedPathItems},
	})
}

// Kind returns KindResolvedComponents
func (*ResolvedComponents) Kind() Kind {
	return KindResolvedComponents
}

var (
	_ Node = (*Components)(nil)
	_ Node = (*ResolvedComponents)(nil)
)
