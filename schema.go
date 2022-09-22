package openapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"strings"

	"github.com/chanced/caps/text"
	"github.com/chanced/jsonx"
	"github.com/chanced/maps"
	"github.com/chanced/uri"
)

// Schema allows the definition of input and output data types. These types can
// be objects, but also primitives and arrays. This object is a superSlice of the
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
	Extensions `json:"-"`
	Location   `json:"-"`

	Schema *uri.URI `json:"$schema,omitempty"`

	// The value of $id is a URI-reference without a fragment that resolves
	// against the Retrieval URI. The resulting URI is the base URI for the
	// schema.
	//
	// https://json-schema.org/understanding-json-schema/structuring.html?highlight=id#id
	ID *uri.URI `json:"$id,omitempty"`

	// A less common way to identify a subschema is to create a named anchor in
	// the schema using the $anchor keyword and using that name in the URI
	// fragment. Anchors must start with a letter followed by any number of
	// letters, digits, -, _, :, or ..
	//
	// https://json-schema.org/understanding-json-schema/structuring.html?highlight=anchor#anchor
	Anchor Text `json:"$anchor,omitempty"`

	DynamicAnchor Text `json:"$dynamicAnchor,omitempty"`

	RecursiveAnchor *bool `json:"$recursiveAnchor,omitempty"`

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

	// The "$dynamicRef" keyword is an applicator that allows for deferring the
	// full resolution until runtime, at which point it is resolved each time it
	// is encountered while evaluating an instance.
	//
	// https://json-schema.org/draft/2020-12/json-schema-core.html#dynamic-ref
	DynamicRef *SchemaRef `json:"$dynamicRef,omitempty"`

	RecursiveRef *SchemaRef `json:"$recursiveRef,omitempty"`

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
	Format Text `json:"format,omitempty"`

	// The const keyword is used to restrict a value to a single value.
	//
	// https://json-schema.org/understanding-json-schema/reference/generic.html?highlight=const#constant-values
	Const jsonx.RawMessage `json:"const,omitempty"`

	Required Texts `json:"required,omitempty"`

	Properties *SchemaMap `json:"properties,omitempty"`

	// The enum keyword is used to restrict a value to a fixed set of values. It
	// must be an array with at least one element, where each element is unique.
	//
	// https://json-schema.org/understanding-json-schema/reference/generic.html?highlight=const#enumerated-values
	Enum Texts `json:"enum,omitempty"`

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
	AllOf *SchemaSlice `json:"allOf,omitempty"`

	// validate against anyOf, the given data must be valid against any (one or
	// more) of the given subschemas.
	//
	// https://json-schema.org/understanding-json-schema/reference/combining.html?highlight=allof#allof
	AnyOf *SchemaSlice `json:"anyOf,omitempty"`

	// alidate against oneOf, the given data must be valid against exactly one of the given subschemas.
	//
	// https://json-schema.org/understanding-json-schema/reference/combining.html?highlight=oneof#oneof
	OneOf *SchemaSlice `json:"oneOf,omitempty"`

	// if, then and else keywords allow the application of a subschema based on
	// the outcome of another schema, much like the if/then/else constructs
	// you’ve probably seen in traditional programming languages.
	//
	// https://json-schema.org/understanding-json-schema/reference/conditionals.html#if-then-else
	If *Schema `json:"if,omitempty"`

	// https://json-schema.org/understanding-json-schema/reference/conditionals.html#if-then-else
	Then *Schema `json:"then,omitempty"`

	// https://json-schema.org/understanding-json-schema/reference/conditionals.html#if-then-else

	Else *Schema `json:"else,omitempty"`

	MinProperties *Number `json:"minProperties,omitempty"`

	MaxProperties *Number `json:"maxProperties,omitempty"`

	PropertyNames *Schema `json:"propertyNames,omitempty"`

	RegexProperties *bool `json:"regexProperties,omitempty"`

	PatternProperties *SchemaMap `json:"patternProperties,omitempty"`

	AdditionalProperties *Schema `json:"additionalProperties,omitempty"`

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
	DependentRequired *Map[Texts] `json:"dependentRequired,omitempty"`

	// The dependentSchemas keyword conditionally applies a subschema when a
	// given property is present. This schema is applied in the same way allOf
	// applies schemas. Nothing is merged or extended. Both schemas apply
	// independently.

	DependentSchemas *SchemaMap `json:"dependentSchemas,omitempty"`

	UnevaluatedProperties *Schema `json:"unevaluatedProperties,omitempty"`

	UniqueItems *bool `json:"uniqueItems,omitempty"`

	// List validation is useful for arrays of arbitrary length where each item
	// matches the same schema. For this kind of array, set the items keyword to
	// a single schema that will be used to validate all of the items in the
	// array.
	//
	// https://json-schema.org/understanding-json-schema/reference/array.html#items
	Items *Schema `json:"items,omitempty"`

	UnevaluatedItems *Schema `json:"unevaluatedItems,omitempty"`

	AdditionalItems *Schema `json:"additionalItems,omitempty"`

	PrefixItems *SchemaSlice `json:"prefixItems,omitempty"`

	Contains *Schema `json:"contains,omitempty"`

	MinContains *Number `json:"minContains,omitempty"`

	MaxContains *Number `json:"maxContains,omitempty"`

	MinLength *Number `json:"minLength,omitempty"`

	MaxLength *Number `json:"maxLength,omitempty"`

	Pattern *Regexp `json:"pattern,omitempty"`

	ContentEncoding Text `json:"contentEncoding,omitempty"`

	ContentMediaType Text `json:"contentMediaType,omitempty"`

	Minimum *Number `json:"minimum,omitempty"`

	ExclusiveMinimum *Number `json:"exclusiveMinimum,omitempty"`

	Maximum *Number `json:"maximum,omitempty"`

	ExclusiveMaximum *Number `json:"exclusiveMaximum,omitempty"`

	MultipleOf *Number `json:"multipleOf,omitempty"`

	Title Text `json:"title,omitempty"`

	Description Text `json:"description,omitempty"`

	Default jsonx.RawMessage `json:"default,omitempty"`

	ReadOnly *bool `json:"readOnly,omitempty"`

	WriteOnly *bool `json:"writeOnly,omitempty"`

	Examples []jsonx.RawMessage `json:"examples,omitempty"`

	Example jsonx.RawMessage `json:"example,omitempty"`

	Deprecated *bool `json:"deprecated,omitempty"`

	ExternalDocs Text `json:"externalDocs,omitempty"`

	// When request bodies or response payloads may be one of a number of
	// different schemas, a discriminator object can be used to aid in
	// serialization, deserialization, and validation. The discriminator is a
	// specific object in a schema which is used to inform the consumer of the
	// document of an alternative schema based on the value associated with it.
	//
	// This object MAY be extended with Specification Extensions.
	//
	// The discriminator object is legal only when using one of the composite
	// keywords oneOf, anyOf, allOf.
	//
	// 3.1:
	//
	// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.0.md#discriminatorObject
	//
	// 3.0:
	//
	// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.2.md#discriminatorObject
	Discriminator *Discriminator `json:"discriminator,omitempty"`

	// This MAY be used only on properties schemas. It has no effect on root
	// schemas. Adds additional metadata to describe the XML representation of
	// this property.
	XML *XML `json:"xml,omitempty"`

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
	Definitions *SchemaMap `json:"$defs,omitempty"`

	Keywords map[Text]jsonx.RawMessage `json:"-"`
}

