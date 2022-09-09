package openapi

import (
	"encoding/json"
)

type (
	// Slice is a slice of Components of type T
	Slice[T Node] []Component[T]
	// Map is a map of Components of type T
	Map[T Node] map[string]Component[T]
)

type Component[T Node] struct {
	Ref    *Reference
	Object T
}

func newComponent[T Node](ref *Reference, obj T) Component[T] {
	return Component[T]{
		Ref:    ref,
		Object: obj,
	}
}

func (c Component[T]) MarshalJSON() ([]byte, error) {
	if c.Ref != nil {
		return json.Marshal(c.Ref)
	}
	if any(c.Object) != nil {
		return c.Object.MarshalJSON()
	}
	return nil, nil
}

func (c *Component[T]) UnmarshalJSON(data []byte) error {
	if isRefJSON(data) {
		var ref Reference
		if err := json.Unmarshal(data, &ref); err != nil {
			return err
		}
		*c = Component[T]{
			Ref: &ref,
		}
		return nil
	}
	var value T
	if err := value.UnmarshalJSON(data); err != nil {
		return err
	}
	*c = Component[T]{
		Object: value,
	}
	return nil
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
