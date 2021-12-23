package openapi

import "github.com/chanced/openapi/yamlutil"

// ResolvedComponents holds a set of resolved reusable objects for different
// aspects of the OAS. All objects defined within the components object will
// have no effect on the API unless they are explicitly referenced from
// properties outside the components object.
type ResolvedComponents struct {
	// An object to hold reusable Schema Objects.
	Schemas *Schemas `json:"schemas,omitempty"`
	// An object to hold reusable Response Objects.
	Responses *Responses `json:"responses,omitempty"`
	// An object to hold reusable Parameter Objects.
	Parameters *Parameters `json:"parameters,omitempty"`
	// An object to hold reusable Example Objects.
	Examples *Examples `json:"examples,omitempty"`
	// An object to hold reusable Request Body Objects.
	RequestBodies *RequestBodies `json:"requestBodies,omitempty"`
	// An object to hold reusable Header Objects.
	Headers *Headers `json:"headers,omitempty"`
	// An object to hold reusable Security Scheme Objects.
	SecuritySchemes *SecuritySchemes `json:"securitySchemes,omitempty"`
	// An object to hold reusable Link Objects.
	Links *Links `json:"links,omitempty"`
	// An object to hold reusable Callback Objects.
	Callbacks *Callbacks `json:"callbacks,omitempty"`
	// An object to hold reusable Path Item Object.
	PathItems  *PathItems `json:"pathItems,omitempty"`
	Extensions `json:"-"`
}

// Components holds a set of reusable objects for different aspects of the OAS.
// All objects defined within the components object will have no effect on the
// API unless they are explicitly referenced from properties outside the
// components object.
type Components struct {
	// An object to hold reusable Schema Objects.
	Schemas *Schemas `json:"schemas,omitempty"`
	// An object to hold reusable Response Objects.
	Responses *Responses `json:"responses,omitempty"`
	// An object to hold reusable Parameter Objects.
	Parameters *Parameters `json:"parameters,omitempty"`
	// An object to hold reusable Example Objects.
	Examples *Examples `json:"examples,omitempty"`
	// An object to hold reusable Request Body Objects.
	RequestBodies *RequestBodies `json:"requestBodies,omitempty"`
	// An object to hold reusable Header Objects.
	Headers *Headers `json:"headers,omitempty"`
	// An object to hold reusable Security Scheme Objects.
	SecuritySchemes *SecuritySchemes `json:"securitySchemes,omitempty"`
	// An object to hold reusable Link Objects.
	Links *Links `json:"links,omitempty"`
	// An object to hold reusable Callback Objects.
	Callbacks *Callbacks `json:"callbacks,omitempty"`
	// An object to hold reusable Path Item Object.
	PathItems  *PathItems `json:"pathItems,omitempty"`
	Extensions `json:"-"`
}
type components Components

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

func (c *Components) resolve(resolver Resolver) (*ResolvedComponents, error) {
	panic("not impl")
}
