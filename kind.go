package openapi

// Kind of Node
type Kind uint16

const (
	// KindNil is the zero value of Kind
	KindNil Kind = iota
	// KindOpenAPI represents *OpenAPI
	KindOpenAPI
	// KindCallback represents *CallbackObj
	KindCallback
	// KindComponents represents *Components
	KindComponents
	// KindContact represents *Contact
	KindContact
	// KindContent represents Content
	KindContent
	// KindDiscriminator represents *Discriminator
	KindDiscriminator
	// KindEncoding represents *Encoding
	KindEncoding
	// KindEncodings represents Encodings
	KindEncodings
	// KindExample represents *ExampleObj
	KindExample
	// KindExamples represents Examples
	KindExamples
	// KindExternalDocs represents *ExternalDocs
	KindExternalDocs
	// KindHeader represents *HeaderObj
	KindHeader
	// KindHeaders represents Headers
	KindHeaders
	// KindInfo represents *Info
	KindInfo
	// KindLicense represents *License
	KindLicense
	// KindLink represents *LinkObj
	KindLink
	// KindLinks represents Links
	KindLinks
	// KindMediaType represents *MediaType
	KindMediaType
	// KindOAuthFlow represents *OAuthFlow
	KindOAuthFlow
	// KindOAuthFlows represents *OAuthFlows
	KindOAuthFlows
	// KindOperation represents *Operation
	KindOperation
	// KindParameter represents *ParameterObj
	KindParameter
	// KindParameterSet represents ParameterSet
	KindParameterSet
	// KindParameters represents Parameters
	KindParameters
	// KindPath represents *PathObj
	KindPath
	// KindPaths represents Paths
	KindPaths
	// KindPathItems represents PathItems
	KindPathItems
	// KindReference represents *Reference
	KindReference
	// KindRequestBodies represents RequestBodies
	KindRequestBodies
	// KindRequestBody represents *RequestBodyObj
	KindRequestBody
	// KindResponse represents *ResponseObj
	KindResponse
	// KindResponses represents Responses
	KindResponses
	// KindSchema represents *SchemaObj
	KindSchema
	// KindSchemaSet represents SchemaSet
	KindSchemaSet
	// KindSchemas represents Schemas
	KindSchemas
	// KindSecurityRequirement represents *SecurityRequirementObj
	KindSecurityRequirement
	// KindSecurityRequirements represents SecurityRquirements
	KindSecurityRequirements
	// KindSecurityScheme represents *SecuritySchemeObj
	KindSecurityScheme
	// KindSecuritySchemes represents SecuritySchemes
	KindSecuritySchemes
	// KindServer represents *Server
	KindServer
	// KindTag represents *Tag
	KindTag
	// KindWebhook represents *WebhookObj
	KindWebhook
	// KindWebhooks represents Webhooks
	KindWebhooks
	// KindXML represents *XML
	KindXML
	// KindResolvedCallback represents *ResolvedCallback
	KindResolvedCallback
	// KindResolvedComponents represents ResolvedComponents
	KindResolvedComponents
	// KindResolvedContent represents ResolvedContent
	KindResolvedContent
	// KindResolvedEncoding represents *ResolvedEncoding
	KindResolvedEncoding
	// KindResolvedEncodings represents ResolvedEncodings
	KindResolvedEncodings
	// KindResolvedExample represents *ResolvedExample
	KindResolvedExample
	// KindResolvedExamples represents ResolvedExamples
	KindResolvedExamples
	// KindResolvedHeader represents *ResolvedHeader
	KindResolvedHeader
	// KindResolvedHeaders represents ResolvedHeaders
	KindResolvedHeaders
	// KindResolvedLink represents *ResolvedLink
	KindResolvedLink
	// KindResolvedLinks represents ResolvedLinks
	KindResolvedLinks
	// KindResolvedMediaType represents *ResolvedMediaType
	KindResolvedMediaType
	// KindResolvedOpenAPI represents *ResolvedOpenAPI
	KindResolvedOpenAPI
	// KindResolvedOperation represents *ResolvedOperation
	KindResolvedOperation
	// KindResolvedParameter represents *ResolvedParameter
	KindResolvedParameter
	// KindResolvedParameterSet represents ParameterSet
	KindResolvedParameterSet
	// KindResolvedParameters represents Parameters
	KindResolvedParameters
	// KindResolvedPath represents *ResolvedPath
	KindResolvedPath
	// KindResolvedPathItems represents ResolvedPathItems
	KindResolvedPathItems
	// KindResolvedPaths represents ResolvedPaths
	KindResolvedPaths
	// KindResolvedRequestBodies represents ResolvedRequestBodies
	KindResolvedRequestBodies
	// KindResolvedRequestBody represents *ResolvedRequestBody
	KindResolvedRequestBody
	// KindResolvedResponse represents *ResolvedResponse
	KindResolvedResponse
	// KindResolvedResponses represents ResolvedResponses
	KindResolvedResponses
	// KindResolvedSchema represents *ResolvedSchema
	KindResolvedSchema
	// KindResolvedSchemaSet represents ResolvedSchemaSet
	KindResolvedSchemaSet
	// KindResolvedSchemas represents ResolvedSchemas
	KindResolvedSchemas
	// KindResolvedSecurityScheme resolved SecurityScheme
	KindResolvedSecurityScheme
	// KindResolvedSecuritySchemes resolved SecuritySchemes
	KindResolvedSecuritySchemes
	// KindResolvedWebhook represents *ResolvedWebhook
	KindResolvedWebhook
	// KindResolvedWebhooks represents ResolvedWebhooks
	KindResolvedWebhooks
)

