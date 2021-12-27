package openapi

// Kind of Node
type Kind uint16

const (
	KindNil                     Kind = iota
	KindOpenAPI                      // KindOpenAPI represents *OpenAPI
	KindCallback                     // KindCallback represents *CallbackObj
	KindComponents                   // KindComponents represents *Components
	KindContact                      // KindContact represents *Contact
	KindDiscriminator                // KindDiscriminator represents *Discriminator
	KindEncoding                     // KindEncoding represents *Encoding
	KindEncodings                    // KindEncoding represents Encodings
	KindExample                      // KindExample represents *ExampleObj
	KindExamples                     // KindExamples represents Examples
	KindExternalDocs                 // KindExternalDocs represents *ExternalDocs
	KindHeader                       // KindHeader represents *HeaderObj
	KindHeaders                      // KindHeaders represents Headers
	KindInfo                         // KindInfo represents *Info
	KindLicense                      // KindLicense represents *License
	KindLink                         // KindLink represents *LinkObj
	KindLinks                        // KindLinks represents Links
	KindMediaType                    // KindMediaType represents *MediaType
	KindOAuthFlow                    // KindOAuthFlow represents *OAuthFlow
	KindOAuthFlows                   // KindOAuthFlows represents *OAuthFlows
	KindOperation                    // KindOperation represents *Operation
	KindParameter                    // KindParameter represents *ParameterObj
	KindParameterSet                 // KindParameterSet represents ParameterSet
	KindParameters                   // KindParameters represents Parameters
	KindPath                         // KindPath represents *PathObj
	KindPaths                        // KindPaths represents Paths
	KindPathItems                    // KindPathItems represents PathItems
	KindReference                    // KindReference represents *Reference
	KindRequestBodies                // KindRequestBodies represents RequestBodies
	KindRequestBody                  // KindRequestBody represents *RequestBodyObj
	KindResponse                     // KindResponse represents *ResponseObj
	KindResponses                    // KindResponses represents Responses
	KindSchema                       // KindSchema represents *SchemaObj
	KindSchemaSet                    // KindSchemaSet represents SchemaSet
	KindSchemas                      // KindSchemas represents Schemas
	KindSecurityRequirement          // KindSecurityRequirement represents *SecurityRequirementObj
	KindSecurityRequirements         // KindSecurityRequirements represents SecurityRquirements
	KindSecurityScheme               // KindSecurityScheme represents *SecuritySchemeObj
	KindSecuritySchemes              // KindSecuritySchemes represents SecuritySchemes
	KindServer                       // KindServer represents *Server
	KindTag                          // KindTag represents *Tag
	KindWebhook                      // KindWebhook represents *WebhookObj
	KindWebhooks                     // KindWebhooks represents Webhooks
	KindXML                          // KindXML represents *XML
	KindResolvedCallback             // KindResolvedCallback represents *ResolvedCallback
	KindResolvedComponents           // KindResolvedComponents resolved Components
	KindResolvedEncoding             // KindResolvedEncoding represents *ResolvedEncoding
	KindResolvedEncodings            // KindResolvedEncoding represents ResolvedEncodings
	KindResolvedExample              // KindResolvedExample represents *ResolvedExample
	KindResolvedExamples             // KindResolvedExamples represents ResolvedExamples
	KindResolvedHeader               // KindResolvedHeader represents *ResolvedHeader
	KindResolvedHeaders              // KindResolvedHeaders represents ResolvedHeaders
	KindResolvedLink                 // KindLink represents *ResolvedLink
	KindResolvedLinks                // KindLinks represents ResolvedLinks
	KindResolvedMediaType            // KindResolvedMediaType represents *ResolvedMediaType
	KindResolvedOpenAPI              // KindResolvedOpenAPI represents *ResolvedOpenAPI
	KindResolvedOperation            // KindResolvedOperation represents *ResolvedOperation
	KindResolvedParameter            // KindResolvedParameter represents *ResolvedParameter
	KindResolvedParameterSet         // KindResolvedParameterSet represents ParameterSet
	KindResolvedParameters           // KindResolvedParameters represents Parameters
	KindResolvedPath                 // KindResolvedPath represents *ResolvedPath
	KindResolvedPathItems            // KindResolvedPathItems represents ResolvedPathItems
	KindResolvedPaths                // KindResolvedPaths represents ResolvedPaths
	KindResolvedRequestBodies        // KindRequestBodies represents ResolvedRequestBodies
	KindResolvedRequestBody          // KindRequestBody represents *ResolvedRequestBody
	KindResolvedResponse             // KindResolvedResponse represents *ResolvedResponse
	KindResolvedResponses            // KindResolvedResponses represents ResolvedResponses
	KindResolvedSchema               // KindResolvedSchema represents *ResolvedSchema
	KindResolvedSchemaSet            // KindResolvedSchemaSet represents ResolvedSchemaSet
	KindResolvedSchemas              // KindResolvedSchemas represents ResolvedSchemas
	KindResolvedSecurityScheme       // KindResolvedSecurityScheme resolved SecurityScheme
	KindResolvedSecuritySchemes      // KindResolvedSecuritySchemes resolved SecuritySchemes
	KindResolvedWebhook              // KindResolvedWebhook represents *ResolvedWebhook
	KindResolvedWebhooks             // KindResolvedWebhooks represents ResolvedWebhooks

)

var kindNames = map[Kind]string{
	KindOpenAPI:                 "OpenAPI",
	KindExternalDocs:            "ExternalDocs",
	KindTag:                     "Tag",
	KindCallback:                "Callback",
	KindComponents:              "Components",
	KindExample:                 "Example",
	KindExamples:                "Examples",
	KindLink:                    "Link",
	KindLinks:                   "Links",
	KindParameters:              "Parameters",
	KindParameter:               "Parameter",
	KindParameterSet:            "ParameterSet",
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
	KindResponse:                "Response",
	KindResponses:               "Responses",
	KindOperation:               "Operation",
	KindHeader:                  "Header",
	KindHeaders:                 "Headers",
	KindLicense:                 "License",
	KindContact:                 "Contact",
	KindEncoding:                "Encoding",
	KindMediaType:               "MediaType",
	KindOAuthFlow:               "OAuthFlow",
	KindOAuthFlows:              "OAuthFlows",
	KindDiscriminator:           "Discriminator",
	KindXML:                     "XML",
	KindResolvedEncoding:        "ResolvedEncoding",
	KindResolvedMediaType:       "ResolvedMediaType",
	KindResolvedHeader:          "ResolvedHeader",
	KindResolvedHeaders:         "ResolvedHeaders",
	KindResolvedOperation:       "ResolvedOperation",
	KindResolvedOpenAPI:         "ResolvedOpenAPI",
	KindResolvedResponse:        "ResolvedResponse",
	KindResolvedResponses:       "ResolvedResponses",
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
	KindResolvedExamples:        "ResolvedExamples",
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
