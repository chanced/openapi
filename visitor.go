package openapi

// type Walker interface {
// Walk(v Visitor) error
// }

// type Visitor interface {
// 	Visit(n Node) (Visitor, error)

// 	VisitDocument(node *Document) (Visitor, error)

// 	VisitCallbacks(node *Callbacks) (Visitor, error)

// 	VisitComponents(node *Components) (Visitor, error)

// 	VisitCallbacksMap(node *CallbacksMap) (Visitor, error)

// 	VisitContact(node *Contact) (Visitor, error)

// 	VisitDiscriminator(node *Discriminator) (Visitor, error)

// 	VisitEncoding(node *Encoding) (Visitor, error)

// 	VisitEncodingMap(node *EncodingMap) (Visitor, error)

// 	VisitExample(node *Example) (Visitor, error)

// 	VisitExternalDocs(node *ExternalDocs) (Visitor, error)

// 	VisitHeader(node *Header) (Visitor, error)

// 	VisitHeaderMap(node *HeaderMap) (Visitor, error)

// 	VisitInfo(node *Info) (Visitor, error)

// 	VisitLicense(node *License) (Visitor, error)

// 	VistLink(node *Link) (Visitor, error)

// 	VistMediaType(node *MediaType) (Visitor, error)

// 	VisitOAuthFlows(node *OAuthFlows) (Visitor, error)

// 	VisitOAuthFlow(node *OAuthFlow) (Visitor, error)

// 	VisitOperation(node *Operation) (Visitor, error)

// 	VisitOperationItem(node *OperationItem) (Visitor, error)

// 	VisitOperationRef(node *OperationRef) (Visitor, error)

// 	VisitParameter(node *Parameter) (Visitor, error)

// 	VisitParameterSlice(node *ParameterSlice) (Visitor, error)

// 	VisitParameterMap(node *ParameterMap) (Visitor, error)

// 	VisitPathItemComponent(node *Component[*PathItem]) (Visitor, error)

// 	VisitPathItem(node *PathItem) (Visitor, error)

// 	// A map of PathItems may not contain references
// 	VisitPathItemObjs(node *PathItemObjs) (Visitor, error)

// 	// Visits a a map of PathItems which may contain references
// 	VisitPathItemMap(node *PathItemMap) (Visitor, error)

// 	VisitPaths(node *Paths) (Visitor, error)

// 	VisitReference(node *Reference) (Visitor, error)

// 	VisitRequestBody(node *RequestBody) (Visitor, error)

// 	VisitRequestBodyMap(node *RequestBodyMap) (Visitor, error)

// 	VisitResponse(node *Response) (Visitor, error)

// 	VisitResponseMap(node *ResponseMap) (Visitor, error)

// 	VisitSchema(node *Schema) (Visitor, error)

// 	VisitSchemaMap(node *SchemaMap) (Visitor, error)

// 	VisitSchemaSlice(node *SchemaSlice) (Visitor, error)

// 	VisitSchemaRef(node *SchemaRef) (Visitor, error)

// 	VisitScope(node *Scope) (Visitor, error)

// 	VisitSecurityRequirement(node *SecurityRequirement) (Visitor, error)

// 	VisitSecurityScheme(node *SecurityScheme) (Visitor, error)

// 	VisitSecuritySchemeMap(node *SecuritySchemeMap) (Visitor, error)

// 	VisitServer(node *Server) (Visitor, error)

// 	VisitServerSlice(node *ServerSlice) (Visitor, error)

// 	VisitServerVariable(node *ServerVariable) (Visitor, error)

// 	VisitServerVariableMap(node *ServerVariableMap) (Visitor, error)

// 	VisitTagSlice(node *TagSlice) (Visitor, error)

// 	VisitTag(node *Tag) (Visitor, error)

// 	VisitXML(node *XML) (Visitor, error)
// }

type BaseVisitor struct{}
