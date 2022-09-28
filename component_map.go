package openapi

import (
	"bytes"
	"encoding/json"
	"reflect"

	"github.com/chanced/jsonx"
	"github.com/chanced/transcode"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v3"
)

// ComponentEntry is an entry in a ComponentMap consisting of a Key/Value pair for
// an object consiting of Component[T]s
type ComponentEntry[V refable] struct {
	Key       Text
	Component *Component[V]
}

// ComponentMap is a pseudo map consisting of Components with type T.
//
// Unlike a regular map, ComponentMap maintains the order of the map's
// fields.
//
// Under the hood, ComponentMap is of a slice of ComponentField[T]
type ComponentMap[T refable] struct {
	Location
	Items []*ComponentEntry[T]
}

func (cm *ComponentMap[T]) nodes() []node {
	if cm == nil {
		return nil
	}
	var edges []node
	for _, item := range cm.Items {
		edges = appendEdges(edges, item.Component)
	}
	return edges
}

func (cm *ComponentMap[T]) Refs() []Ref {
	if cm == nil {
		return nil
	}
	var refs []Ref
	for _, item := range cm.Items {
		refs = append(refs, item.Component.Refs()...)
	}
	return refs
}

func (cm ComponentMap[T]) Map() map[Text]*Component[T] {
	m := make(map[Text]*Component[T], len(cm.Items))
	for _, item := range cm.Items {
		m[item.Key] = item.Component
	}
	return m
}

func (*ComponentMap[T]) Kind() Kind {
	var t T
	return t.mapKind()
}
func (*ComponentMap[T]) mapKind() Kind   { return KindUndefined }
func (*ComponentMap[T]) sliceKind() Kind { return KindUndefined }

func (cm *ComponentMap[T]) UnmarshalJSON(data []byte) error {
	if !jsonx.IsObject(data) {
		return &json.UnmarshalTypeError{
			Value:  jsonx.TypeOf(data).String(),
			Type:   reflect.TypeOf(cm),
			Struct: "ComponentMap",
		}
	}
	var err error
	*cm = ComponentMap[T]{
		Items: make([]*ComponentEntry[T], 0),
	}
	gjson.ParseBytes(data).ForEach(func(key, value gjson.Result) bool {
		var comp Component[T]
		err = comp.UnmarshalJSON([]byte(value.Raw))
		cm.Items = append(cm.Items, &ComponentEntry[T]{
			Key:       Text(key.String()),
			Component: &comp,
		})
		return err == nil
	})
	return err
}

func (cm *ComponentMap[T]) isNil() bool {
	return cm == nil
}

// MarshalJSON marshals JSON
func (cm ComponentMap[T]) MarshalJSON() ([]byte, error) {
	b := bytes.Buffer{}
	b.WriteByte('{')
	for _, field := range cm.Items {
		if b.Len() > 1 {
			b.WriteByte(',')
		}
		jsonx.EncodeAndWriteString(&b, field.Key)
		b.WriteByte(':')
		cb, err := field.Component.MarshalJSON()
		if err != nil {
			return nil, err
		}
		b.Write(cb)
	}
	b.WriteByte('}')
	return b.Bytes(), nil
}

func (cm *ComponentMap[T]) Get(key Text) *Component[T] {
	if cm == nil || cm.Items == nil {
		return nil
	}
	for _, v := range cm.Items {
		if v.Key == key {
			return v.Component
		}
	}
	return nil
}

// Set sets the value of the key in the ComponentMap
func (cm *ComponentMap[T]) Set(key Text, value *Component[T]) {
	if cm == nil {
		*cm = ComponentMap[T]{}
	}
	for i, v := range cm.Items {
		if v.Key == key {
			cm.Items[i] = &ComponentEntry[T]{
				Key:       key,
				Component: value,
			}
		}
	}
	cm.Items = append(cm.Items, &ComponentEntry[T]{
		Key:       key,
		Component: value,
	})
}

func (cm *ComponentMap[T]) Del(key Text) {
	for i, v := range cm.Items {
		if v.Key == key {
			cm.Items = append(cm.Items[:i], cm.Items[i+1:]...)
			return
		}
	}
}

func (cm *ComponentMap[T]) MarshalYAML() (interface{}, error) {
	j, err := cm.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return transcode.YAMLFromJSON(j)
}

// UnmarshalYAML implements yaml.Unmarshaler
func (cm *ComponentMap[T]) UnmarshalYAML(value *yaml.Node) error {
	v, err := yaml.Marshal(value)
	if err != nil {
		return err
	}
	j, err := transcode.JSONFromYAML(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(j, cm)
}

func (cm *ComponentMap[T]) setLocation(loc Location) error {
	if cm == nil {
		return nil
	}
	cm.Location = loc
	for _, kv := range cm.Items {
		if err := kv.Component.setLocation(loc.AppendLocation(kv.Key.String())); err != nil {
			return err
		}
	}

	return nil
}

func (cm *ComponentMap[T]) Anchors() (*Anchors, error) {
	if cm == nil {
		return nil, nil
	}
	var anchors *Anchors
	var err error

	for _, kv := range cm.Items {
		if anchors, err = kv.Component.Anchors(); err != nil {
			return nil, err
		}
	}
	return anchors, nil
}

var _ node = (*ComponentMap[*Response])(nil)

// func (cm ComponentMap[T]) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if err := ptr.Validate(); err != nil {
// 		return nil, err
// 	}
// 	return cm.resolveNodeByPointer(ptr)
// }

// func (c *ComponentMap[T]) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if ptr.IsRoot() {
// 		return c, nil
// 	}
// 	nxt, tok, _ := ptr.Next()
// 	n := c.Get(Text(tok))

// 	if nxt.IsRoot() {
// 		if n == nil {
// 			return nil, newErrNotFound(c.AbsoluteLocation(), tok)
// 		}
// 		if n.Reference != nil {
// 			return n.Reference, nil
// 		}
// 		if !n.Object.isNil() {
// 			return n.Object, nil
// 		}
// 		return nil, newErrNotFound(c.Location.AbsoluteLocation(), tok)
// 	}
// 	if n == nil {
// 		return nil, newErrNotFound(c.Location.AbsoluteLocation(), tok)
// 	}
// 	return n.resolveNodeByPointer(nxt)
// }
