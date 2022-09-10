package openapi

type kind uint

const (
	kindUnknown        kind = iota
	kindDocument            // OpenAPI Document
	kindExample             // Example
	kindSchema              // Schema
	kindHeader              // Header
	kindLink                // Link
	kindPath                // Path
	kindResponse            // Response
	kindParameter           // Parameter
	kindRequestBody         // RequestBody
	kindCallbacks           // Callbacks
	kindSecurityScheme      // SecurityScheme
	kindOperation           // Operation
	kindLicense             // License
	kindTag                 // Tag
	kindPaths               // Paths
	kindMediaType           // MediaType
	kindInfo                // Info
	kindContact             // Contact
	kindEncoding            // Encoding
	kindExternalDocs        // ExternalDocs
	kindReference           // Reference
)
