package model

import (
	"net/http"
	"time"
)

// BehaviorSet represents a collection of behaviors with a default response behavior.
type BehaviorSet struct {
	DefaultBehavior *ResponseBehavior
	Behaviors       []*Behavior
}

// ResponseBehavior defines the structure of a response behavior.
type ResponseBehavior struct {
	Delay      *time.Duration
	StatusCode *uint16
	Body       *string
	Headers    map[string]string
	Cookies    []*http.Cookie
	Redirect   *string
	SSE        bool
}

// Behavior defines the structure of a behavior with an associated HTTP method and URL.
type Behavior struct {
	*ResponseBehavior
	Method HTTPMethod
	URL    string
	Repeat *uint
}
