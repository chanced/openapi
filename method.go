package openapi

import "net/http"

type Method string

const (
	MethodGet     Method = http.MethodGet
	MethodPut     Method = http.MethodPut
	MethodPost    Method = http.MethodPost
	MethodDelete  Method = http.MethodDelete
	MethodOptions Method = http.MethodOptions
	MethodHead    Method = http.MethodHead
	MethodPatch   Method = http.MethodPatch
	MethodTrace   Method = http.MethodTrace
)
