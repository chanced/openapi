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

type SchemaItem struct {
	Key    Text
	Schema *Schema
}

func (si *SchemaItem) Clone() SchemaItem {
	if si == nil {
		return SchemaItem{}
	}
	return SchemaItem{
		Key:    si.Key,
		Schema: si.Schema.Clone(),
	}
}

// SchemaMap is a pseudo, ordered map ofASew3 Schemas
//
// Under the hood, SchemaMap is a slice of SchemaEntry
type SchemaMap struct {
	Location
	Items []SchemaItem
}

func (sm *SchemaMap) Nodes() []Node {
	if sm == nil {
		return nil
	}
	return downcastNodes(sm.nodes())
}

func (sm *SchemaMap) nodes() []node {
	if sm == nil {
		return nil
	}
	var edges []node
	for _, e := range sm.Items {
		edges = append(edges, e.Schema)
	}
	return edges
}

func (sm *SchemaMap) Refs() []Ref {
	if sm == nil {
		return nil
	}
	var refs []Ref
	for _, e := range sm.Items {
		refs = append(refs, e.Schema.Refs()...)
	}
	return refs
}

func (*SchemaMap) Kind() Kind      { return KindSchemaMap }
func (*SchemaMap) sliceKind() Kind { return KindUndefined }
func (*SchemaMap) mapKind() Kind   { return KindUndefined }
func (sm *SchemaMap) isNil() bool  { return sm == nil }
func (sm *SchemaMap) Anchors() (*Anchors, error) {
	if sm == nil {
		return nil, nil
	}
	var anchors *Anchors
	var err error
	for _, e := range sm.Items {
		if anchors, err = e.Schema.Anchors(); err != nil {
			return nil, err
		}
	}
	return anchors, nil
}

// func (sm *SchemaMap) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if err := ptr.Validate(); err != nil {
// 		return nil, err
// 	}
// 	return sm.resolveNodeByPointer(ptr)
// }

// func (sm *SchemaMap) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if ptr.IsRoot() {
// 		return sm, nil
// 	}
// 	tok, _ := ptr.NextToken()
// 	v := sm.Get(Text(tok))
// 	if v == nil {
// 		return nil, newErrNotFound(sm.AbsoluteLocation(), tok)
// 	}
// 	return v.resolveNodeByPointer(ptr)
// }

func (sm *SchemaMap) Set(key Text, s *Schema) {
	se := SchemaItem{
		Key:    key,
		Schema: s,
	}
	for i, v := range sm.Items {
		if v.Key == key {
			sm.Items[i] = se
			return
		}
	}
	sm.Items = append(sm.Items, se)
}

func (sm *SchemaMap) setLocation(loc Location) error {
	if sm == nil {
		return nil
	}
	sm.Location = loc

	for _, e := range sm.Items {
		err := e.Schema.setLocation(loc.AppendLocation(e.Key.String()))
		if err != nil {
			return err
		}
	}
	return nil
}

func (sm SchemaMap) Get(key Text) *Schema {
	for _, v := range sm.Items {
		if v.Key == key {
			return v.Schema
		}
	}
	return nil
}

func (sm SchemaMap) MarshalJSON() ([]byte, error) {
	b := bytes.Buffer{}
	b.WriteByte('{')
	var err error
	var s []byte
	for _, v := range sm.Items {
		if b.Len() > 1 {
			b.WriteByte(',')
		}
		jsonx.EncodeAndWriteString(&b, v.Key.String())
		b.WriteByte(':')
		s, err = v.Schema.MarshalJSON()
		if err != nil {
			return nil, err
		}
		b.Write(s)
	}
	b.WriteByte('}')
	return b.Bytes(), nil
}

func (sm *SchemaMap) UnmarshalJSON(data []byte) error {
	t := jsonx.TypeOf(data)
	if t != jsonx.TypeObject {
		return &json.UnmarshalTypeError{Value: t.String(), Type: reflect.TypeOf(sm)}
	}
	*sm = SchemaMap{}
	var err error
	gjson.ParseBytes(data).ForEach(func(key, value gjson.Result) bool {
		var s Schema
		err = json.Unmarshal([]byte(value.Raw), &s)
		sm.Items = append(sm.Items, SchemaItem{Key: Text(key.String()), Schema: &s})
		return err == nil
	})
	return err
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Marshaler interface
func (sm SchemaMap) MarshalYAML() (interface{}, error) {
	j, err := sm.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return transcode.YAMLFromJSON(j)
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Unmarshaler interface
func (sm *SchemaMap) UnmarshalYAML(value *yaml.Node) error {
	j, err := transcode.YAMLFromJSON([]byte(value.Value))
	if err != nil {
		return err
	}
	return json.Unmarshal(j, sm)
}

func (sm *SchemaMap) Clone() *SchemaMap {
	if sm == nil {
		return nil
	}
	m := make([]SchemaItem, len(sm.Items))
	for i, v := range sm.Items {
		m[i] = v.Clone()
	}
	return &SchemaMap{
		Location: Location{
			absolute: *sm.absolute.Clone(),
			relative: sm.relative,
		},
		Items: m,
	}
}

var _ node = (*SchemaMap)(nil)
