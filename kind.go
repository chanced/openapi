package openapi

type Kind uint16

const (
	KindUndefined                Kind = iota
	KindDocument                      // *Document
	KindComponents                    // *Components
	KindExample                       // *Example
	KindExampleMap                    // *ExampleMap
	KindExampleComponent              // *Component[*Example]
	KindSchema                        // *Schema
	KindSchemaSlice                   // *SchemaSlice
	KindSchemaMap                     // *SchemaMap
	KindSchemaRef                     // *SchemaRef
	KindDiscriminator                 // *Discriminator
	KindHeader                        // *Header
	KindHeaderMap                     // *HeaderMap
	KindHeaderSlice                   // *HeaderSlice
	KindHeaderComponent               // *Component[*Header]
	KindLink                          // *Link
	KindLinkComponent                 // *Component[*Link]
	KindLinkMap                       // *LinkMap
	KindResponse                      // *Response
	KindResponseMap                   // *ResponseMap
	KindResponseComponent             // *Component[*Response]
	KindParameter                     // *Parameter
	KindParameterComponent            // *Component[*Parameter]
	KindParameterSlice                // *ParameterSlice
	KindParameterMap                  // *ParameterMap
	KindPaths                         // *Paths
	KindPathItem                      // *PathItem
	KindPathItemComponent             // *Component[*PathItem]
	KindPathItemMap                   // *PathItemMap
	KindRequestBody                   // *RequestBody
	KindRequestBodyMap                // *RequestBodyMap
	KindRequestBodyComponent          // *Component[*RequestBody]
	KindCallbacks                     // *Callbacks
	KindCallbacksComponent            // *Component[*Callbacks]
	KindCallbacksMap                  // *CallbacksMap
	KindSecurityRequirementSlice      // *SecurityRequirements
	KindSecurityRequirement           // *SecurityRequirement
	KindSecurityRequirementItem       // *SecurityRequirementItem
	KindSecurityScheme                // *SecurityScheme
	KindSecuritySchemeComponent       // *Component[*SecurityScheme]
	KindSecuritySchemeMap             // *SecuritySchemeMap
	KindOperation                     // *Operation
	KindOperationRef                  // *OperationRef
	KindLicense                       // *License
	KindTag                           // *Tag
	KindTagSlice                      // *TagSlice
	KindMediaType                     // *MediaType
	KindMediaTypeMap                  // *MediaTypeMap
	KindInfo                          // *Info
	KindContact                       // *Contact
	KindEncoding                      // *Encoding
	KindEncodingMap                   // *EncodingMap
	KindExternalDocs                  // *ExternalDocs
	KindReference                     // *Reference
	KindServer                        // *Server
	KindServerComponent               // *Component[*Server]
	KindServerSlice                   // *ServerSlice
	KindServerVariable                // *ServerVariable
	KindServerVariableMap             // *ServerVariableMap
	KindOAuthFlow                     // *OAuthFlow
	KindOAuthFlows                    // *OAuthFlows
	KindXML                           // *XML
	KindScope                         // *Scope
	KindScopes                        // *Scopes
)

func (k Kind) String() string {
	switch k {
	case KindUndefined:
		return "Undefined"
	case KindDocument:
		return "Document"
	case KindComponents:
		return "Components"
	case KindExample:
		return "Example"
	case KindExampleMap:
		return "ExampleMap"
	case KindExampleComponent:
		return "ExampleComponent"
	case KindSchema:
		return "Schema"
	case KindSchemaSlice:
		return "SchemaSlice"
	case KindSchemaMap:
		return "SchemaMap"
	case KindSchemaRef:
		return "SchemaRef"
	case KindDiscriminator:
		return "Discriminator"
	case KindHeader:
		return "Header"
	case KindHeaderMap:
		return "HeaderMap"
	case KindHeaderSlice:
		return "HeaderSlice"
	case KindHeaderComponent:
		return "HeaderComponent"
	case KindLink:
		return "Link"
	case KindLinkComponent:
		return "LinkComponent"
	case KindLinkMap:
		return "LinkMap"
	case KindResponse:
		return "Response"
	case KindResponseMap:
		return "ResponseMap"
	case KindResponseComponent:
		return "ResponseComponent"
	case KindParameter:
		return "Parameter"
	case KindParameterComponent:
		return "ParameterComponent"
	case KindParameterSlice:
		return "ParameterSlice"
	case KindParameterMap:
		return "ParameterMap"
	case KindPaths:
		return "Paths"
	case KindPathItem:
		return "PathItem"
	case KindPathItemComponent:
		return "PathItemComponent"
	case KindPathItemMap:
		return "PathItemMap"
	case KindRequestBody:
		return "RequestBody"
	case KindRequestBodyMap:
		return "RequestBodyMap"
	case KindRequestBodyComponent:
		return "RequestBodyComponent"
	case KindCallbacks:
		return "Callbacks"
	case KindCallbacksComponent:
		return "CallbacksComponent"
	case KindCallbacksMap:
		return "CallbacksMap"
	case KindSecurityRequirementSlice:
		return "SecurityRequirementSlice"
	case KindSecurityRequirement:
		return "SecurityRequirement"
	case KindSecurityRequirementItem:
		return "SecurityRequirementItem"
	case KindSecurityScheme:
		return "SecurityScheme"
	case KindSecuritySchemeComponent:
		return "SecuritySchemeComponent"
	case KindSecuritySchemeMap:
		return "SecuritySchemeMap"
	case KindOperation:
		return "Operation"
	case KindOperationRef:
		return "OperationRef"
	case KindLicense:
		return "License"
	case KindTag:
		return "Tag"
	case KindTagSlice:
		return "TagSlice"
	case KindMediaType:
		return "MediaType"
	case KindMediaTypeMap:
		return "MediaTypeMap"
	case KindInfo:
		return "Info"
	case KindContact:
		return "Contact"
	case KindEncoding:
		return "Encoding"
	case KindEncodingMap:
		return "EncodingMap"
	case KindExternalDocs:
		return "ExternalDocs"
	case KindReference:
		return "Reference"
	case KindServer:
		return "Server"
	case KindServerComponent:
		return "ServerComponent"
	case KindServerSlice:
		return "ServerSlice"
	case KindServerVariable:
		return "ServerVariable"
	case KindServerVariableMap:
		return "ServerVariableMap"
	case KindOAuthFlow:
		return "OAuthFlow"
	case KindOAuthFlows:
		return "OAuthFlows"
	case KindXML:
		return "XML"
	case KindScope:
		return "Scope"
	case KindScopes:
		return "Scopes"
	default:
		return "Invalid"
	}
}

func objSliceKind(n node) Kind {
	if sn, ok := n.(objSlicedNode); ok {
		return sn.objSliceKind()
	}
	return KindUndefined
}
