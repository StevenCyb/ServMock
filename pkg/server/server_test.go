package server

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/StevenCyb/ServMock/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestServerStartupShutdown(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		t.Skip("Skipping server startup/shutdown test in GitHub Actions environment")
	}

	t.Parallel()

	server := New(":8080", &model.BehaviorSet{})
	errorChan := server.Start()

	time.Sleep(500 * time.Millisecond)
	resp, err := http.Get("http://localhost:8080/")
	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
	if resp.Body != nil {
		resp.Body.Close()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	assert.NoError(t, server.Shutdown(ctx))

	select {
	case err := <-errorChan:
		if err != nil && err.Error() != "http: Server closed" {
			assert.NoError(t, err, "Server should start and shutdown without error")
		}
	case <-time.After(1 * time.Second):
		// No error received, assume normal shutdown
	}
}
