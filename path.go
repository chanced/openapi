package openapi

import (
	"encoding/json"
	"strings"

	"github.com/chanced/openapi/yamlutil"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v2"
)

// Path can either be a Path or a Reference
type Path interface {
	Node
	ResolvePath(func(ref string) (*PathObj, error)) (*PathObj, error)
}

// PathValue is relative path to an individual endpoint. The path is appended
// (no relative URL resolution) to the expanded URL from the Server Object's url
// field in order to construct the full URL. PathValue templating is allowed. When
// matching URLs, concrete (non-templated) paths would be matched before their
// templated counterparts. Templated paths with the same hierarchy but different
// templated names MUST NOT exist as they are identical. In case of ambiguous
// matching, it's up to the tooling to decide which one to use.
type PathValue string

func (pv PathValue) String() string {
	str := string(pv)
	if len(pv) == 0 {
		return "/"
	}
	if pv[0] != '/' {
		return "/" + str
	}
	return str
}

// // Params returns all params in the path
// func (pv PathValue) Params() []string {
// }

// MarshalJSON Marshals PathEntry to JSON
func (pv PathValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(pv.String())
}

// MarshalYAML Marshals PathEntry to YAML
func (pv PathValue) MarshalYAML() ([]byte, error) {
	return yaml.Marshal(pv.String())
}

func (pv PathValue) MarshalText() ([]byte, error) {
	return []byte(pv), nil
}

func (pv *PathValue) UnmarshalText(txt []byte) error {
	*pv = PathValue(txt)
	return nil
}

// Paths holds the relative paths to the individual endpoints and their
// operations. The path is appended to the URL from the Server Object in order
// to construct the full URL. The Paths MAY be empty, due to Access Control List
// (ACL) constraints.
type Paths struct {
	Items      map[PathValue]*PathObj `json:"-"`
	Extensions `json:"-"`
}

func (Paths) Kind() Kind { return KindPaths }

func (ps *Paths) Len() int {
	if ps == nil || ps.Items == nil {
		return 0
	}
	return len(ps.Items)
}

func (ps *Paths) Get(key string) (*PathObj, bool) {
	if ps == nil || ps.Items == nil {
		return nil, false
	}
	v, ok := ps.Items[PathValue(key)]
	return v, ok
}

func (ps *Paths) Set(key string, val *PathObj) {
	if ps == nil || ps.Items == nil {
		*ps = Paths{
			Items: map[PathValue]*PathObj{
				PathValue(key): val,
			},
		}
		return
	}
	ps.Items[PathValue(key)] = val
}

func (ps *Paths) Nodes() Nodes {
	if ps.Len() == 0 {
		return nil
	}
	nl := make(Nodes, ps.Len())
	for i, v := range ps.Items {
		nl.maybeAdd(i, v, KindPath)
	}
	return nl
}

// MarshalJSON marshals JSON
func (p Paths) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{}, len(p.Items)+len(p.Extensions))
	for k, v := range p.Items {
		m[k.String()] = v
	}
	for k, v := range p.Extensions {
		m[k] = v
	}
	return json.Marshal(m)
}

// UnmarshalJSON unmarshals JSON data into p
func (p *Paths) UnmarshalJSON(data []byte) error {
	*p = Paths{
		Items:      map[PathValue]*PathObj{},
		Extensions: Extensions{},
	}
	var err error
	gjson.ParseBytes(data).ForEach(func(key, value gjson.Result) bool {
		if strings.HasPrefix(key.String(), "x-") {
			p.SetEncodedExtension(key.String(), []byte(value.Raw))
		} else {
			var v PathObj
			err = json.Unmarshal([]byte(value.Raw), &v)
			p.Items[PathValue(key.String())] = &v
		}
		return err == nil
	})
	return err
}

// PathItems is a map of Paths that can either be a Path or a Reference
type PathItems map[string]Path

func (pi *PathItems) Len() int {
	if pi == nil || *pi == nil {
		return 0
	}
	return len(*pi)
}

func (pi *PathItems) Get(key string) (Path, bool) {
	if pi.Len() == 0 {
		return nil, false
	}
	v, ok := (*pi)[key]
	return v, ok
}

