package openapi

type Kind uint16

const (
	KindUndefined         Kind = iota
	KindDocument               // *Document
	KindExample                // *Example
	KindExampleMap             // *ExampleMap
	KindSchema                 // *Schema
	KindSchemaSlice            // *SchemaSlice
	KindSchemaMap              // *SchemaMap
	KindHeader                 // *Header
	KindHeaderMap              // HeaderMap
	KindLink                   // *Link
	KindLinkMap                // LinkMap
	KindResponse               // *Response
	KindResponseMap            // ResponseMap
	KindParameter              // *Parameter
	KindParameterSlice         // ParameterSlice
	KindParameterMap           // ParameterMap
	KindPaths                  // Paths
	KindPathItem               // *PathItem
	KindPathItemMap            // PathItemMap
	KindRequestBody            // RequestBody
	KindRequestBodyMap         // RequestBodyMap
	KindCallbacks              // *Callbacks
	KindCallbacksMap           // CallbacksMap
	KindSecurityScheme         // SecurityScheme
	KindSecuritySchemeMap      // SecuritySchemeMap
	KindOperation              // *Operation
	KindLicense                // *License
	KindTag                    // *Tag
	KindTagSlice               // TagSlice
	KindMediaType              // *MediaType
	KindMediaTypeMap           // MediaTypeMap
	KindInfo                   // *Info
	KindContact                // *Contact
	KindEncoding               // *Encoding
	KindEncodingMap            // EncodingMap
	KindExternalDocs           // *ExternalDocs
	KindReference              // *Reference
	KindServer                 // *Server
	KindServerSlice            // ServerSlice
	KindServerVariable         // *ServerVariable
	KindServerVariableMap      // ServerVariableMap
	KindOAuthFlow              // *OAuthFlow
	KindOAuthFlows             // *OAuthFlows
	KindXML                    // *XML

)
