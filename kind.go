package openapi

// Kind of Node
type Kind uint16

const (
	KindOpenAPI                 Kind = iota // KindOpenAPI represents *OpenAPI
	KindExternalDocs                        // KindExternalDocs represents *ExternalDocs
	KindCallback                            // KindCallback represents *CallbackObj
	KindComponents                          // KindComponents represents *Components
	KindExample                             // KindExample represents *ExampleObj
	KindExamples                            // KindExamples represents Examples
	KindLink                                // KindLink represents *LinkObj
	KindLinks                               // KindLinks represents Links
	KindParameters                          // KindParameters represents Parameters
	KindParameter                           // KindParameter represents *ParameterObj
	KindParameterSet                        // KindParameterSet represents ParameterSet
	KindInfo                                // KindInfo represents *Info
	KindPath                                // KindPath represents *PathObj
	KindPathItems                           // KindPathItems represents Paths
	KindReference                           // KindReference represents *Reference
	KindRequestBody                         // KindRequestBody represents *RequestBodyObj
	KindRequestBodies                       // KindRequestBodies represents RequestBodies
	KindSecurityRequirement                 // KindSecurityRequirement represents *SecurityRequirementObj
	KindSecurityRequirements                // KindSecurityRequirements represents SecurityRquirements
	KindSecurityScheme                      // KindSecurityScheme represents *SecuritySchemeObj
	KindSecuritySchemes                     // KindSecuritySchemes represents SecuritySchemes
	KindSchema                              // KindSchema represents *SchemaObj
	KindSchemaSet                           // KindSchemaSet represents SchemaSet
	KindSchemas                             // KindSchemas represents Schemas
	KindServer                              // KindServer represents *Server
	KindWebhook                             // KindWebhook represents *WebhookObj
	KindWebhooks                            // KindWebhooks represents Webhooks
	KindResponse                            // KindResponse represents *ResponseObj
	KindResponses                           // KindResponses represents Responses
	KindOperation                           // KindOperation represents *Operation
	KindResolvedOperation                   // KindResolvedOperation represents *ResolvedOperation
	KindResolvedOpenAPI                     // KindResolvedOpenAPI represents *ResolvedOpenAPI
	KindResolvedResponse                    // KindResolvedResponse represents *ResolvedResponse
	KindResolvedResponses                   // KindResolvedResponses represents ResolvedResponses
	KindResolvedRequestBody                 // KindRequestBody represents *ResolvedRequestBody
	KindResolvedRequestBodies               // KindRequestBodies represents ResolvedRequestBodies
	KindResolvedParameters                  // KindResolvedParameters represents Parameters
	KindResolvedParameterSet                // KindResolvedParameterSet represents ParameterSet
	KindResolvedParameter                   // KindResolvedParameter represents *ResolvedParameter
	KindResolvedLink                        // KindLink represents *ResolvedLink
	KindResolvedLinks                       // KindLinks represents ResolvedLinks
	KindResolvedCallback                    // KindResolvedCallback represents *ResolvedCallback
	KindResolvedComponents                  // KindResolvedComponents resolved Components
	KindResolvedExample                     // KindResolvedExample represents *ResolvedExample
	KindResolvedExamples                    // KindResolvedExamples represents ResolvedExamples
	KindResolvedPath                        // KindResolvedPath represents *ResolvedPath
	KindResolvedPathItems                   // KindResolvedPathItems represents ResolvedPathItems
	KindResolvedPaths                       // KindResolvedPaths represents ResolvedPaths
	KindResolvedSecurityScheme              // KindResolvedSecurityScheme resolved SecurityScheme
	KindResolvedSecuritySchemes             // KindResolvedSecuritySchemes resolved SecuritySchemes
	KindResolvedSchema                      // KindResolvedSchema represents *ResolvedSchema
	KindResolvedSchemaSet                   // KindResolvedSchemaSet represents ResolvedSchemaSet
	KindResolvedSchemas                     // KindResolvedSchemas represents ResolvedSchemas
)

var kindNames = map[Kind]string{}

func (k Kind) String() string {
	if s, ok := kindNames[k]; ok {
		return s
	}
	return ""
}
