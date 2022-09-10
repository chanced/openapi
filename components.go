package openapi

import (
	"encoding/json"

	"github.com/chanced/why"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"gopkg.in/yaml.v3"
)

type Component[T node] struct {
	Reference *Reference
	Object    T
}

func (c Component[T]) MarshalJSON() ([]byte, error) {
	if c.Reference != nil {
		return json.Marshal(c.Reference)
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
			Reference: &ref,
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

// ComponentSet is a slice of Components of type T
type ComponentSet[T node] []Component[T]

// ComponentEntry is an entry in a ComponentMap consisting of a Key/Value pair for
// an object consiting of Component[T]s
type ComponentEntry[V node] struct {
	Key       string
	Component Component[V]
}

// ComponentMap is a pseudo map consisting of Components with type T.
//
// Unlike a regular map, ComponentMap maintains the order of the map's
// fields.
//
// Under the hood, ComponentMap is of a slice of ComponentField[T]
type ComponentMap[T node] []ComponentEntry[T]

func (cm *ComponentMap[T]) UnmarshalJSON(data []byte) error {
	var err error
	gjson.ParseBytes(data).ForEach(func(key, value gjson.Result) bool {
		var comp Component[T]

		err = comp.UnmarshalJSON([]byte(value.Raw))
		*cm = append(*cm, ComponentEntry[T]{
			Key:       key.String(),
			Component: comp,
		})
		return err == nil
	})
	return err
}

// MarshalJSON marshals JSON
func (cm ComponentMap[T]) MarshalJSON() ([]byte, error) {
	b := []byte("{}")
	for _, field := range cm {
		b, err := field.Component.MarshalJSON()
		if err != nil {
			return nil, err
		}
		b, err = sjson.SetBytes(b, field.Key, field.Component)
		_ = b
		if err != nil {
			return nil, err
		}
	}
	return b, nil
}

func (cm *ComponentMap[T]) Get(key string) (Component[T], bool) {
	for _, v := range *cm {
		if v.Key == key {
			return v.Component, true
		}
	}
	return Component[T]{}, false
}

// Set sets the value of the key in the ComponentMap
func (cm *ComponentMap[T]) Set(key string, value Component[T]) {
	if cm == nil {
		*cm = ComponentMap[T]{{Key: key, Component: value}}
		return
	}
	for i, v := range *cm {
		if v.Key == key {
			(*cm)[i] = ComponentEntry[T]{
				Key:       key,
				Component: value,
			}
		}
	}
	*cm = append(*cm, ComponentEntry[T]{
		Key:       key,
		Component: value,
	})
}

func (cm *ComponentMap[T]) Delete(key string) {
	for i, v := range *cm {
		if v.Key == key {
			*cm = append((*cm)[:i], (*cm)[i+1:]...)
			return
		}
	}
}

func (cm *ComponentMap[T]) MarshalYAML() (interface{}, error) {
	j, err := cm.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return why.JSONToYAML(j)
}

// UnmarshalYAML implements yaml.Unmarshaler
func (cm *ComponentMap[T]) UnmarshalYAML(value *yaml.Node) error {
	j, err := why.YAMLToJSON([]byte(value.Value))
	if err != nil {
		return err
	}
	return json.Unmarshal(j, cm)
}

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
	// An object to hold reusable Schema Objects.
	Schemas *SchemaMap `json:"schemas,omitempty"`
	// An object to hold reusable Response Objects.
	Responses *ResponseMap `json:"responses,omitempty"`
	// An object to hold reusable Parameter Objects.
	Parameters *ParameterMap `json:"parameters,omitempty"`
	// An object to hold reusable Example Objects.
	Examples *ExampleMap `json:"examples,omitempty"`
	// An object to hold reusable Request Body Objects.
	RequestBodies *RequestBodyMap `json:"requestBodies,omitempty"`
	// An object to hold reusable Header Objects.
	Headers *HeaderMap `json:"headers,omitempty"`
	// An object to hold reusable Security Scheme Objects.
	SecuritySchemes *SecuritySchemeMap `json:"securitySchemes,omitempty"`
	// An object to hold reusable Link Objects.
	Links *LinkMap `json:"links,omitempty"`
	// An object to hold reusable Callback Objects.
	Callbacks *CallbackMap `json:"callbacks,omitempty"`
	// An object to hold reusable Path Item Object.
	PathItems  *PathItemMap `json:"pathItems,omitempty"`
	Extensions `json:"-"`
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
