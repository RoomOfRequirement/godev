package mmap

import (
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	data := []byte("abcdefg123456")

	f, err := ioutil.TempFile(".", "test")
	assert.NoError(t, err)

	fWithData, err := ioutil.TempFile(".", "test_with_data")
	assert.NoError(t, err)
	n, err := fWithData.Write(data)
	assert.NoError(t, err)
	assert.Equal(t, len(data), n)
	defer func() {
		_ = os.Remove(f.Name())
		_ = os.Remove(fWithData.Name())
	}()
	err = f.Close()
	assert.NoError(t, err)
	err = fWithData.Close()
	assert.NoError(t, err)

	// wrong filename
	m, err := New("nope", RDWR, 0, 0)
	assert.Error(t, err)
	assert.Nil(t, m)

	// without data
	m, err = New(f.Name(), RDWR, 0, 0)
	assert.NoError(t, err)
	rn, err := m.ReadAt(nil, 0)
	assert.Errorf(t, err, "mmap: closed")
	assert.Equal(t, 0, rn)
	err = m.Close()
	assert.NoError(t, err)

	// with data
	m, err = New(fWithData.Name(), RDWR, ANON, 0)
	assert.NoError(t, err)
	err = m.Close()
	assert.NoError(t, err)
	m, err = New(fWithData.Name(), COPY, 0, 0)
	assert.NoError(t, err)
	err = m.Close()
	assert.NoError(t, err)
	m, err = New(fWithData.Name(), EXEC, 0, 0)
	assert.NoError(t, err)
	err = m.Close()
	assert.NoError(t, err)

	m, err = New(fWithData.Name(), RDONLY, 0, 0)
	assert.NoError(t, err)
	assert.Equal(t, n, m.Len())
	for i := range data {
		assert.Equal(t, data[i], m.At(i))
	}
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
	err = m.Close()
	assert.NoError(t, err)

	// non-zero offset
	ps := os.Getpagesize()
	data = make([]byte, 2*ps, 2*ps)
	n = len(data)
	// pagesize error
	m, err = New(f.Name(), RDONLY, 0, 5)
	assert.Errorf(t, err, "mmap: offset must be a multiple of os pagesize")
	assert.Nil(t, m)
	// offset exceeds file size, m.data is nil
	m, err = New(f.Name(), RDONLY, 0, int64(3*ps))
	assert.NoError(t, err)
	assert.Nil(t, m.data)
	// normal
	m, err = New(f.Name(), RDONLY, 0, int64(ps))
	assert.NoError(t, err)
	assert.Equal(t, 0, m.Len())

	// lock/unlock
	err = m.Lock()
	assert.NoError(t, err)
	err = m.Unlock()
	assert.NoError(t, err)
	err = m.Close()
	assert.NoError(t, err)
}
