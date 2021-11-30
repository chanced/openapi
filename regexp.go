package openapi

import (
	"encoding/json"
	"regexp"
)

// Regexp is a wrapper around *regexp.Regexp to allow for marshinaling/unmarshaling
type Regexp struct {
	*regexp.Regexp
}

// MarshalJSON unmarshals data into sr
func (sr *Regexp) MarshalJSON() ([]byte, error) {
	if sr.IsNil() {
		return []byte{}, nil
	}
	return json.Marshal(sr.Regexp.String())
}

// UnmarshalJSON unmarshals data into sr
func (sr *Regexp) UnmarshalJSON(data []byte) error {
	var expr string
	var err error
	if err = json.Unmarshal(data, &expr); err != nil {
		return err
	}
	sr.Regexp, err = regexp.Compile(expr)
	return err
}

// IsNil returns true if either sr or sr.Regexp is nil
func (sr *Regexp) IsNil() bool {
	return sr == nil || sr.Regexp == nil
}
