package openapi

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/chanced/jay"
	"github.com/chanced/jsonpointer"
	"github.com/chanced/uri"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type SchemaEntry struct {
	Key    string
	Schema *Schema
}

// SchemaMap is a psuedo, ordered map of Schemas
//
// Under the hood, SchemaMap is a slice of SchemaEntry
type SchemaMap []SchemaEntry

func (sm *SchemaMap) Set(key string, s *Schema) {
	se := SchemaEntry{
		Key:    key,
		Schema: s,
	}
	for i, v := range *sm {
		if v.Key == key {
			(*sm)[i] = se
			return
		}
	}
	(*sm) = append((*sm), se)
}

func (sm SchemaMap) setLocation(loc Location) error {
	if sm == nil {
		return nil
	}
	for _, e := range sm {
		err := e.Schema.setLocation(loc.Append(e.Key))
		if err != nil {
			return err
		}
	}
	return nil
}

func (sm SchemaMap) Get(key string) *Schema {
	for _, v := range sm {
		if v.Key == key {
			return v.Schema
		}
	}
	return nil
}

func (sm *SchemaMap) MarshalJSON() ([]byte, error) {
	b := []byte("{}")
	var err error
	for _, v := range *sm {
		b, err = json.Marshal(v.Schema)
		if err != nil {
			return b, err
		}
		b, err = sjson.SetBytes(b, v.Key, b)
		if err != nil {
			return b, err
		}
	}
	return b, nil
}

func (sm *SchemaMap) UnmarshalJSON(data []byte) error {
	t := jay.TypeOf(data)
	if t != jay.TypeObject {
		return &json.UnmarshalTypeError{Value: t.String(), Type: reflect.TypeOf(sm)}
	}
	*sm = make(SchemaMap, 0)
	var err error
	gjson.ParseBytes(data).ForEach(func(key, value gjson.Result) bool {
		var s Schema
		err = json.Unmarshal([]byte(value.Raw), &s)
		*sm = append(*sm, SchemaEntry{Key: key.String(), Schema: &s})
		return err == nil
	})
	return err
}

type SchemaSet []*Schema

