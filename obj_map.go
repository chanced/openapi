package openapi

import (
	"bytes"
	"encoding/json"
	"reflect"

	"github.com/chanced/jsonpointer"
	"github.com/chanced/jsonx"
	"github.com/tidwall/gjson"
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
func (*ObjMap[T]) mapKind() Kind   { return KindUndefined }
func (*ObjMap[T]) sliceKind() Kind { return KindUndefined }

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

func (om *ObjMap[T]) edges() []node {
	if om == nil {
		return nil
	}
	edges := make([]node, 0, len(om.Items))
	for _, item := range om.Items {
		edges = appendEdges(edges, item.Value)
	}
	return edges
}

func (om *ObjMap[T]) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
	if err := ptr.Validate(); err != nil {
		return nil, err
	}
	return om.resolveNodeByPointer(ptr)
}

func (om *ObjMap[T]) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
	if ptr.IsRoot() {
		return om, nil
	}
	tok, _ := ptr.NextToken()
	v := om.Get(Text(tok))
	if v.isNil() {
		return nil, newErrNotFound(om.Location.AbsolutePath(), tok)
	}
	return nil, nil
}

func (om ObjMap[T]) setLocation(loc Location) error {
	for _, kv := range om.Items {
		if err := kv.Value.setLocation(loc); err != nil {
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
	for i, kv := range om.Items {
		if kv.Key == key {
			om.Items[i] = ObjMapEntry[T]{
				Location: om.Location.Append(key.String()),
				Key:      key,
				Value:    obj,
			}
			return
		}
	}
	om.Items = append(om.Items, ObjMapEntry[T]{
		Location: om.Location.Append(key.String()),
		Key:      key,
		Value:    obj,
	})
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
