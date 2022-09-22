package openapi

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/chanced/jsonx"
	"github.com/chanced/maps"
	"github.com/tidwall/gjson"
)

// Extensions for OpenAPI
//
// While the OpenAPI Specification tries to accommodate most use cases,
// additional data can be added to extend the specification at certain points.
//
// The extensions properties are implemented as patterned fields that are always
// prefixed by "x-".
//
// Field Pattern    Type    Description ^x-  Any Allows extensions to the
// OpenAPI Schema. The field name MUST begin with x-, for example,
// x-internal-id. Field names beginning x-oai- and x-oas- are reserved for uses
// defined by the OpenAPI Initiative. The value can be null, a primitive, an
// array or an object. The extensions may or may not be supported by the
// available tooling, but those may be extended as well to add requested support
// (if tools are internal or open-sourced).
//
// Security Filtering Some objects in the OpenAPI Specification MAY be declared
// and remain empty, or be completely removed, even though they are inherently
// the core of the API documentation.
//
// The reasoning is to allow an additional layer of access control over the
// documentation. While not part of the specification itself, certain libraries
// MAY choose to allow access to parts of the documentation based on some form
// of authentication/authorization.
//
// Two examples of this:
//
// The Paths Object MAY be present but empty. It may be counterintuitive, but
// this may tell the viewer that they got to the right place, but can't access
// any documentation. They would still have access to at least the Info Object
// which may contain additional information regarding authentication. The Path
// Item Object MAY be empty. In this case, the viewer will be aware that the
// path exists, but will not be able to see any of its operations or parameters.
// This is different from hiding the path itself from the Paths Object, because
// the user will be aware of its existence. This allows the documentation
// provider to finely control what the viewer can see.
type Extensions map[Text]jsonx.RawMessage

type extended interface {
	exts() Extensions
}

type extender interface {
	setExts(Extensions)
}

// Decode decodes all extensions into dst.
func (e Extensions) DecodeExtensions(dst interface{}) error {
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dst)
}

// DecodeExtension decodes extension at key into dst.
func (e Extensions) DecodeExtension(key Text, dst interface{}) error {
	if !key.HasPrefix("x-") {
		key = "x-" + key
	}
	return json.Unmarshal(e[key], dst)
}

func (e Extensions) exts() Extensions { return e }

func (e *Extensions) setExts(v Extensions) { *e = v }

// SetExtension encodes val and sets the result to key
func (e *Extensions) SetExtension(key Text, val interface{}) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	e.SetRawExtension(key, data)
	return nil
}

// SetRawExtension sets the raw JSON encoded val to key
func (e *Extensions) SetRawExtension(key Text, val []byte) {
	if !key.HasPrefix("x-") {
		key = "x-" + key
	}
	(*e)[key] = val
}

// Extension returns an extension by name
func (e Extensions) Extension(key Text) (interface{}, bool) {
	if !key.HasPrefix("x-") {
		key = "x-" + key
	}
	v, exists := e[key]
	return v, exists
}

// IsExtensionKey returns true if the key starts with "x-"
func IsExtensionKey(key Text) bool {
	return key.HasPrefix("x-")
}

func unmarshalExtendedJSON(data []byte, dst extender) error {
	ev := Extensions{}
	if err := json.Unmarshal(data, dst); err != nil {
		return err
	}
	gjson.ParseBytes(data).ForEach(func(key, value gjson.Result) bool {
		if IsExtensionKey(Text(key.String())) {
			ev[Text(key.String())] = jsonx.RawMessage(value.Raw)
		}
		return true
	})
	dst.setExts(ev)
	return nil
}

func marshalExtendedJSON(dst extended) ([]byte, error) {
	data, err := json.Marshal(dst)
	if err != nil {
		return nil, err
	}
	if !jsonx.IsObject(data) {
		// this shouldn't happen
		return nil, fmt.Errorf("openapi: cannot marshal extensions into non-object")
	}

	b := bytes.Buffer{}
	b.Write(data[:len(data)-1])
	return marshalExtensionsInto(&b, dst.exts())
}

func marshalExtensionsInto(b *bytes.Buffer, e Extensions) ([]byte, error) {
	var err error
	for _, kv := range maps.SortByKeys(e) {
		if b.Len() > 1 {
			b.WriteByte(',')
		}
		jsonx.EncodeAndWriteString(b, kv.Key)
		b.WriteByte(':')
		b.Write(kv.Value)
		if err != nil {
			return nil, err
		}
	}
	b.WriteByte('}')
	return b.Bytes(), nil
}
