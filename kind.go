package openapi

type Kind uint

const (
	KindUnknown        Kind = iota
	KindOpenAPI             // OpenAPI
	KindExample             // Example
	KindSchema              // Schema
	KindHeader              // Header
	KindLink                // Link
	KindPath                // Path
	KindResponse            // Response
	KindParameter           // Parameter
	KindRequestBody         // RequestBody
	KindCallback            // Callback
	KindSecurityScheme      // SecurityScheme
)