func (ss SchemaSet) setLocation(loc Location) error {
	if ss == nil {
		return nil
	}
	for i, s := range ss {
		err := s.setLocation(loc.Append(strconv.Itoa(i)))
		if err != nil {
			return err
		}
	}
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
	// Always will be assigned if the schema value is a boolean
	Always *bool `json:"-"`

	Schema *uri.URI `json:"$schema,omitempty"`
	// The value of $id is a URI-reference without a fragment that resolves
	// against the Retrieval URI. The resulting URI is the base URI for the
	// schema.
	//
	// https://json-schema.org/understanding-json-schema/structuring.html?highlight=id#id
	ID *uri.URI `json:"$id,omitempty"`
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
	Ref *SchemaRef `json:"$ref,omitempty"`
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
	Definitions SchemaMap `json:"$defs,omitempty"`
	// The format keyword allows for basic semantic identification of certain Kinds of string values that are commonly used. For example, because JSON doesn’t have a “DateTime” type, dates need to be encoded as strings. format allows the schema author to indicate that the string value should be interpreted as a date. By default, format is just an annotation and does not effect validation.
	//
	// Optionally, validator implementations can provide a configuration option to
	// enable format to function as an assertion rather than just an annotation.
	// That means that validation will fail if, for example, a value with a date
	// format isn’t in a form that can be parsed as a date. This can allow values to
	// be constrained beyond what the other tools in JSON *SchemaObj, including Regular
	// Expressions can do.
	//
	// https://json-schema.org/understanding-json-schema/reference/string.html#format
	Format        Text `json:"format,omitempty"`
	DynamicAnchor Text `json:"$dynamicAnchor,omitempty"`
	// The "$dynamicRef" keyword is an applicator that allows for deferring the
	// full resolution until runtime, at which point it is resolved each time it
	// is encountered while evaluating an instance.
	//
	// https://json-schema.org/draft/2020-12/json-schema-core.html#dynamic-ref
	DynamicRef *SchemaRef `json:"$dynamicRef,omitempty"`
	// A less common way to identify a subschema is to create a named anchor in
	// the schema using the $anchor keyword and using that name in the URI
	// fragment. Anchors must start with a letter followed by any number of
	// letters, digits, -, _, :, or ..
	//
	// https://json-schema.org/understanding-json-schema/structuring.html?highlight=anchor#anchor
	Anchor Text `json:"$anchor,omitempty"`
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
	Comments Text `json:"$comment,omitempty"`

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
	Else                 *Schema   `json:"else,omitempty"`
	MinProperties        *Number   `json:"minProperties,omitempty"`
	MaxProperties        *Number   `json:"maxProperties,omitempty"`
	Required             []string  `json:"required,omitempty"`
	Properties           SchemaMap `json:"properties,omitempty"`
	PropertyNames        *Schema   `json:"propertyNames,omitempty"`
	RegexProperties      *bool     `json:"regexProperties,omitempty"`
	PatternProperties    SchemaMap `json:"patternProperties,omitempty"`
	AdditionalProperties *Schema   `json:"additionalProperties,omitempty"`
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
	DependentSchemas      SchemaMap `json:"dependentSchemas,omitempty"`
	UnevaluatedProperties *Schema   `json:"unevaluatedProperties,omitempty"`
	UniqueObjs            *bool     `json:"uniqueObjs,omitempty"`
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
	ContentEncoding  Text              `json:"contentEncoding,omitempty"`
	ContentMediaType Text              `json:"contentMediaType,omitempty"`
	Minimum          *Number           `json:"minimum,omitempty"`
	ExclusiveMinimum *Number           `json:"exclusiveMinimum,omitempty"`
	Maximum          *Number           `json:"maximum,omitempty"`
	ExclusiveMaximum *Number           `json:"exclusiveMaximum,omitempty"`
	MultipleOf       *Number           `json:"multipleOf,omitempty"`
	Title            Text              `json:"title,omitempty"`
	Description      Text              `json:"description,omitempty"`
	Default          json.RawMessage   `json:"default,omitempty"`
	ReadOnly         *bool             `json:"readOnly,omitempty"`
	WriteOnly        *bool             `json:"writeOnly,omitempty"`
	Examples         []json.RawMessage `json:"examples,omitempty"`
	Example          json.RawMessage   `json:"example,omitempty"`
	Deprecated       *bool             `json:"deprecated,omitempty"`
	ExternalDocs     Text              `json:"externalDocs,omitempty"`
	// Deprecated: renamed to dynamicAnchor
	RecursiveAnchor *bool `json:"$recursiveAnchor,omitempty"`
	// Deprecated: renamed to dynamicRef
	RecursiveRef *SchemaRef `json:"$recursiveRef,omitempty"`

	Discriminator *Discriminator `json:"discriminator,omitempty"`
	// This MAY be used only on properties schemas. It has no effect on root
	// schemas. Adds additional metadata to describe the XML representation of
	// this property.
	XML        *XML `json:"xml,omitempty"`
	Extensions `json:"-"`
	Keywords   map[string]json.RawMessage `json:"-"`
	Location   *Location                  `json:"-"`
}

// MarshalJSON marshals JSON
func (s Schema) MarshalJSON() ([]byte, error) {
	type schema Schema
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
	t := jay.TypeOf(data)
	switch t {
	case jay.TypeBool:
		return s.unmarshalJSONBool(data)
	case jay.TypeObject:
		return s.unmarshalJSONObj(data)
	default:
		return &json.UnmarshalTypeError{Value: t.String(), Type: reflect.TypeOf(s)}
	}
}

