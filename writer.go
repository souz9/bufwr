package bufwr

import (
	"bytes"
	"io"
	"sync"
	"time"
)

type writer struct {
	w    io.Writer
	size int

	m  sync.Mutex
	b  bytes.Buffer
	ts time.Time
}

// New creates a new buffered writer that implements io.Write.
// The writer accumulates input data in a buffer with size of size,
// until it is filled up or delayFlush time is pass, and then writes to
// underlying writer w.
func New(size int, delayFlush time.Duration, w io.Writer) *writer {
	wr := &writer{w: w, size: size}
	go wr.flusher(delayFlush)
	return wr
}

// Write writes data to the internal buffer.
// It actually does not return any error because writing to the underlying
// writer commonly happens asynchronously.
func (w *writer) Write(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, nil
	}

	w.m.Lock()
	if w.b.Len()+len(data) > w.size {
		w.flush()
	}
	if w.b.Len() == 0 {
		w.ts = time.Now() // Reset flush timer
	}
	w.b.Write(data)
	w.m.Unlock()

	return len(data), nil
}

func (w *writer) flush() {
	// w.m must be locked here
	if w.b.Len() > 0 {
		w.w.Write(w.b.Bytes())
		w.b.Reset()
	}
}

func (w *writer) flusher(delay time.Duration) {
	for {
		w.m.Lock()

		sleep := delay
		if w.b.Len() > 0 {
			elapsed := time.Since(w.ts)
			if elapsed >= delay {
				w.flush()
			} else {
				sleep = delay - elapsed
			}
		}

		w.m.Unlock()
		time.Sleep(sleep)
	}
}
