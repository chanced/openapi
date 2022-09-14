package openapi

import (
	"bytes"
	"encoding/json"
	"reflect"

	"github.com/chanced/jsonx"
	"github.com/tidwall/gjson"
)

type KeyValue[V any] struct {
	Key   Text
	Value V
}

type Map[T any] struct {
	Items []KeyValue[T]
}

func (m Map[T]) Get(key Text) (T, bool) {
	for _, v := range m.Items {
		if v.Key == key {
			return v.Value, true
		}
	}
	var t T
	return t, false
}

func (m Map[T]) Has(key Text) bool {
	for _, v := range m.Items {
		if v.Key == key {
			return true
		}
	}
	return false
}

func (m *Map[T]) Set(key Text, value T) {
	if m == nil {
		*m = Map[T]{}
	}
	for i, v := range m.Items {
		if v.Key == key {
			m.Items[i].Value = value
			return
		}
	}
	m.Items = append(m.Items, KeyValue[T]{Key: key, Value: value})
}

func (m *Map[T]) Del(key Text) {
	if m == nil {
		return
	}
	for i, v := range m.Items {
		if v.Key == key {
			m.Items = append(m.Items[:i], m.Items[i+1:]...)
			return
		}
	}
}

func (m Map[T]) MarshalJSON() ([]byte, error) {
	b := bytes.Buffer{}
	b.WriteByte('{')
	var err error
	var s []byte
	for _, v := range m.Items {
		if b.Len() > 1 {
			b.WriteByte(',')
		}
		jsonx.EncodeAndWriteString(&b, v.Key.String())
		b.WriteByte(':')
		s, err = json.Marshal(v.Value)
		if err != nil {
			return nil, err
		}
		b.Write(s)
	}
	b.WriteByte('}')
	return b.Bytes(), nil
}

func (m *Map[T]) UnmarshalJSON(data []byte) error {
	*m = Map[T]{}
	var v KeyValue[T]
	if !jsonx.IsObject(data) {
		return &json.UnmarshalTypeError{Value: jsonx.TypeOf(data).String(), Type: reflect.TypeOf(v), Struct: "Map"}
	}
	var err error
	gjson.ParseBytes(data).ForEach(func(key, value gjson.Result) bool {
		var t T
		if err = json.Unmarshal([]byte(value.Raw), &t); err != nil {
			return false
		}
		v = KeyValue[T]{Key: Text(key.String()), Value: t}
		m.Items = append(m.Items, v)
		return true
	})
	return err
}
