package openapi

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/chanced/jsonx"
	"github.com/tidwall/gjson"
)

type JSONObjEntry struct {
	Key   Text
	Value jsonx.RawMessage
}

type OrderedJSONObj []JSONObjEntry

func (j OrderedJSONObj) MarshalJSON() ([]byte, error) {
	b := strings.Builder{}
	b.WriteByte('{')
	for _, e := range j {
		if b.Len() > 1 {
			b.WriteByte(',')
		}
		b.WriteString("\"" + e.Key.String() + "\":")
		if e.Value == nil {
			b.WriteString("null")
		} else {
			b.Write(e.Value)
		}
	}
	b.WriteByte('}')
	return []byte(b.String()), nil
}

func (j *OrderedJSONObj) UnmarshalJSON(data []byte) error {
	t := jsonx.TypeOf(data)
	var v OrderedJSONObj
	switch t {
	case jsonx.TypeNull:
		*j = nil
		return nil
	case jsonx.TypeObject:
		gjson.ParseBytes(data).ForEach(func(key, value gjson.Result) bool {
			v = append(v, JSONObjEntry{
				Key:   Text(key.String()),
				Value: jsonx.RawMessage(value.Raw),
			})
			return true
		})
		return nil
	default:
		return &json.UnmarshalTypeError{Value: t.String(), Type: reflect.TypeOf(jsonx.TypeObject)}
	}
}

func (j OrderedJSONObj) Get(key Text) jsonx.RawMessage {
	for _, v := range j {
		if v.Key == key {
			return v.Value
		}
	}
	return nil
}

// Has returns true if key exists in j
func (j OrderedJSONObj) Has(key Text) bool {
	for _, v := range j {
		if v.Key == Text(key) {
			return true
		}
	}
	return false
}

func (j OrderedJSONObj) Map() map[string]jsonx.RawMessage {
	m := make(map[string]jsonx.RawMessage, len(j))
	for _, v := range j {
		m[v.Key.String()] = v.Value
	}
	return m
}

// Set concrete object to lp. To add JSON, use SetEncoded
func (j *OrderedJSONObj) Set(key Text, value interface{}) error {
	var data []byte
	var ok bool
	var err error
	if data, ok = value.([]byte); !ok {
		data, err = json.Marshal(value)
		if err != nil {
			return err
		}
	}
	for i, v := range *j {
		if v.Key == key {
			(*j)[i] = JSONObjEntry{
				Key:   Text(key),
				Value: data,
			}
		}
	}
	return nil
}

// DecodeValue decodes a given parameter by key.
func (j OrderedJSONObj) DecodeValue(key Text, dst interface{}) error {
	if g := j.Get(key); g != nil {
		return json.Unmarshal(g, dst)
	} else {
		return json.Unmarshal([]byte("null"), dst)
	}
}

// Decode decodes all of j into dst
//
// For field-level decoding, use DecodeValue
func (j OrderedJSONObj) Decode(dst interface{}) error {
	b, err := json.Marshal(j.Map())
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dst)
}
