package uri

import (
	"encoding"
	"errors"
	"net/url"
	"strings"
)

// Error reports an error and the operation and URI that caused it.
type Error url.Error

func (e *Error) Unwrap() error   { return (*url.Error)(e) }
func (e *Error) Error() string   { return errfmt.Replace((*url.Error)(e).Error()) }
func (e *Error) Timeout() bool   { return (*url.Error)(e).Timeout() }
func (e *Error) Temporary() bool { return (*url.Error)(e).Temporary() }

type EscapeError url.EscapeError

func (e EscapeError) Error() string { return errfmt.Replace(url.EscapeError(e).Error()) }
func (e EscapeError) Unwrap() error { return url.EscapeError(e) }

type InvalidHostError url.InvalidHostError

func (e InvalidHostError) Error() string { return errfmt.Replace(url.InvalidHostError(e).Error()) }
func (e InvalidHostError) Unwrap() error { return url.InvalidHostError(e) }

func errOrNil(err error) error {
	if err == nil {
		return nil
	}
	var escapeErr url.EscapeError
	if errors.As(err, &escapeErr) {
		return EscapeError(escapeErr)
	}

	var invalidHostErr url.InvalidHostError
	if errors.As(err, &invalidHostErr) {
		return InvalidHostError(invalidHostErr)
	}

	var e *url.Error
	if errors.As(err, &e) {
		e := Error(*e)
		return &e
	}

	return errors.New(errfmt.Replace(err.Error()))
}

var errfmt = strings.NewReplacer("net/url", "uri", "URL", "URI")

// QueryUnescape does the inverse transformation of QueryEscape,
// converting each 3-byte encoded substring of the form "%AB" into the
// hex-decoded byte 0xAB.
// It returns an error if any % is not followed by two hexadecimal
// digits.
func QueryUnescape(s string) (string, error) {
	res, err := url.QueryUnescape(s)
	return res, errOrNil(err)
}

// PathUnescape does the inverse transformation of PathEscape,
// converting each 3-byte encoded substring of the form "%AB" into the
// hex-decoded byte 0xAB. It returns an error if any % is not followed
// by two hexadecimal digits.
//
// PathUnescape is identical to QueryUnescape except that it does not
// unescape '+' to ' ' (space).
func PathUnescape(s string) (string, error) {
	res, err := url.PathUnescape(s)
	return res, errOrNil(err)
}

// QueryEscape escapes the string so it can be safely placed
// inside a URI query.
func QueryEscape(s string) string {
	return url.QueryEscape(s)
}

// PathEscape escapes the string so it can be safely placed inside a URI path segment,
// replacing special characters (including /) with %XX sequences as needed.
func PathEscape(s string) string {
	return url.PathEscape(s)
}

// A URI represents a parsed URI reference. It is a newtype of url.URL.
//
// The general form represented is:
//
//	[scheme:][//[userinfo@]host][/]path[?query][#fragment]
//
// URIs that do not start with a slash after the scheme are interpreted as:
//
//	scheme:opaque[?query][#fragment]
//
// Note that the Path field is stored in decoded form: /%47%6f%2f becomes /Go/.
// A consequence is that it is impossible to tell which slashes in the Path were
// slashes in the raw URI and which were %2f. This distinction is rarely important,
// but when it is, the code should use RawPath, an optional field which only gets
// set if the default encoding is different from Path.
//
// URI's String method uses the EscapedPath method to obtain the path. See the
// EscapedPath method for more details.
type URI url.URL

// User returns a Userinfo containing the provided username
// and no password set.
func User(username string) *url.Userinfo {
	return url.User(username)
}

// UserPassword returns a Userinfo containing the provided username
// and password.
//
// This functionality should only be used with legacy web sites.
// RFC 2396 warns that interpreting Userinfo this way
// “is NOT RECOMMENDED, because the passing of authentication
// information in clear text (such as URI) has proven to be a
// security risk in almost every case where it has been used.”
func UserPassword(username, password string) *Userinfo {
	ui := Userinfo(*url.UserPassword(username, password))
	return &ui
}

