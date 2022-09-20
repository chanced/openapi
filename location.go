package openapi

import (
	"github.com/chanced/jsonpointer"
	"github.com/chanced/uri"
)

// TODO: relToRes needs to be a slice

func NewLocation(uri uri.URI) (Location, error) {
	ptr, err := jsonpointer.Parse(uri.Fragment)
	if err != nil {
		return Location{}, err
	}
	loc := Location{
		absolute: uri,
		relative: ptr,
	}
	return loc, nil
}

type Location struct {
	absolute uri.URI
	relative jsonpointer.Pointer
}

func (l Location) String() string {
	return l.absolute.String()
}

func (l Location) AbsoluteLocation() uri.URI {
	return l.absolute
}

// RelativeLocation returns a jsonpointer.Pointer of the path from the
// containing resource file.
func (l Location) RelativeLocation() jsonpointer.Pointer {
	return l.relative
}

func (l Location) AppendLocation(p string) Location {
	l.relative = l.relative.AppendString(p)
	l.absolute.Fragment = l.relative.String()
	l.absolute.RawFragment = l.relative.String()
	return l
}

func (l Location) withURI(uri *uri.URI) (Location, error) {
	l.absolute = *uri
	if len(l.absolute.Fragment) > 0 {
		var err error
		l.relative, err = jsonpointer.Parse(l.absolute.Fragment)
		if err != nil {
			return l, err
		}
	}
	// we dont know what this is yet
	l.relative = ""
	return l, nil
}

func (l Location) location() Location {
	return l
}

func (l Location) IsRelativeTo(uri *uri.URI) bool {
	if uri == nil {
		return false
	}
	a := l.absolute
	a.Fragment = ""
	a.RawFragment = ""
	u := *uri
	u.Fragment = ""
	u.RawFragment = ""

	return a.String() == u.String()
}
