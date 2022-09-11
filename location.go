package openapi

import (
	"github.com/chanced/jsonpointer"
	"github.com/chanced/uri"
)

type Location struct {
	absolute             uri.URI
	relativeFromDocument jsonpointer.Pointer
	relativeFromResource jsonpointer.Pointer
}

func NewLocation(uri *uri.URI) (Location, error) {
	ptr, err := jsonpointer.Parse(uri.Fragment)
	if err != nil {
		return Location{}, err
	}
	loc := Location{
		absolute:             *uri,
		relativeFromDocument: ptr,
	}
	return loc, nil
}

func (l Location) String() string {
	return l.absolute.String()
}

func (l Location) Absolute() uri.URI {
	return l.absolute
}

// RelativeFromDocument returns a jsonpointer.Pointer of the relative path from
// the root OpenAPI Document to the current Node.
func (l Location) RelativeFromDocument() jsonpointer.Pointer {
	return l.relativeFromDocument
}

// RelativeFromResource returns a jsonpointer.Pointer of the path from the
// resource file that it is in.
func (l Location) RelativeFromResource() jsonpointer.Pointer {
	return l.relativeFromResource
}

func (l Location) Append(p string) Location {
	l.relativeFromDocument = l.relativeFromDocument.AppendString(p)
	l.absolute.Fragment = l.relativeFromDocument.String()
	l.absolute.RawFragment = l.relativeFromDocument.String()
	return l
}

func (l Location) WithURI(uri *uri.URI) (Location, error) {
	l.absolute = *uri
	if len(l.absolute.Fragment) > 0 {
		var err error
		l.relativeFromDocument, err = jsonpointer.Parse(l.absolute.Fragment)
		if err != nil {
			return l, err
		}
	}
	l.relativeFromDocument = jsonpointer.Root
	return l, nil
}
