package openapi

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/chanced/dynamic"
	"github.com/chanced/openapi/yamlutil"
	"github.com/tidwall/sjson"
	"gopkg.in/yaml.v2"
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

// SchemaKind indicates whether the Schema is a SchemaObj, Reference, or Boolean
type SchemaKind uint8

const (
	// SchemaKindObj = *SchemaObj
	SchemaKindObj SchemaKind = iota
	// SchemaKindBool = *Boolean
	SchemaKindBool
)

// Schema can either be a SchemaObj, Reference, or Boolean
type Schema interface {
	SchemaKind() SchemaKind
	IsRef() bool
}

// Schemas is a map of Schemas
type Schemas map[string]Schema

// UnmarshalJSON unmarshals JSON
func (s *Schemas) UnmarshalJSON(data []byte) error {
	var dm map[string]json.RawMessage
	if err := json.Unmarshal(data, &dm); err != nil {
		return err
	}
	res := make(Schemas, len(dm))

	for k, d := range dm {
		v, err := unmarshalSchemaJSON(d)
		if err != nil {
			return err
		}
		res[k] = v
	}
	*s = res
	return nil
}

// SchemaObj allows the definition of input and output data types. These types can
// be objects, but also primitives and arrays. This object is a superset of the
// [JSON SchemaObj Specification Draft
// 2020-12](https://tools.ietf.org/html/draft-bhutton-json-schema-00).
//
// For more information about the properties, see [JSON SchemaObj
// Core](https://tools.ietf.org/html/draft-bhutton-json-schema-00) and [JSON
// SchemaObj
// Validation](https://tools.ietf.org/html/draft-bhutton-json-schema-validation-00).
//
// Unless stated otherwise, the property definitions follow those of JSON SchemaObj
// and do not add any additional semantics. Where JSON SchemaObj indicates that
// behavior is defined by the application (e.g. for annotations), OAS also
// defers the definition of semantics to the application consuming the OpenAPI
// document.
//
// The OpenAPI SchemaObj Object
// [dialect](https://tools.ietf.org/html/draft-bhutton-json-schema-00#section-4.3.3)
// is defined as requiring the [OAS base vocabulary](#baseVocabulary), in
// addition to the vocabularies as specified in the JSON SchemaObj draft 2020-12
// [general purpose
// meta-schema](https://tools.ietf.org/html/draft-bhutton-json-schema-00#section-8).
//
// The OpenAPI SchemaObj Object dialect for this version of the specification is
// identified by the URI `https://spec.openapis.org/oas/3.1/dialect/base` (the
// <a name="dialectSchemaId"></a>"OAS dialect schema id").
//
// The following properties are taken from the JSON SchemaObj specification but
// their definitions have been extended by the OAS:
//
// - description - [CommonMark syntax](https://spec.commonmark.org/) MAY be used
// for rich text representation. - format - See [Data Type
// Formats](#dataTypeFormat) for further details. While relying on JSON SchemaObj's
// defined formats, the OAS offers a few additional predefined formats.
//
// In addition to the JSON SchemaObj properties comprising the OAS dialect, the
// SchemaObj Object supports keywords from any other vocabularies, or entirely
// arbitrary properties.
// A SchemaObj represents compiled version of json-schema.
type SchemaObj struct {
	Schema string `json:"$schema,omitempty"`
	// The value of $id is a URI-reference without a fragment that resolves
	// against the Retrieval URI. The resulting URI is the base URI for the
	// schema.
	//
	// https://json-schema.org/understanding-json-schema/structuring.html?highlight=id#id
	ID string `json:"$id,omitempty"`
	// At its core, JSON Schema defines the following basic types:
	//
	// 	"string", "number", "integer", "object", "array", "boolean", "null"
	//
	// https://json-schema.org/understanding-json-schema/reference/type.html#type
	Type Types `json:"type,omitempty"`
	// The "$ref" keyword is an applicator that is used to reference a
	// statically identified schema. Its results are the results of the
	// referenced schema. [CREF5]
	//
	// The value of the "$ref" keyword MUST be a string which is a
	// URI-Reference. Resolved against the current URI base, it produces the URI
	// of the schema to apply. This resolution is safe to perform on schema
	// load, as the process of evaluating an instance cannot change how the
	// reference resolves.
	//
	// https://json-schema.org/draft/2020-12/json-schema-core.html#ref
	//
	// https://json-schema.org/understanding-json-schema/structuring.html?highlight=ref#ref
	Ref string `json:"$ref,omitempty"`
	// The "$defs" keyword reserves a location for schema authors to inline
	// re-usable JSON Schemas into a more general schema. The keyword does not
	// directly affect the validation result.
	//
	// This keyword's value MUST be an object. Each member value of this object
	// MUST be a valid JSON Schema.
	//
	// https://json-schema.org/draft/2020-12/json-schema-core.html#defs
	//
	// https://json-schema.org/understanding-json-schema/structuring.html?highlight=defs#defs
	Definitions Schemas `json:"$defs,omitempty"`
	// The format keyword allows for basic semantic identification of certain kinds of string values that are commonly used. For example, because JSON doesn’t have a “DateTime” type, dates need to be encoded as strings. format allows the schema author to indicate that the string value should be interpreted as a date. By default, format is just an annotation and does not effect validation.
	//
	// Optionally, validator implementations can provide a configuration option to
	// enable format to function as an assertion rather than just an annotation.
	// That means that validation will fail if, for example, a value with a date
	// format isn’t in a form that can be parsed as a date. This can allow values to
	// be constrained beyond what the other tools in JSON Schema, including Regular
	// Expressions can do.
	//
	// https://json-schema.org/understanding-json-schema/reference/string.html#format
	Format        string `json:"format,omitempty"`
	DynamicAnchor string `json:"$dynamicAnchor,omitempty"`
	// The "$dynamicRef" keyword is an applicator that allows for deferring the
	// full resolution until runtime, at which point it is resolved each time it
	// is encountered while evaluating an instance.
	//
	// https://json-schema.org/draft/2020-12/json-schema-core.html#dynamic-ref
	DynamicRef string `json:"$dynamicRef,omitempty"`
	// A less common way to identify a subschema is to create a named anchor in
	// the schema using the $anchor keyword and using that name in the URI
	// fragment. Anchors must start with a letter followed by any number of
	// letters, digits, -, _, :, or ..
	//
	// https://json-schema.org/understanding-json-schema/structuring.html?highlight=anchor#anchor
	Anchor string `json:"$anchor,omitempty"`
	// The const keyword is used to restrict a value to a single value.
	//
	// https://json-schema.org/understanding-json-schema/reference/generic.html?highlight=const#constant-values
	Const json.RawMessage `json:"const,omitempty"`
	// The enum keyword is used to restrict a value to a fixed set of values. It
	// must be an array with at least one element, where each element is unique.
	//
	// https://json-schema.org/understanding-json-schema/reference/generic.html?highlight=const#enumerated-values
	Enum []string `json:"enum,omitempty"`
	// The $comment keyword is strictly intended for adding comments to a
	// schema. Its value must always be a string. Unlike the annotations title,
	// description, and examples, JSON schema implementations aren’t allowed to
	// attach any meaning or behavior to it whatsoever, and may even strip them
	// at any time. Therefore, they are useful for leaving notes to future
	// editors of a JSON schema, but should not be used to communicate to users
	// of the schema.
	//
	// https://json-schema.org/understanding-json-schema/reference/generic.html?highlight=const#comments
	Comments string `json:"$comment,omitempty"`

	// The not keyword declares that an instance validates if it doesn’t
	// validate against the given subschema.
	//
	// https://json-schema.org/understanding-json-schema/reference/combining.html?highlight=not#not
	Not Schema `json:"not,omitempty"`
	// validate against allOf, the given data must be valid against all of the
	// given subschemas.
	//
	// https://json-schema.org/understanding-json-schema/reference/combining.html?highlight=anyof#anyof
	AllOf SchemaSet `json:"allOf,omitempty"`
	// validate against anyOf, the given data must be valid against any (one or
	// more) of the given subschemas.
	//
	// https://json-schema.org/understanding-json-schema/reference/combining.html?highlight=allof#allof
	AnyOf SchemaSet `json:"anyOf,omitempty"`
	// alidate against oneOf, the given data must be valid against exactly one of the given subschemas.
	//
	// https://json-schema.org/understanding-json-schema/reference/combining.html?highlight=oneof#oneof
	OneOf SchemaSet `json:"oneOf,omitempty"`
	// if, then and else keywords allow the application of a subschema based on
	// the outcome of another schema, much like the if/then/else constructs
	// you’ve probably seen in traditional programming languages.
	//
	// https://json-schema.org/understanding-json-schema/reference/conditionals.html#if-then-else
	If Schema `json:"if,omitempty"`
	// https://json-schema.org/understanding-json-schema/reference/conditionals.html#if-then-else
	Then Schema `json:"then,omitempty"`
	// https://json-schema.org/understanding-json-schema/reference/conditionals.html#if-then-else
	Else                 Schema   `json:"else,omitempty"`
	MinProperties        *int     `json:"minProperties,omitempty"`
	MaxProperties        *int     `json:"maxProperties,omitempty"`
	Required             []string `json:"required,omitempty"`
	Properties           Schemas  `json:"properties,omitempty"`
	PropertyNames        Schema   `json:"propertyNames,omitempty"`
	RegexProperties      *bool    `json:"regexProperties,omitempty"`
	PatternProperties    Schemas  `json:"patternProperties,omitempty"`
	AdditionalProperties Schema   `json:"additionalProperties,omitempty"`
	// The dependentRequired keyword conditionally requires that certain
	// properties must be present if a given property is present in an object.
	// For example, suppose we have a schema representing a customer. If you
	// have their credit card number, you also want to ensure you have a billing
	// address. If you don’t have their credit card number, a billing address
	// would not be required. We represent this dependency of one property on
	// another using the dependentRequired keyword. The value of the
	// dependentRequired keyword is an object. Each entry in the object maps
	// from the name of a property, p, to an array of strings listing properties
	// that are required if p is present.
	DependentRequired map[string][]string `json:"dependentRequired,omitempty"`
	// The dependentSchemas keyword conditionally applies a subschema when a
	// given property is present. This schema is applied in the same way allOf
	// applies schemas. Nothing is merged or extended. Both schemas apply
	// independently.
	DependentSchemas      Schemas `json:"dependentSchemas,omitempty"`
	UnevaluatedProperties Schema  `json:"unevaluatedProperties,omitempty"`
	UniqueObjs            *bool   `json:"uniqueObjs,omitempty"`
	// List validation is useful for arrays of arbitrary length where each item
	// matches the same schema. For this kind of array, set the items keyword to
	// a single schema that will be used to validate all of the items in the
	// array.
	Items            Schema          `json:"items,omitempty"`
	UnevaluatedObjs  Schema          `json:"unevaluatedObjs,omitempty"`
	AdditionalObjs   Schema          `json:"additionalObjs,omitempty"`
	PrefixObjs       SchemaSet       `json:"prefixObjs,omitempty"`
	Contains         Schema          `json:"contains,omitempty"`
	MinContains      *Number         `json:"minContains,omitempty"`
	MaxContains      *Number         `json:"maxContains,omitempty"`
	MinLength        *Number         `json:"minLength,omitempty"`
	MaxLength        *Number         `json:"maxLength,omitempty"`
	Pattern          *Regexp         `json:"pattern,omitempty"`
	ContentEncoding  string          `json:"contentEncoding,omitempty"`
	ContentMediaType string          `json:"contentMediaType,omitempty"`
	Minimum          *Number         `json:"minimum,omitempty"`
	ExclusiveMinimum *Number         `json:"exclusiveMinimum,omitempty"`
	Maximum          *Number         `json:"maximum,omitempty"`
	ExclusiveMaximum *Number         `json:"exclusiveMaximum,omitempty"`
	MultipleOf       *Number         `json:"multipleOf,omitempty"`
	Title            string          `json:"title,omitempty"`
	Description      string          `json:"description,omitempty"`
	Default          json.RawMessage `json:"default,omitempty"`
	ReadOnly         *bool           `json:"readOnly,omitempty"`
	WriteOnly        *bool           `json:"writeOnly,omitempty"`
	Examples         json.RawMessage `json:"examples,omitempty"`
	Deprecated       *bool           `json:"deprecated,omitempty"`
	ExternalDocs     string          `json:"externalDocs,omitempty"`
	// Deprecated: renamed to dynamicAnchor
	RecursiveAnchor *bool `json:"$recursiveAnchor,omitempty"`
	// Deprecated: renamed to dynamicRef
	RecursiveRef string `json:"$recursiveRef,omitempty"`

	Discriminator *Discriminator `json:"discriminator,omitempty"`
	// This MAY be used only on properties schemas. It has no effect on root
	// schemas. Adds additional metadata to describe the XML representation of
	// this property.
	XML        *XML `json:"xml,omitempty"`
	Extensions `json:"-"`
	Keywords   map[string]json.RawMessage `json:"-"`
}

