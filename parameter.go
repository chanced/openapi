package openapi

import (
	"encoding/json"

	"github.com/chanced/jsonx"
	"github.com/chanced/transcode"
	"gopkg.in/yaml.v3"
)

// ParameterMap is a map of Parameter
type ParameterMap = ComponentMap[*Parameter]

// ParameterSlice is list of parameters that are applicable for a given operation.
// If a parameter is already defined at the Path Item, the new definition will
// override it but can never remove it. The list MUST NOT include duplicated
// parameters. A unique parameter is defined by a combination of a name and
// location. The list can use the Reference Object to link to parameters that
// are defined at the OpenAPI Object's components/parameters.
//
// Can either be a Parameter or a Reference
type ParameterSlice = ComponentSlice[*Parameter]

/*
* Path Parameters

* Path parameters support the following style values:
* - simple – (default) comma-separated values. Corresponds to the {param_name} URI template.
* - label – dot-prefixed values, also known as label expansion. Corresponds to the {.param_name} URI template.
* - matrix – semicolon-prefixed values, also known as path-style expansion. Corresponds to the {;param_name} URI template.
*
* The label and matrix styles are sometimes used with partial path parameters, such as /users{id}, because the parameter values get prefixed.
*
* +----------+---------+---------------+------------------------+------------------------+----------------------------------------------------+
* |  style   | explode | URI template  | Primitive value id = 5 |  Array id = [3, 4, 5]  | Object id = {"role": "admin", "firstName": "Alex"} |
* +----------+---------+---------------+------------------------+------------------------+----------------------------------------------------+
* | simple * | false * | /users/{id}   | /users/5               | /users/3,4,5           | /users/role,admin,firstName,Alex                   |
* | simple   | true    | /users/{id*}  | /users/5               | /users/3,4,5           | /users/role=admin,firstName=Alex                   |
* | label    | false   | /users/{.id}  | /users/.5              | /users/.3,4,5          | /users/.role,admin,firstName,Alex                  |
* | label    | true    | /users/{.id*} | /users/.5              | /users/.3.4.5          | /users/.role=admin.firstName=Alex                  |
* | matrix   | false   | /users/{;id}  | /users/;id=5           | /users/;id=3,4,5       | /users/;id=role,admin,firstName,Alex               |
* | matrix   | true    | /users/{;id*} | /users/;id=5           | /users/;id=3;id=4;id=5 | /users/;role=admin;firstName=Alex                  |
* +----------+---------+---------------+------------------------+------------------------+----------------------------------------------------+
 */

/*
* Query Parameters
*
* - form – (default) ampersand-separated values, also known as form-style query expansion.
*   		Corresponds to the {?param_name} URI template.
* - spaceDelimited – space-separated array values. Same as collectionFormat: ssv in
*   		OpenAPI 2.0. Has effect only for non-exploded arrays (explode: false),
*   		that is, the space separates the array values if the array is a single
*   		parameter, as in arr=a b c.
* - pipeDelimited – pipeline-separated array values. Same as collectionFormat: pipes in
*   		OpenAPI 2.0. Has effect only for non-exploded arrays (explode: false),
*   		that is, the pipe separates the array values if the array is a single
*   		parameter, as in arr=a|b|c.
* - deepObject – a simple way of rendering nested objects using form parameters
*   		(applies to objects only).
*
* The default serialization method is style: form and explode: true. This
* corresponds to  collectionFormat: multi from OpenAPI 2.0. Given the path
* users with a query parameter id, the query string is serialized as follows:
*
* +----------+---------+---------------+------------------------+------------------------+----------------------------------------------------+
* |  style   | explode | URI template  | Primitive value id = 5 |  Array id = [3, 4, 5]  | Object id = {"role": "admin", "firstName": "Alex"} |
* +----------+---------+---------------+------------------------+------------------------+----------------------------------------------------+
* | simple * | false * | /users/{id}   | /users/5               | /users/3,4,5           | /users/role,admin,firstName,Alex                   |
* | simple   | true    | /users/{id*}  | /users/5               | /users/3,4,5           | /users/role=admin,firstName=Alex                   |
* | label    | false   | /users/{.id}  | /users/.5              | /users/.3,4,5          | /users/.role,admin,firstName,Alex                  |
* | label    | true    | /users/{.id*} | /users/.5              | /users/.3.4.5          | /users/.role=admin.firstName=Alex                  |
* | matrix   | false   | /users/{;id}  | /users/;id=5           | /users/;id=3,4,5       | /users/;id=role,admin,firstName,Alex               |
* | matrix   | true    | /users/{;id*} | /users/;id=5           | /users/;id=3;id=4;id=5 | /users/;role=admin;firstName=Alex                  |
* +----------+---------+---------------+------------------------+------------------------+----------------------------------------------------+
 */

