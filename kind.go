package openapi

type Kind uint16

const (
	KindUndefined         Kind = iota
	KindDocument               // *Document
	KindExample                // *Example
	KindExampleMap             // *ExampleMap
	KindSchema                 // *Schema
	KindSchemaSet              // *SchemaSet
	KindSchemaMap              // *SchemaMap
	KindHeader                 // *Header
	KindHeaderMap              // HeaderMap
	KindLink                   // *Link
	KindLinkMap                // LinkMap
	KindResponse               // *Response
	KindResponseMap            // ResponseMap
	KindParameter              // *Parameter
	KindParameterSet           // ParameterSet
	KindParameterMap           // ParameterMap
	KindPaths                  // Paths
	KindPathItem               // *PathItem
	KindPathItemMap            // PathItemMap
	KindRequestBody            // RequestBody
	KindRequestBodyMap         // RequestBodyMap
	KindCallbacks              // *Callbacks
	KindCallbackMap            // CallbackMap
	KindSecurityScheme         // SecurityScheme
	KindSecuritySchemeMap      // SecuritySchemeMap
	KindOperation              // *Operation
	KindLicense                // *License
	KindTag                    // *Tag
	KindMediaType              // *MediaType
	KindInfo                   // *Info
	KindContact                // *Contact
	KindEncoding               // *Encoding
	KindExternalDocs           // *ExternalDocs
	KindReference              // *Reference
)