type schema SchemaObj

// Detail returns a ptr to the Schema
func (s SchemaObj) Detail() *SchemaObj {
	return &s
}

// MarshalJSON marshals JSON
func (s SchemaObj) MarshalJSON() ([]byte, error) {
	data, err := marshalExtendedJSON(schema(s))
	if s.Keywords != nil {
		for k, v := range s.Keywords {
			data, err = sjson.SetBytes(data, k, v)
			if err != nil {
				return data, err
			}
		}
	}
	return data, err
}

// SchemaKind returns SchemaKindObj
func (s *SchemaObj) SchemaKind() SchemaKind {
	return SchemaKindObj
}

// UnmarshalJSON unmarshals JSON
func (s *SchemaObj) UnmarshalJSON(data []byte) error {
	sv, err := unmarshalSchemaJSON(data)
	v, ok := sv.(*SchemaObj)
	if !ok {
		return errors.New("invalid schema")
	}
	*s = *v
	return err
}

// MarshalYAML first marshals and unmarshals into JSON and then marshals into
// YAML
func (s SchemaObj) MarshalYAML() (interface{}, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	var v interface{}
	err = json.Unmarshal(b, &v)
	return v, err
}

// UnmarshalYAML unmarshals yaml into s
func (s *SchemaObj) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return yamlutil.Unmarshal(unmarshal, s)
}

