package openapi

const (
	// InQuery - Parameters that are appended to the URL. For example, in
	// /items?id=###, the query parameter is id.
	InQuery In = "query"
	// InHeader - Custom headers that are expected as part of the request. Note
	// that RFC7230 states header names are case insensitive.
	InHeader In = "header"
	// InCookie -  Used to pass a specific cookie value to the API.
	InCookie In = "cookie"
	// InPath - Used together with Path Templating, where the parameter value is
	// actually part of the operation's URL. This does not include the host or
	// base path of the API. For example, in /items/{itemId}, the path parameter
	// is itemId.
	InPath In = "path"
)

// In is a location where a paremeter may be located in a request.
type In string

func (in In) String() string {
	return string(in)
}