// The Userinfo type is an immutable encapsulation of username and
// password details for a URI. An existing Userinfo value is guaranteed
// to have a username set (potentially empty, as allowed by RFC 2396),
// and optionally a password.
type Userinfo = url.Userinfo

// Parse parses a raw url into a URI structure.
//
// The url may be relative (a path, without a host) or absolute
// (starting with a scheme). Trying to parse a hostname and path
// without a scheme is invalid but may not necessarily return an
// error, due to parsing ambiguities.
func Parse(rawURL string) (*URI, error) {
	url, err := url.Parse(rawURL)
	if url == nil {
		return nil, errOrNil(err)
	}
	uri := URI(*url)
	return &uri, errOrNil(err)
}

// ParseRequestURI parses a raw url into a URI structure. It assumes that
// url was received in an HTTP request, so the url is interpreted
// only as an absolute URI or an absolute path.
// The string url is assumed not to have a #fragment suffix.
// (Web browsers strip #fragment before sending the URL to a web server.)
func ParseRequestURI(rawURL string) (*URI, error) {
	url, err := url.ParseRequestURI(rawURL)
	if url == nil {
		return nil, errOrNil(err)
	}
	uri := URI(*url)
	return &uri, errOrNil(err)
}

// EscapedPath returns the escaped form of u.Path.
// In general there are multiple possible escaped forms of any path.
// EscapedPath returns u.RawPath when it is a valid escaping of u.Path.
// Otherwise EscapedPath ignores u.RawPath and computes an escaped
// form on its own.
// The String and RequestURI methods use EscapedPath to construct
// their results.
// In general, code should call EscapedPath instead of
// reading u.RawPath directly.
func (u *URI) EscapedPath() string {
	return (*url.URL)(u).EscapedPath()
}

// EscapedFragment returns the escaped form of u.Fragment.
// In general there are multiple possible escaped forms of any fragment.
// EscapedFragment returns u.RawFragment when it is a valid escaping of u.Fragment.
// Otherwise EscapedFragment ignores u.RawFragment and computes an escaped
// form on its own.
// The String method uses EscapedFragment to construct its result.
// In general, code should call EscapedFragment instead of
// reading u.RawFragment directly.
func (u *URI) EscapedFragment() string {
	return (*url.URL)(u).EscapedFragment()
}

// String reassembles the URI into a valid URI string.
// The general form of the result is one of:
//
//	scheme:opaque?query#fragment
//	scheme://userinfo@host/path?query#fragment
//
// If u.Opaque is non-empty, String uses the first form;
// otherwise it uses the second form.
// Any non-ASCII characters in host are escaped.
// To obtain the path, String uses u.EscapedPath().
//
// In the second form, the following rules apply:
//   - if u.Scheme is empty, scheme: is omitted.
//   - if u.User is nil, userinfo@ is omitted.
//   - if u.Host is empty, host/ is omitted.
//   - if u.Scheme and u.Host are empty and u.User is nil,
//     the entire scheme://userinfo@host/ is omitted.
//   - if u.Host is non-empty and u.Path begins with a /,
//     the form host/path does not add its own /.
//   - if u.RawQuery is empty, ?query is omitted.
//   - if u.Fragment is empty, #fragment is omitted.
func (u *URI) String() string {
	return (*url.URL)(u).String()
}

// Redacted is like String but replaces any password with "xxxxx".
// Only the password in u.URL is redacted.
func (u *URI) Redacted() string {
	return (*url.URL)(u).Redacted()
}

// Values maps a string key to a list of values.
// It is typically used for query parameters and form values.
// Unlike in the http.Header map, the keys in a Values map
// are case-sensitive.
type Values = url.Values

// ParseQuery parses the URL-encoded query string and returns
// a map listing the values specified for each key.
// ParseQuery always returns a non-nil map containing all the
// valid query parameters found; err describes the first decoding error
// encountered, if any.
//
// Query is expected to be a list of key=value settings separated by ampersands.
// A setting without an equals sign is interpreted as a key set to an empty
// value.
// Settings containing a non-URL-encoded semicolon are considered invalid.
func ParseQuery(query string) (Values, error) {
	q, err := url.ParseQuery(query)
	if q == nil {
		return nil, errOrNil(err)
	}
	return Values(q), errOrNil(err)
}

