package openapi

const (
	// InQuery - Parameters that are appended to the URL. For example, in
	// /items?id=###, the query parameter is id.
	InQuery Text = "query"
	// InHeader - Custom headers that are expected as part of the request. Note
	// that RFC7230 states header names are case insensitive.
	InHeader Text = "header"
	// InCookie -  Used to pass a specific cookie value to the API.
	InCookie Text = "cookie"
	// InPath - Used together with Path Templating, where the parameter value is
	// actually part of the operation's URL. This does not include the host or
	// base path of the API. For example, in /items/{itemId}, the path parameter
	// is itemId.
	InPath Text = "path"
)

type In = Text
