package openapi

import (
	"encoding/json"

	"github.com/chanced/openapi/yamlutil"
)

// Webhooks is a map of Webhooks that can either be a Webhook or a Reference
type Webhooks map[string]Webhook

// Kind returns KindWebhooks
func (Webhooks) Kind() Kind {
	return KindWebhooks
}

func (ws *Webhooks) Len() int {
	if ws == nil || *ws == nil {
		return 0
	}
	return len(*ws)
}

func (ws *Webhooks) Get(key string) (Webhook, bool) {
	if ws == nil || *ws == nil {
		return nil, false
	}
	v, ok := (*ws)[key]
	return v, ok
}

func (ws *Webhooks) Set(key string, val Webhook) {
	if *ws == nil {
		*ws = Webhooks{
			key: val,
		}
		return
	}
	(*ws)[key] = val
}

func (ws Webhooks) Nodes() Nodes {
	if len(ws) == 0 {
		return nil
	}
	n := make(Nodes, len(ws))
	for k, v := range ws {
		n.maybeAdd(k, v, KindWebhook)
	}
	if len(n) == 0 {
		return nil
	}
	return n
}

// Webhook can either be a WebhookObj or a Reference
type Webhook interface {
	Node
	ResolveWebhook(func(ref string) (*WebhookObj, error)) (*WebhookObj, error)
}

type webhook pathobj

// WebhookObj is a PathObj
type WebhookObj PathObj

func (w *WebhookObj) Nodes() Nodes {
	return makeNodes(nodes{
		"get":        {w.Get, KindOperation},
		"put":        {w.Put, KindOperation},
		"post":       {w.Post, KindOperation},
		"delete":     {w.Delete, KindOperation},
		"options":    {w.Options, KindOperation},
		"head":       {w.Head, KindOperation},
		"patch":      {w.Patch, KindOperation},
		"trace":      {w.Trace, KindOperation},
		"servers":    {w.Servers, KindServers},
		"parameters": {w.Parameters, KindParameterSet},
	})
}

// Kind returns KindWebhook
func (*WebhookObj) Kind() Kind {
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

func (rw *ResolvedWebhook) Nodes() Nodes {
	return makeNodes(nodes{
		"get":        {rw.Get, KindResolvedOperation},
		"put":        {rw.Put, KindResolvedOperation},
		"post":       {rw.Post, KindResolvedOperation},
		"delete":     {rw.Delete, KindResolvedOperation},
		"options":    {rw.Options, KindResolvedOperation},
		"head":       {rw.Head, KindResolvedOperation},
		"patch":      {rw.Patch, KindResolvedOperation},
		"trace":      {rw.Trace, KindResolvedOperation},
		"servers":    {rw.Servers, KindServers},
		"parameters": {rw.Parameters, KindResolvedParameterSet},
	})
}

// Kind returns KindResolvedWebhook
func (*ResolvedWebhook) Kind() Kind {
	return KindResolvedWebhook
}

// ResolvedWebhooks is a map of *ResolvedWebhook
type ResolvedWebhooks map[string]*ResolvedWebhook

// Kind returns KindResolvedWebhooks
func (ResolvedWebhooks) Kind() Kind {
	return KindResolvedWebhooks
}

func (ws *ResolvedWebhooks) Len() int {
	if ws == nil || *ws == nil {
		return 0
	}
	return len(*ws)
}

func (ws *ResolvedWebhooks) Get(key string) (*ResolvedWebhook, bool) {
	if ws == nil || *ws == nil {
		return nil, false
	}
	v, ok := (*ws)[key]
	return v, ok
}

func (ws *ResolvedWebhooks) Set(key string, val *ResolvedWebhook) {
	if *ws == nil {
		*ws = ResolvedWebhooks{
			key: val,
		}
		return
	}
	(*ws)[key] = val
}

func (ws ResolvedWebhooks) Nodes() Nodes {
	if len(ws) == 0 {
		return nil
	}
	n := make(Nodes, len(ws))
	for k, v := range ws {
		n.maybeAdd(k, v, KindResolvedWebhook)
	}
	if len(n) == 0 {
		return nil
	}
	return n
}

// // Kind returns KindResolvedPath
// func (*ResolvedPath) Kind() Kind {
// 	return KindResolvedPath
// }

// // ResolvedPathItems is a map of resolved Path objects
// type ResolvedPathItems map[string]*ResolvedPath

// // Kind returns KindResolvedPathItems
// func (ResolvedPathItems) Kind() Kind {
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

// func (*ResolvedPaths) Kind() Kind {
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