func (pi *PathItems) Set(key string, val Path) {
	if *pi == nil {
		*pi = PathItems{
			key: val,
		}
		return
	}
	(*pi)[key] = val
}

func (pi *PathItems) Delete(key string) {
	if pi == nil || *pi == nil {
		return
	}
	delete(*pi, key)
}

func (pi PathItems) Nodes() Nodes {
	if pi.Len() == 0 {
		return nil
	}
	nodes := make(Nodes, pi.Len())

	for k, v := range pi {
		nodes.maybeAdd(k, v, KindPath)
	}
	return nodes
}

// Kind returns KindPathItems
func (PathItems) Kind() Kind {
	return KindPathItems
}

// UnmarshalJSON unmarshals JSON data into pi
func (pi *PathItems) UnmarshalJSON(data []byte) error {
	var rd map[string]json.RawMessage
	err := json.Unmarshal(data, &rd)
	if err != nil {
		return err
	}
	res := PathItems{}
	for k, d := range rd {
		if isRefJSON(data) {
			var v Reference
			if err = json.Unmarshal(d, &v); err != nil {
				return err
			}
			res[k] = &v
		} else {
			var v PathObj
			if err = json.Unmarshal(d, &v); err != nil {
				return err
			}
			res[k] = &v
		}
	}
	*pi = res
	return nil
}

// UnmarshalYAML unmarshals YAML data into pi
func (pi *PathItems) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return yamlutil.Unmarshal(unmarshal, pi)
}

// MarshalYAML marshals pi into YAML
func (pi PathItems) MarshalYAML() (interface{}, error) {
	b, err := json.Marshal(pi)
	if err != nil {
		return nil, err
	}
	var v interface{}
	err = json.Unmarshal(b, &v)
	return v, err
}

// PathObj describes the operations available on a single path. A PathObj Item MAY
// be empty, due to ACL constraints. The path itself is still exposed to the
// documentation viewer but they will not know which operations and parameters
// are available.
type PathObj struct {
	// Allows for a referenced definition of this path item. The referenced
	// structure MUST be in the form of a Path Item Object. In case a Path Item
	// Object field appears both in the defined object and the referenced
	// object, the behavior is undefined. See the rules for resolving Relative
	// References.
	Ref string `json:"$ref,omitempty"`
	// An optional, string summary, intended to apply to all operations in this path.
	Summary string `json:"summary,omitempty"`
	// An optional, string description, intended to apply to all operations in
	// this path. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`
	// A definition of a GET operation on this path.
	Get *Operation `json:"get,omitempty"`
	// A definition of a PUT operation on this path.
	Put *Operation `json:"put,omitempty"`
	// A definition of a POST operation on this path.
	Post *Operation `json:"post,omitempty"`
	// A definition of a DELETE operation on this path.
	Delete *Operation `json:"delete,omitempty"`
	// A definition of a OPTIONS operation on this path.
	Options *Operation `json:"options,omitempty"`
	// A definition of a HEAD operation on this path.
	Head *Operation `json:"head,omitempty"`
	// A definition of a PATCH operation on this path.
	Patch *Operation `json:"patch,omitempty"`
	// A definition of a TRACE operation on this path.
	Trace *Operation `json:"trace,omitempty"`
	// An alternative server array to service all operations in this path.
	Servers Servers `json:"servers,omitempty"`
	// A list of parameters that are applicable for all the operations described
	// under this path. These parameters can be overridden at the operation
	// level, but cannot be removed there. The list MUST NOT include duplicated
	// parameters. A unique parameter is defined by a combination of a name and
	// location. The list can use the Reference Object to link to parameters
	// that are defined at the OpenAPI Object's components/parameters.
	Parameters *ParameterSet `json:"parameters,omitempty"`
	Extensions `json:"-"`
}

