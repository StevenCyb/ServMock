package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/StevenCyb/ServMock/pkg/model"
	"github.com/stretchr/testify/assert"
)

type mockServer struct {
	Server
}

func newTestServer(behaviors []*model.Behavior, defaultBehavior *model.ResponseBehavior) *mockServer {
	return &mockServer{
		Server{
			behaviorSet: &model.BehaviorSet{
				Behaviors:       behaviors,
				DefaultBehavior: defaultBehavior,
			},
		},
	}
}

func TestHandleRequest_Body(t *testing.T) {
	body := "hello world"
	beh := &model.Behavior{
		Method:           http.MethodGet,
		URL:              "/body",
		ResponseBehavior: &model.ResponseBehavior{Body: &body},
	}
	ts := newTestServer([]*model.Behavior{beh}, nil)
	r := httptest.NewRequest(http.MethodGet, "/body", nil)
	w := httptest.NewRecorder()
	ts.handleRequest(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, body, w.Body.String())
}

func TestHandleRequest_Headers(t *testing.T) {
	headers := map[string]string{"X-Test": "val"}
	beh := &model.Behavior{
		Method:           http.MethodGet,
		URL:              "/headers",
		ResponseBehavior: &model.ResponseBehavior{Headers: headers},
	}
	ts := newTestServer([]*model.Behavior{beh}, nil)
	r := httptest.NewRequest(http.MethodGet, "/headers", nil)
	w := httptest.NewRecorder()
	ts.handleRequest(w, r)
	assert.Equal(t, "val", w.Header().Get("X-Test"))
}

func TestHandleRequest_Cookies(t *testing.T) {
	cookie := &http.Cookie{Name: "foo", Value: "bar"}
	beh := &model.Behavior{
		Method:           http.MethodGet,
		URL:              "/cookie",
		ResponseBehavior: &model.ResponseBehavior{Cookies: []*http.Cookie{cookie}},
	}
	ts := newTestServer([]*model.Behavior{beh}, nil)
	r := httptest.NewRequest(http.MethodGet, "/cookie", nil)
	w := httptest.NewRecorder()
	ts.handleRequest(w, r)
	assert.Contains(t, w.Header().Get("Set-Cookie"), "foo=bar")
}

func TestHandleRequest_Redirect(t *testing.T) {
	redirect := "/new"
	beh := &model.Behavior{
		Method:           http.MethodGet,
		URL:              "/redirect",
		ResponseBehavior: &model.ResponseBehavior{Redirect: &redirect},
	}
	ts := newTestServer([]*model.Behavior{beh}, nil)
	r := httptest.NewRequest(http.MethodGet, "/redirect", nil)
	w := httptest.NewRecorder()
	ts.handleRequest(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/new", w.Header().Get("Location"))
}

func TestHandleRequest_SSE(t *testing.T) {
	body := "chunk1\nchunk2"
	beh := &model.Behavior{
		Method:           http.MethodGet,
		URL:              "/sse",
		ResponseBehavior: &model.ResponseBehavior{SSE: true, Body: &body},
	}
	ts := newTestServer([]*model.Behavior{beh}, nil)
	r := httptest.NewRequest(http.MethodGet, "/sse", nil)
	w := httptest.NewRecorder()
	ts.handleRequest(w, r)
	assert.Equal(t, "text/event-stream", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Body.String(), "data: chunk1")
	assert.Contains(t, w.Body.String(), "data: chunk2")
}

func TestHandleRequest_Delay(t *testing.T) {
	d := 10 * time.Millisecond
	beh := &model.Behavior{
		Method:           http.MethodGet,
		URL:              "/delay",
		ResponseBehavior: &model.ResponseBehavior{Delay: &d, Body: nil},
	}
	ts := newTestServer([]*model.Behavior{beh}, nil)
	r := httptest.NewRequest(http.MethodGet, "/delay", nil)
	w := httptest.NewRecorder()
	start := time.Now()
	ts.handleRequest(w, r)
	elapsed := time.Since(start)
	assert.GreaterOrEqual(t, elapsed, d)
}

func TestHandleRequest_StatusCodeOverride(t *testing.T) {
	status := uint16(201)
	def := &model.ResponseBehavior{StatusCode: &status}
	beh := &model.Behavior{
		Method:           http.MethodGet,
		URL:              "/status",
		ResponseBehavior: &model.ResponseBehavior{StatusCode: &status},
	}
	ts := newTestServer([]*model.Behavior{beh}, def)
	r := httptest.NewRequest(http.MethodGet, "/status", nil)
	w := httptest.NewRecorder()
	ts.handleRequest(w, r)
	assert.Equal(t, 201, w.Code)
}

func TestHandleRequest_NotFound_DefaultBehavior(t *testing.T) {
	body := "not found"
	def := &model.ResponseBehavior{Body: &body}
	ts := newTestServer([]*model.Behavior{}, def)
	r := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	w := httptest.NewRecorder()
	ts.handleRequest(w, r)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, body, w.Body.String())
}

func TestHandleRequest_NotFound_NoDefault(t *testing.T) {
	ts := newTestServer([]*model.Behavior{}, nil)
	r := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	w := httptest.NewRecorder()
	ts.handleRequest(w, r)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestHandleRequest_RepeatBehavior_Removal(t *testing.T) {
	body := "repeat once"
	repeat := uint(1)
	beh := &model.Behavior{
		Method:           http.MethodGet,
		URL:              "/repeat",
		Repeat:           &repeat,
		ResponseBehavior: &model.ResponseBehavior{Body: &body},
	}
	ts := newTestServer([]*model.Behavior{beh}, nil)

	// First call: should match and decrement repeat to 0, then remove behavior
	r1 := httptest.NewRequest(http.MethodGet, "/repeat", nil)
	w1 := httptest.NewRecorder()
	ts.handleRequest(w1, r1)
	assert.Equal(t, http.StatusOK, w1.Code)
	assert.Equal(t, body, w1.Body.String())
	assert.Empty(t, ts.behaviorSet.Behaviors, "Behavior should be removed after repeat reaches 0")

	// Second call: should not match, returns NotFound
	r2 := httptest.NewRequest(http.MethodGet, "/repeat", nil)
	w2 := httptest.NewRecorder()
	ts.handleRequest(w2, r2)
	assert.Equal(t, http.StatusNotFound, w2.Code)
}

func TestHandleRequest_RepeatBehavior_Decrement(t *testing.T) {
	body := "repeat twice"
	repeat := uint(2)
	beh := &model.Behavior{
		Method:           http.MethodGet,
		URL:              "/repeat2",
		Repeat:           &repeat,
		ResponseBehavior: &model.ResponseBehavior{Body: &body},
	}
	ts := newTestServer([]*model.Behavior{beh}, nil)

	// First call: should match and decrement repeat to 1
	r1 := httptest.NewRequest(http.MethodGet, "/repeat2", nil)
	w1 := httptest.NewRecorder()
	ts.handleRequest(w1, r1)
	assert.Equal(t, http.StatusOK, w1.Code)
	assert.Equal(t, body, w1.Body.String())
	assert.Len(t, ts.behaviorSet.Behaviors, 1, "Behavior should still exist after first call")

	// Second call: should match and decrement repeat to 0, then remove behavior
	r2 := httptest.NewRequest(http.MethodGet, "/repeat2", nil)
	w2 := httptest.NewRecorder()
	ts.handleRequest(w2, r2)
	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Equal(t, body, w2.Body.String())
	assert.Empty(t, ts.behaviorSet.Behaviors, "Behavior should be removed after second call")
}
