package watcher

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWatcher_FileChangeTriggersListener(t *testing.T) {
	f, err := os.CreateTemp("", "watcher_test_*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())

	triggered := false
	listener := func(path string) {
		triggered = true
	}

	w := NewWatcher(f.Name(), 250)
	w.RegisterListener(listener)
	w.Start()
	defer w.Stop()

	time.Sleep(1 * time.Second)

	_, err = f.WriteString("change")
	assert.NoError(t, err, "failed to write to file")
	f.Sync()

	assert.Eventually(t, func() bool { return triggered }, 500*time.Millisecond, 20*time.Millisecond, "Listener was not triggered on file change")
}

func TestWatcher_StartStopIdempotent(t *testing.T) {
	f, err := os.CreateTemp("", "watcher_test_*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())

	w := NewWatcher(f.Name(), 10)
	w.Start()
	w.Start()
	w.Stop()
	w.Stop()
}

func TestWatcher_NoListenerNoPanic(t *testing.T) {
	f, err := os.CreateTemp("", "watcher_test_*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())

	w := NewWatcher(f.Name(), 10)
	w.Start()
	defer w.Stop()

	if _, err := f.WriteString("change"); err != nil {
		t.Fatalf("failed to write to file: %v", err)
	}
	f.Sync()

	time.Sleep(50 * time.Millisecond)
}