// IsStrings returns false
func (s *SchemaObj) IsStrings() bool {
	return false
}

// IsBool returns false
func (s *SchemaObj) IsBool() bool {
	return false
}

// IsRef returns true if s.Ref is set
func (s *SchemaObj) IsRef() bool {
	return s.Ref != ""
}

// SetKeyword encodes and sets the keyword key to the encoded value
func (s *SchemaObj) SetKeyword(key string, value interface{}) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.SetEncodedKeyword(key, b)

}

// SetEncodedKeyword sets the keyword key to value
func (s *SchemaObj) SetEncodedKeyword(key string, value []byte) error {
	if strings.HasPrefix(key, "x-") {
		return errors.New("keyword keys may not start with \"x-\"")
	}
	s.Keywords[key] = value
	return nil
}

// DecodeKeyword unmarshals the keyword's raw data into dst
func (s *SchemaObj) DecodeKeyword(key string, dst interface{}) error {
	return json.Unmarshal(s.Keywords[key], dst)
}

// DecodeKeywords unmarshals all keywords raw data into dst
func (s *SchemaObj) DecodeKeywords(dst interface{}) error {
	data, err := json.Marshal(s.Keywords)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}

// SchemaSet is a slice of *Schema
type SchemaSet []Schema

//  UnmarshalJSON unmarshals JSON
func (s *SchemaSet) UnmarshalJSON(data []byte) error {
	var j []dynamic.JSON
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	res := make(SchemaSet, len(j))
	for i, d := range j {
		v, err := unmarshalSchemaJSON(d)
		if err != nil {
			return err
		}
		res[i] = v
	}
	*s = res
	return nil
}

