package bufwr

import (
	"bytes"
	"io"
	"sync"
	"time"
)

type Writer struct {
	w    io.Writer
	size int

	m  sync.Mutex
	b  bytes.Buffer
	ts time.Time

	onError func(error)
}

// New creates a new buffered writer that implements io.Write.
// The writer accumulates input data in a buffer with size of size,
// until it is filled up or delayFlush time is pass, and then writes to
// underlying writer w.
func New(size int, delayFlush time.Duration, w io.Writer) *Writer {
	wr := &Writer{w: w, size: size}
	go wr.flusher(delayFlush)
	return wr
}

// OnError sets the handler function that will be called if a error
// occurs while write to the underlying writer.
func (w *Writer) OnError(handler func(error)) {
	w.onError = handler
}

// Write writes data to the internal buffer.
// It actually does not return any error because writing to the underlying
// writer most of the time happens asynchronously; see OnError method.
func (w *Writer) Write(data []byte) (int, error) {
	if len(data) == 0 {
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

func (w *Writer) flush() {
	// w.m must be locked here
	if w.b.Len() > 0 {
		_, err := w.w.Write(w.b.Bytes())
		if err != nil && w.onError != nil {
			w.onError(err)
		}
		w.b.Reset()
	}
}

func (w *Writer) flusher(delay time.Duration) {
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
