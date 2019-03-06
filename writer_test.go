package bufwr

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestWriter(t *testing.T) {
	b := &bytes.Buffer{}
	w := New(20, 100*time.Millisecond, b)

	w.Write([]byte("MessageA"))
	assert.Empty(t, b.Bytes())

	w.Write([]byte("MessageB"))
	assert.Empty(t, b.Bytes())

	w.Write([]byte("MessageC"))
	// There is not enough free space in the internal buffer, should flush here
	assert.Equal(t, []byte("MessageAMessageB"), b.Bytes())

	// Should flust the rest part after the specified delay
	b.Reset()
	time.Sleep(200 * time.Millisecond)
	assert.Equal(t, []byte("MessageC"), b.Bytes())
}