/*
*
* Header Parameters
*
* Header parameters always use the simple style, that is, comma-separated
* values. This corresponds to the {param_name} URI template. An optional
* explode keyword controls the object serialization. Given the request header
* named X-MyHeader, the header value is serialized as follows:
*
* +----------+---------+--------------+--------------------------------+------------------------------+------------------------------------------------------------+
* |  style   | explode | URI template | Primitive value X-MyHeader = 5 | Array X-MyHeader = [3, 4, 5] | Object X-MyHeader = {"role": "admin", "firstName": "Alex"} |
* +----------+---------+--------------+--------------------------------+------------------------------+------------------------------------------------------------+
* | simple * | false * | {id}         | X-MyHeader: 5                  | X-MyHeader: 3,4,5            | X-MyHeader: role,admin,firstName,Alex                      |
* | simple   | true    | {id*}        | X-MyHeader: 5                  | X-MyHeader: 3,4,5            | X-MyHeader: role=admin,firstName=Alex                      |
* +----------+---------+--------------+--------------------------------+------------------------------+------------------------------------------------------------+
 */

/*
* Cookie Parameters
*
* Cookie parameters always use the form style. An optional explode keyword
* controls the array and object serialization. Given the cookie named id, the
* cookie value is serialized as follows:
*
* +--------+---------+--------------+------------------------+----------------------+----------------------------------------------------+
* | style  | explode | URI template | Primitive value id = 5 | Array id = [3, 4, 5] | Object id = {"role": "admin", "firstName": "Alex"} |
* +--------+---------+--------------+------------------------+----------------------+----------------------------------------------------+
* | form * | true *  |              | Cookie: id=5           |                      |                                                    |
* | form   | false   | id={id}      | Cookie: id=5           | Cookie: id=3,4,5     | Cookie: id=role,admin,firstName,Alex               |
* +--------+---------+--------------+------------------------+----------------------+----------------------------------------------------+
 */

/*
* +-----------------------+------------------------------------------------------------------------------------------------------------------------------------------+
* |        Keyword        |                                                          URI Template Modifier                                                           |
* +-----------------------+------------------------------------------------------------------------------------------------------------------------------------------+
* | style: simple         | none                                                                                                                                     |
* | style: label          | . prefix                                                                                                                                 |
* | style: matrix         | ; prefix                                                                                                                                 |
* | style: form           | ? or & prefix (depending on the parameter position in the query string)                                                                  |
* | style: pipeDelimited  | ? or & prefix (depending on the parameter position in the query string) – but using pipes | instead of commas , to join the array values |
* | style: spaceDelimited | ? or & prefix (depending on the parameter position in the query string) – but using spaces instead of commas , to join the array values  |
* | explode: false        | none                                                                                                                                     |
* | explode: true         | * suffix                                                                                                                                 |
* | allowReserved: false  | none                                                                                                                                     |
* | allowReserved: true   | + prefix                                                                                                                                 |
* +-----------------------+------------------------------------------------------------------------------------------------------------------------------------------+
 */

// Parameter describes a single operation parameter.
//
// A unique parameter is defined by a combination of a name and location.
type Parameter struct {
	Extensions `json:"-"`
	Location   `json:"-"`

	// The name of the parameter. Parameter names are case sensitive:
	//
	// - If In is "path", the name field MUST correspond to a template
	// expression occurring within the path field in the Paths Object.
	// See Path Templating for further information.
	//
	// - If In is "header" and the name field is "Accept", "Content-Type"
	// or "Authorization", the parameter definition SHALL be ignored.
	//
	// - For all other cases, the name corresponds to the parameter name
	// used by the in property.
	//
	//  *required*
	Name Text `json:"name"`

	// The location of the parameter. Possible values are "query", "header",
	// "path" or "cookie".
	//
	//  *required*
	In In `json:"in"`

	// A brief description of the parameter. This could contain examples of use.
	// CommonMark syntax MAY be used for rich text representation.
	Description Text `json:"description,omitempty"`

	// Determines whether this parameter is mandatory. If the parameter location
	// is "path", this property is REQUIRED and its value MUST be true.
	// Otherwise, the property MAY be included and its default value is false.
	Required *bool `json:"required,omitempty"`

	// Specifies that a parameter is deprecated and SHOULD be transitioned out
	// of usage. Default value is false.
	Deprecated bool `json:"deprecated,omitempty"`

	// Sets the ability to pass empty-valued parameters. This is valid only for
	// query parameters and allows sending a parameter with an empty value.
	// Default value is false. If style is used, and if behavior is n/a (cannot
	// be serialized), the value of allowEmptyValue SHALL be ignored. Use of
	// this property is NOT RECOMMENDED, as it is likely to be removed in a
	// later revision.
	AllowEmptyValue bool `json:"allowEmptyValue,omitempty"`

	// Describes how the parameter value will be serialized depending on the
	// type of the parameter value.
	// Default values (based on value of in):
	//  - for query - form;
	// 	- for path - simple;
	// 	- for header - simple;
	// 	- for cookie - form.
	Style Text `json:"style,omitempty"`

	// When this is true, parameter values of type array or object generate
	// separate parameters for each value of the array or key-value pair of the
	// map. For other types of parameters this property has no effect. When
	// style is form, the default value is true. For all other styles, the
	// default value is false.
	Explode bool `json:"explode,omitempty"`

	// Determines whether the parameter value SHOULD allow reserved characters,
	// as defined by RFC3986 :/?#[]@!$&'()*+,;= to be included without
	// percent-encoding. This property only applies to parameters with an in
	// value of query. The default value is false.
	AllowReserved bool `json:"allowReserved,omitempty"`

	// The schema defining the type used for the parameter.
	Schema *Schema `json:"schema,omitempty"`

	// Examples of the parameter's potential value. Each example SHOULD
	// contain a value in the correct format as specified in the parameter
	// encoding. The examples field is mutually exclusive of the example
	// field. Furthermore, if referencing a schema that contains an example,
	// the examples value SHALL override the example provided by the schema.
	Examples *ExampleMap `json:"examples,omitempty"`

	Example jsonx.RawMessage `json:"example,omitempty"`

	// For more complex scenarios, the content property can define the media
	// type and schema of the parameter. A parameter MUST contain either a
	// schema property, or a content property, but not both. When example or
	// examples are provided in conjunction with the schema object, the example
	// MUST follow the prescribed serialization strategy for the parameter.
	Content *ContentMap `json:"content,omitempty"`
}

