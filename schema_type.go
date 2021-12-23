package openapi

import (
	"encoding/json"

	"github.com/chanced/dynamic"
)

const (
	// TypeString = string
	//
	// https://json-schema.org/understanding-json-schema/reference/string.html#string
	TypeString SchemaType = "string"
	// TypeNumber = number
	//
	// https://json-schema.org/understanding-json-schema/reference/numeric.html#number
	TypeNumber SchemaType = "number"
	// TypeInteger = integer
	//
	// https://json-schema.org/understanding-json-schema/reference/numeric.html#integer
	TypeInteger SchemaType = "integer"
	// TypeObject = object
	//
	// https://json-schema.org/understanding-json-schema/reference/object.html#object
	TypeObject SchemaType = "object"
	// TypeArray = array
	//
	// https://json-schema.org/understanding-json-schema/reference/array.html#array
	TypeArray SchemaType = "array"
	// TypeBoolean = boolean
	//
	// https://json-schema.org/understanding-json-schema/reference/boolean.html#boolean
	TypeBoolean SchemaType = "boolean"
	// TypeNull = null
	//
	// https://json-schema.org/understanding-json-schema/reference/null.html#null
	TypeNull SchemaType = "null"
)

// SchemaType restricts to a JSON Schema specific type
//
// https://json-schema.org/understanding-json-schema/reference/type.html#type
type SchemaType string

func (t SchemaType) String() string {
	return string(t)
}

// Types is a set of Types. A single Type marshals/unmarshals into a string
// while 2+ marshals into an array.
type Types []SchemaType
type types Types

// ContainsString returns true if TypeString is present
func (t Types) ContainsString() bool {
	return t.Contains(TypeString)
}

// ContainsNumber returns true if TypeNumber is present
func (t Types) ContainsNumber() bool {
	return t.Contains(TypeNumber)
}

// ContainsInteger returns true if TypeInteger is present
func (t Types) ContainsInteger() bool {
	return t.Contains(TypeInteger)
}

// ContainsObject returns true if TypeObject is present
func (t Types) ContainsObject() bool {
	return t.Contains(TypeObject)
}

// ContainsArray returns true if TypeArray is present
func (t Types) ContainsArray() bool {
	return t.Contains(TypeArray)
}

// ContainsBoolean returns true if TypeBoolean is present
func (t Types) ContainsBoolean() bool {
	return t.Contains(TypeBoolean)
}

// ContainsNull returns true if TypeNull is present
func (t Types) ContainsNull() bool {
	return t.Contains(TypeNull)
}

// IsSingle returns true if len(t) == 1
func (t Types) IsSingle() bool {
	return len(t) == 1
}

// IsEmpty returns true if len(t) == 0
func (t SchemaType) IsEmpty() bool {
	return len(t) == 0
}

// Len returns len(t)
func (t Types) Len() int {
	return len(t)
}

// Contains returns true if t contains typ
func (t Types) Contains(typ SchemaType) bool {
	for _, v := range t {
		if v == typ {
			return true
		}
	}
	return false
}

// Add adds typ if not present
func (t *Types) Add(typ SchemaType) Types {
	if !t.Contains(typ) {
		*t = append(*t, typ)
	}
	return *t
}

// Remove removes typ if present
func (t *Types) Remove(typ SchemaType) Types {
	for i, v := range *t {
		if typ == v {
			copy((*t)[i:], (*t)[i+1:])
			(*t)[len(*t)-1] = ""
			*t = (*t)[:len(*t)-1]
		}
	}
	return *t
}

// MarshalJSON marshals JSON
func (t Types) MarshalJSON() ([]byte, error) {
	switch len(t) {
	case 1:
		return json.Marshal(t[0].String())
	default:
		return json.Marshal(types(t))
	}
}

// UnmarshalJSON unmarshals JSON
func (t *Types) UnmarshalJSON(data []byte) error {
	d := dynamic.JSON(data)
	if d.IsString() {
		var v SchemaType
		err := json.Unmarshal(data, &v)
		*t = Types{v}
		return err
	}
	var v types
	err := json.Unmarshal(data, &v)
	*t = Types(v)
	return err
}
