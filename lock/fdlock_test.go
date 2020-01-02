package lock

import (
	"github.com/stretchr/testify/assert"
	"os"
	"syscall"
	"testing"
)

func TestFdLock(t *testing.T) {
	fd, err := syscall.Open("./tmp.txt", syscall.O_CREAT|syscall.O_RDONLY, 0750)
	assert.NoError(t, err)
	defer func() {
		syscall.Close(fd)
		os.Remove("./tmp.txt")
	}()

	fl, err := NewFdLock(fd)
	assert.NoError(t, err)

	err = fl.Lock()
	assert.NoError(t, err)

	err = fl.Unlock()
	assert.NoError(t, err)

	err = fl.RLock()
	assert.NoError(t, err)

	err = fl.RUnlock()
	assert.NoError(t, err)

	err = fl.Close()
	assert.NoError(t, err)

	fl, err = NewFdLock(-1)
	assert.Nil(t, fl)
	assert.Error(t, err, "bad file descriptor")

	fd, err = syscall.Open("./tmp.txt", syscall.O_CREAT|syscall.O_RDONLY, 0750)
	assert.NoError(t, err)
	fl, err = NewFdLock(fd)
	assert.NoError(t, err)
	err = syscall.Close(fd)
	assert.NoError(t, err)

	err = fl.Lock()
	assert.Error(t, err, "bad file descriptor")

	err = fl.Unlock()
	assert.Error(t, err, "bad file descriptor")

	err = fl.RLock()
	assert.Error(t, err, "bad file descriptor")

	err = fl.RUnlock()
	assert.Error(t, err, "bad file descriptor")

	err = fl.Close()
	assert.Error(t, err, "bad file descriptor")
}
