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
	KindParameterSet                        // KindParameterSet represents ParameterSet
	KindInfo                                // KindInfo represents *Info
	KindPath                                // KindPath represents *PathObj
	KindPathItems                           // KindPathItems represents Paths
	KindReference                           // KindReference represents *Reference
	KindRequestBody                         // KindRequestBody represents*RequestBodyObj
	KindRequestBodies                       // KindRequestBodies representsRequestBodies map
	KindSecurityRequirement                 // KindSecurityRequirement represents*SecurityRequirementObj
	KindSecurityRequirements                // KindSecurityRequirements represents aSecurityRquirements
	KindSecurityScheme                      // KindSecurityScheme represents*SecuritySchemeObj
	KindSecuritySchemes                     // KindSecuritySchemes represents SecuritySchemes
	KindSchema                              // KindSchema represents*SchemaObj
	KindSchemaSet                           // KindSchemaSet representsSchemaSet
	KindSchemas                             // KindSchemas represents Schemas
	KindServer                              // KindServer represents *Server
	KindWebhook                             // KindWebhook represents*WebhookObj
	KindWebhooks                            // KindWebhooks representsWebhooks map
	KindResolvedRequestBody                 // KindRequestBody represents*ResolvedRequestBody
	KindResolvedRequestBodies               // KindRequestBodies represents ResolvedRequestBodies
	KindResolvedParameters                  // KindResolvedParameters representsParameters map
	KindResolvedParameterSet                // KindResolvedParameterSet representsParameterSet
	KindResolvedParameter                   // KindResolvedParameter represents*ResolvedParameter
	KindResolvedLink                        // KindLink represents*ResolvedLink
	KindResolvedLinks                       // KindLinks representsResolvedLinks map
	KindResolvedCallback                    // KindResolvedCallback represents*ResolvedCallback
	KindResolvedComponents                  // KindResolvedComponents resolved Components
	KindResolvedExample                     // KindResolvedExample represents*ResolvedExample
	KindResolvedOpenAPI                     // KindResolvedOpenAPI represents*ResolvedOpenAPI
	KindResolvedPath                        // KindResolvedPath represents*ResolvedPath
	KindResolvedPathItems                   // KindResolvedPathItems representsResolvedPathItems map
	KindResolvedPaths                       // KindResolvedPaths representsResolvedPaths map
	KindResolvedSecurityScheme              // KindResolvedSecurityScheme resolved SecurityScheme
	KindResolvedSecuritySchemes             // KindResolvedSecuritySchemes resolved SecuritySchemes
	KindResolvedSchema                      // KindResolvedSchema represents*ResolvedSchema
	KindResolvedSchemaSet                   // KindResolvedSchemaSet representsResolvedSchemaSet
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
