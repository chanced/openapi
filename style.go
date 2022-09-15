package openapi

const (
	// StyleForm for
	StyleForm Text = "form"
	// StyleSimple comma-separated values. Corresponds to the
	// {param_name} URI template.
	StyleSimple Text = "simple"
	// StyleMatrix is semicolon-prefixed values, also known as path-style
	// expansion. Corresponds to the {;param_name} URI template.
	StyleMatrix Text = "matrix"
	// StyleLabel dot-prefixed values, also known as label expansion.
	// Corresponds to the {.param_name} URI template.
	StyleLabel Text = "label"
	// StyleDeepObject a simple way of rendering nested objects using
	// form parameters (applies to objects only).
	StyleDeepObject Text = "deepObject"
	// StylePipeDelimited is pipeline-separated array values.
	//
	// Same as collectionFormat: pipes in OpenAPI 2.0. Has effect only for
	// non-exploded arrays (explode: false), that is, the pipe separates the
	// array values if the array is a single parameter, as in
	// 	arr=a|b|c
	StylePipeDelimited Text = "pipeDelimited"
)
