package openapi

import (
	"github.com/chanced/jsonpointer"
	"github.com/chanced/uri"
)

type Location struct {
	absolute uri.URI
	relative jsonpointer.Pointer
}

func NewLocation(uri *uri.URI) (Location, error) {
	ptr, err := jsonpointer.Parse(uri.Fragment)
	if err != nil {
		return Location{}, err
	}
	loc := Location{
		absolute: *uri,
		relative: ptr,
	}
	return loc, nil
}

func (l Location) String() string {
	return l.absolute.String()
}

func (l Location) Absolute() uri.URI {
	return l.absolute
}

func (l Location) Relative() jsonpointer.Pointer {
	return l.relative
}

func (l Location) Append(p string) Location {
	l.relative = l.relative.AppendString(p)
	l.absolute.Fragment = l.relative.String()
	l.absolute.RawFragment = l.relative.String()
	return l
}