func unmarshalSchemaJSON(data []byte) (Schema, error) {
	var str string
	l := len(data)
	if l >= 4 && l <= 5 {
		str = string(data)
	}
	switch {
	case str == "true":
		return Boolean(true), nil
	case str == "false":
		return Boolean(false), nil
	default:
		return unmarshalSchemaObjJSON(data)
	}
}

func unmarshalSchemaObjJSON(data []byte) (Schema, error) {
	var err error
	exts := Extensions{}
	kw := make(map[string]json.RawMessage)
	var dst partialschema
	if err = json.Unmarshal(data, &dst); err != nil {
		return nil, err
	}
	var jm map[string]json.RawMessage
	if err = json.Unmarshal(data, &jm); err != nil {
		return nil, err
	}

	for key, d := range jm {
		if strings.HasPrefix(key, "x-") {
			exts[key] = d
		} else if set, isSchema := schemaFieldSetters[key]; isSchema {
			var v Schema
			v, err = unmarshalSchemaJSON(d)
			if err != nil {
				return nil, err
			}
			set(&dst, v)
		} else if _, isfield := jsfields[key]; !isfield {
			kw[key] = d
		}
	}
	res := SchemaObj(dst)
	res.Keywords = kw
	res.Extensions = exts
	return &res, err
}

