package openapi

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/chanced/jsonpointer"
	"github.com/chanced/uri"
)

var (
	ErrEmptyRef      = errors.New("openapi: empty $ref")
	ErrNotFound      = fmt.Errorf("openapi: component not found")
	ErrNotResolvable = errors.New("openapi: pointer path not resolvable")

	ErrMissingOpenAPIVersion = errors.New("openapi: missing openapi version")

	// ErrInvalidSemVer is returned a version is found to be invalid when
	// being parsed.
	ErrInvalidSemVer = errors.New("invalid semantic version")

	// ErrInvalidMetadata is returned when the metadata of a semver is an invalid format
	ErrInvalidSemVerMetadata = errors.New("invalid semantic version metadata string")

	// ErrInvalidPrerelease is returned when the pre-release of a semver is an invalid format
	ErrInvalidSemVerPrerelease = errors.New("invalid semantic version prerelease string")

	ErrInvalidResolution = errors.New("openapi: invalid resolution")
)

type Error struct {
	Err         error
	ResourceURI uri.URI
}

func NewError(err error, resource uri.URI) error {
	return &Error{
		Err:         err,
		ResourceURI: resource,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: [%q]", e.Err, e.ResourceURI.String())
}

func (e *Error) Unwrap() error {
	return e.Err
}

func newErrNotFound(uri uri.URI, tok jsonpointer.Token) error {
	return NewError(fmt.Errorf("%w: %q", ErrNotFound, tok), uri)
}

func newErrNotResolvable(uri uri.URI, tok jsonpointer.Token) error {
	return NewError(fmt.Errorf("%w: %q", ErrNotResolvable, tok), uri)
}

type UnsupportedVersionError struct {
	Version string  `json:"version"`
	Errs    []error `json:"errors"`
}

func (e *UnsupportedVersionError) Error() string {
	b := strings.Builder{}
	b.WriteString("openapi: unsupported version:")
	b.WriteString(fmt.Sprintf(" %q:", e.Version))
	for _, err := range e.Errs {
		b.WriteString(fmt.Sprintf("\n- %s", err))
	}
	return b.String()
}

func (e *UnsupportedVersionError) As(target interface{}) bool {
	for _, v := range e.Errs {
		if errors.As(v, target) {
			return true
		}
	}
	return false
}

func (e *UnsupportedVersionError) Is(err error) bool {
	for _, v := range e.Errs {
		if errors.Is(v, err) {
			return true
		}
	}
	return false
}

type ValidationError struct {
	Kind Kind
	Err  error
	URI  uri.URI
}

func NewValidationError(err error, kind Kind, resource uri.URI) error {
	var ve *ValidationError
	if errors.As(err, &ve) {
		return ve
	}

	return NewError(&ValidationError{
		Kind: kind,
		Err:  err,
		URI:  resource,
	}, resource)
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("openapi: error validating %s %s: %s", e.Kind, e.URI, e.Err)
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}

type SemVerError struct {
	Value string
	Err   error
	URI   uri.URI
}

func (e SemVerError) Error() string {
	return fmt.Sprintf("openapi: error parsing semver %q for %s: %v", e.Value, e.URI, e.Err)
}

func (e SemVerError) Unwrap() error {
	return e.Err
}

func (e SemVerError) Is(err error) bool {
	if errors.Is(e.Err, err) {
		return true
	}
	switch {
	case errors.Is(err, ErrInvalidSemVer):
		return errors.Is(err, semver.ErrInvalidSemVer)
	case errors.Is(e.Err, ErrInvalidSemVerMetadata):
		return errors.Is(err, semver.ErrInvalidMetadata)
	case errors.Is(e.Err, ErrInvalidSemVerPrerelease):
		return errors.Is(err, semver.ErrInvalidPrerelease)
	}
	return false
}

func NewSemVerError(err error, value string, uri uri.URI) error {
	if err == nil {
		return nil
	}
	var sv *SemVerError
	if errors.As(err, &sv) {
		return sv
	}
	return NewError(&SemVerError{
		Value: value,
		Err:   translateSemVerErr(err),
		URI:   uri,
	}, uri)
}

func translateSemVerErr(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, semver.ErrInvalidSemVer):
		return ErrInvalidSemVer
	case errors.Is(err, semver.ErrInvalidMetadata):
		return ErrInvalidSemVerMetadata
	case errors.Is(err, semver.ErrInvalidPrerelease):
		return ErrInvalidSemVerPrerelease
	}
	return err
}

type ResolutionError struct {
	URI      uri.URI
	Expected Kind
	Actual   Kind
	RefType  RefType
}

func (e *ResolutionError) Error() string {
	return fmt.Sprintf("%v cannot resolve %s to %s for %s: %s", ErrInvalidResolution, e.Actual, e.Expected, e.RefType, e.URI)
}

func (e *ResolutionError) Unwrap() error {
	return ErrInvalidResolution
}

func NewResolutionError(r Ref, expected, actual Kind) error {
	return &ResolutionError{
		URI:      r.AbsoluteLocation(),
		Actual:   actual,
		Expected: expected,
		RefType:  r.RefType(),
	}
}
