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
	StatusCode *uint16
	Body       *string
	Delay      *time.Duration
	Headers    map[string]string
	Cookies    []*http.Cookie
	Redirect   *string
	Stream     bool
}

// Behavior defines the structure of a behavior with an associated HTTP method and URL.
type Behavior struct {
	*ResponseBehavior
	Method HttpMethod
	URL    string
	Repeat *uint
}
