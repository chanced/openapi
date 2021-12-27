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

// Schemas is a map of Schemas
type Schemas map[string]*SchemaObj

// Kind returns KindSchemas
func (Schemas) Kind() Kind {
	return KindSchemas
}

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
	// Always will be assigned if the schema value is a boolean
	Always *bool  `json:"-"`
	Schema string `json:"$schema,omitempty"`
	// The value of $id is a URI-reference without a fragment that resolves
	// against the Retrieval URI. The resulting URI is the base URI for the
	// schema.
	//
	// https://json-schema.org/understanding-json-schema/structuring.html?highlight=id#id
	ID string `json:"$id,omitempty"`
	// At its core, JSON *SchemaObj defines the following basic types:
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
	// MUST be a valid JSON *SchemaObj.
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
	// be constrained beyond what the other tools in JSON *SchemaObj, including Regular
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
	Comment string `json:"$comment,omitempty"`

	// The not keyword declares that an instance validates if it doesn’t
	// validate against the given subschema.
	//
	// https://json-schema.org/understanding-json-schema/reference/combining.html?highlight=not#not
	Not *SchemaObj `json:"not,omitempty"`
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
	If *SchemaObj `json:"if,omitempty"`
	// https://json-schema.org/understanding-json-schema/reference/conditionals.html#if-then-else
	Then *SchemaObj `json:"then,omitempty"`
	// https://json-schema.org/understanding-json-schema/reference/conditionals.html#if-then-else
	Else                 *SchemaObj `json:"else,omitempty"`
	MinProperties        *int       `json:"minProperties,omitempty"`
	MaxProperties        *int       `json:"maxProperties,omitempty"`
	Required             []string   `json:"required,omitempty"`
	Properties           Schemas    `json:"properties,omitempty"`
	PropertyNames        *SchemaObj `json:"propertyNames,omitempty"`
	RegexProperties      *bool      `json:"regexProperties,omitempty"`
	PatternProperties    Schemas    `json:"patternProperties,omitempty"`
	AdditionalProperties *SchemaObj `json:"additionalProperties,omitempty"`
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
	DependentSchemas      Schemas    `json:"dependentSchemas,omitempty"`
	UnevaluatedProperties *SchemaObj `json:"unevaluatedProperties,omitempty"`
	UniqueObjs            *bool      `json:"uniqueObjs,omitempty"`
	// List validation is useful for arrays of arbitrary length where each item
	// matches the same schema. For this kind of array, set the items keyword to
	// a single schema that will be used to validate all of the items in the
	// array.
	Items            *SchemaObj      `json:"items,omitempty"`
	UnevaluatedObjs  *SchemaObj      `json:"unevaluatedObjs,omitempty"`
	AdditionalObjs   *SchemaObj      `json:"additionalObjs,omitempty"`
	PrefixObjs       SchemaSet       `json:"prefixObjs,omitempty"`
	Contains         *SchemaObj      `json:"contains,omitempty"`
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

// Detail returns a ptr to the *SchemaObj
func (s SchemaObj) Detail() *SchemaObj {
	return &s
}

