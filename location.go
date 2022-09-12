package openapi

import (
	"github.com/chanced/jsonpointer"
	"github.com/chanced/uri"
)

// TODO: relToRes needs to be a slice

func NewLocation(uri *uri.URI) (Location, error) {
	ptr, err := jsonpointer.Parse(uri.Fragment)
	if err != nil {
		return Location{}, err
	}
	loc := Location{
		absolute: *uri,
		relToRes: ptr,
	}
	return loc, nil
}

type Location struct {
	absolute uri.URI
	relToDoc jsonpointer.Pointer
	relToRes jsonpointer.Pointer
}

func (l Location) String() string {
	return l.absolute.String()
}

func (l Location) AbsoluteLocation() uri.URI {
	return l.absolute
}

// LocationRelativeToDocument returns a jsonpointer.Pointer of the relative path from
// the root OpenAPI Document to the current Node.
func (l Location) FirstPathRelativeToDocument() jsonpointer.Pointer {
	return l.relToDoc
}

// LocationRelativeToResource returns a jsonpointer.Pointer of the path from the
// resource file that it is in.
func (l Location) PathRelativeToResource() jsonpointer.Pointer {
	return l.relToRes
}

func (l Location) Append(p string) Location {
	l.relToDoc = l.relToDoc.AppendString(p)
	l.absolute.Fragment = l.relToDoc.String()
	l.absolute.RawFragment = l.relToDoc.String()
	return l
}

func (l Location) withURI(uri *uri.URI) (Location, error) {
	l.absolute = *uri
	if len(l.absolute.Fragment) > 0 {
		var err error
		l.relToRes, err = jsonpointer.Parse(l.absolute.Fragment)
		if err != nil {
			return l, err
		}
	}
	// we dont know what this is yet
	l.relToDoc = ""
	return l, nil
}

func (l Location) location() Location {
	return l
}
