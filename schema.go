package openapi

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/chanced/dynamic"
	"github.com/tidwall/sjson"
)

// Schemas is a map of Schemas
type Schemas map[string]*Schema

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

// Schema allows the definition of input and output data types. These types can
// be objects, but also primitives and arrays. This object is a superset of the
// [JSON Schema Specification Draft
// 2020-12](https://tools.ietf.org/html/draft-bhutton-json-schema-00).
//
// For more information about the properties, see [JSON Schema
// Core](https://tools.ietf.org/html/draft-bhutton-json-schema-00) and [JSON
// Schema
// Validation](https://tools.ietf.org/html/draft-bhutton-json-schema-validation-00).
//
// Unless stated otherwise, the property definitions follow those of JSON Schema
// and do not add any additional semantics. Where JSON Schema indicates that
// behavior is defined by the application (e.g. for annotations), OAS also
// defers the definition of semantics to the application consuming the OpenAPI
// document.
//
// The OpenAPI Schema Object
// [dialect](https://tools.ietf.org/html/draft-bhutton-json-schema-00#section-4.3.3)
// is defined as requiring the [OAS base vocabulary](#baseVocabulary), in
// addition to the vocabularies as specified in the JSON Schema draft 2020-12
// [general purpose
// meta-schema](https://tools.ietf.org/html/draft-bhutton-json-schema-00#section-8).
//
// The OpenAPI Schema Object dialect for this version of the specification is
// identified by the URI `https://spec.openapis.org/oas/3.1/dialect/base` (the
// <a name="dialectSchemaId"></a>"OAS dialect schema id").
//
// The following properties are taken from the JSON Schema specification but
// their definitions have been extended by the OAS:
//
// - description - [CommonMark syntax](https://spec.commonmark.org/) MAY be used
// for rich text representation. - format - See [Data Type
// Formats](#dataTypeFormat) for further details. While relying on JSON Schema's
// defined formats, the OAS offers a few additional predefined formats.
//
// In addition to the JSON Schema properties comprising the OAS dialect, the
// Schema Object supports keywords from any other vocabularies, or entirely
// arbitrary properties.
// A Schema represents compiled version of json-schema.
type Schema struct {
	RefResolved          *Schema `json:"-"`
	RecursiveRefResolved *Schema `json:"-"`
	DynamicRefResolved   *Schema `json:"-"`
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
	Comments string `json:"$comment,omitempty"`

	// The not keyword declares that an instance validates if it doesn’t
	// validate against the given subschema.
	//
	// https://json-schema.org/understanding-json-schema/reference/combining.html?highlight=not#not
	Not *Schema `json:"not,omitempty"`
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
	If *Schema `json:"if,omitempty"`
	// https://json-schema.org/understanding-json-schema/reference/conditionals.html#if-then-else
	Then *Schema `json:"then,omitempty"`
	// https://json-schema.org/understanding-json-schema/reference/conditionals.html#if-then-else
	Else                 *Schema  `json:"else,omitempty"`
	MinProperties        *int     `json:"minProperties,omitempty"`
	MaxProperties        *int     `json:"maxProperties,omitempty"`
	Required             []string `json:"required,omitempty"`
	Properties           Schemas  `json:"properties,omitempty"`
	PropertyNames        *Schema  `json:"propertyNames,omitempty"`
	RegexProperties      *bool    `json:"regexProperties,omitempty"`
	PatternProperties    Schemas  `json:"patternProperties,omitempty"`
	AdditionalProperties *Schema  `json:"additionalProperties,omitempty"`
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
	UnevaluatedProperties *Schema `json:"unevaluatedProperties,omitempty"`
	UniqueObjs            *bool   `json:"uniqueObjs,omitempty"`
	// List validation is useful for arrays of arbitrary length where each item
	// matches the same schema. For this kind of array, set the items keyword to
	// a single schema that will be used to validate all of the items in the
	// array.
	Items            *Schema           `json:"items,omitempty"`
	UnevaluatedObjs  *Schema           `json:"unevaluatedObjs,omitempty"`
	AdditionalObjs   *Schema           `json:"additionalObjs,omitempty"`
	PrefixObjs       SchemaSet         `json:"prefixObjs,omitempty"`
	Contains         *Schema           `json:"contains,omitempty"`
	MinContains      *Number           `json:"minContains,omitempty"`
	MaxContains      *Number           `json:"maxContains,omitempty"`
	MinLength        *Number           `json:"minLength,omitempty"`
	MaxLength        *Number           `json:"maxLength,omitempty"`
	Pattern          *Regexp           `json:"pattern,omitempty"`
	ContentEncoding  string            `json:"contentEncoding,omitempty"`
	ContentMediaType string            `json:"contentMediaType,omitempty"`
	Minimum          *Number           `json:"minimum,omitempty"`
	ExclusiveMinimum *Number           `json:"exclusiveMinimum,omitempty"`
	Maximum          *Number           `json:"maximum,omitempty"`
	ExclusiveMaximum *Number           `json:"exclusiveMaximum,omitempty"`
	MultipleOf       *Number           `json:"multipleOf,omitempty"`
	Title            string            `json:"title,omitempty"`
	Description      string            `json:"description,omitempty"`
	Default          json.RawMessage   `json:"default,omitempty"`
	ReadOnly         *bool             `json:"readOnly,omitempty"`
	WriteOnly        *bool             `json:"writeOnly,omitempty"`
	Examples         []json.RawMessage `json:"examples,omitempty"`
	Example          json.RawMessage   `json:"example,omitempty"`
	Deprecated       *bool             `json:"deprecated,omitempty"`
	ExternalDocs     string            `json:"externalDocs,omitempty"`
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

type schema Schema

// Detail returns a ptr to the *SchemaObj
func (s Schema) Detail() *Schema {
	return &s
}

// MarshalJSON marshals JSON
func (s Schema) MarshalJSON() ([]byte, error) {
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

// UnmarshalJSON unmarshals JSON
func (s *Schema) UnmarshalJSON(data []byte) error {
	sv, err := unmarshalSchemaJSON(data)
	*s = *sv
	return err
}

// IsStrings returns false
func (s *Schema) IsStrings() bool {
	return false
}

// IsBool returns false
func (s *Schema) IsBool() bool {
	return false
}

// IsRef returns true if s.Ref is set
func (s *Schema) IsRef() bool {
	return s.Ref != ""
}

// SetKeyword encodes and sets the keyword key to the encoded value
func (s *Schema) SetKeyword(key string, value interface{}) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.SetEncodedKeyword(key, b)
}

// SetEncodedKeyword sets the keyword key to value
func (s *Schema) SetEncodedKeyword(key string, value []byte) error {
	if strings.HasPrefix(key, "x-") {
		return errors.New("keyword keys may not start with \"x-\"")
	}
	s.Keywords[key] = value
	return nil
}

// DecodeKeyword unmarshals the keyword's raw data into dst
func (s *Schema) DecodeKeyword(key string, dst interface{}) error {
	return json.Unmarshal(s.Keywords[key], dst)
}

// DecodeKeywords unmarshals all keywords raw data into dst
func (s *Schema) DecodeKeywords(dst interface{}) error {
	data, err := json.Marshal(s.Keywords)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}

// SchemaSet is a slice of **SchemaObj
type SchemaSet []*Schema

// UnmarshalJSON unmarshals JSON
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

func unmarshalSchemaJSON(data []byte) (*Schema, error) {
	var str string
	l := len(data)
	if l >= 4 && l <= 5 {
		str = string(data)
	}
	switch {
	case str == "true":
		t := true
		return &Schema{Always: &t}, nil
	case str == "false":
		f := false
		return &Schema{Always: &f}, nil
	default:
		return unmarshalSchemaObjJSON(data)
	}
}

func unmarshalSchemaObjJSON(data []byte) (*Schema, error) {
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
			var v *Schema
			v, err = unmarshalSchemaJSON(d)
			if err != nil {
				return nil, err
			}
			set(&dst, v)
		} else if _, isfield := jsfields[key]; !isfield {
			kw[key] = d
		}
	}
	res := Schema(dst)
	res.Keywords = kw
	res.Extensions = exts
	return &res, err
}