// MarshalJSON marshals JSON
func (s SchemaObj) MarshalJSON() ([]byte, error) {
	if s.Always != nil {
		return json.Marshal(s.Always)
	}
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

// Kind returns KindSchema
func (s *SchemaObj) Kind() Kind { return KindSchema }

// ResolveSchema resolves *SchemaObj by returning s
func (s *SchemaObj) ResolveSchema(func(ref string) (*SchemaObj, error)) (*SchemaObj, error) {
	return s, nil
}

// UnmarshalJSON unmarshals JSON
func (s *SchemaObj) UnmarshalJSON(data []byte) error {
	sv, err := unmarshalSchemaJSON(data)
	*s = *sv
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

// IsBool returns true if s.Always is not nil
func (s *SchemaObj) IsBool() bool {
	return s.Always != nil
}

// IsRef returns true if s.Ref is set
func (s *SchemaObj) IsRef() bool {
	return s.Ref != ""
}

// HasField returns true if json representation of the field exists on the
// SchemaObj
//
// - if field is a member of the SchemaObj struct, its value is returned.
//
// - if field starts with "x-" then s.Extensinons is checked
//
// - if field is not a member field and does not have the extensions prefix,
// s.Keywords is checked
func (s *SchemaObj) HasField(field string) bool {
	if _, ok := schemaobjfields[field]; ok {
		return true
	}
	if strings.HasPrefix(field, "x-") {
		_, ok := s.Extensions[field]
		return ok
	}
	_, ok := s.Keywords[field]
	return ok
}

// ValueOf returns the value of the field if it exists and a bool indicating
// whether or not it does.
//
// - if field is a member of the SchemaObj struct, its value is returned.
//
// - if field starts with "x-" then s.Extensinons is checked
//
// - if field is not a member field and does not have the extensions prefix,
// s.Keywords is checked
func (s *SchemaObj) ValueOf(field string) (interface{}, bool) {
	if fn, ok := schemaobjfields[field]; ok {
		return fn(s), true
	}
	if strings.HasPrefix(field, "x-") {
		v, ok := s.Extensions[field]
		return v, ok
	}
	v, ok := s.Keywords[field]
	return v, ok
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

// SchemaSet is a slice of **SchemaObj
type SchemaSet []*SchemaObj

func (s SchemaSet) Len() int {
	return len(s)
}

// Kind returns KindSchemaSet
func (s SchemaSet) Kind() Kind {
	return KindSchemaSet
}

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

func unmarshalSchemaJSON(data []byte) (*SchemaObj, error) {
	var str string
	l := len(data)
	if l >= 4 && l <= 5 {
		str = string(data)
	}
	switch {
	case str == "true":
		t := true
		return &SchemaObj{Always: &t}, nil
	case str == "false":
		f := false
		return &SchemaObj{Always: &f}, nil
	default:
		return unmarshalSchemaObjJSON(data)
	}
}

func unmarshalSchemaObjJSON(data []byte) (*SchemaObj, error) {
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
			var v *SchemaObj
			v, err = unmarshalSchemaJSON(d)
			if err != nil {
				return nil, err
			}
			set(&dst, v)
		} else if _, isfield := schemaobjfields[key]; !isfield {
			kw[key] = d
		}
	}
	res := SchemaObj(dst)
	res.Keywords = kw
	res.Extensions = exts
	return &res, err
}

type partialschema struct {
	Always                *bool               `json:"-"`
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
	Comment               string              `json:"$comment,omitempty"`
	Not                   *SchemaObj          `json:"-"`
	AllOf                 SchemaSet           `json:"allOf,omitempty"`
	AnyOf                 SchemaSet           `json:"anyOf,omitempty"`
	OneOf                 SchemaSet           `json:"oneOf,omitempty"`
	If                    *SchemaObj          `json:"-"`
	Then                  *SchemaObj          `json:"-"`
	Else                  *SchemaObj          `json:"-"`
	MinProperties         *int                `json:"minProperties,omitempty"`
	MaxProperties         *int                `json:"maxProperties,omitempty"`
	Required              []string            `json:"required,omitempty"`
	Properties            Schemas             `json:"properties,omitempty"`
	PropertyNames         *SchemaObj          `json:"-"`
	RegexProperties       *bool               `json:"regexProperties,omitempty"`
	PatternProperties     Schemas             `json:"patternProperties,omitempty"`
	AdditionalProperties  *SchemaObj          `json:"-"`
	DependentRequired     map[string][]string `json:"dependentRequired,omitempty"`
	DependentSchemas      Schemas             `json:"dependentSchemas,omitempty"`
	UnevaluatedProperties *SchemaObj          `json:"-"`
	UniqueObjs            *bool               `json:"uniqueObjs,omitempty"`
	Items                 *SchemaObj          `json:"-"`
	UnevaluatedObjs       *SchemaObj          `json:"-"`
	AdditionalObjs        *SchemaObj          `json:"-"`
	PrefixObjs            SchemaSet           `json:"prefixObjs,omitempty"`
	Contains              *SchemaObj          `json:"-"`
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

type ResolvedSchemas map[string]*ResolvedSchema

func (rs ResolvedSchemas) Kind() Kind {
	return KindSchema
}

func (rs *ResolvedSchemas) Set(key string, val *ResolvedSchema) {
	if *rs == nil {
		*rs = ResolvedSchemas{key: val}
	} else {
		(*rs)[key] = val
	}
}

func (rs ResolvedSchemas) Get(name string) *ResolvedSchema {
	if v, ok := rs[name]; ok {
		return v
	}
	return nil
}

// ResolvedSchemaSet is a slice of *ResolvedSchemas
type ResolvedSchemaSet []*ResolvedSchema

// Kind returns KindResolvedSchemaSet
func (ResolvedSchemaSet) Kind() Kind {
	return KindResolvedSchemaSet
}

type ResolvedSchema struct {
	Always                *bool               `json:"-"`
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
	Comment               string              `json:"$comment,omitempty"`
	Not                   *ResolvedSchema     `json:"-"`
	AllOf                 ResolvedSchemaSet   `json:"allOf,omitempty"`
	AnyOf                 ResolvedSchemaSet   `json:"anyOf,omitempty"`
	OneOf                 ResolvedSchemaSet   `json:"oneOf,omitempty"`
	If                    *ResolvedSchema     `json:"-"`
	Then                  *ResolvedSchema     `json:"-"`
	Else                  *ResolvedSchema     `json:"-"`
	MinProperties         *int                `json:"minProperties,omitempty"`
	MaxProperties         *int                `json:"maxProperties,omitempty"`
	Required              []string            `json:"required,omitempty"`
	Properties            Schemas             `json:"properties,omitempty"`
	PropertyNames         *SchemaObj          `json:"-"`
	RegexProperties       *bool               `json:"regexProperties,omitempty"`
	PatternProperties     ResolvedSchemas     `json:"patternProperties,omitempty"`
	AdditionalProperties  *ResolvedSchema     `json:"-"`
	DependentRequired     map[string][]string `json:"dependentRequired,omitempty"`
	DependentSchemas      ResolvedSchemas     `json:"dependentSchemas,omitempty"`
	UnevaluatedProperties *SchemaObj          `json:"-"`
	UniqueObjs            *bool               `json:"uniqueObjs,omitempty"`
	Items                 *ResolvedSchema     `json:"-"`
	UnevaluatedObjs       *ResolvedSchema     `json:"-"`
	AdditionalObjs        *ResolvedSchema     `json:"-"`
	PrefixObjs            SchemaSet           `json:"prefixObjs,omitempty"`
	Contains              *ResolvedSchema     `json:"-"`
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

func (*ResolvedSchema) Kind() Kind {
	return KindResolvedSchema
}

var schemaFieldSetters = map[string]func(s *partialschema, v *SchemaObj){
	"not":                   func(s *partialschema, v *SchemaObj) { s.Not = v },
	"if":                    func(s *partialschema, v *SchemaObj) { s.If = v },
	"then":                  func(s *partialschema, v *SchemaObj) { s.Then = v },
	"else":                  func(s *partialschema, v *SchemaObj) { s.Else = v },
	"propertyNames":         func(s *partialschema, v *SchemaObj) { s.PropertyNames = v },
	"additionalProperties":  func(s *partialschema, v *SchemaObj) { s.AdditionalProperties = v },
	"unevaluatedProperties": func(s *partialschema, v *SchemaObj) { s.UnevaluatedProperties = v },
	"items":                 func(s *partialschema, v *SchemaObj) { s.Items = v },
	"contains":              func(s *partialschema, v *SchemaObj) { s.Contains = v },
	"unevaluatedObjs":       func(s *partialschema, v *SchemaObj) { s.UnevaluatedObjs = v },
	"additionalObjs":        func(s *partialschema, v *SchemaObj) { s.AdditionalObjs = v },
}

var schemaobjfields = map[string]func(s *SchemaObj) interface{}{
	"$schema":               func(s *SchemaObj) interface{} { return s.Schema },
	"$id":                   func(s *SchemaObj) interface{} { return s.ID },
	"type":                  func(s *SchemaObj) interface{} { return s.Type },
	"$ref":                  func(s *SchemaObj) interface{} { return s.Ref },
	"$defs":                 func(s *SchemaObj) interface{} { return s.Definitions },
	"format":                func(s *SchemaObj) interface{} { return s.Format },
	"$dynamicAnchor":        func(s *SchemaObj) interface{} { return s.DynamicAnchor },
	"$dynamicRef":           func(s *SchemaObj) interface{} { return s.DynamicRef },
	"$anchor":               func(s *SchemaObj) interface{} { return s.Anchor },
	"const":                 func(s *SchemaObj) interface{} { return s.Const },
	"enum":                  func(s *SchemaObj) interface{} { return s.Enum },
	"$comment":              func(s *SchemaObj) interface{} { return s.Comment },
	"not":                   func(s *SchemaObj) interface{} { return s.Not },
	"allOf":                 func(s *SchemaObj) interface{} { return s.AllOf },
	"anyOf":                 func(s *SchemaObj) interface{} { return s.AnyOf },
	"oneOf":                 func(s *SchemaObj) interface{} { return s.OneOf },
	"if":                    func(s *SchemaObj) interface{} { return s.If },
	"then":                  func(s *SchemaObj) interface{} { return s.Then },
	"else":                  func(s *SchemaObj) interface{} { return s.Else },
	"minProperties":         func(s *SchemaObj) interface{} { return s.MinProperties },
	"maxProperties":         func(s *SchemaObj) interface{} { return s.MaxProperties },
	"required":              func(s *SchemaObj) interface{} { return s.Required },
	"properties":            func(s *SchemaObj) interface{} { return s.Properties },
	"propertyNames":         func(s *SchemaObj) interface{} { return s.PropertyNames },
	"regexProperties":       func(s *SchemaObj) interface{} { return s.RegexProperties },
	"patternProperties":     func(s *SchemaObj) interface{} { return s.PatternProperties },
	"additionalProperties":  func(s *SchemaObj) interface{} { return s.AdditionalProperties },
	"dependentRequired":     func(s *SchemaObj) interface{} { return s.DependentRequired },
	"dependentSchemas":      func(s *SchemaObj) interface{} { return s.DependentSchemas },
	"unevaluatedProperties": func(s *SchemaObj) interface{} { return s.UnevaluatedProperties },
	"uniqueObjs":            func(s *SchemaObj) interface{} { return s.UniqueObjs },
	"items":                 func(s *SchemaObj) interface{} { return s.Items },
	"unevaluatedObjs":       func(s *SchemaObj) interface{} { return s.UnevaluatedObjs },
	"additionalObjs":        func(s *SchemaObj) interface{} { return s.AdditionalObjs },
	"prefixObjs":            func(s *SchemaObj) interface{} { return s.PrefixObjs },
	"contains":              func(s *SchemaObj) interface{} { return s.Contains },
	"minContains":           func(s *SchemaObj) interface{} { return s.MinContains },
	"maxContains":           func(s *SchemaObj) interface{} { return s.MaxContains },
	"minLength":             func(s *SchemaObj) interface{} { return s.MinLength },
	"maxLength":             func(s *SchemaObj) interface{} { return s.MaxLength },
	"pattern":               func(s *SchemaObj) interface{} { return s.Pattern },
	"contentEncoding":       func(s *SchemaObj) interface{} { return s.ContentEncoding },
	"contentMediaType":      func(s *SchemaObj) interface{} { return s.ContentMediaType },
	"minimum":               func(s *SchemaObj) interface{} { return s.Minimum },
	"exclusiveMinimum":      func(s *SchemaObj) interface{} { return s.ExclusiveMinimum },
	"maximum":               func(s *SchemaObj) interface{} { return s.Maximum },
	"exclusiveMaximum":      func(s *SchemaObj) interface{} { return s.ExclusiveMaximum },
	"multipleOf":            func(s *SchemaObj) interface{} { return s.MultipleOf },
	"title":                 func(s *SchemaObj) interface{} { return s.Title },
	"description":           func(s *SchemaObj) interface{} { return s.Description },
	"default":               func(s *SchemaObj) interface{} { return s.Default },
	"readOnly":              func(s *SchemaObj) interface{} { return s.ReadOnly },
	"writeOnly":             func(s *SchemaObj) interface{} { return s.WriteOnly },
	"examples":              func(s *SchemaObj) interface{} { return s.Examples },
	"deprecated":            func(s *SchemaObj) interface{} { return s.Deprecated },
	"externalDocs":          func(s *SchemaObj) interface{} { return s.ExternalDocs },
	"$recursiveAnchor":      func(s *SchemaObj) interface{} { return s.RecursiveAnchor },
	"$recursiveRef":         func(s *SchemaObj) interface{} { return s.RecursiveRef },
	"discriminator":         func(s *SchemaObj) interface{} { return s.Discriminator },
	"xml":                   func(s *SchemaObj) interface{} { return s.XML },
}

var resolvedSchemafields = map[string]func(s *ResolvedSchema) interface{}{
	"$schema":               func(s *ResolvedSchema) interface{} { return s.Schema },
	"$id":                   func(s *ResolvedSchema) interface{} { return s.ID },
	"type":                  func(s *ResolvedSchema) interface{} { return s.Type },
	"$ref":                  func(s *ResolvedSchema) interface{} { return s.Ref },
	"$defs":                 func(s *ResolvedSchema) interface{} { return s.Definitions },
	"format":                func(s *ResolvedSchema) interface{} { return s.Format },
	"$dynamicAnchor":        func(s *ResolvedSchema) interface{} { return s.DynamicAnchor },
	"$dynamicRef":           func(s *ResolvedSchema) interface{} { return s.DynamicRef },
	"$anchor":               func(s *ResolvedSchema) interface{} { return s.Anchor },
	"const":                 func(s *ResolvedSchema) interface{} { return s.Const },
	"enum":                  func(s *ResolvedSchema) interface{} { return s.Enum },
	"$comment":              func(s *ResolvedSchema) interface{} { return s.Comment },
	"not":                   func(s *ResolvedSchema) interface{} { return s.Not },
	"allOf":                 func(s *ResolvedSchema) interface{} { return s.AllOf },
	"anyOf":                 func(s *ResolvedSchema) interface{} { return s.AnyOf },
	"oneOf":                 func(s *ResolvedSchema) interface{} { return s.OneOf },
	"if":                    func(s *ResolvedSchema) interface{} { return s.If },
	"then":                  func(s *ResolvedSchema) interface{} { return s.Then },
	"else":                  func(s *ResolvedSchema) interface{} { return s.Else },
	"minProperties":         func(s *ResolvedSchema) interface{} { return s.MinProperties },
	"maxProperties":         func(s *ResolvedSchema) interface{} { return s.MaxProperties },
	"required":              func(s *ResolvedSchema) interface{} { return s.Required },
	"properties":            func(s *ResolvedSchema) interface{} { return s.Properties },
	"propertyNames":         func(s *ResolvedSchema) interface{} { return s.PropertyNames },
	"regexProperties":       func(s *ResolvedSchema) interface{} { return s.RegexProperties },
	"patternProperties":     func(s *ResolvedSchema) interface{} { return s.PatternProperties },
	"additionalProperties":  func(s *ResolvedSchema) interface{} { return s.AdditionalProperties },
	"dependentRequired":     func(s *ResolvedSchema) interface{} { return s.DependentRequired },
	"dependentSchemas":      func(s *ResolvedSchema) interface{} { return s.DependentSchemas },
	"unevaluatedProperties": func(s *ResolvedSchema) interface{} { return s.UnevaluatedProperties },
	"uniqueObjs":            func(s *ResolvedSchema) interface{} { return s.UniqueObjs },
	"items":                 func(s *ResolvedSchema) interface{} { return s.Items },
	"unevaluatedObjs":       func(s *ResolvedSchema) interface{} { return s.UnevaluatedObjs },
	"additionalObjs":        func(s *ResolvedSchema) interface{} { return s.AdditionalObjs },
	"prefixObjs":            func(s *ResolvedSchema) interface{} { return s.PrefixObjs },
	"contains":              func(s *ResolvedSchema) interface{} { return s.Contains },
	"minContains":           func(s *ResolvedSchema) interface{} { return s.MinContains },
	"maxContains":           func(s *ResolvedSchema) interface{} { return s.MaxContains },
	"minLength":             func(s *ResolvedSchema) interface{} { return s.MinLength },
	"maxLength":             func(s *ResolvedSchema) interface{} { return s.MaxLength },
	"pattern":               func(s *ResolvedSchema) interface{} { return s.Pattern },
	"contentEncoding":       func(s *ResolvedSchema) interface{} { return s.ContentEncoding },
	"contentMediaType":      func(s *ResolvedSchema) interface{} { return s.ContentMediaType },
	"minimum":               func(s *ResolvedSchema) interface{} { return s.Minimum },
	"exclusiveMinimum":      func(s *ResolvedSchema) interface{} { return s.ExclusiveMinimum },
	"maximum":               func(s *ResolvedSchema) interface{} { return s.Maximum },
	"exclusiveMaximum":      func(s *ResolvedSchema) interface{} { return s.ExclusiveMaximum },
	"multipleOf":            func(s *ResolvedSchema) interface{} { return s.MultipleOf },
	"title":                 func(s *ResolvedSchema) interface{} { return s.Title },
	"description":           func(s *ResolvedSchema) interface{} { return s.Description },
	"default":               func(s *ResolvedSchema) interface{} { return s.Default },
	"readOnly":              func(s *ResolvedSchema) interface{} { return s.ReadOnly },
	"writeOnly":             func(s *ResolvedSchema) interface{} { return s.WriteOnly },
	"examples":              func(s *ResolvedSchema) interface{} { return s.Examples },
	"deprecated":            func(s *ResolvedSchema) interface{} { return s.Deprecated },
	"externalDocs":          func(s *ResolvedSchema) interface{} { return s.ExternalDocs },
	"$recursiveAnchor":      func(s *ResolvedSchema) interface{} { return s.RecursiveAnchor },
	"$recursiveRef":         func(s *ResolvedSchema) interface{} { return s.RecursiveRef },
	"discriminator":         func(s *ResolvedSchema) interface{} { return s.Discriminator },
	"xml":                   func(s *ResolvedSchema) interface{} { return s.XML },
}

var (
	_ Node             = (*SchemaObj)(nil)
	_ Node             = (SchemaSet)(nil)
	_ Node             = (Schemas)(nil)
	_ Node             = (*ResolvedSchema)(nil)
	_ Node             = (ResolvedSchemas)(nil)
	_ Node             = (ResolvedSchemaSet)(nil)
	_ json.Marshaler   = (*SchemaObj)(nil)
	_ json.Unmarshaler = (*SchemaObj)(nil)
	_ yaml.Unmarshaler = (*SchemaObj)(nil)
	_ yaml.Marshaler   = (*SchemaObj)(nil)
)