var schemaFieldSetters = map[string]func(s *partialschema, v Schema){
	"not":                   func(s *partialschema, v Schema) { s.Not = v },
	"if":                    func(s *partialschema, v Schema) { s.If = v },
	"then":                  func(s *partialschema, v Schema) { s.Then = v },
	"else":                  func(s *partialschema, v Schema) { s.Else = v },
	"propertyNames":         func(s *partialschema, v Schema) { s.PropertyNames = v },
	"additionalProperties":  func(s *partialschema, v Schema) { s.AdditionalProperties = v },
	"unevaluatedProperties": func(s *partialschema, v Schema) { s.UnevaluatedProperties = v },
	"items":                 func(s *partialschema, v Schema) { s.Items = v },
	"contains":              func(s *partialschema, v Schema) { s.Contains = v },
	"unevaluatedObjs":       func(s *partialschema, v Schema) { s.UnevaluatedObjs = v },
	"additionalObjs":        func(s *partialschema, v Schema) { s.AdditionalObjs = v },
}

var jsfields = map[string]struct{}{
	"$schema":               {},
	"$id":                   {},
	"type":                  {},
	"$ref":                  {},
	"$defs":                 {},
	"format":                {},
	"$dynamicAnchor":        {},
	"$dynamicRef":           {},
	"$anchor":               {},
	"const":                 {},
	"enum":                  {},
	"$comment":              {},
	"not":                   {},
	"allOf":                 {},
	"anyOf":                 {},
	"oneOf":                 {},
	"if":                    {},
	"then":                  {},
	"else":                  {},
	"minProperties":         {},
	"maxProperties":         {},
	"required":              {},
	"properties":            {},
	"propertyNames":         {},
	"regexProperties":       {},
	"patternProperties":     {},
	"additionalProperties":  {},
	"dependentRequired":     {},
	"dependentSchemas":      {},
	"unevaluatedProperties": {},
	"uniqueObjs":            {},
	"items":                 {},
	"unevaluatedObjs":       {},
	"additionalObjs":        {},
	"prefixObjs":            {},
	"contains":              {},
	"minContains":           {},
	"maxContains":           {},
	"minLength":             {},
	"maxLength":             {},
	"pattern":               {},
	"contentEncoding":       {},
	"contentMediaType":      {},
	"minimum":               {},
	"exclusiveMinimum":      {},
	"maximum":               {},
	"exclusiveMaximum":      {},
	"multipleOf":            {},
	"title":                 {},
	"description":           {},
	"default":               {},
	"readOnly":              {},
	"writeOnly":             {},
	"examples":              {},
	"deprecated":            {},
	"externalDocs":          {},
	"$recursiveAnchor":      {},
	"$recursiveRef":         {},
	"discriminator":         {},
	"xml":                   {},
}

