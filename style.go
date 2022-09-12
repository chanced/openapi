package openapi

const (
	// StyleForm for
	StyleForm Style = "form"
	// StyleSimple comma-separated values. Corresponds to the
	// {param_name} URI template.
	StyleSimple Style = "simple"
	// StyleMatrix is semicolon-prefixed values, also known as path-style
	// expansion. Corresponds to the {;param_name} URI template.
	StyleMatrix Style = "matrix"
	// StyleLabel dot-prefixed values, also known as label expansion.
	// Corresponds to the {.param_name} URI template.
	StyleLabel Style = "label"
	// StyleDeepObject a simple way of rendering nested objects using
	// form parameters (applies to objects only).
	StyleDeepObject Style = "deepObject"
	// StylePipeDelimited is pipeline-separated array values.
	//
	// Same as collectionFormat: pipes in OpenAPI 2.0. Has effect only for
	// non-exploded arrays (explode: false), that is, the pipe separates the
	// array values if the array is a single parameter, as in
	// 	arr=a|b|c
	StylePipeDelimited Style = "pipeDelimited"
)

// Style describes how the parameter value will be serialized depending
// on the type of the parameter value.
type Style string

func (s Style) String() string {
	return string(s)
}

func (s Style) Text() Text {
	return Text(s.String())
}

// IsForm reports whether or not s equals "form"
func (s Style) IsForm() bool {
	return s == StyleForm
}

// IsSimple reports whether or not s equals "simple"
func (s Style) IsSimple() bool {
	return s == StyleSimple
}

// IsMatrix reports whether or not s equals "matrix"
func (s Style) IsMatrix() bool {
	return s == StyleMatrix
}

// IsLabel reports whether or not s equals "label"
func (s Style) IsLabel() bool {
	return s == StyleLabel
}

// IsDeepObject reports whether or not s equals "deepObject"
func (s Style) IsDeepObject() bool {
	return s == StyleDeepObject
}

// IsPipeDelimited reports whether or not s equals "pipeDelimited"
func (s Style) IsPipeDelimited() bool {
	return s == StylePipeDelimited
}

// IsValid reports whether or not s is a valid Style
//
// Valid styles are:
//
//	"simple"  "matrix"  "label"  "deepObject"  "pipeDelimited"
func (s Style) IsValid() bool {
	switch s {
	case StyleForm, StyleSimple, StyleMatrix, StyleLabel, StyleDeepObject, StylePipeDelimited:
		return true
	}
	return false
}
