package mmap

import (
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	data := []byte("abcdefg123456")
	dir, err := ioutil.TempDir(".", "tmp")
	assert.NoError(t, err)
	f, err := ioutil.TempFile(dir, "test")
	assert.NoError(t, err)

	fWithData, err := ioutil.TempFile(dir, "test_with_data")
	assert.NoError(t, err)
	n, err := fWithData.Write(data)
	assert.NoError(t, err)
	assert.Equal(t, len(data), n)
	defer func() {
		_ = os.RemoveAll(dir)
	}()
	err = f.Close()
	assert.NoError(t, err)
	err = fWithData.Close()
	assert.NoError(t, err)

	// not exist filename -> create file
	m, err := New(filepath.Join(dir, "/nope"), 10, RDWR, 0, 0)
	assert.NoError(t, err)
	assert.NotNil(t, m)
	err = m.Close()
	assert.NoError(t, err)
	// not exist dir -> fail with error
	m, err = New(filepath.Join(dir, "/nope/nope"), 10, RDWR, 0, 0)
	assert.Error(t, err)
	assert.Nil(t, m)
	// size == 0
	m, err = New(filepath.Join(dir, "/nope"), 0, RDWR, 0, 0)
	assert.Errorf(t, err, "mmap: can NOT map region size == 0")
	assert.Nil(t, m)

	// without data
	m, err = New(f.Name(), n, RDWR, 0, 0)
	assert.NoError(t, err)
	// ReadAt
	rn, err := m.ReadAt(nil, 0)
	assert.Errorf(t, err, "mmap: closed")
	assert.Equal(t, 0, rn)
	// WriteAt
	wn, err := m.WriteAt(nil, 0)
	assert.Errorf(t, err, "mmap: closed")
	assert.Equal(t, 0, wn)
	err = m.Close()
	assert.NoError(t, err)

	// with data
	m, err = New(fWithData.Name(), n, RDWR, ANON, 0)
	assert.NoError(t, err)
	err = m.Close()
	assert.NoError(t, err)
	m, err = New(fWithData.Name(), n, COPY, 0, 0)
	assert.NoError(t, err)
	err = m.Close()
	assert.NoError(t, err)
	m, err = New(fWithData.Name(), n, EXEC, 0, 0)
	assert.NoError(t, err)
	err = m.Close()
	assert.NoError(t, err)

	// RDWR for both ReadAt and WriteAt
	m, err = New(fWithData.Name(), n, RDWR, 0, 0)
	assert.NoError(t, err)
	assert.Equal(t, n, m.Len())
	for i := range data {
		assert.Equal(t, data[i], m.At(i))
	}
	// ReadAt
	buf := make([]byte, n)
	rn, err = m.ReadAt(buf, -1)
	assert.Errorf(t, err, "mmap: invalid ReadAt offset: -1")
	assert.Equal(t, 0, rn)
	rn, err = m.ReadAt(buf, 0)
	assert.NoError(t, err)
	assert.Equal(t, n, rn)
	rn, err = m.ReadAt(buf, 1)
	assert.Errorf(t, err, io.EOF.Error())
	assert.Equal(t, n-1, rn)
	// WriteAt
	// change data
	data = []byte("123456abcdefg")
	wn, err = m.WriteAt(nil, -1)
	assert.Errorf(t, err, "mmap: invalid WriteAt offset: -1")
	assert.Equal(t, 0, wn)
	wn, err = m.WriteAt(data, 0)
	assert.NoError(t, err)
	assert.Equal(t, n, wn)
	// change on file
	fData, err := ioutil.ReadFile(fWithData.Name())
	assert.NoError(t, err)
	assert.Equal(t, data, fData)
	wn, err = m.WriteAt(data, 1)
	assert.Errorf(t, err, io.ErrShortWrite.Error())
	assert.Equal(t, n-1, wn)
	err = m.Close()
	assert.NoError(t, err)

	// non-zero offset
	ps := os.Getpagesize()
	data = make([]byte, 2*ps, 2*ps)
	n = len(data)
	// pagesize error
	m, err = New(f.Name(), n, RDONLY, 0, 5)
	assert.Errorf(t, err, "mmap: offset must be a multiple of os pagesize")
	assert.Nil(t, m)
	// offset exceeds file size, m.data is nil
	m, err = New(f.Name(), n, RDONLY, 0, int64(3*ps))
	assert.NoError(t, err)
	assert.Nil(t, m.data)
	// normal
	m, err = New(f.Name(), n, RDONLY, 0, int64(ps))
	assert.NoError(t, err)
	assert.Equal(t, n, m.Len())
	err = m.Close()
	assert.NoError(t, err)
}
