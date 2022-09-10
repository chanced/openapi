package openapi

type Kind uint

const (
	KindUnknown        Kind = iota
	KindDocument            // OpenAPI Document
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
	KindOperation           // Operation
	KindLicense             // License
	KindTag                 // Tag
	KindPaths               // Paths
	KindMediaType           // MediaType
	KindInfo                // Info
	KindContact             // Contact
	KindEncoding            // Encoding
	KindExternalDocs        // ExternalDocs
	KindReference           // Reference
)
