package openapi

import (
	"encoding/json"

	"github.com/chanced/openapi/yamlutil"
)

// Path can either be a Path or a Reference
type Webhook interface {
	Node
	ResolveWebhook(func(ref string) (*PathObj, error)) (*WebhookObj, error)
}

type WebhookObj PathObj

// KindPath returns KindWebhook
func (p *WebhookObj) Kind() Kind {
	return KindWebhook
}

type webhook PathObj

// ResolveWebhook resolves WebhookObj by returning itself. resolve is  not called.
func (w *WebhookObj) ResolveWebhook(func(ref string) (*WebhookObj, error)) (*WebhookObj, error) {
	return w, nil
}

// MarshalJSON marshals p into JSON
func (w WebhookObj) MarshalJSON() ([]byte, error) {
	return marshalExtendedJSON(webhook(w))
}

// UnmarshalJSON unmarshals json into p
func (w *WebhookObj) UnmarshalJSON(data []byte) error {
	var v webhook
	if err := unmarshalExtendedJSON(data, &v); err != nil {
		return err
	}
	*w = WebhookObj(v)
	return nil
}

// MarshalYAML first marshals and unmarshals into JSON and then marshals into
// YAML
func (w WebhookObj) MarshalYAML() (interface{}, error) {
	return yamlutil.Marshal(w)
}

// UnmarshalYAML unmarshals yaml into s
func (w *WebhookObj) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return yamlutil.Unmarshal(unmarshal, w)
}

// Webhooks is a map of Webhooks that can either be a Webhook or a Reference
type Webhooks map[string]Webhook

// Kind returns KindWebhooks
func (Webhooks) Kind() Kind {
	return KindWebhooks
}

// UnmarshalJSON unmarshals JSON data into rp
func (ws *Webhooks) UnmarshalJSON(data []byte) error {
	var rd map[string]json.RawMessage
	err := json.Unmarshal(data, &rd)
	if err != nil {
		return err
	}
	res := make(Webhooks)
	for k, d := range rd {
		if isRefJSON(data) {
			var v Reference
			if err = json.Unmarshal(d, &v); err != nil {
				return err
			}
			res[k] = &v
		} else {
			var v WebhookObj
			if err = json.Unmarshal(d, &v); err != nil {
				return err
			}
			res[k] = &v
		}
	}
	*ws = res
	return nil
}

// // UnmarshalYAML unmarshals YAML data into rp
// func (rp *PathItems) UnmarshalYAML(unmarshal func(interface{}) error) error {
// 	return yamlutil.Unmarshal(unmarshal, rp)
// }

// // MarshalYAML marshals rp into YAML
// func (rp PathItems) MarshalYAML() (interface{}, error) {
// 	b, err := json.Marshal(rp)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var v interface{}
// 	err = json.Unmarshal(b, &v)
// 	return v, err
// }

// func unmarshalPathJSON(data []byte) (Path, error) {
// 	if isRefJSON(data) {
// 		return unmarshalReferenceJSON(data)
// 	}
// 	var p pathobj
// 	err := json.Unmarshal(data, &p)
// 	v := PathObj(p)
// 	return &v, err
// }

// // ResolvedPath is a Path Object which has beeen resolved. It describes the
// // operations available on a single path. A PathObj Item MAY be empty, due to
// // ACL constraints. The path itself is still exposed to the documentation viewer
// // but they will not know which operations and parameters are available.
// type ResolvedPath struct {
// 	// Allows for a referenced definition of this path item. The referenced
// 	// structure MUST be in the form of a Path Item Object. In case a Path Item
// 	// Object field appears both in the defined object and the referenced
// 	// object, the behavior is undefined. See the rules for resolving Relative
// 	// References.
// 	Ref string `json:"$ref,omitempty"`
// 	// An optional, string summary, intended to apply to all operations in this path.
// 	Summary string `json:"summary,omitempty"`
// 	// An optional, string description, intended to apply to all operations in
// 	// this path. CommonMark syntax MAY be used for rich text representation.
// 	Description string `json:"description,omitempty"`
// 	// A definition of a GET operation on this path.
// 	Get *Operation `json:"get,omitempty"`
// 	// A definition of a PUT operation on this path.
// 	Put *Operation `json:"put,omitempty"`
// 	// A definition of a POST operation on this path.
// 	Post *Operation `json:"post,omitempty"`
// 	// A definition of a DELETE operation on this path.
// 	Delete *Operation `json:"delete,omitempty"`
// 	// A definition of a OPTIONS operation on this path.
// 	Options *Operation `json:"options,omitempty"`
// 	// A definition of a HEAD operation on this path.
// 	Head *Operation `json:"head,omitempty"`
// 	// A definition of a PATCH operation on this path.
// 	Patch *Operation `json:"patch,omitempty"`
// 	// A definition of a TRACE operation on this path.
// 	Trace *Operation `json:"trace,omitempty"`
// 	// An alternative server array to service all operations in this path.
// 	Servers []*Server `json:"servers,omitempty"`
// 	// A list of parameters that are applicable for all the operations described
// 	// under this path. These parameters can be overridden at the operation
// 	// level, but cannot be removed there. The list MUST NOT include duplicated
// 	// parameters. A unique parameter is defined by a combination of a name and
// 	// location. The list can use the Reference Object to link to parameters
// 	// that are defined at the OpenAPI Object's components/parameters.
// 	Parameters *ResolvedParameterSet `json:"parameters,omitempty"`
// 	Extensions `json:"-"`
// }

// // Kind returns KindResolvedPath
// func (rp *ResolvedPath) Kind() Kind {
// 	return KindResolvedPath
// }

// // ResolvedPathItems is a map of resolved Path objects
// type ResolvedPathItems map[string]*ResolvedPath

// // Kind returns KindResolvedPathItems
// func (rpi ResolvedPathItems) Kind() Kind {
// 	return KindResolvedPathItems
// }

// // Paths holds the relative paths to the individual endpoints and their
// // operations. The path is appended to the URL from the Server Object in order
// // to construct the full URL. The Paths MAY be empty, due to Access Control List
// // (ACL) constraints.
// type ResolvedPaths struct {
// 	Items      map[PathValue]*ResolvedPath `json:"-"`
// 	Extensions `json:"-"`
// }

// func (rp *ResolvedPaths) Kind() Kind {
// 	return KindResolvedPaths
// }

// // MarshalJSON marshals JSON
// func (p ResolvedPaths) MarshalJSON() ([]byte, error) {
// 	m := make(map[string]interface{}, len(p.Items)+len(p.Extensions))
// 	for k, v := range p.Items {
// 		m[k.String()] = v
// 	}
// 	for k, v := range p.Extensions {
// 		m[k] = v
// 	}
// 	return json.Marshal(m)
// }

// var _ Node = (*ResolvedPath)(nil)
// var _ Node = (*PathObj)(nil)
// var _ Node = (*PathItems)(nil)
// var _ Node = (*ResolvedPathItems)(nil)
// var _ Node = (PathItems)(nil)
// var _ Node = (ResolvedPathItems)(nil)
