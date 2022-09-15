package openapi

import "github.com/chanced/jsonpointer"

const (
	// SecuritySchemeTypeAPIKey = "apiKey"
	SecuritySchemeTypeAPIKey Text = "apiKey"
	// SecuritySchemeTypeHTTP = "http"
	SecuritySchemeTypeHTTP Text = "http"
	// SecuritySchemeTypeMutualTLS = mutualTLS
	SecuritySchemeTypeMutualTLS Text = "mutualTLS"
	// SecuritySchemeTypeOAuth2 = oauth2
	SecuritySchemeTypeOAuth2 Text = "oauth2"
	// SecuritySchemeTypeOpenIDConnect = "openIdConnect"
	SecuritySchemeTypeOpenIDConnect Text = "openIdConnect"
)

// SecuritySchemeMap is a map of SecurityScheme
type SecuritySchemeMap = ComponentMap[*SecurityScheme]

// SecurityScheme defines a security scheme that can be used by the operations.
type SecurityScheme struct {
	Extensions `json:"-"`
	Location   `json:"-"`

	// The type of the security scheme.
	//
	// *required
	Type Text `json:"type,omitempty"`

	// Any description for security scheme. CommonMark syntax MAY be used for
	// rich text representation.
	Description Text `json:"description,omitempty"`
	// The name of the header, query or cookie parameter to be used.
	//
	// Applies to: API Key
	//
	// 	*required*
	Name Text `json:"name,omitempty"`
	// The location of the API key. Valid values are "query", "header" or "cookie".
	//
	// Applies to: APIKey
	//
	// 	*required*
	In In `json:"in,omitempty"`
	// The name of the HTTP Authorization scheme to be used in the Authorization
	// header as defined in RFC7235. The values used SHOULD be registered in the
	// IANA Authentication Scheme registry.
	//
	// 	*required*
	Scheme Text `json:"scheme,omitempty"`

	// http ("bearer")  A hint to the client to identify how the bearer token is
	// formatted. Bearer tokens are usually generated by an authorization
	// server, so this information is primarily for documentation purposes.
	BearerFormat Text `json:"bearerFormat,omitempty"`

	// An object containing configuration information for the flow types supported.
	//
	// 	*required*
	Flows *OAuthFlows `json:"flows,omitempty"`

	// OpenId Connect URL to discover OAuth2 configuration values. This MUST be
	// in the form of a URL. The OpenID Connect standard requires the use of
	// TLS.
	//
	// 	*required*
	OpenIDConnectURL Text `json:"openIdConnect,omitempty"`
}

func (*SecurityScheme) IsRef() bool { return false }

func (ss *SecurityScheme) Edges() []Node {
	if ss == nil {
		return nil
	}
	return downcastNodes(ss.edges())
}

func (ss *SecurityScheme) edges() []node {
	if ss == nil {
		return nil
	}
	return appendEdges(nil, ss.Flows)
}

func (ss *SecurityScheme) Refs() []Ref {
	if ss == nil {
		return nil
	}
	return ss.Flows.Refs()
}
func (ss *SecurityScheme) isNil() bool { return ss == nil }

func (ss *SecurityScheme) Anchors() (*Anchors, error) {
	if ss == nil {
		return nil, nil
	}

	return ss.Flows.Anchors()
}

func (ss *SecurityScheme) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
	if err := ptr.Validate(); err != nil {
		return nil, err
	}
	return ss.resolveNodeByPointer(ptr)
}

func (ss *SecurityScheme) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
	if ptr.IsRoot() {
		return ss, nil
	}
	nxt, tok, _ := ptr.Next()
	switch nxt {
	case "flows":
		if nxt.IsRoot() {
			return ss.Flows, nil
		}
		if ss.Flows == nil {
			return nil, newErrNotFound(ss.AbsoluteLocation(), tok)
		}
		return ss.Flows.resolveNodeByPointer(nxt)
	default:
		return nil, newErrNotResolvable(ss.AbsoluteLocation(), tok)
	}
}

func (s *SecurityScheme) setLocation(loc Location) error {
	if s == nil {
		return nil
	}
	s.Location = loc
	return s.Flows.setLocation(loc.Append("flows"))
}

// UnmarshalJSON unmarshals JSON
func (ss *SecurityScheme) UnmarshalJSON(data []byte) error {
	type securityscheme SecurityScheme

	var v securityscheme
	err := unmarshalExtendedJSON(data, &v)
	*ss = SecurityScheme(v)
	return err
}

// MarshalJSON marshals JSON
func (ss SecurityScheme) MarshalJSON() ([]byte, error) {
	type securityscheme SecurityScheme

	return marshalExtendedJSON(securityscheme(ss))
}

func (*SecurityScheme) Kind() Kind      { return KindSecurityScheme }
func (*SecurityScheme) mapKind() Kind   { return KindSecuritySchemeMap }
func (*SecurityScheme) sliceKind() Kind { return KindUndefined }

var (
	_ node   = (*SecurityScheme)(nil)
	_ Walker = (*SecurityScheme)(nil)
)