var schemaFieldSetters = map[string]func(s *partialschema, v *Schema){
	"not":                   func(s *partialschema, v *Schema) { s.Not = v },
	"if":                    func(s *partialschema, v *Schema) { s.If = v },
	"then":                  func(s *partialschema, v *Schema) { s.Then = v },
	"else":                  func(s *partialschema, v *Schema) { s.Else = v },
	"propertyNames":         func(s *partialschema, v *Schema) { s.PropertyNames = v },
	"additionalProperties":  func(s *partialschema, v *Schema) { s.AdditionalProperties = v },
	"unevaluatedProperties": func(s *partialschema, v *Schema) { s.UnevaluatedProperties = v },
	"items":                 func(s *partialschema, v *Schema) { s.Items = v },
	"contains":              func(s *partialschema, v *Schema) { s.Contains = v },
	"unevaluatedObjs":       func(s *partialschema, v *Schema) { s.UnevaluatedObjs = v },
	"additionalObjs":        func(s *partialschema, v *Schema) { s.AdditionalObjs = v },
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
	RefResolved           *Schema             `json:"-"`
	RecursiveRefResolved  *Schema             `json:"-"`
	DynamicRefResolved    *Schema             `json:"-"`
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
	Comments              string              `json:"$comment,omitempty"`
	Not                   *Schema             `json:"-"`
	AllOf                 SchemaSet           `json:"allOf,omitempty"`
	AnyOf                 SchemaSet           `json:"anyOf,omitempty"`
	OneOf                 SchemaSet           `json:"oneOf,omitempty"`
	If                    *Schema             `json:"-"`
	Then                  *Schema             `json:"-"`
	Else                  *Schema             `json:"-"`
	MinProperties         *int                `json:"minProperties,omitempty"`
	MaxProperties         *int                `json:"maxProperties,omitempty"`
	Required              []string            `json:"required,omitempty"`
	Properties            Schemas             `json:"properties,omitempty"`
	PropertyNames         *Schema             `json:"-"`
	RegexProperties       *bool               `json:"regexProperties,omitempty"`
	PatternProperties     Schemas             `json:"patternProperties,omitempty"`
	AdditionalProperties  *Schema             `json:"-"`
	DependentRequired     map[string][]string `json:"dependentRequired,omitempty"`
	DependentSchemas      Schemas             `json:"dependentSchemas,omitempty"`
	UnevaluatedProperties *Schema             `json:"-"`
	UniqueObjs            *bool               `json:"uniqueObjs,omitempty"`
	Items                 *Schema             `json:"-"`
	UnevaluatedObjs       *Schema             `json:"-"`
	AdditionalObjs        *Schema             `json:"-"`
	PrefixObjs            SchemaSet           `json:"prefixObjs,omitempty"`
	Contains              *Schema             `json:"-"`
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
	Examples              []json.RawMessage   `json:"examples,omitempty"`
	Example               json.RawMessage     `json:"example,omitempty"`
	Deprecated            *bool               `json:"deprecated,omitempty"`
	ExternalDocs          string              `json:"externalDocs,omitempty"`
	RecursiveAnchor       *bool               `json:"$recursiveAnchor,omitempty"`
	RecursiveRef          string              `json:"$recursiveRef,omitempty"`
	Discriminator         *Discriminator      `json:"discriminator,omitempty"`
	XML                   *XML                `json:"xml,omitempty"`
	Extensions            `json:"-"`
	Keywords              map[string]json.RawMessage `json:"-"`
}

var (
	_ json.Marshaler   = (*Schema)(nil)
	_ json.Unmarshaler = (*Schema)(nil)
)
