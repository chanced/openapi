package openapi

import "fmt"

type DuplicateAnchorError struct {
	A *Anchor
	B *Anchor
}

func (dae *DuplicateAnchorError) Error() string {
	return fmt.Sprintf("duplicate anchor: %s", dae.A.Name)
}

type AnchorType uint8

const (
	AnchorTypeUndefined AnchorType = iota
	AnchorTypeRegular              // $anchor
	AnchorTypeRecursive            // $recursiveAnchor
	AnchorTypeDynamic              // $dynamicAnchor
)

type Anchor struct {
	Location
	In   *Schema
	Name Text
	Type AnchorType
}

type Anchors struct {
	Standard  map[Text]Anchor // $anchor
	Recursive *Anchor         // $recursiveAnchor
	Dynamic   map[Text]Anchor // $dynamicAnchor
}

func (a *Anchors) merge(b *Anchors, err error) (*Anchors, error) {
	if err != nil {
		return nil, err
	}
	if b == nil {
		return a, nil
	}

	// we do not merge recursive anchors as they must be at the root of the
	// document. This method is only called when merging schemas from nested
	// components, so we can, and should, drop them from result if not coming
	// from a.

	if a == nil {
		return &Anchors{
			Standard: b.Standard,
			Dynamic:  b.Dynamic,
		}, nil
	}
	for k, bv := range b.Standard {
		if av, ok := a.Standard[k]; ok {
			return nil, &DuplicateAnchorError{&av, &bv}
		}
		a.Standard[k] = bv
	}

	for k, bv := range b.Dynamic {
		if av, ok := a.Dynamic[k]; ok {
			return nil, &DuplicateAnchorError{&av, &bv}
		}
		a.Dynamic[k] = bv
	}

	return a, nil
}
