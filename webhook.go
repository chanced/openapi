package openapi

import (
	"encoding/json"

	"github.com/chanced/openapi/yamlutil"
)

// Path can either be a Path or a Reference
type Webhook interface {
	Node
	ResolveWebhook(func(ref string) (*WebhookObj, error)) (*WebhookObj, error)
}

type webhook pathobj

// Webhook is a PathObj
type WebhookObj PathObj

// KindPath returns KindWebhook
func (p *WebhookObj) Kind() Kind {
	return KindWebhook
}

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

// UnmarshalYAML unmarshals YAML data into rp
func (ws *Webhooks) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return yamlutil.Unmarshal(unmarshal, ws)
}

// MarshalYAML marshals rp into YAML
func (ws Webhooks) MarshalYAML() (interface{}, error) {
	b, err := json.Marshal(ws)
	if err != nil {
		return nil, err
	}
	var v interface{}
	err = json.Unmarshal(b, &v)
	return v, err
}

// ResolvedWebhook is a Webhook that has been fully resolved
type ResolvedWebhook ResolvedPath

func (*ResolvedWebhook) Kind() Kind {
	return KindResolvedWebhook
}

// ResolvedWebhooks is a map of *ResolvedWebhook
type ResolvedWebhooks map[string]*ResolvedWebhook

func (ResolvedWebhooks) Kind() Kind {
	return KindResolvedWebhooks
}

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

var (
	_ Node = (*WebhookObj)(nil)
	_ Node = (*Webhooks)(nil)
	_ Node = (*ResolvedWebhook)(nil)
	_ Node = (ResolvedWebhooks)(nil)
)