// IsAbs reports whether the URL is absolute.
// Absolute means that it has a non-empty scheme.
func (u *URI) IsAbs() bool {
	return (*url.URL)(u).IsAbs()
}

// Parse parses a URL in the context of the receiver. The provided URL
// may be relative or absolute. Parse returns nil, err on parse
// failure, otherwise its return value is the same as ResolveReference.
func (u *URI) Parse(ref string) (*URI, error) {
	url, err := (*url.URL)(u).Parse(ref)
	if url == nil {
		return nil, errOrNil(err)
	}
	uri := URI(*url)
	return &uri, errOrNil(err)
}

// ResolveReference resolves a URI reference to an absolute URI from
// an absolute base URI u, per RFC 3986 Section 5.2. The URI reference
// may be relative or absolute. ResolveReference always returns a new
// URL instance, even if the returned URL is identical to either the
// base or reference. If ref is an absolute URL, then ResolveReference
// ignores base and returns a copy of ref.
func (u *URI) ResolveReference(ref *URI) *URI {
	url := (*url.URL)(u).ResolveReference((*url.URL)(ref))
	uri := URI(*url)
	return &uri
}

// Query parses RawQuery and returns the corresponding values.
// It silently discards malformed value pairs.
// To check errors use ParseQuery.
func (u *URI) Query() Values {
	v, _ := ParseQuery(u.RawQuery)
	return v
}

// RequestURI returns the encoded path?query or opaque?query
// string that would be used in an HTTP request for u.
func (u *URI) RequestURI() string {
	return (*url.URL)(u).RequestURI()
}

// Hostname returns u.Host, stripping any valid port number if present.
//
// If the result is enclosed in square brackets, as literal IPv6 addresses are,
// the square brackets are removed from the result.
func (u *URI) Hostname() string {
	return (*url.URL)(u).Hostname()
}

// Port returns the port part of u.Host, without the leading colon.
//
// If u.Host doesn't contain a valid numeric port, Port returns an empty string.
func (u *URI) Port() string {
	return (*url.URL)(u).Port()
}

// Marshaling interface implementations.
// Would like to implement MarshalText/UnmarshalText but that will change the JSON representation of URLs.

func (u *URI) MarshalBinary() (text []byte, err error) {
	return (*url.URL)(u).MarshalBinary()
}

func (u *URI) UnmarshalBinary(text []byte) error {
	u1, err := Parse(string(text))
	if err != nil {
		return errOrNil(err)
	}
	*u = *u1
	return nil
}

// MarshalText implements encoding.TextMarshaler
func (u *URI) MarshalText() (text []byte, err error) {
	return []byte(u.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler
func (u *URI) UnmarshalText(text []byte) error {
	p, err := Parse(string(text))
	if err != nil {
		return errOrNil(err)
	}
	*u = *p
	return nil
}

// URL returns the
func (u *URI) URL() *url.URL {
	return (*url.URL)(u)
}

// JoinPath returns a new URL with the provided path elements joined to
// any existing path and the resulting path cleaned of any ./ or ../ elements.
// Any sequences of multiple / characters will be reduced to a single /.
func (u *URI) JoinPath(elem ...string) *URI {
	url := (*url.URL)(u).JoinPath(elem...)
	uri := URI(*url)
	return &uri
}

// JoinPath returns a URL string with the provided path elements joined to
// the existing path of base and the resulting path cleaned of any ./ or ../ elements.
func JoinPath(base string, elem ...string) (result string, err error) {
	url, err := Parse(base)
	if err != nil {
		return
	}
	result = url.JoinPath(elem...).String()
	return
}

var (
	_ encoding.TextMarshaler   = (*URI)(nil)
	_ encoding.TextUnmarshaler = (*URI)(nil)
)
