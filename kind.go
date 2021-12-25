package openapi

// Kind of Node
type Kind uint16

const (
	KindInvalid                 Kind = iota // KindInvalid is an invalid node.
	KindOpenAPI                             // KindOpenAPI represents *OpenAPI
	KindCallback                            // KindCallback represents *CallbackObj
	KindComponents                          // KindComponents represents *Components
	KindExample                             // KindExample represents *ExampleObj
	KindLink                                // KindLink represents *LinkObj
	KindLinks                               // KindLinks represents Links
	KindParameters                          // KindParameters represents Parameters
	KindParameter                           // KindParameter represents *ParameterObj
	KindInfo                                // KindInfo represents *Info
	KindPath                                // KindPath represents *PathObj
	KindPathItems                           // KindPathItems represents Paths
	KindReference                           // KindReference represents *Reference
	KindRequestBody                         // KindRequestBody represents a *RequestBodyObj
	KindRequestBodies                       // KindRequestBodies represents a RequestBodies map
	KindSecurityRequirement                 // KindSecurityRequirement represents a *SecurityRequirementObj
	KindSecurityRequirements                // KindSecurityRequirements represents aSecurityRquirements
	KindSecurityScheme                      // KindSecurityScheme represents a *SecuritySchemeObj
	KindSecuritySchemes                     // KindSecuritySchemes represents SecuritySchemes
	KindSchema                              // KindSchema represents a *SchemaObj
	KindSchemaSet                           // KindSchemaSet represents a SchemaSet
	KindSchemas                             // KindSchemas represents Schemas
	KindServer                              // KindServer represents *Server
	KindWebhook                             // KindWebhook represents a *WebhookObj
	KindWebhooks                            // KindWebhooks represents a Webhooks map
	KindResolvedRequestBody                 // KindRequestBody represents a *ResolvedRequestBody
	KindResolvedRequestBodies               // KindRequestBodies represents a
	KindResolvedParameters                  // KindResolvedParameters represents a Parameters map
	KindResolvedParameterSet                // KindResolvedParameterSet represents a ParameterSet
	KindResolvedParameter                   // KindResolvedParameter represents a *ResolvedParameter
	KindResolvedLink                        // KindLink represents a *ResolvedLink
	KindResolvedLinks                       // KindLinks represents a ResolvedLinks map
	KindResolvedCallback                    // KindResolvedCallback represents a *ResolvedCallback
	KindResolvedComponents                  // KindResolvedComponents resolved Components
	KindResolvedExample                     // KindResolvedExample represents a *ResolvedExample
	KindResolvedOpenAPI                     // KindResolvedOpenAPI represents a *ResolvedOpenAPI
	KindResolvedPath                        // KindResolvedPath represents a *ResolvedPath
	KindResolvedPathItems                   // KindResolvedPathItems represents a ResolvedPathItems map
	KindResolvedPaths                       // KindResolvedPaths represents a ResolvedPaths map
	KindResolvedSecurityScheme              // KindResolvedSecurityScheme resolved SecurityScheme
	KindResolvedSecuritySchemes             // KindResolvedSecuritySchemes resolved SecuritySchemes
	KindResolvedSchema                      // KindResolvedSchema represents a *ResolvedSchema
	KindResolvedSchemaSet                   // KindResolvedSchemaSet represents a ResolvedSchemaSet
	KindResolvedSchemas                     // KindResolvedSchemas represents ResolvedSchemas
)

var kindNames = map[Kind]string{
	KindOpenAPI:                 "OpenAPI",
	KindCallback:                "Callback",
	KindComponents:              "Components",
	KindExample:                 "Example",
	KindLink:                    "Link",
	KindLinks:                   "Links",
	KindParameters:              "Parameters",
	KindParameter:               "Parameter",
	KindInfo:                    "Info",
	KindPath:                    "Path",
	KindPathItems:               "PathItems",
	KindReference:               "Reference",
	KindRequestBody:             "RequestBody",
	KindRequestBodies:           "RequestBodies",
	KindSecurityRequirement:     "SecurityRequirement",
	KindSecurityRequirements:    "SecurityRequirements",
	KindSecurityScheme:          "SecurityScheme",
	KindSecuritySchemes:         "SecuritySchemes",
	KindSchema:                  "Schema",
	KindSchemaSet:               "SchemaSet",
	KindSchemas:                 "Schemas",
	KindServer:                  "Server",
	KindWebhook:                 "Webhook",
	KindWebhooks:                "Webhooks",
	KindResolvedRequestBody:     "ResolvedRequestBody",
	KindResolvedRequestBodies:   "ResolvedRequestBodies",
	KindResolvedParameters:      "ResolvedParameters",
	KindResolvedParameterSet:    "ResolvedParameterSet",
	KindResolvedParameter:       "ResolvedParameter",
	KindResolvedLink:            "ResolvedLink",
	KindResolvedLinks:           "ResolvedLinks",
	KindResolvedCallback:        "ResolvedCallback",
	KindResolvedComponents:      "ResolvedComponents",
	KindResolvedExample:         "ResolvedExample",
	KindResolvedOpenAPI:         "ResolvedOpenAPI",
	KindResolvedPath:            "ResolvedPath",
	KindResolvedPathItems:       "ResolvedPathItems",
	KindResolvedPaths:           "ResolvedPaths",
	KindResolvedSecurityScheme:  "ResolvedSecurityScheme",
	KindResolvedSecuritySchemes: "ResolvedSecuritySchemes",
	KindResolvedSchema:          "ResolvedSchema",
	KindResolvedSchemaSet:       "ResolvedSchemaSet",
	KindResolvedSchemas:         "ResolvedSchemas",
}

func (k Kind) String() string {
	if s, ok := kindNames[k]; ok {
		return s
	}
	return ""
}