type partialschema struct {
	Schema                string              `json:"$schema,omitempty"`
	ID                    string              `json:"$id,omitempty"`
	Type                  Types               `json:"type,omitempty"`
	Ref                   string              `json:"$ref,omitempty"`
	Definitions           Schemas             `json:"$defs,omitempty"`
	Format                string              `json:"format,omitempty"`
	DynamicAnchor         string              `json:"$dynamicAnchor,omitempty"`
	DynamicRef            string              `json:"$dynamicRef,omitempty"`
	Anchor                string              `json:"$anchor,omitempty"`
	Const                 json.RawMessage     `json:"const,omitempty"`
	Enum                  []string            `json:"enum,omitempty"`
	Comments              string              `json:"$comment,omitempty"`
	Not                   Schema              `json:"-"`
	AllOf                 SchemaSet           `json:"allOf,omitempty"`
	AnyOf                 SchemaSet           `json:"anyOf,omitempty"`
	OneOf                 SchemaSet           `json:"oneOf,omitempty"`
	If                    Schema              `json:"-"`
	Then                  Schema              `json:"-"`
	Else                  Schema              `json:"-"`
	MinProperties         *int                `json:"minProperties,omitempty"`
	MaxProperties         *int                `json:"maxProperties,omitempty"`
	Required              []string            `json:"required,omitempty"`
	Properties            Schemas             `json:"properties,omitempty"`
	PropertyNames         Schema              `json:"-"`
	RegexProperties       *bool               `json:"regexProperties,omitempty"`
	PatternProperties     Schemas             `json:"patternProperties,omitempty"`
	AdditionalProperties  Schema              `json:"-"`
	DependentRequired     map[string][]string `json:"dependentRequired,omitempty"`
	DependentSchemas      Schemas             `json:"dependentSchemas,omitempty"`
	UnevaluatedProperties Schema              `json:"-"`
	UniqueObjs            *bool               `json:"uniqueObjs,omitempty"`
	Items                 Schema              `json:"-"`
	UnevaluatedObjs       Schema              `json:"-"`
	AdditionalObjs        Schema              `json:"-"`
	PrefixObjs            SchemaSet           `json:"prefixObjs,omitempty"`
	Contains              Schema              `json:"-"`
	MinContains           *Number             `json:"minContains,omitempty"`
	MaxContains           *Number             `json:"maxContains,omitempty"`
	MinLength             *Number             `json:"minLength,omitempty"`
	MaxLength             *Number             `json:"maxLength,omitempty"`
	Pattern               *Regexp             `json:"pattern,omitempty"`
	ContentEncoding       string              `json:"contentEncoding,omitempty"`
	ContentMediaType      string              `json:"contentMediaType,omitempty"`
	Minimum               *Number             `json:"minimum,omitempty"`
	ExclusiveMinimum      *Number             `json:"exclusiveMinimum,omitempty"`
	Maximum               *Number             `json:"maximum,omitempty"`
	ExclusiveMaximum      *Number             `json:"exclusiveMaximum,omitempty"`
	MultipleOf            *Number             `json:"multipleOf,omitempty"`
	Title                 string              `json:"title,omitempty"`
	Description           string              `json:"description,omitempty"`
	Default               json.RawMessage     `json:"default,omitempty"`
	ReadOnly              *bool               `json:"readOnly,omitempty"`
	WriteOnly             *bool               `json:"writeOnly,omitempty"`
	Examples              json.RawMessage     `json:"examples,omitempty"`
	Deprecated            *bool               `json:"deprecated,omitempty"`
	ExternalDocs          string              `json:"externalDocs,omitempty"`
	RecursiveAnchor       *bool               `json:"$recursiveAnchor,omitempty"`
	RecursiveRef          string              `json:"$recursiveRef,omitempty"`
	Discriminator         *Discriminator      `json:"discriminator,omitempty"`
	XML                   *XML                `json:"xml,omitempty"`
	Extensions            `json:"-"`
	Keywords              map[string]json.RawMessage `json:"-"`
}

var _ json.Marshaler = (*SchemaObj)(nil)
var _ json.Unmarshaler = (*SchemaObj)(nil)
var _ yaml.Unmarshaler = (*SchemaObj)(nil)
var _ yaml.Marshaler = (*SchemaObj)(nil)
