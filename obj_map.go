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

type ObjMapEntry[T node] struct {
	Location
	Key   Text
	Value T
}

// ObjMap is a map of OpenAPI Objects of type T
type ObjMap[T node] struct {
	Location
	Items []ObjMapEntry[T]
}

func (*ObjMap[T]) Kind() Kind {
	var t T
	return t.Kind()
}
func (*ObjMap[T]) mapKind() Kind { return KindUndefined }
func (*ObjMap[T]) sliceKind() Kind {
	var t T
	return objSliceKind(t)
}

func (om *ObjMap[T]) Refs() []Ref {
	if om == nil {
		return nil
	}
	refs := []Ref{}
	for _, item := range om.Items {
		refs = append(refs, item.Value.Refs()...)
	}
	return refs
}

func (om *ObjMap[T]) nodes() []node {
	if om == nil {
		return nil
	}
	edges := make([]node, 0, len(om.Items))
	for _, item := range om.Items {
		edges = appendEdges(edges, item.Value)
	}
	return edges
}

func (om *ObjMap[T]) setLocation(loc Location) error {
	if om == nil {
		return nil
	}
	om.Location = loc
	for _, kv := range om.Items {
		if err := kv.Value.setLocation(loc.AppendLocation(string(kv.Key))); err != nil {
			return err
		}
	}
	return nil
}

func (om *ObjMap[T]) Get(key Text) T {
	var t T
	for _, kv := range om.Items {
		if kv.Key == key {
			t = kv.Value
			break
		}
	}
	return t
}

func (om *ObjMap[T]) Set(key Text, obj T) {
	if om == nil || om.Items == nil {
		*om = ObjMap[T]{
			Items: []ObjMapEntry[T]{},
		}
	}
	for i, kv := range om.Items {
		if kv.Key == key {
			om.Items[i] = ObjMapEntry[T]{
				Location: om.AppendLocation(key.String()),
				Key:      key,
				Value:    obj,
			}
			return
		}
	}
	om.Items = append(om.Items, ObjMapEntry[T]{
		Location: om.AppendLocation(key.String()),
		Key:      key,
		Value:    obj,
	})
}

func (om *ObjMap[T]) Del(key Text) {
}

func (om *ObjMap[T]) UnmarshalJSON(data []byte) error {
	var t T
	var m ObjMap[T]
	*om = m

	if !jsonx.IsObject(data) {
		return &json.UnmarshalTypeError{
			Value:  jsonx.TypeOf(data).String(),
			Type:   reflect.TypeOf(t),
			Struct: "PathItemMap",
		}
	}
	var pi T
	var err error
	gjson.ParseBytes(data).ForEach(func(key, value gjson.Result) bool {
		if err = json.Unmarshal([]byte(value.Raw), &pi); err != nil {
			return false
		}
		m.Items = append(m.Items, ObjMapEntry[T]{Key: Text(key.String()), Value: pi})
		return true
	})
	*om = m
	return err
}

func (om *ObjMap[T]) MarshalJSON() ([]byte, error) {
	var err error
	b := bytes.Buffer{}
	var j []byte
	_ = j
	b.WriteByte('{')
	for _, entry := range om.Items {
		if b.Len() > 1 {
			b.WriteByte(',')
		}
		jsonx.EncodeAndWriteString(&b, entry.Key)
		b.WriteByte(':')
		j, err = entry.Value.MarshalJSON()
		if err != nil {
			return nil, err
		}
		b.Write(j)
	}
	b.WriteByte('}')
	return b.Bytes(), err
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Marshaler interface
func (om ObjMap[T]) MarshalYAML() (interface{}, error) {
	j, err := om.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return transcode.YAMLFromJSON(j)
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Unmarshaler interface
func (om *ObjMap[T]) UnmarshalYAML(value *yaml.Node) error {
	j, err := transcode.YAMLFromJSON([]byte(value.Value))
	if err != nil {
		return err
	}
	return json.Unmarshal(j, om)
}

func (om *ObjMap[T]) Anchors() (*Anchors, error) {
	if om == nil {
		return nil, nil
	}
	var anchors *Anchors
	var err error
	for _, e := range om.Items {
		anchors, err = anchors.merge(e.Value.Anchors())
		if err != nil {
			return anchors, err
		}
	}
	//∆
	return anchors, nil
}

func (om *ObjMap[T]) isNil() bool { return om == nil }

var _ (node) = (*ObjMap[*Server])(nil)

// func (om *ObjMap[T]) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if err := ptr.Validate(); err != nil {
// 		return nil, err
// 	}
// 	return om.resolveNodeByPointer(ptr)
// }

// func (om *ObjMap[T]) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if ptr.IsRoot() {
// 		return om, nil
// 	}
// 	tok, _ := ptr.NextToken()
// 	v := om.Get(Text(tok))
// 	if v.isNil() {
// 		return nil, newErrNotFound(om.Location.AbsoluteLocation(), tok)
// 	}
// 	return nil, nil
// }
