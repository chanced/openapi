package openapi

import "net/http"

const (
	MethodGet     = Text(http.MethodGet)
	MethodPut     = Text(http.MethodPut)
	MethodPost    = Text(http.MethodPost)
	MethodDelete  = Text(http.MethodDelete)
	MethodOptions = Text(http.MethodOptions)
	MethodHead    = Text(http.MethodHead)
	MethodPatch   = Text(http.MethodPatch)
	MethodTrace   = Text(http.MethodTrace)
)
