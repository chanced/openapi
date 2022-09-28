package openapi

import (
	"encoding/json"

	"github.com/chanced/transcode"
	"gopkg.in/yaml.v3"
)

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

func (f *OAuthFlow) Refs() []Ref {
	if f == nil {
		return nil
	}
	return f.Scopes.Refs()
}

func (f *OAuthFlow) Anchors() (*Anchors, error) {
	if f == nil {
		return nil, nil
	}
	return f.Scopes.Anchors()
}

// func (f *OAuthFlow) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if err := ptr.Validate(); err != nil {
// 		return nil, err
// 	}
// 	return f.resolveNodeByPointer(ptr)
// }

//	func (f *OAuthFlow) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
//		if ptr.IsRoot() {
//			return f, nil
//		}
//		nxt, tok, _ := ptr.Next()
//		switch tok {
//		case "scopes":
//			if f.Scopes == nil {
//				return nil, newErrNotFound(f.Location.AbsoluteLocation(), tok)
//			}
//			return f.Scopes.resolveNodeByPointer(nxt)
//		default:
//			return nil, newErrNotResolvable(f.Location.AbsoluteLocation(), tok)
//		}
//	}
func (r *OAuthFlow) Nodes() []Node {
	if r == nil {
		return nil
	}
	return downcastNodes(r.nodes())
}

func (r *OAuthFlow) nodes() []node {
	if r == nil {
		return nil
	}
	return appendEdges(nil, r.Scopes)
}

func (*OAuthFlow) Kind() Kind      { return KindOAuthFlow }
func (*OAuthFlow) mapKind() Kind   { return KindUndefined }
func (*OAuthFlow) sliceKind() Kind { return KindUndefined }

func (o *OAuthFlow) setLocation(loc Location) error {
	if o == nil {
		return nil
	}
	o.Location = loc
	o.Scopes.setLocation(loc.AppendLocation("scopes"))
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

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Marshaler interface
func (o OAuthFlow) MarshalYAML() (interface{}, error) {
	j, err := o.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return transcode.YAMLFromJSON(j)
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Unmarshaler interface
func (o *OAuthFlow) UnmarshalYAML(value *yaml.Node) error {
	v, err := yaml.Marshal(value)
	if err != nil {
		return err
	}
	j, err := transcode.JSONFromYAML(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(j, o)
}

func (o *OAuthFlow) isNil() bool { return o == nil }

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

func (f *OAuthFlows) Nodes() []Node {
	if f == nil {
		return nil
	}
	return downcastNodes(f.nodes())
}

func (f *OAuthFlows) nodes() []node {
	if f == nil {
		return nil
	}
	return appendEdges(nil, f.Implicit, f.Password, f.ClientCredentials, f.AuthorizationCode)
}

func (f *OAuthFlows) Refs() []Ref {
	if f == nil {
		return nil
	}
	var refs []Ref
	refs = append(refs, f.Implicit.Refs()...)
	refs = append(refs, f.Password.Refs()...)
	refs = append(refs, f.ClientCredentials.Refs()...)
	refs = append(refs, f.AuthorizationCode.Refs()...)
	return refs
}

func (f *OAuthFlows) isNil() bool { return f == nil }
func (f *OAuthFlows) Anchors() (*Anchors, error) {
	if f == nil {
		return nil, nil
	}
	var anchors *Anchors
	var err error
	if anchors, err = anchors.merge(f.Implicit.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(f.Password.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(f.ClientCredentials.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(f.AuthorizationCode.Anchors()); err != nil {
		return nil, err
	}
	return anchors, nil
}

// func (f *OAuthFlows) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if err := ptr.Validate(); err != nil {
// 		return nil, err
// 	}
// 	return f.resolveNodeByPointer(ptr)
// }

// func (f *OAuthFlows) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if ptr.IsRoot() {
// 		return f, nil
// 	}
// 	nxt, tok, _ := ptr.Next()
// 	switch tok {
// 	case "implicit":
// 		if f.Implicit == nil {
// 			return nil, newErrNotFound(f.Location.AbsoluteLocation(), tok)
// 		}
// 		return f.Implicit.resolveNodeByPointer(nxt)
// 	case "password":
// 		if f.Password == nil {
// 			return nil, newErrNotFound(f.Location.AbsoluteLocation(), tok)
// 		}
// 		return f.Password.resolveNodeByPointer(nxt)
// 	case "clientCredentials":
// 		if f.ClientCredentials == nil {
// 			return nil, newErrNotFound(f.Location.AbsoluteLocation(), tok)
// 		}
// 		return f.ClientCredentials.resolveNodeByPointer(nxt)
// 	case "authorizationCode":
// 		if f.AuthorizationCode == nil {
// 			return nil, newErrNotFound(f.Location.AbsoluteLocation(), tok)
// 		}
// 		return f.AuthorizationCode.resolveNodeByPointer(nxt)
// 	default:
// 		return nil, newErrNotResolvable(f.Location.AbsoluteLocation(), tok)
// 	}
// }

func (*OAuthFlows) Kind() Kind      { return KindOAuthFlows }
func (*OAuthFlows) mapKind() Kind   { return KindUndefined }
func (*OAuthFlows) sliceKind() Kind { return KindUndefined }

func (o *OAuthFlows) setLocation(loc Location) error {
	if o == nil {
		return nil
	}
	o.Location = loc
	if err := o.Implicit.setLocation(loc.AppendLocation("implicit")); err != nil {
		return err
	}
	if err := o.Password.setLocation(loc.AppendLocation("password")); err != nil {
		return err
	}
	if err := o.ClientCredentials.setLocation(loc.AppendLocation("clientCredentials")); err != nil {
		return err
	}
	if err := o.AuthorizationCode.setLocation(loc.AppendLocation("authorizationCode")); err != nil {
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

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Marshaler interface
func (f OAuthFlows) MarshalYAML() (interface{}, error) {
	j, err := f.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return transcode.YAMLFromJSON(j)
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Unmarshaler interface
func (f *OAuthFlows) UnmarshalYAML(value *yaml.Node) error {
	v, err := yaml.Marshal(value)
	if err != nil {
		return err
	}
	j, err := transcode.JSONFromYAML(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(j, f)
}

var (
	_ node = (*OAuthFlow)(nil)

	_ node = (*OAuthFlows)(nil)
)
