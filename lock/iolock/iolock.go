// Package iolock is a general lock wrapper for io.Reader, io.Writer and io.ReadWriter
package iolock

import (
	"io"
	"sync"
)

// Reader wrapper with mutex
type Reader struct {
	sync.Mutex
	reader io.Reader
}

// NewReader creates a new io.Reader with mutex
func NewReader(r io.Reader) *Reader {
	return &Reader{reader: r}
}

// Read to meet io.Reader interface
func (r *Reader) Read(p []byte) (n int, err error) {
	r.Lock()
	defer r.Unlock()
	return r.reader.Read(p)
}

// Writer wrapper with mutex
type Writer struct {
	sync.Mutex
	writer io.Writer
}

// NewWriter creates a new io.Writer with mutex
func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

// Write to meet io.Writer interface
func (w *Writer) Write(p []byte) (n int, err error) {
	w.Lock()
	defer w.Unlock()
	return w.writer.Write(p)
}

// ReadWriter wrapper with mutex
type ReadWriter struct {
	sync.Mutex
	rw io.ReadWriter
}

// NewReadWriter creates a new io.ReadWriter with mutex
func NewReadWriter(rw io.ReadWriter) *ReadWriter {
	return &ReadWriter{rw: rw}
}

// Read to meet io.Reader interface
func (rw *ReadWriter) Read(p []byte) (n int, err error) {
	rw.Lock()
	defer rw.Unlock()
	return rw.rw.Read(p)
}

// Write to meet io.Writer interface
func (rw *ReadWriter) Write(p []byte) (n int, err error) {
	rw.Lock()
	defer rw.Unlock()
	return rw.rw.Write(p)
}

// ReadWriterRW wrapper with rw_mutex, Read has rlock, Write has mutex
type ReadWriterRW struct {
	sync.RWMutex
	rw io.ReadWriter
}

// NewReadWriterRW creates a new io.ReadWriter with rw_mutex
func NewReadWriterRW(rw io.ReadWriter) *ReadWriterRW {
	return &ReadWriterRW{rw: rw}
}

// Read to meet io.Reader interface
func (rw *ReadWriterRW) Read(p []byte) (n int, err error) {
	rw.RLock()
	defer rw.RUnlock()
	return rw.rw.Read(p)
}

// Write to meet io.Writer interface
func (rw *ReadWriterRW) Write(p []byte) (n int, err error) {
	rw.Lock()
	defer rw.Unlock()
	return rw.rw.Write(p)
}
