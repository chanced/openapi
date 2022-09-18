package openapi

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/chanced/jsonpointer"
	"github.com/chanced/transcode"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v3"
)

type PathItemEntry struct {
	Key      string
	PathItem *PathItem
}

type PathItemObjs = ObjMap[*PathItem]

// PathItemMap is a map of Paths that can either be a Path or a Reference
type PathItemMap = ComponentMap[*PathItem]

// Paths holds the relative paths to the individual endpoints and their
// operations. The path is appended to the URL from the Server Object in order
// to construct the full URL. The Paths MAY be empty, due to Access Control List
// (ACL) constraints.
type Paths struct {
	Location   `json:"-"`
	Extensions `json:"-"`

	// Items are the Path
	Items *PathItemObjs `json:"-"`
}

func (p *Paths) Edges() []Node {
	if p == nil {
		return nil
	}
	return downcastNodes(p.edges())
}

func (p *Paths) edges() []node {
	if p == nil {
		return nil
	}
	return appendEdges(nil, p.Items)
}

func (*Paths) ref() Ref { return nil }

func (p *Paths) Refs() []Ref {
	if p == nil {
		return nil
	}
	return p.Items.Refs()
}

func (p *Paths) Anchors() (*Anchors, error) {
	if p == nil {
		return nil, nil
	}
	return p.Items.Anchors()
}

func (p *Paths) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
	if err := ptr.Validate(); err != nil {
		return nil, err
	}
	return p.resolveNodeByPointer(ptr)
}

func (p *Paths) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
	if ptr.IsRoot() {
		return p, nil
	}
	nxt, tok, _ := ptr.Next()
	v := p.Items.Get(Text(tok))
	if v == nil {
		return nil, newErrNotFound(p.Location.AbsolutePath(), tok)
	}
	return v.resolveNodeByPointer(nxt)
}

func (p *Paths) isNil() bool { return p == nil }

func (*Paths) Kind() Kind      { return KindPaths }
func (*Paths) mapKind() Kind   { return KindUndefined }
func (*Paths) sliceKind() Kind { return KindUndefined }

func (p *Paths) setLocation(loc Location) error {
	if p == nil {
		return nil
	}
	p.Location = loc
	return p.Items.setLocation(loc)
}

// MarshalJSON marshals JSON
func (p Paths) MarshalJSON() ([]byte, error) {
	j, err := p.Items.MarshalJSON()
	if err != nil {
		return nil, err
	}
	b := bytes.Buffer{}
	// removing the last } as marshalExtensionsInto execpts a buffer without it
	b.Write(j[:len(j)-1])
	return marshalExtensionsInto(&b, p.Extensions)
}

// UnmarshalJSON unmarshals JSON data into p
func (p *Paths) UnmarshalJSON(data []byte) error {
	*p = Paths{
		Extensions: Extensions{},
	}
	var err error
	gjson.ParseBytes(data).ForEach(func(key, value gjson.Result) bool {
		if strings.HasPrefix(key.String(), "x-") {
			p.SetRawExtension(key.String(), []byte(value.Raw))
		} else {
			var v PathItem
			err = json.Unmarshal([]byte(value.Raw), &v)
			p.Items.Set(Text(key.String()), &v)
		}
		return err == nil
	})
	return err
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Marshaler interface
func (p Paths) MarshalYAML() (interface{}, error) {
	j, err := p.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return transcode.YAMLFromJSON(j)
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Unmarshaler interface
func (p *Paths) UnmarshalYAML(value *yaml.Node) error {
	j, err := transcode.YAMLFromJSON([]byte(value.Value))
	if err != nil {
		return err
	}
	return json.Unmarshal(j, p)
}

var _ node = (*Paths)(nil)
