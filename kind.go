package openapi

type Kind uint16

const (
	KindUndefined               Kind = iota
	KindDocument                     // *Document
	KindExample                      // *Example
	KindExampleMap                   // *ExampleMap
	KindExampleComponent             // *Component[*Example]
	KindSchema                       // *Schema
	KindSchemaSlice                  // *SchemaSlice
	KindSchemaMap                    // *SchemaMap
	KindSchemaRef                    // *SchemaRef
	KindHeader                       // *Header
	KindHeaderMap                    // *HeaderMap
	KindHeaderComponent              // *Component[*Header]
	KindLink                         // *Link
	KindLinkComponent                // *Component[*Link]
	KindLinkMap                      // *LinkMap
	KindResponse                     // *Response
	KindResponseMap                  // *ResponseMap
	KindResponseComponent            // *Component[*Response]
	KindParameter                    // *Parameter
	KindParameterComponent           // *Component[*Parameter]
	KindParameterSlice               // *ParameterSlice
	KindParameterMap                 // *ParameterMap
	KindPaths                        // *Paths
	KindPathItem                     // *PathItem
	KindPathItemComponent            // *Component[*PathItem]
	KindPathItemMap                  // *PathItemMap
	KindRequestBody                  // *RequestBody
	KindRequestBodyMap               // *RequestBodyMap
	KindRequestBodyComponent         // *Component[*RequestBody]
	KindCallbacks                    // *Callbacks
	KindCallbacksComponent           // *Component[*Callbacks]
	KindCallbacksMap                 // *CallbacksMap
	KindSecurityRequirements         // *SecurityRequirements
	KindSecurityRequirement          // *SecurityRequirement
	KindSecurityRequirementItem      // *SecurityRequirementItem
	KindSecurityScheme               // *SecurityScheme
	KindSecuritySchemeComponent      // *Component[*SecurityScheme]
	KindSecuritySchemeMap            // *SecuritySchemeMap
	KindOperation                    // *Operation
	KindLicense                      // *License
	KindTag                          // *Tag
	KindTagSlice                     // *TagSlice
	KindMediaType                    // *MediaType
	KindMediaTypeMap                 // *MediaTypeMap
	KindInfo                         // *Info
	KindContact                      // *Contact
	KindEncoding                     // *Encoding
	KindEncodingMap                  // *EncodingMap
	KindExternalDocs                 // *ExternalDocs
	KindReference                    // *Reference
	KindServer                       // *Server
	KindServerSlice                  // *ServerSlice
	KindServerVariable               // *ServerVariable
	KindServerVariableMap            // *ServerVariableMap
	KindOAuthFlow                    // *OAuthFlow
	KindOAuthFlows                   // *OAuthFlows
	KindXML                          // *XML
	KindScope                        // *Scope
	KindScopes                       // *Scopes
)
