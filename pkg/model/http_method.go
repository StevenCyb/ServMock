package model

import "strings"

// HTTPMethod represents the HTTP methods used in behaviors.
type HTTPMethod string

const (
	MethodGet     HTTPMethod = "GET"
	MethodPost    HTTPMethod = "POST"
	MethodPut     HTTPMethod = "PUT"
	MethodDelete  HTTPMethod = "DELETE"
	MethodPatch   HTTPMethod = "PATCH"
	MethodHead    HTTPMethod = "HEAD"
	MethodOptions HTTPMethod = "OPTIONS"
)

// HTTPMethodFromString converts a string to an HttpMethod.
func HTTPMethodFromString(method string) (HTTPMethod, bool) {
	method = strings.ToUpper(method)
	switch method {
	case "GET":
		return MethodGet, true
	case "POST":
		return MethodPost, true
	case "PUT":
		return MethodPut, true
	case "DELETE":
		return MethodDelete, true
	case "PATCH":
		return MethodPatch, true
	case "HEAD":
		return MethodHead, true
	case "OPTIONS":
		return MethodOptions, true
	}
	return MethodGet, false
}