func (p *PathObj) Nodes() Nodes {
	return makeNodes(nodes{
		"get":        {p.Get, KindOperation},
		"put":        {p.Put, KindOperation},
		"post":       {p.Post, KindOperation},
		"delete":     {p.Delete, KindOperation},
		"options":    {p.Options, KindOperation},
		"head":       {p.Head, KindOperation},
		"patch":      {p.Patch, KindOperation},
		"trace":      {p.Trace, KindOperation},
		"servers":    {p.Servers, KindServers},
		"parameters": {p.Parameters, KindParameterSet},
	})
}

// Kind returns KindPath
func (*PathObj) Kind() Kind {
	return KindPath
}

type pathobj PathObj

// ResolvePath resolves PathObj by returning itself. resolve is  not called.
func (p *PathObj) ResolvePath(func(ref string) (*PathObj, error)) (*PathObj, error) {
	return p, nil
}

// MarshalJSON marshals p into JSON
func (p PathObj) MarshalJSON() ([]byte, error) {
	return marshalExtendedJSON(pathobj(p))
}

// UnmarshalJSON unmarshals json into p
func (p *PathObj) UnmarshalJSON(data []byte) error {
	var v pathobj
	if err := unmarshalExtendedJSON(data, &v); err != nil {
		return err
	}
	*p = PathObj(v)
	return nil
}

// MarshalYAML first marshals and unmarshals into JSON and then marshals into
// YAML
func (p PathObj) MarshalYAML() (interface{}, error) {
	return yamlutil.Marshal(p)
}

// UnmarshalYAML unmarshals yaml into s
func (p *PathObj) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return yamlutil.Unmarshal(unmarshal, p)
}

func unmarshalPathJSON(data []byte) (Path, error) {
	if isRefJSON(data) {
		return unmarshalReferenceJSON(data)
	}
	var p pathobj
	err := json.Unmarshal(data, &p)
	v := PathObj(p)
	return &v, err
}

// ResolvedPath is a Path Object which has beeen resolved. It describes the
// operations available on a single path. A PathObj Item MAY be empty, due to
// ACL constraints. The path itself is still exposed to the documentation viewer
// but they will not know which operations and parameters are available.
type ResolvedPath struct {
	// Allows for a referenced definition of this path item. The referenced
	// structure MUST be in the form of a Path Item Object. In case a Path Item
	// Object field appears both in the defined object and the referenced
	// object, the behavior is undefined. See the rules for resolving Relative
	// References.
	Ref string `json:"$ref,omitempty"`
	// An optional, string summary, intended to apply to all operations in this path.
	Summary string `json:"summary,omitempty"`
	// An optional, string description, intended to apply to all operations in
	// this path. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`
	// A definition of a GET operation on this path.
	Get *ResolvedOperation `json:"get,omitempty"`
	// A definition of a PUT operation on this path.
	Put *ResolvedOperation `json:"put,omitempty"`
	// A definition of a POST operation on this path.
	Post *ResolvedOperation `json:"post,omitempty"`
	// A definition of a DELETE operation on this path.
	Delete *ResolvedOperation `json:"delete,omitempty"`
	// A definition of a OPTIONS operation on this path.
	Options *ResolvedOperation `json:"options,omitempty"`
	// A definition of a HEAD operation on this path.
	Head *ResolvedOperation `json:"head,omitempty"`
	// A definition of a PATCH operation on this path.
	Patch *ResolvedOperation `json:"patch,omitempty"`
	// A definition of a TRACE operation on this path.
	Trace *ResolvedOperation `json:"trace,omitempty"`
	// An alternative server array to service all operations in this path.
	Servers Servers `json:"servers,omitempty"`
	// A list of parameters that are applicable for all the operations described
	// under this path. These parameters can be overridden at the operation
	// level, but cannot be removed there. The list MUST NOT include duplicated
	// parameters. A unique parameter is defined by a combination of a name and
	// location. The list can use the Reference Object to link to parameters
	// that are defined at the OpenAPI Object's components/parameters.
	Parameters *ResolvedParameterSet `json:"parameters,omitempty"`
	Extensions `json:"-"`
}