var kindNames = map[Kind]string{
	KindNil:                     "nil",
	KindOpenAPI:                 "OpenAPI",
	KindCallback:                "Callback",
	KindComponents:              "Components",
	KindContact:                 "Contact",
	KindContent:                 "Content",
	KindDiscriminator:           "Discriminator",
	KindEncoding:                "Encoding",
	KindEncodings:               "Encodings",
	KindExample:                 "Example",
	KindExamples:                "Examples",
	KindExternalDocs:            "ExternalDocs",
	KindHeader:                  "Header",
	KindHeaders:                 "Headers",
	KindInfo:                    "Info",
	KindLicense:                 "License",
	KindLink:                    "Link",
	KindLinks:                   "Links",
	KindMediaType:               "MediaType",
	KindOAuthFlow:               "OAuthFlow",
	KindOAuthFlows:              "OAuthFlows",
	KindOperation:               "Operation",
	KindParameter:               "Parameter",
	KindParameterSet:            "ParameterSet",
	KindParameters:              "Parameters",
	KindPath:                    "Path",
	KindPaths:                   "Paths",
	KindPathItems:               "PathItems",
	KindReference:               "Reference",
	KindRequestBodies:           "RequestBodies",
	KindRequestBody:             "RequestBody",
	KindResponse:                "Response",
	KindResponses:               "Responses",
	KindSchema:                  "Schema",
	KindSchemaSet:               "SchemaSet",
	KindSchemas:                 "Schemas",
	KindSecurityRequirement:     "SecurityRequirement",
	KindSecurityRequirements:    "SecurityRequirements",
	KindSecurityScheme:          "SecurityScheme",
	KindSecuritySchemes:         "SecuritySchemes",
	KindServer:                  "Server",
	KindTag:                     "Tag",
	KindWebhook:                 "Webhook",
	KindWebhooks:                "Webhooks",
	KindXML:                     "XML",
	KindResolvedCallback:        "ResolvedCallback",
	KindResolvedComponents:      "ResolvedComponents",
	KindResolvedContent:         "ResolvedContent",
	KindResolvedEncoding:        "ResolvedEncoding",
	KindResolvedEncodings:       "ResolvedEncodings",
	KindResolvedExample:         "ResolvedExample",
	KindResolvedExamples:        "ResolvedExamples",
	KindResolvedHeader:          "ResolvedHeader",
	KindResolvedHeaders:         "ResolvedHeaders",
	KindResolvedLink:            "ResolvedLink",
	KindResolvedLinks:           "ResolvedLinks",
	KindResolvedMediaType:       "ResolvedMediaType",
	KindResolvedOpenAPI:         "ResolvedOpenAPI",
	KindResolvedOperation:       "ResolvedOperation",
	KindResolvedParameter:       "ResolvedParameter",
	KindResolvedParameterSet:    "ResolvedParameterSet",
	KindResolvedParameters:      "ResolvedParameters",
	KindResolvedPath:            "ResolvedPath",
	KindResolvedPathItems:       "ResolvedPathItems",
	KindResolvedPaths:           "ResolvedPaths",
	KindResolvedRequestBodies:   "ResolvedRequestBodies",
	KindResolvedRequestBody:     "ResolvedRequestBody",
	KindResolvedResponse:        "ResolvedResponse",
	KindResolvedResponses:       "ResolvedResponses",
	KindResolvedSchema:          "ResolvedSchema",
	KindResolvedSchemaSet:       "ResolvedSchemaSet",
	KindResolvedSchemas:         "ResolvedSchemas",
	KindResolvedSecurityScheme:  "ResolvedSecurityScheme",
	KindResolvedSecuritySchemes: "ResolvedSecuritySchemes",
	KindResolvedWebhook:         "ResolvedWebhook",
	KindResolvedWebhooks:        "ResolvedWebhooks",
}

func (k Kind) String() string {
	if s, ok := kindNames[k]; ok {
		return s
	}
	return ""
}
