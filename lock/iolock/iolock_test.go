package iolock

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	// Reader
	b := make([]byte, 5)
	reader := bytes.NewReader([]byte("test"))
	r := NewReader(reader)
	n, err := r.Read(b)
	assert.Equal(t, 4, n)
	assert.NoError(t, err)

	// Writer
	b = make([]byte, 5)
	writer := bytes.NewBuffer(b)
	w := NewWriter(writer)
	n, err = w.Write([]byte("test"))
	assert.Equal(t, 4, n)
	assert.NoError(t, err)

	// ReadWriter
	var readWriter bytes.Buffer
	rw := NewReadWriter(&readWriter)
	n, err = rw.Write([]byte("test"))
	assert.Equal(t, 4, n)
	assert.NoError(t, err)
	assert.Equal(t, "test", readWriter.String())
	b = make([]byte, 5)
	n, err = rw.Read(b)
	assert.Equal(t, 4, n)
	assert.NoError(t, err)

	// ReadWriterRW
	rwRW := NewReadWriterRW(&readWriter)
	n, err = rwRW.Write([]byte("test"))
	assert.Equal(t, 4, n)
	assert.NoError(t, err)
	assert.Equal(t, "test", readWriter.String())
	b = make([]byte, 5)
	n, err = rwRW.Read(b)
	assert.Equal(t, 4, n)
	assert.NoError(t, err)
}