func (p *Parameter) Nodes() []Node {
	if p == nil {
		return nil
	}
	return downcastNodes(p.nodes())
}

func (p *Parameter) nodes() []node {
	if p == nil {
		return nil
	}
	return appendEdges(nil, p.Schema, p.Examples, p.Content)
}

func (p *Parameter) Refs() []Ref {
	if p == nil {
		return nil
	}
	var refs []Ref
	if p.Schema != nil {
		refs = append(refs, p.Schema.Refs()...)
	}
	if p.Examples != nil {
		refs = append(refs, p.Examples.Refs()...)
	}
	if p.Content != nil {
		refs = append(refs, p.Content.Refs()...)
	}
	return refs
}

func (p *Parameter) Anchors() (*Anchors, error) {
	if p == nil {
		return nil, nil
	}
	var anchors *Anchors
	var err error
	if anchors, err = anchors.merge(p.Schema.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(p.Content.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(p.Examples.Anchors()); err != nil {
		return nil, err
	}
	return anchors, nil
}

// func (p *Parameter) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if err := ptr.Validate(); err != nil {
// 		return nil, err
// 	}
// 	return p.resolveNodeByPointer(ptr)
// }

// func (p *Parameter) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if ptr.IsRoot() {
// 		return p, nil
// 	}
// 	nxt, tok, _ := ptr.Next()
// 	switch tok {
// 	case "schema":
// 		if p.Schema == nil {
// 			return nil, newErrNotFound(p.AbsoluteLocation(), tok)
// 		}
// 		return p.Schema.resolveNodeByPointer(nxt)
// 	case "content":
// 		if p.Content == nil {
// 			return nil, newErrNotFound(p.AbsoluteLocation(), tok)
// 		}
// 		return p.Content.resolveNodeByPointer(nxt)
// 	case "examples":
// 		if p.Examples == nil {
// 			return nil, newErrNotFound(p.AbsoluteLocation(), tok)
// 		}
// 		return p.Examples.resolveNodeByPointer(nxt)
// 	default:
// 		return nil, newErrNotResolvable(p.AbsoluteLocation(), tok)
// 	}
// }

// MarshalJSON marshals h into JSON
func (p Parameter) MarshalJSON() ([]byte, error) {
	type parameter Parameter
	return marshalExtendedJSON(parameter(p))
}

// UnmarshalJSON unmarshals json into p
func (p *Parameter) UnmarshalJSON(data []byte) error {
	type parameter Parameter
	var v parameter

	if err := unmarshalExtendedJSON(data, &v); err != nil {
		return err
	}
	*p = Parameter(v)
	return nil
}

func (*Parameter) Kind() Kind      { return KindParameter }
func (*Parameter) mapKind() Kind   { return KindParameterMap }
func (*Parameter) sliceKind() Kind { return KindParameterSlice }

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Marshaler interface
func (p Parameter) MarshalYAML() (interface{}, error) {
	j, err := p.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return transcode.YAMLFromJSON(j)
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Unmarshaler interface
func (p *Parameter) UnmarshalYAML(value *yaml.Node) error {
	v, err := yaml.Marshal(value)
	if err != nil {
		return err
	}
	j, err := transcode.JSONFromYAML(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(j, p)
}

func (p *Parameter) setLocation(loc Location) error {
	if p == nil {
		return nil
	}
	p.Location = loc
	if err := p.Schema.setLocation(loc.AppendLocation("schema")); err != nil {
		return err
	}
	if err := p.Content.setLocation(loc.AppendLocation("content")); err != nil {
		return err
	}
	if err := p.Examples.setLocation(loc.AppendLocation("examples")); err != nil {
		return err
	}

	return nil
}
func (p *Parameter) isNil() bool { return p == nil }

func (*Parameter) refable() {}

var _ node = (*Parameter)(nil)