func (s *Schema) Nodes() []Node {
	if s == nil {
		return nil
	}
	return downcastNodes(s.nodes())
}

func (s *Schema) nodes() []node {
	return appendEdges(nil, s.Ref,
		s.DynamicRef,
		s.RecursiveRef,
		s.Properties,
		s.Not,
		s.AllOf,
		s.AnyOf,
		s.OneOf,
		s.If,
		s.Then,
		s.Else,
		s.PropertyNames,
		s.PatternProperties,
		s.AdditionalProperties,
		s.DependentSchemas,
		s.UnevaluatedProperties,
		s.Items,
		s.UnevaluatedItems,
		s.AdditionalItems,
		s.PrefixItems,
		s.Contains,
		s.Discriminator,
		s.XML,
		s.Definitions,
	)
}

func (s *Schema) Refs() []Ref {
	if s == nil {
		return nil
	}
	var refs []Ref
	if s.Ref != nil {
		refs = append(refs, s.Ref)
	}
	if s.DynamicRef != nil {
		refs = append(refs, s.DynamicRef)
	}
	if s.RecursiveRef != nil {
		refs = append(refs, s.RecursiveRef)
	}
	refs = append(refs, s.Definitions.Refs()...)
	refs = append(refs, s.Not.Refs()...)
	refs = append(refs, s.AllOf.Refs()...)
	refs = append(refs, s.AnyOf.Refs()...)
	refs = append(refs, s.OneOf.Refs()...)
	refs = append(refs, s.If.Refs()...)
	refs = append(refs, s.Then.Refs()...)
	refs = append(refs, s.Else.Refs()...)
	refs = append(refs, s.Properties.Refs()...)
	refs = append(refs, s.PropertyNames.Refs()...)
	refs = append(refs, s.PatternProperties.Refs()...)
	refs = append(refs, s.AdditionalProperties.Refs()...)
	refs = append(refs, s.DependentSchemas.Refs()...)
	refs = append(refs, s.UnevaluatedProperties.Refs()...)
	refs = append(refs, s.Items.Refs()...)
	refs = append(refs, s.UnevaluatedItems.Refs()...)
	refs = append(refs, s.AdditionalItems.Refs()...)
	refs = append(refs, s.PrefixItems.Refs()...)
	refs = append(refs, s.Contains.Refs()...)
	refs = append(refs, s.XML.Refs()...)

	return refs
}

