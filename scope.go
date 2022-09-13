package openapi

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/chanced/jsonpointer"
	"github.com/chanced/jsonx"
	"github.com/tidwall/gjson"
)

type Scope struct {
	Location `json:"-"`
	Key      Text `json:"-"`
	Value    Text `json:"-"`
}

func (*Scope) Anchors() (*Anchors, error) { return nil, nil }
func (*Scope) Kind() Kind                 { return KindScope }
func (*Scope) mapKind() Kind              { return KindUndefined }
func (*Scope) sliceKind() Kind            { return KindUndefined }

func (s *Scope) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
	if err := ptr.Validate(); err != nil {
		return nil, err
	}
	return s.resolveNodeByPointer(ptr)
}

func (s *Scope) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
	if ptr.IsRoot() {
		return s, nil
	}
	tok, _ := ptr.NextToken()
	return nil, newErrNotResolvable(s.AbsoluteLocation(), tok)
}

func (s Scope) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Value)
}

func (s *Scope) UnmarshalJSON(data []byte) error {
	*s = Scope{}
	if len(data) == 0 {
		return nil
	}
	t := jsonx.TypeOf(data)
	switch t {
	case jsonx.TypeString:
		var v Text
		err := json.Unmarshal(data, &v)
		if err != nil {
			return err
		}
		s.Value = v
		return nil
	default:
		var v map[Text]Text
		err := json.Unmarshal(data, &v)
		if err != nil {
			return err
		}
		if len(v) > 1 {
			return fmt.Errorf("can not unmarshal more than a single key/value pair into a Scope")
		}
		for k, v := range v {
			s.Key = k
			s.Value = v
			break
		}
		return nil
	}
}

func (s *Scope) setLocation(loc Location) error {
	s.Location = loc
	return nil
}

func (s Scope) String() string {
	return s.Value.String()
}

func (s Scope) Text() Text {
	return s.Value
}

type Scopes struct {
	Location `json:"-"`

	Items []*Scope `json:"-"`
}

func (*Scopes) Kind() Kind      { return KindScopes }
func (*Scopes) mapKind() Kind   { return KindUndefined }
func (*Scopes) sliceKind() Kind { return KindUndefined }

func (s *Scopes) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
	if err := ptr.Validate(); err != nil {
		return nil, err
	}
	return s.resolveNodeByPointer(ptr)
}

func (s *Scopes) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
	if ptr.IsRoot() {
		return s, nil
	}
	tok, _ := ptr.NextToken()
	if tok == "" {
		return s, nil
	}
	tk := Text(tok)
	for _, v := range s.Items {
		if v.Key == tk {
			return v.ResolveNodeByPointer(ptr)
		}
	}
	return nil, newErrNotFound(s.AbsoluteLocation(), tok)
}

func (s *Scopes) setLocation(loc Location) error {
	if s == nil {
		return nil
	}
	s.Location = loc
	for _, item := range s.Items {
		item.Location = loc.Append(item.Key.String())
	}
	return nil
}

func (s Scopes) Get(key Text) *Scope {
	for _, v := range s.Items {
		if v.Key == key {
			return v
		}
	}
	return nil
}

func (s Scopes) Has(key Text) bool {
	for _, v := range s.Items {
		if v.Key == key {
			return true
		}
	}
	return false
}

func (s *Scopes) Set(key Text, value Text) {
	if s == nil {
		*s = Scopes{}
	}
	for _, v := range s.Items {
		if v.Key == key {
			v.Value = value
			return
		}
	}
	s.Items = append(s.Items, &Scope{
		Key:   key,
		Value: value,
	})
}

func (s *Scopes) Map() map[Text]Text {
	if s == nil || s.Items == nil {
		return nil
	}
	m := make(map[Text]Text, len(s.Items))
	for _, v := range s.Items {
		m[v.Key] = v.Value
	}
	return m
}

func (s Scopes) MarshalJSON() ([]byte, error) {
	b := strings.Builder{}
	b.WriteByte('{')
	for _, v := range s.Items {
		if b.Len() > 1 {
			b.WriteByte(',')
		}
		key, err := json.Marshal(v.Key)
		if err != nil {
			return nil, err
		}
		b.Write(key)
		b.WriteByte(':')
		value, err := json.Marshal(v.Value)
		if err != nil {
			return nil, err
		}
		b.Write(value)
	}
	b.WriteByte('}')
	return []byte(b.String()), nil
}

func (s *Scopes) UnmarshalJSON(data []byte) error {
	*s = Scopes{}
	if len(data) == 0 {
		return nil
	}
	if !jsonx.IsObject(data) {
		return &json.UnmarshalTypeError{
			Value:  jsonx.TypeOf(data).String(),
			Type:   reflect.TypeOf(s),
			Struct: "Scopes",
		}
	}
	var err error
	var v string
	gjson.ParseBytes(data).ForEach(func(key, value gjson.Result) bool {
		err = json.Unmarshal([]byte(value.Raw), &v)
		if err != nil {
			return false
		}
		s.Items = append(s.Items, &Scope{
			Key:   Text(key.String()),
			Value: Text(v),
		})
		return true
	})
	return nil
}

var (
	_ node = (*Scope)(nil)
	_ node = (*Scopes)(nil)
)