func (rp *ResolvedPath) Nodes() Nodes {
	return makeNodes(nodes{
		"get":        {rp.Get, KindResolvedOperation},
		"put":        {rp.Put, KindResolvedOperation},
		"post":       {rp.Post, KindResolvedOperation},
		"delete":     {rp.Delete, KindResolvedOperation},
		"options":    {rp.Options, KindResolvedOperation},
		"head":       {rp.Head, KindResolvedOperation},
		"patch":      {rp.Patch, KindResolvedOperation},
		"trace":      {rp.Trace, KindResolvedOperation},
		"servers":    {rp.Servers, KindServers},
		"parameters": {rp.Parameters, KindResolvedParameterSet},
	})
}

// Kind returns KindResolvedPath
func (*ResolvedPath) Kind() Kind {
	return KindResolvedPath
}

// ResolvedPathItems is a map of resolved Path objects
type ResolvedPathItems map[string]*ResolvedPath

func (rpi *ResolvedPathItems) Len() int {
	if rpi == nil || *rpi == nil {
		return 0
	}
	return len(*rpi)
}

func (rpi *ResolvedPathItems) Get(key string) (*ResolvedPath, bool) {
	if rpi.Len() == 0 {
		return nil, false
	}
	v, ok := (*rpi)[key]
	return v, ok
}

func (rpi *ResolvedPathItems) Set(key string, val *ResolvedPath) {
	if *rpi == nil {
		*rpi = ResolvedPathItems{
			key: val,
		}
		return
	}
	(*rpi)[key] = val
}

func (rpi *ResolvedPathItems) Delete(key string) {
	if rpi == nil || *rpi == nil {
		return
	}
	delete(*rpi, key)
}

func (rpi ResolvedPathItems) Nodes() Nodes {
	if rpi.Len() == 0 {
		return nil
	}
	n := make(Nodes, rpi.Len())
	for k, v := range rpi {
		n.maybeAdd(k, v, KindResolvedPath)
	}
	return n
}

// Kind returns KindResolvedPathItems
func (ResolvedPathItems) Kind() Kind {
	return KindResolvedPathItems
}

// ResolvedPaths holds the relative paths to the individual endpoints and their
// operations. The path is appended to the URL from the Server Object in order
// to construct the full URL. The Paths MAY be empty, due to Access Control List
// (ACL) constraints.
type ResolvedPaths struct {
	Items      map[PathValue]*ResolvedPath `json:"-"`
	Extensions `json:"-"`
}

func (rps *ResolvedPaths) Len() int {
	if rps == nil || rps.Items == nil {
		return 0
	}
	return len(rps.Items)
}

func (rps *ResolvedPaths) Get(key string) (*ResolvedPath, bool) {
	if rps == nil || rps.Items == nil {
		return nil, false
	}
	v, ok := rps.Items[PathValue(key)]
	return v, ok
}

func (rps *ResolvedPaths) Set(key string, val *ResolvedPath) {
	if rps == nil || rps.Items == nil {
		*rps = ResolvedPaths{
			Items: map[PathValue]*ResolvedPath{
				PathValue(key): val,
			},
		}
		return
	}
	rps.Items[PathValue(key)] = val
}

func (rps *ResolvedPaths) Nodes() Nodes {
	if rps.Len() == 0 {
		return nil
	}
	nl := make(Nodes, rps.Len())
	for i, v := range rps.Items {
		nl.maybeAdd(i, v, KindResolvedPath)
	}
	return nl
}

// Kind returns KindResolvedPaths
func (*ResolvedPaths) Kind() Kind {
	return KindResolvedPaths
}

// MarshalJSON marshals JSON
func (rp ResolvedPaths) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{}, len(rp.Items)+len(rp.Extensions))
	for k, v := range rp.Items {
		m[k.String()] = v
	}
	for k, v := range rp.Extensions {
		m[k] = v
	}
	return json.Marshal(m)
}

var (
	_ Node = (*PathObj)(nil)
	_ Node = (*PathItems)(nil)
	_ Node = (PathItems)(nil)
	_ Node = (*Paths)(nil)
	_ Node = (*ResolvedPath)(nil)
	_ Node = (*ResolvedPaths)(nil)
	_ Node = (ResolvedPathItems)(nil)
)