func (s *Schema) unmarshalJSONBool(data []byte) error {
	var b bool
	err := json.Unmarshal(data, &b)
	*s = Schema{Always: &b}
	return err
}

func (s *Schema) unmarshalJSONObj(data []byte) error {
	res := Schema{
		Extensions: make(Extensions),
		Keywords:   make(map[string]json.RawMessage),
	}

	d := map[string]json.RawMessage{}
	err := json.Unmarshal(data, &d)
	if err != nil {
		return err
	}
	fields := res.fields()
	for k, v := range d {
		if f, ok := fields[k]; ok {
			err = json.Unmarshal(v, f)
			if err != nil {
				return err
			}
		} else if strings.HasPrefix(k, "x-") {
			res.Extensions[k] = v
		} else {
			res.Keywords[k] = v
		}
	}
	return nil
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
	return s.Ref != nil
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

func (s *Schema) fields() map[string]interface{} {
	return map[string]interface{}{
		"$schema":               &s.Schema,
		"$id":                   &s.ID,
		"type":                  &s.Type,
		"$ref":                  &s.Ref,
		"$defs":                 &s.Definitions,
		"format":                &s.Format,
		"$dynamicAnchor":        &s.DynamicAnchor,
		"$dynamicRef":           &s.DynamicRef,
		"$anchor":               &s.Anchor,
		"const":                 &s.Const,
		"enum":                  &s.Enum,
		"$comment":              &s.Comments,
		"not":                   &s.Not,
		"allOf":                 &s.AllOf,
		"anyOf":                 &s.AnyOf,
		"oneOf":                 &s.OneOf,
		"if":                    &s.If,
		"then":                  &s.Then,
		"else":                  &s.Else,
		"minProperties":         &s.MinProperties,
		"maxProperties":         &s.MaxProperties,
		"required":              &s.Required,
		"properties":            &s.Properties,
		"propertyNames":         &s.PropertyNames,
		"regexProperties":       &s.RegexProperties,
		"patternProperties":     &s.PatternProperties,
		"additionalProperties":  &s.AdditionalProperties,
		"dependentRequired":     &s.DependentRequired,
		"dependentSchemas":      &s.DependentSchemas,
		"unevaluatedProperties": &s.UnevaluatedProperties,
		"uniqueObjs":            &s.UniqueObjs,
		"items":                 &s.Items,
		"unevaluatedObjs":       &s.UnevaluatedObjs,
		"additionalObjs":        &s.AdditionalObjs,
		"prefixObjs":            &s.PrefixObjs,
		"contains":              &s.Contains,
		"minContains":           &s.MinContains,
		"maxContains":           &s.MaxContains,
		"minLength":             &s.MinLength,
		"maxLength":             &s.MaxLength,
		"pattern":               &s.Pattern,
		"contentEncoding":       &s.ContentEncoding,
		"contentMediaType":      &s.ContentMediaType,
		"minimum":               &s.Minimum,
		"exclusiveMinimum":      &s.ExclusiveMinimum,
		"maximum":               &s.Maximum,
		"exclusiveMaximum":      &s.ExclusiveMaximum,
		"multipleOf":            &s.MultipleOf,
		"title":                 &s.Title,
		"description":           &s.Description,
		"default":               &s.Default,
		"readOnly":              &s.ReadOnly,
		"writeOnly":             &s.WriteOnly,
		"examples":              &s.Examples,
		"example":               &s.Example,
		"deprecated":            &s.Deprecated,
		"externalDocs":          &s.ExternalDocs,
		"$recursiveAnchor":      &s.RecursiveAnchor,
		"$recursiveRef":         &s.RecursiveRef,
		"discriminator":         &s.Discriminator,
		"xml":                   &s.XML,
	}
}

type SchemaRef struct {
	Ref      *uri.URI `json:"-"`
	Resolved *Schema  `json:"-"`
}

func (sr *SchemaRef) setLocation(l Location) error {
	if sr == nil {
		return nil
	}
	if sr.Resolved != nil {
		if sr.Ref != nil {
			nl, err := NewLocation(sr.Ref)
			if err != nil {
				return err
			}
			sr.Resolved.setLocation(nl)
			return nil
		}
		return sr.Resolved.setLocation(l)
	}
	return nil
}

func (sr *SchemaRef) UnmarshalJSON(data []byte) error {
	if jay.IsString(data) {
		var u uri.URI
		if err := json.Unmarshal(data, &u); err != nil {
			return err
		}
		sr.Ref = &u
		return nil
	}

	var s Schema
	err := json.Unmarshal(data, &s)
	sr.Resolved = &s
	return err
}

func (sr *SchemaRef) MarshalJSON() ([]byte, error) {
	return json.Marshal(sr.Ref)
}

// init implements node
func (*Schema) init(ctx context.Context, resolver *resolver) error {
	panic("unimplemented")
}

// kind implements node
func (*Schema) Kind() Kind { return KindSchema }

// resolve implements node
func (*Schema) resolve(ctx context.Context, resolver *resolver, p jsonpointer.Pointer) (node, error) {
	panic("unimplemented")
}

// setLocation implements node
func (s *Schema) setLocation(loc Location) error {
	if s.Location != nil {
		return nil
	}
	s.Location = &loc

	if err := s.Ref.setLocation(loc.Append("ref")); err != nil {
		return err
	}
	if err := s.Definitions.setLocation(loc.Append("definitions")); err != nil {
		return err
	}
	if err := s.DynamicRef.setLocation(loc.Append("dynamicRef")); err != nil {
		return err
	}
	if err := s.Not.setLocation(loc.Append("not")); err != nil {
		return err
	}
	if err := s.AllOf.setLocation(loc.Append("allOf")); err != nil {
		return err
	}
	if err := s.AnyOf.setLocation(loc.Append("anyOf")); err != nil {
		return err
	}
	if err := s.OneOf.setLocation(loc.Append("oneOf")); err != nil {
		return err
	}
	if err := s.If.setLocation(loc.Append("if")); err != nil {
		return err
	}
	if err := s.Then.setLocation(loc.Append("then")); err != nil {
		return err
	}
	if err := s.Else.setLocation(loc.Append("else")); err != nil {
		return err
	}
	if err := s.Properties.setLocation(loc.Append("properties")); err != nil {
		return err
	}
	if err := s.PropertyNames.setLocation(loc.Append("propertyNames")); err != nil {
		return err
	}
	if err := s.PatternProperties.setLocation(loc.Append("patternProperties")); err != nil {
		return err
	}
	if err := s.AdditionalProperties.setLocation(loc.Append("additionalProperties")); err != nil {
		return err
	}
	if err := s.DependentSchemas.setLocation(loc.Append("dependentSchemas")); err != nil {
		return err
	}
	if err := s.UnevaluatedProperties.setLocation(loc.Append("unevaluatedProperties")); err != nil {
		return err
	}
	if err := s.Items.setLocation(loc.Append("items")); err != nil {
		return err
	}
	if err := s.UnevaluatedObjs.setLocation(loc.Append("unevaluatedObjs")); err != nil {
		return err
	}
	if err := s.AdditionalObjs.setLocation(loc.Append("additionalObjs")); err != nil {
		return err
	}
	if err := s.PrefixObjs.setLocation(loc.Append("prefixObjs")); err != nil {
		return err
	}
	if err := s.Contains.setLocation(loc.Append("contains")); err != nil {
		return err
	}
	if err := s.RecursiveRef.setLocation(loc.Append("recursiveRef")); err != nil {
		return err
	}
	if err := s.Discriminator.setLocation(loc.Append("discriminator")); err != nil {
		return err
	}
	if err := s.XML.setLocation(loc.Append("xml")); err != nil {
		return err
	}
	return nil
}

var _ node = (*Schema)(nil)
