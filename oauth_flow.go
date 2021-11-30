package openapi

import "github.com/chanced/openapi/yamlutil"

// OAuthFlows allows configuration of the supported OAuth Flows.
type OAuthFlows struct {
	// Configuration for the OAuth Implicit flow
	Implicit *OAuthFlow `json:"implicit,omitempty"`
	// Configuration for the OAuth Resource Owner Password flow
	Password *OAuthFlow `json:"password,omitempty"`
	// Configuration for the OAuth Client Credentials flow. Previously called
	// application in OpenAPI 2.0.
	ClientCredentials *OAuthFlow `json:"clientCredentials,omitempty"`
	// Configuration for the OAuth Authorization Code flow. Previously called accessCode in OpenAPI 2.0.
	AuthorizationCode *OAuthFlow `json:"authorizationCode,omitempty"`
	Extensions        `json:"-"`
}

type oauthflows OAuthFlows

// MarshalJSON marshals json
func (oaf OAuthFlows) MarshalJSON() ([]byte, error) {
	return marshalExtendedJSON(oauthflows(oaf))
}

// UnmarshalJSON unmarshals json
func (oaf *OAuthFlows) UnmarshalJSON(data []byte) error {
	var v oauthflows
	if err := unmarshalExtendedJSON(data, &v); err != nil {
		return err
	}
	*oaf = OAuthFlows(v)
	return nil
}

// MarshalYAML marshals YAML
func (oaf OAuthFlows) MarshalYAML() (interface{}, error) {
	return yamlutil.Marshal(oaf)
}

// UnmarshalYAML unmarshals YAML
func (oaf *OAuthFlows) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return yamlutil.Unmarshal(unmarshal, oaf)
}

// OAuthFlow configuration details for a supported OAuth Flow
type OAuthFlow struct {
	// The authorization URL to be used for this flow. This MUST be in the form
	// of a URL. The OAuth2 standard requires the use of TLS.
	//
	// Applies to: OAuth2 ("implicit", "authorizationCode")
	//
	// 	*required*
	AuthorizationURL string `json:"authorizationUrl,omitempty"`
	// The token URL to be used for this flow. This MUST be in the form of a
	// URL. The OAuth2 standard requires the use of TLS.
	//
	// Applies to: OAuth2Flow ("password", "clientCredentials", "authorizationCode")
	//
	// 	*required*
	TokenURL string `json:"tokenUrl,omitempty"`
	// The URL to be used for obtaining refresh tokens. This MUST be in the form
	// of a URL. The OAuth2 standard requires the use of TLS.
	RefreshURL string `json:"refreshUrl,omitempty"`
	// The available scopes for the OAuth2 security scheme. A map between the
	// scope name and a short description for it. The map MAY be empty.
	//
	// 	*required*
	Scopes     map[string]string `json:"scopes"`
	Extensions `json:"-"`
}

type oauthflow OAuthFlow

// MarshalJSON marshals json
func (o OAuthFlow) MarshalJSON() ([]byte, error) {
	return marshalExtendedJSON(oauthflow(o))
}

// UnmarshalJSON unmarshals json
func (o *OAuthFlow) UnmarshalJSON(data []byte) error {
	var v oauthflow
	if err := unmarshalExtendedJSON(data, &v); err != nil {
		return err
	}
	*o = OAuthFlow(v)
	return nil
}

// MarshalYAML marshals YAML
func (o OAuthFlow) MarshalYAML() (interface{}, error) {
	return yamlutil.Marshal(o)
}

// UnmarshalYAML unmarshals YAML
func (o *OAuthFlow) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return yamlutil.Unmarshal(unmarshal, o)
}
