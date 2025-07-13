package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/StevenCyb/ServMock/pkg/model"
)

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	matchingBehavior, statusCode := s.findMatchingBehavior(r)
	if matchingBehavior == nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	if matchingBehavior.Delay != nil {
		time.Sleep(*matchingBehavior.Delay)
	}

	if matchingBehavior.StatusCode != nil {
		statusCode = int(*s.behaviorSet.DefaultBehavior.StatusCode)
	}
	w.WriteHeader(statusCode)

	if matchingBehavior.Headers != nil {
		for key, value := range matchingBehavior.Headers {
			w.Header().Set(key, value)
		}
	}

	if matchingBehavior.Cookies != nil {
		for _, cookie := range matchingBehavior.Cookies {
			http.SetCookie(w, cookie)
		}
	}

	if matchingBehavior.Redirect != nil {
		http.Redirect(w, r, *matchingBehavior.Redirect, statusCode)
		return
	}

	if matchingBehavior.SSE {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
			return
		}

		if matchingBehavior.Body != nil {
			for _, chunk := range strings.Split(*matchingBehavior.Body, "\n") {
				fmt.Fprintf(w, "data: %s\n\n", chunk)
				flusher.Flush()
			}
		}

		flusher.Flush()
	} else if matchingBehavior.Body != nil {
		w.Write([]byte(*matchingBehavior.Body))
	}
}

func (s *Server) findMatchingBehavior(r *http.Request) (*model.ResponseBehavior, int) {
	var matchingBehavior *model.ResponseBehavior
	var statusCode = http.StatusOK

	for i, behavior := range s.behaviorSet.Behaviors {
		if behavior.Method == model.HttpMethod(r.Method) && behavior.URL == r.URL.Path {
			if behavior.Repeat != nil {
				*behavior.Repeat--
				if *behavior.Repeat <= 0 {
					s.behaviorSet.Behaviors = append(s.behaviorSet.Behaviors[:i], s.behaviorSet.Behaviors[i+1:]...)
				}
			}

			matchingBehavior = behavior.ResponseBehavior
			break
		}
	}

	if matchingBehavior == nil {
		statusCode = http.StatusNotFound
		if s.behaviorSet.DefaultBehavior != nil {
			matchingBehavior = s.behaviorSet.DefaultBehavior
		}
	}
	return matchingBehavior, statusCode
}
