package openapi

import (
	"encoding/json"
	"strconv"

	"github.com/chanced/transcodefmt"
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

func (c *Component[T]) setLocation(loc Location) error {
	if c == nil {
		return nil
	}
	if c.Reference != nil {
		return c.Reference.setLocation(loc)
	} else if (any)(c.Object) != nil {
		return c.Object.setLocation(loc)
	}
	return nil
}

// ComponentSlice is a slice of Components of type T
type ComponentSlice[T node] []Component[T]

func (cs ComponentSlice[T]) setLocation(loc Location) error {
	for i, c := range cs {
		if err := c.setLocation(loc.Append(strconv.Itoa(i))); err != nil {
			return err
		}
	}
	return nil
}

// ComponentEntry is an entry in a ComponentMap consisting of a Key/Value pair for
// an object consiting of Component[T]s
type ComponentEntry[V node] struct {
	Key       Text
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
			Key:       Text(key.String()),
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
		b, err = sjson.SetBytes(b, field.Key.String(), field.Component)
		_ = b
		if err != nil {
			return nil, err
		}
	}
	return b, nil
}

func (cm *ComponentMap[T]) Get(key Text) (Component[T], bool) {
	for _, v := range *cm {
		if v.Key == key {
			return v.Component, true
		}
	}
	return Component[T]{}, false
}

// Set sets the value of the key in the ComponentMap
func (cm *ComponentMap[T]) Set(key Text, value Component[T]) {
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

func (cm *ComponentMap[T]) Del(key Text) {
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
	return transcodefmt.JSONToYAML(j)
}

// UnmarshalYAML implements yaml.Unmarshaler
func (cm *ComponentMap[T]) UnmarshalYAML(value *yaml.Node) error {
	j, err := transcodefmt.YAMLToJSON([]byte(value.Value))
	if err != nil {
		return err
	}
	return json.Unmarshal(j, cm)
}

func (cm ComponentMap[T]) setLocation(loc Location) error {
	for _, kv := range cm {
		if err := kv.Component.setLocation(loc); err != nil {
			return err
		}
	}
	return nil
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

	Extensions `json:"-"`
	Location   *Location `json:"-"`
}

func (c *Components) setLocation(loc Location) error {
	if c == nil {
		return nil
	}
	c.Location = &loc
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
