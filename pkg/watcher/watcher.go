package watcher

import (
	"os"
	"time"
)

// Watcher monitors a file for changes at a specified interval.
type Watcher struct {
	path        string
	interval    int
	listener    func(string)
	running     bool
	stopChan    chan struct{}
	lastModTime int64
}

// NewWatcher creates a new Watcher for the given path and polling interval (ms).
func NewWatcher(path string, interval int) *Watcher {
	return &Watcher{
		path:     path,
		interval: interval,
		stopChan: make(chan struct{}),
	}
}

// RegisterListener sets the callback to be called when the file changes.
func (w *Watcher) RegisterListener(listener func(string)) {
	w.listener = listener
}

// Start begins polling for file changes.
func (w *Watcher) Start() {
	if w.running {
		return
	}
	w.running = true
	if w.listener != nil {
		w.listener(w.path)
	}
	go w.poll()
}

// Stop halts the watcher.
func (w *Watcher) Stop() {
	if !w.running {
		return
	}
	w.running = false
	close(w.stopChan)
}

func (w *Watcher) poll() {
	ticker := time.NewTicker(time.Duration(w.interval) * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			w.checkFile()
		case <-w.stopChan:
			return
		}
	}
}

func (w *Watcher) checkFile() {
	fi, err := os.Stat(w.path)
	if err != nil {
		return
	}
	modTime := fi.ModTime().UnixNano()
	if w.lastModTime == 0 {
		w.lastModTime = modTime
		return
	}
	if modTime != w.lastModTime {
		if w.listener != nil {
			w.listener(w.path)
		}
		w.lastModTime = modTime
	}
}
