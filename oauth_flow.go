package openapi

import "github.com/chanced/jsonpointer"

// OAuthFlows allows configuration of the supported OAuth Flows.
type OAuthFlows struct {
	Extensions `json:"-"`
	Location   `json:"-"`

	// Configuration for the OAuth Implicit flow
	Implicit *OAuthFlow `json:"implicit,omitempty"`
	// Configuration for the OAuth Resource Owner Password flow
	Password *OAuthFlow `json:"password,omitempty"`
	// Configuration for the OAuth Client Credentials flow. Previously called
	// application in OpenAPI 2.0.
	ClientCredentials *OAuthFlow `json:"clientCredentials,omitempty"`
	// Configuration for the OAuth Authorization Code flow. Previously called accessCode in OpenAPI 2.0.
	AuthorizationCode *OAuthFlow `json:"authorizationCode,omitempty"`
}

func (f *OAuthFlows) Resolve(ptr jsonpointer.Pointer) (Node, error) {
	if err := ptr.Validate(); err != nil {
		return nil, err
	}
	return f.resolve(ptr)
}

func (f *OAuthFlows) resolve(ptr jsonpointer.Pointer) (Node, error) {
	if ptr.IsRoot() {
		return f, nil
	}
	nxt, tok, _ := ptr.Next()
	var n node
	switch tok {
	case "implicit":
		n = f.Implicit
	case "password":
		n = f.Password
	case "clientCredentials":
		n = f.ClientCredentials
	case "authorizationCode":
		n = f.AuthorizationCode
	default:
		return nil, newErrNotResolvable(f.Location.AbsoluteLocation(), tok)
	}
	if nxt.IsRoot() {
		return n, nil
	}

	if n == nil {
		return nil, newErrNotFound(f.Location.AbsoluteLocation(), tok)
	}
	return n.resolve(nxt)
}

func (*OAuthFlows) Kind() Kind      { return KindOAuthFlows }
func (*OAuthFlows) mapKind() Kind   { return KindUndefined }
func (*OAuthFlows) sliceKind() Kind { return KindUndefined }

func (o *OAuthFlows) setLocation(loc Location) error {
	if o == nil {
		return nil
	}
	o.Location = loc
	if err := o.Implicit.setLocation(loc.Append("implicit")); err != nil {
		return err
	}
	if err := o.Password.setLocation(loc.Append("password")); err != nil {
		return err
	}
	if err := o.ClientCredentials.setLocation(loc.Append("clientCredentials")); err != nil {
		return err
	}
	if err := o.AuthorizationCode.setLocation(loc.Append("authorizationCode")); err != nil {
		return err
	}
	return nil
}

// MarshalJSON marshals json
func (o OAuthFlows) MarshalJSON() ([]byte, error) {
	type oauthflows OAuthFlows

	return marshalExtendedJSON(oauthflows(o))
}

// UnmarshalJSON unmarshals json
func (f *OAuthFlows) UnmarshalJSON(data []byte) error {
	type oauthflows OAuthFlows
	var v oauthflows
	if err := unmarshalExtendedJSON(data, &v); err != nil {
		return err
	}
	*f = OAuthFlows(v)
	return nil
}

// OAuthFlow configuration details for a supported OAuth Flow
type OAuthFlow struct {
	Extensions `json:"-"`
	Location   `json:"-"`

	// The authorization URL to be used for this flow. This MUST be in the form
	// of a URL. The OAuth2 standard requires the use of TLS.
	//
	// Applies to: OAuth2 ("implicit", "authorizationCode")
	//
	// 	*required*
	AuthorizationURL Text `json:"authorizationUrl,omitempty"`
	// The token URL to be used for this flow. This MUST be in the form of a
	// URL. The OAuth2 standard requires the use of TLS.
	//
	// Applies to: OAuth2Flow ("password", "clientCredentials", "authorizationCode")
	//
	// 	*required*
	TokenURL Text `json:"tokenUrl,omitempty"`
	// The URL to be used for obtaining refresh tokens. This MUST be in the form
	// of a URL. The OAuth2 standard requires the use of TLS.
	RefreshURL Text `json:"refreshUrl,omitempty"`
	// The available scopes for the OAuth2 security scheme. A map between the
	// scope name and a short description for it. The map MAY be empty.
	//
	// 	*required*
	Scopes *Scopes `json:"scopes"`
}

func (f *OAuthFlow) Resolve(ptr jsonpointer.Pointer) (Node, error) {
	if err := ptr.Validate(); err != nil {
		return nil, err
	}
	return f.resolve(ptr)
}

func (f *OAuthFlow) resolve(ptr jsonpointer.Pointer) (Node, error) {
	if ptr.IsRoot() {
		return f, nil
	}
	nxt, tok, _ := ptr.Next()
	var n node
	switch tok {
	case "scopes":
		n = f.Scopes
	default:
		return nil, newErrNotResolvable(f.Location.AbsoluteLocation(), tok)
	}
	if nxt.IsRoot() {
		return n, nil
	}

	if n == nil {
		return nil, newErrNotFound(f.Location.AbsoluteLocation(), tok)
	}
	return n.resolve(nxt)
}

func (*OAuthFlow) Kind() Kind      { return KindOAuthFlow }
func (*OAuthFlow) mapKind() Kind   { return KindUndefined }
func (*OAuthFlow) sliceKind() Kind { return KindUndefined }

func (o *OAuthFlow) setLocation(loc Location) error {
	if o == nil {
		return nil
	}
	o.Location = loc
	o.Scopes.setLocation(loc.Append("scopes"))
	return nil
}

// MarshalJSON marshals json
func (o OAuthFlow) MarshalJSON() ([]byte, error) {
	type oauthflow OAuthFlow

	return marshalExtendedJSON(oauthflow(o))
}

// UnmarshalJSON unmarshals json
func (o *OAuthFlow) UnmarshalJSON(data []byte) error {
	type oauthflow OAuthFlow
	var v oauthflow
	if err := unmarshalExtendedJSON(data, &v); err != nil {
		return err
	}
	*o = OAuthFlow(v)
	return nil
}

var (
	_ node = (*OAuthFlow)(nil)
	_ node = (*OAuthFlows)(nil)
)