func (s *Schema) Anchors() (*Anchors, error) {
	if s == nil {
		return nil, nil
	}
	anchors := &Anchors{
		Standard: make(map[text.Text]Anchor),
		Dynamic:  make(map[text.Text]Anchor),
	}
	if s.Anchor != "" {
		anchors.Standard[s.Anchor] = Anchor{
			Location: s.Location.AppendLocation("$anchor"),
			In:       s,
			Name:     s.Anchor,
			Type:     AnchorTypeRegular,
		}
	}
	if s.DynamicAnchor != "" {
		anchors.Dynamic[s.DynamicAnchor] = Anchor{
			Location: s.Location.AppendLocation("$dynamicAnchor"),
			In:       s,
			Name:     s.DynamicAnchor,
			Type:     AnchorTypeDynamic,
		}
	}
	if s.RecursiveAnchor != nil {
		anchors.Recursive = &Anchor{
			Location: s.Location.AppendLocation("$recursiveAnchor"),
			In:       s,
			Name:     "",
			Type:     AnchorTypeRecursive,
		}
	}
	var err error

	if anchors, err = anchors.merge(s.Ref.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(s.Definitions.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(s.DynamicRef.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(s.Not.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(s.AllOf.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(s.AnyOf.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(s.OneOf.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(s.If.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(s.Then.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(s.Else.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(s.Properties.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(s.PropertyNames.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(s.PatternProperties.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(s.AdditionalProperties.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(s.DependentSchemas.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(s.UnevaluatedProperties.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(s.Items.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(s.UnevaluatedItems.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(s.AdditionalItems.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(s.PrefixItems.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(s.Contains.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(s.RecursiveRef.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(s.XML.Anchors()); err != nil {
		return nil, err
	}

	return anchors, nil
}

// MarshalJSON marshals JSON
func (s Schema) MarshalJSON() ([]byte, error) {
	type schema Schema
	b := bytes.Buffer{}
	data, err := json.Marshal(schema(s))
	if err != nil {
		return nil, err
	}
	// trimming the last }
	b.Write(data[:len(data)-1])

	if len(s.Keywords) == 0 && len(s.Extensions) == 0 && b.Len() < 10 {
		bs := b.String()
		switch bs {
		case "{":
			return []byte("true"), nil
		case `{"not":true`:
			return []byte("false"), nil
		}
	}
	if s.Keywords != nil {
		for _, kv := range maps.SortByKeys(s.Keywords) {
			if b.Len() > 2 {
				b.WriteString(",")
			}
			jsonx.EncodeAndWriteString(&b, kv.Key)
			b.WriteByte(':')
			if kv.Value != nil {
				bb, err := json.Marshal(kv.Value)
				if err != nil {
					return nil, err
				}
				b.Write(bb)
			}
		}
	}
	b.WriteByte('}')
	return b.Bytes(), err
}

// UnmarshalJSON unmarshals JSON
func (s *Schema) UnmarshalJSON(data []byte) error {
	t := jsonx.TypeOf(data)
	switch t {
	case jsonx.TypeBool:
		return s.unmarshalJSONBool(data)
	case jsonx.TypeObject:
		return s.unmarshalJSONObj(data)
	default:
		return &json.UnmarshalTypeError{Value: t.String(), Type: reflect.TypeOf(s)}
	}
}

func (s *Schema) unmarshalJSONBool(data []byte) error {
	if jsonx.IsTrue(data) {
		*s = Schema{}
		return nil
	} else {
		*s = Schema{Not: &Schema{}}
		return nil
	}
}

func (s *Schema) unmarshalJSONObj(data []byte) error {
	res := Schema{}

	d := map[Text]jsonx.RawMessage{}
	err := json.Unmarshal(data, &d)
	if err != nil {
		return err
	}
	fields := res.fields()
	for k, v := range d {
		if f, ok := fields[k.String()]; ok {
			err = json.Unmarshal(v, f)
			if err != nil {
				return err
			}
		} else if strings.HasPrefix(k.String(), "x-") {
			if res.Extensions == nil {
				res.Extensions = Extensions{}
			}
			res.Extensions[k] = v
		} else {
			if res.Keywords == nil {
				res.Keywords = make(map[Text]jsonx.RawMessage)
			}
			res.Keywords[k] = v
		}
	}
	if err != nil {
		return err
	}
	if res.Ref != nil {
		res.Ref.SchemaRefKind = SchemaRefTypeRef
	}
	if res.DynamicRef != nil {
		res.DynamicRef.SchemaRefKind = SchemaRefTypeDynamic
	}
	if res.RecursiveRef != nil {
		res.RecursiveRef.SchemaRefKind = SchemaRefTypeRecursive
	}
	*s = res

	return nil
}

// SetKeyword marshals value and sets the encoded json to key in Keywords
//
// If setting the value as []byte, it should be in the form of json.RawMessage
// or jsonx.RawMessage as both types implement json.Marshaler
func (s *Schema) SetKeyword(key Text, value interface{}) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.setEncodedKeyword(key, b)
}

// SetEncodedKeyword sets the keyword key to value
func (s *Schema) setEncodedKeyword(key Text, value []byte) error {
	if key.HasPrefix("x-") {
		return errors.New("keyword keys may not start with \"x-\"")
	}
	s.Keywords[Text(key)] = value
	return nil
}

// DecodeKeyword unmarshals the keyword's raw data into dst
func (s *Schema) DecodeKeyword(key Text, dst interface{}) error {
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
		"uniqueItems":           &s.UniqueItems,
		"items":                 &s.Items,
		"unevaluatedItems":      &s.UnevaluatedItems,
		"additionalItems":       &s.AdditionalItems,
		"prefixItems":           &s.PrefixItems,
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

func (*Schema) Kind() Kind      { return KindSchema }
func (*Schema) mapKind() Kind   { return KindSchemaMap }
func (*Schema) sliceKind() Kind { return KindSchemaSlice }

func (s *Schema) setLocation(loc Location) error {
	if s == nil {
		return nil
	}
	s.Location = loc

	if err := s.Ref.setLocation(loc.AppendLocation("$ref")); err != nil {
		return err
	}
	if err := s.Definitions.setLocation(loc.AppendLocation("$defs")); err != nil {
		return err
	}
	if err := s.DynamicRef.setLocation(loc.AppendLocation("$dynamicRef")); err != nil {
		return err
	}
	if err := s.Not.setLocation(loc.AppendLocation("not")); err != nil {
		return err
	}
	if err := s.AllOf.setLocation(loc.AppendLocation("allOf")); err != nil {
		return err
	}
	if err := s.AnyOf.setLocation(loc.AppendLocation("anyOf")); err != nil {
		return err
	}
	if err := s.OneOf.setLocation(loc.AppendLocation("oneOf")); err != nil {
		return err
	}
	if err := s.If.setLocation(loc.AppendLocation("if")); err != nil {
		return err
	}
	if err := s.Then.setLocation(loc.AppendLocation("then")); err != nil {
		return err
	}
	if err := s.Else.setLocation(loc.AppendLocation("else")); err != nil {
		return err
	}
	if err := s.Properties.setLocation(loc.AppendLocation("properties")); err != nil {
		return err
	}
	if err := s.PropertyNames.setLocation(loc.AppendLocation("propertyNames")); err != nil {
		return err
	}
	if err := s.PatternProperties.setLocation(loc.AppendLocation("patternProperties")); err != nil {
		return err
	}
	if err := s.AdditionalProperties.setLocation(loc.AppendLocation("additionalProperties")); err != nil {
		return err
	}
	if err := s.DependentSchemas.setLocation(loc.AppendLocation("dependentSchemas")); err != nil {
		return err
	}

	if err := s.UnevaluatedProperties.setLocation(loc.AppendLocation("unevaluatedProperties")); err != nil {
		return err
	}
	if err := s.Items.setLocation(loc.AppendLocation("items")); err != nil {
		return err
	}
	if err := s.UnevaluatedItems.setLocation(loc.AppendLocation("unevaluatedItems")); err != nil {
		return err
	}
	if err := s.AdditionalItems.setLocation(loc.AppendLocation("additionalItems")); err != nil {
		return err
	}
	if err := s.PrefixItems.setLocation(loc.AppendLocation("prefixItems")); err != nil {
		return err
	}
	if err := s.Contains.setLocation(loc.AppendLocation("contains")); err != nil {
		return err
	}
	if err := s.RecursiveRef.setLocation(loc.AppendLocation("$recursiveRef")); err != nil {
		return err
	}
	if err := s.Discriminator.setLocation(loc.AppendLocation("discriminator")); err != nil {
		return err
	}
	if err := s.XML.setLocation(loc.AppendLocation("xml")); err != nil {
		return err
	}
	return nil
}

// Clone returns a deep copy of Schema. This is to avoid overriding the initial
// Schema when dealing with $dynamicRef and $recursiveRef.
func (s *Schema) Clone() *Schema {
	if s == nil {
		return nil
	}
	var recAnc *bool
	if s.RecursiveAnchor != nil {
		*recAnc = *s.RecursiveAnchor
	}
	var cnst jsonx.RawMessage
	if s.Const != nil {
		cnst = make(jsonx.RawMessage, len(s.Const))
		copy(cnst, s.Const)
	}

	var required text.Texts
	if s.Required != nil {
		required = make(text.Texts, len(s.Required))
		copy(required, s.Required)
	}
	var example jsonx.RawMessage
	if s.Example != nil {
		example = make(jsonx.RawMessage, len(s.Example))
		copy(example, s.Example)
	}
	var examples []jsonx.RawMessage
	if s.Examples != nil {
		examples = make([]jsonx.RawMessage, len(s.Examples))
		copy(examples, s.Examples)
	}
	var enum text.Texts
	if s.Enum != nil {
		enum = make(text.Texts, len(s.Enum))
		copy(enum, s.Enum)
	}
	var minprops *jsonx.Number
	if s.MinProperties != nil {
		v := *s.MinProperties
		minprops = &v
	}
	var maxprops *jsonx.Number
	if s.MaxProperties != nil {
		v := *s.MaxProperties
		maxprops = &v
	}
	var regexpProps *bool
	if s.RegexProperties != nil {
		v := *s.RegexProperties
		regexpProps = &v
	}
	var depReq *Map[Texts]
	if s.DependentRequired != nil {
		i := make([]KeyValue[Texts], len(s.DependentRequired.Items))
		copy(i, s.DependentRequired.Items)
		depReq = &Map[Texts]{Items: i}
	}
	var uniqItems *bool
	if s.UniqueItems != nil {
		v := *s.UniqueItems
		uniqItems = &v
	}
	var minContains *jsonx.Number
	if s.MinContains != nil {
		v := *s.MinContains
		minContains = &v
	}
	var maxContains *jsonx.Number
	if s.MaxContains != nil {
		v := *s.MaxContains
		maxContains = &v
	}
	var minLen *jsonx.Number
	if s.MinLength != nil {
		v := *s.MinLength
		minLen = &v
	}
	var maxLen *jsonx.Number
	if s.MaxLength != nil {
		v := *s.MaxLength
		maxLen = &v
	}
	var min *jsonx.Number
	if s.Minimum != nil {
		v := *s.Minimum
		min = &v
	}
	var max *jsonx.Number
	if s.Maximum != nil {
		v := *s.Maximum
		max = &v
	}

	var exclMin *jsonx.Number
	if s.ExclusiveMinimum != nil {
		v := *s.ExclusiveMinimum
		exclMin = &v
	}
	var exclMax *jsonx.Number
	if s.ExclusiveMaximum != nil {
		v := *s.ExclusiveMaximum
		exclMax = &v
	}
	var multipleOf *jsonx.Number
	if s.MultipleOf != nil {
		v := *s.MultipleOf
		multipleOf = &v
	}
	var readonly *bool
	if s.ReadOnly != nil {
		v := *s.ReadOnly
		readonly = &v
	}
	var writeOnly *bool
	if s.WriteOnly != nil {
		v := *s.WriteOnly
		writeOnly = &v
	}
	var deprecated *bool
	if s.Deprecated != nil {
		v := *s.Deprecated
		deprecated = &v
	}
	var k map[Text]jsonx.RawMessage
	if s.Keywords != nil {
		k = make(map[Text]jsonx.RawMessage, len(s.Keywords))
		for key, value := range s.Keywords {
			k[key] = value
		}
	}
	var id *uri.URI
	if s.ID != nil {
		id = s.ID.Clone()
	}
	var pattern *Regexp
	if s.Pattern != nil {
		pattern = &Regexp{s.Pattern.Copy()}
	}
	cloned := &Schema{
		RecursiveAnchor:       recAnc,
		Const:                 cnst,
		Required:              required,
		Enum:                  enum,
		Example:               example,
		Examples:              examples,
		MinProperties:         minprops,
		MaxProperties:         maxprops,
		RegexProperties:       regexpProps,
		DependentRequired:     depReq,
		UniqueItems:           uniqItems,
		MinContains:           minContains,
		MaxContains:           maxContains,
		MinLength:             minLen,
		MaxLength:             maxLen,
		Minimum:               min,
		Maximum:               max,
		ExclusiveMinimum:      exclMin,
		ExclusiveMaximum:      exclMax,
		MultipleOf:            multipleOf,
		ReadOnly:              readonly,
		WriteOnly:             writeOnly,
		Deprecated:            deprecated,
		Keywords:              k,
		Schema:                s.Schema,
		ID:                    id,
		Title:                 s.Title,
		Description:           s.Description,
		Default:               s.Default,
		ExternalDocs:          s.ExternalDocs,
		Format:                s.Format,
		ContentMediaType:      s.ContentMediaType,
		Discriminator:         s.Discriminator.Clone(),
		XML:                   s.XML.Clone(),
		Definitions:           s.Definitions.Clone(),
		Anchor:                s.Anchor,
		DynamicAnchor:         s.DynamicAnchor,
		Ref:                   s.Ref.Clone(),
		Type:                  s.Type.Clone(),
		DynamicRef:            s.DynamicRef.Clone(),
		Not:                   s.Not.Clone(),
		AllOf:                 s.AllOf.Clone(),
		AnyOf:                 s.AnyOf.Clone(),
		RecursiveRef:          s.RecursiveRef.Clone(),
		OneOf:                 s.OneOf.Clone(),
		Properties:            s.Properties.Clone(),
		Comments:              s.Comments,
		PropertyNames:         s.PropertyNames.Clone(),
		PatternProperties:     s.PatternProperties.Clone(),
		If:                    s.If.Clone(),
		Then:                  s.Then.Clone(),
		Else:                  s.Else.Clone(),
		AdditionalProperties:  s.AdditionalProperties.Clone(),
		DependentSchemas:      s.DependentSchemas.Clone(),
		UnevaluatedProperties: s.UnevaluatedProperties.Clone(),
		Items:                 s.Items.Clone(),
		UnevaluatedItems:      s.UnevaluatedItems.Clone(),
		AdditionalItems:       s.AdditionalItems.Clone(),
		PrefixItems:           s.PrefixItems.Clone(),
		Contains:              s.Contains.Clone(),
		Pattern:               pattern,
		ContentEncoding:       s.ContentEncoding,
		Extensions:            cloneExtensions(s.Extensions),
		Location:              s.Location,
	}
	return cloned
}

func cloneExtensions(e Extensions) Extensions {
	if e == nil {
		return nil
	}
	a := make(Extensions, len(e))
	for k, v := range e {
		a[k] = v
	}
	return a
}

func (s *Schema) isNil() bool { return s == nil }

var _ node = (*Schema)(nil)

// func (s *Schema) ResolveByAnchor(anchor Text) (*Schema, error) {
// }

// func (s *Schema) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if err := ptr.Validate(); err != nil {
// 		return nil, err
// 	}
// 	return s.resolveNodeByPointer(ptr)
// }

// func (s *Schema) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if ptr.IsRoot() {
// 		return s, nil
// 	}
// 	nxt, tok, _ := ptr.Next()

// 	switch tok {
// 	case "ref":
// 		if s.Ref == nil {
// 			return nil, newErrNotFound(s.AbsoluteLocation(), tok)
// 		}
// 		return s.Ref.resolveNodeByPointer(nxt)
// 	case "definitions":
// 		if s.Definitions == nil {
// 			return nil, newErrNotFound(s.AbsoluteLocation(), tok)
// 		}
// 		return s.Definitions.resolveNodeByPointer(nxt)
// 	case "dynamicRef":
// 		if s.DynamicRef == nil {
// 			return nil, newErrNotFound(s.AbsoluteLocation(), tok)
// 		}
// 		return s.DynamicRef.resolveNodeByPointer(nxt)
// 	case "not":
// 		if s.Not == nil {
// 			return nil, newErrNotFound(s.AbsoluteLocation(), tok)
// 		}
// 		return s.Not.resolveNodeByPointer(nxt)
// 	case "allOf":
// 		if s.AllOf == nil {
// 			return nil, newErrNotFound(s.AbsoluteLocation(), tok)
// 		}
// 		return s.AllOf.resolveNodeByPointer(nxt)
// 	case "anyOf":
// 		if s.AnyOf == nil {
// 			return nil, newErrNotFound(s.AbsoluteLocation(), tok)
// 		}
// 		return s.AnyOf.resolveNodeByPointer(nxt)
// 	case "oneOf":
// 		if s.OneOf == nil {
// 			return nil, newErrNotFound(s.AbsoluteLocation(), tok)
// 		}
// 		return s.OneOf.resolveNodeByPointer(nxt)
// 	case "if":
// 		if s.If == nil {
// 			return nil, newErrNotFound(s.AbsoluteLocation(), tok)
// 		}
// 		return s.If.resolveNodeByPointer(nxt)
// 	case "then":
// 		if s.Then == nil {
// 			return nil, newErrNotFound(s.AbsoluteLocation(), tok)
// 		}
// 		return s.Then.resolveNodeByPointer(nxt)
// 	case "else":
// 		if s.Else == nil {
// 			return nil, newErrNotFound(s.AbsoluteLocation(), tok)
// 		}
// 		return s.Else.resolveNodeByPointer(nxt)
// 	case "properties":
// 		if s.Properties == nil {
// 			return nil, newErrNotFound(s.AbsoluteLocation(), tok)
// 		}
// 		return s.Properties.resolveNodeByPointer(nxt)
// 	case "propertyNames":
// 		if s.PropertyNames == nil {
// 			return nil, newErrNotFound(s.AbsoluteLocation(), tok)
// 		}
// 		return s.PropertyNames.resolveNodeByPointer(nxt)
// 	case "patternProperties":
// 		if s.PatternProperties == nil {
// 			return nil, newErrNotFound(s.AbsoluteLocation(), tok)
// 		}
// 		return s.PatternProperties.resolveNodeByPointer(nxt)
// 	case "additionalProperties":
// 		if s.AdditionalProperties == nil {
// 			return nil, newErrNotFound(s.AbsoluteLocation(), tok)
// 		}
// 		return s.AdditionalProperties.resolveNodeByPointer(nxt)
// 	case "dependentSchemas":
// 		if s.DependentSchemas == nil {
// 			return nil, newErrNotFound(s.AbsoluteLocation(), tok)
// 		}
// 		return s.DependentSchemas.resolveNodeByPointer(nxt)
// 	case "unevaluatedProperties":
// 		if s.UnevaluatedProperties == nil {
// 			return nil, newErrNotFound(s.AbsoluteLocation(), tok)
// 		}
// 		return s.UnevaluatedProperties.resolveNodeByPointer(nxt)
// 	case "items":
// 		if s.Items == nil {
// 			return nil, newErrNotFound(s.AbsoluteLocation(), tok)
// 		}
// 		return s.Items.resolveNodeByPointer(nxt)
// 	case "unevaluatedItems":
// 		if s.UnevaluatedItems == nil {
// 			return nil, newErrNotFound(s.AbsoluteLocation(), tok)
// 		}
// 		return s.UnevaluatedItems.resolveNodeByPointer(nxt)
// 	case "additionalItems":
// 		if s.AdditionalItems == nil {
// 			return nil, newErrNotFound(s.AbsoluteLocation(), tok)
// 		}
// 		return s.AdditionalItems.resolveNodeByPointer(nxt)
// 	case "prefixItems":
// 		if s.PrefixItems == nil {
// 			return nil, newErrNotFound(s.AbsoluteLocation(), tok)
// 		}
// 		return s.PrefixItems.resolveNodeByPointer(nxt)
// 	case "contains":
// 		if s.Contains == nil {
// 			return nil, newErrNotFound(s.AbsoluteLocation(), tok)
// 		}
// 		return s.Contains.resolveNodeByPointer(nxt)
// 	case "recursiveRef":
// 		if s.RecursiveRef == nil {
// 			return nil, newErrNotFound(s.AbsoluteLocation(), tok)
// 		}
// 		return s.RecursiveRef.resolveNodeByPointer(nxt)
// 	case "xml":
// 		if s.XML == nil {
// 			return nil, newErrNotFound(s.AbsoluteLocation(), tok)
// 		}
// 		return s.XML.resolveNodeByPointer(nxt)
// 	default:
// 		return nil, newErrNotResolvable(s.Location.AbsoluteLocation(), tok)
// 	}
// }
