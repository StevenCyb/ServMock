package model

import "strings"

// HttpMethod represents the HTTP methods used in behaviors.
type HttpMethod string

const (
	MethodGet     HttpMethod = "GET"
	MethodPost    HttpMethod = "POST"
	MethodPut     HttpMethod = "PUT"
	MethodDelete  HttpMethod = "DELETE"
	MethodPatch   HttpMethod = "PATCH"
	MethodHead    HttpMethod = "HEAD"
	MethodOptions HttpMethod = "OPTIONS"
)

// HttpMethodFromString converts a string to an HttpMethod.
func HttpMethodFromString(method string) (HttpMethod, bool) {
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
