package openapi

import (
	"encoding/json"
	"reflect"

	"github.com/chanced/jsonx"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type ObjMapEntry[T node] struct {
	Key    Text
	Object T
}

// ObjMap is a map of OpenAPI Objects of type T
type ObjMap[T node] []ObjMapEntry[T]

func (om ObjMap[T]) setLocation(loc Location) error {
	for _, kv := range om {
		if err := kv.Object.setLocation(loc); err != nil {
			return err
		}
	}
	return nil
}

func (om *ObjMap[T]) Get(key Text) T {
	var t T
	for _, kv := range *om {
		if kv.Key == key {
			t = kv.Object
			break
		}
	}
	return t
}

func (om *ObjMap[T]) Set(key Text, obj T) {
	for i, kv := range *om {
		if kv.Key == key {
			(*om)[i] = ObjMapEntry[T]{key, obj}
			return
		}
	}
	*om = append(*om, ObjMapEntry[T]{key, obj})
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
		if err = t.UnmarshalJSON([]byte(value.Raw)); err != nil {
			return false
		}
		m = append(m, ObjMapEntry[T]{Key: Text(key.String()), Object: pi})
		return true
	})
	return err
}

func (om *ObjMap[T]) MarshalJSON() ([]byte, error) {
	data := []byte("{}")
	var err error
	var j []byte
	for _, entry := range *om {
		j, err = entry.Object.MarshalJSON()
		if err != nil {
			return nil, err
		}
		data, err = sjson.SetRawBytesOptions(data, entry.Key.String(), j, &sjson.Options{ReplaceInPlace: true})
		if err != nil {
			return nil, err
		}
	}
	return data, err
}
