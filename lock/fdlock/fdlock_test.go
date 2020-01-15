package fdlock

import (
	"github.com/stretchr/testify/assert"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestFdLock(t *testing.T) {
	fd, err := syscall.Open("./tmp.txt", syscall.O_CREAT|syscall.O_RDONLY, 0750)
	assert.NoError(t, err)
	defer func() {
		_ = syscall.Close(fd)
		_ = os.Remove("./tmp.txt")
	}()

	// normal case
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

	// error case
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

func TestFdLock_Try(t *testing.T) {
	bfd, err := syscall.Open("./tmp.txt", syscall.O_CREAT|syscall.O_RDONLY, 0750)
	assert.NoError(t, err)
	fd, err := syscall.Open("./tmp.txt", syscall.O_CREAT|syscall.O_RDONLY, 0750)
	assert.NoError(t, err)
	fd1, err := syscall.Open("./tmp.txt", syscall.O_CREAT|syscall.O_RDONLY, 0750)
	assert.NoError(t, err)
	defer func() {
		_ = syscall.Close(bfd)
		_ = syscall.Close(fd)
		_ = syscall.Close(fd1)
		_ = os.Remove("./tmp.txt")
	}()

	// error case
	bfl, err := NewFdLock(bfd)
	assert.NoError(t, err)
	err = bfl.Close()

	assert.NoError(t, err)
	ok, err := bfl.TryLock()
	assert.False(t, ok)
	assert.Error(t, err, "bad file descriptor")

	fl, err := NewFdLock(fd)
	assert.NoError(t, err)
	fl1, err := NewFdLock(fd1)
	assert.NoError(t, err)

	// Lock Lock
	ok, err = fl.TryLock()
	assert.True(t, ok)
	assert.NoError(t, err)
	ok, err = fl1.TryLock()
	assert.False(t, ok)
	assert.NoError(t, err)
	err = fl.Unlock()
	assert.NoError(t, err)
	ok, err = fl1.TryLock()
	assert.True(t, ok)
	assert.NoError(t, err)
	err = fl1.Unlock()
	assert.NoError(t, err)

	// RLock Lock
	ok, err = fl.TryRLock()
	assert.True(t, ok)
	assert.NoError(t, err)
	ok, err = fl1.TryLock()
	assert.False(t, ok)
	assert.NoError(t, err)
	err = fl.RUnlock()
	assert.NoError(t, err)
	ok, err = fl1.TryLock()
	assert.True(t, ok)
	assert.NoError(t, err)
	err = fl1.Unlock()
	assert.NoError(t, err)

	// Lock RLock
	ok, err = fl.TryLock()
	assert.True(t, ok)
	assert.NoError(t, err)
	ok, err = fl1.TryRLock()
	assert.False(t, ok)
	assert.NoError(t, err)
	err = fl.Unlock()
	assert.NoError(t, err)
	ok, err = fl1.TryRLock()
	assert.True(t, ok)
	assert.NoError(t, err)
	err = fl1.RUnlock()
	assert.NoError(t, err)

	// RLock RLock
	ok, err = fl.TryRLock()
	assert.True(t, ok)
	assert.NoError(t, err)
	ok, err = fl1.TryRLock()
	assert.True(t, ok)
	assert.NoError(t, err)
	err = fl.RUnlock()
	assert.NoError(t, err)
	err = fl1.RUnlock()
	assert.NoError(t, err)
}

func TestFdLock_TryWithTimeout(t *testing.T) {
	fd, err := syscall.Open("./tmp.txt", syscall.O_CREAT|syscall.O_RDONLY, 0750)
	assert.NoError(t, err)
	fd1, err := syscall.Open("./tmp.txt", syscall.O_CREAT|syscall.O_RDONLY, 0750)
	assert.NoError(t, err)
	defer func() {
		_ = syscall.Close(fd)
		_ = syscall.Close(fd1)
		_ = os.Remove("./tmp.txt")
	}()

	fl, err := NewFdLock(fd)
	assert.NoError(t, err)
	fl1, err := NewFdLock(fd1)
	assert.NoError(t, err)

	// Lock
	err = fl.Lock()
	assert.NoError(t, err)
	ok, err := fl1.TryLockWithTimeout(1 * time.Millisecond)
	assert.False(t, ok)
	assert.NoError(t, nil)

	time.AfterFunc(1*time.Millisecond, func() {
		_ = fl.Unlock()
	})
	ok, err = fl1.TryLockWithTimeout(2 * time.Millisecond)
	assert.True(t, ok)
	assert.NoError(t, nil)
	err = fl1.Unlock()
	assert.NoError(t, err)

	// RLock
	err = fl.Lock()
	assert.NoError(t, err)
	ok, err = fl1.TryRLockWithTimeout(1 * time.Millisecond)
	assert.False(t, ok)
	assert.NoError(t, nil)

	time.AfterFunc(1*time.Millisecond, func() {
		_ = fl.Unlock()
	})
	ok, err = fl1.TryRLockWithTimeout(2 * time.Millisecond)
	assert.True(t, ok)
	assert.NoError(t, nil)
}

func TestFdLock_Retry(t *testing.T) {
	fd, err := syscall.Open("./tmp.txt", syscall.O_CREAT|syscall.O_RDONLY, 0750)
	assert.NoError(t, err)
	fd1, err := syscall.Open("./tmp.txt", syscall.O_CREAT|syscall.O_RDONLY, 0750)
	assert.NoError(t, err)
	defer func() {
		_ = syscall.Close(fd)
		_ = syscall.Close(fd1)
		_ = os.Remove("./tmp.txt")
	}()

	fl, err := NewFdLock(fd)
	assert.NoError(t, err)
	fl1, err := NewFdLock(fd1)
	assert.NoError(t, err)

	// Lock Lock
	err = fl.Lock()
	assert.NoError(t, err)
	cancel := make(chan struct{})
	time.AfterFunc(1*time.Millisecond, func() {
		close(cancel)
	})
	ok, err := fl1.RetryLock(cancel, 1*time.Microsecond)
	assert.False(t, ok)
	assert.NoError(t, nil)

	cancel = make(chan struct{})
	time.AfterFunc(1*time.Millisecond, func() {
		_ = fl.Unlock()
	})
	ok, err = fl1.RetryLock(cancel, 1*time.Microsecond)
	assert.True(t, ok)
	assert.NoError(t, nil)
	close(cancel)
	err = fl1.Unlock()
	assert.NoError(t, err)

	// Lock RLock
	err = fl.Lock()
	assert.NoError(t, err)
	cancel = make(chan struct{})
	time.AfterFunc(1*time.Millisecond, func() {
		close(cancel)
	})
	ok, err = fl1.RetryRLock(cancel, 1*time.Microsecond)
	assert.False(t, ok)
	assert.NoError(t, nil)

	cancel = make(chan struct{})
	time.AfterFunc(1*time.Millisecond, func() {
		_ = fl.Unlock()
	})
	ok, err = fl1.RetryRLock(cancel, 1*time.Microsecond)
	assert.True(t, ok)
	assert.NoError(t, nil)
	close(cancel)
	err = fl1.RUnlock()
	assert.NoError(t, err)

	// RLock RLock
	err = fl.RLock()
	assert.NoError(t, err)
	cancel = make(chan struct{})
	time.AfterFunc(1*time.Millisecond, func() {
		close(cancel)
	})
	ok, err = fl1.RetryRLock(cancel, 1*time.Microsecond)
	assert.True(t, ok)
	assert.NoError(t, nil)
	err = fl1.RUnlock()
	assert.NoError(t, err)
}

func TestFdLock_RetryWithTimeout(t *testing.T) {
	fd, err := syscall.Open("./tmp.txt", syscall.O_CREAT|syscall.O_RDONLY, 0750)
	assert.NoError(t, err)
	fd1, err := syscall.Open("./tmp.txt", syscall.O_CREAT|syscall.O_RDONLY, 0750)
	assert.NoError(t, err)
	defer func() {
		_ = syscall.Close(fd)
		_ = syscall.Close(fd1)
		_ = os.Remove("./tmp.txt")
	}()

	fl, err := NewFdLock(fd)
	assert.NoError(t, err)
	fl1, err := NewFdLock(fd1)
	assert.NoError(t, err)

	// Lock
	cancel := make(chan struct{})

	err = fl.Lock()
	assert.NoError(t, err)
	ok, err := fl1.RetryLockWithTimeout(cancel, 1*time.Microsecond, 1*time.Millisecond)
	assert.False(t, ok)
	assert.NoError(t, nil)
	err = fl.Unlock()
	assert.NoError(t, err)

	err = fl.Lock()
	assert.NoError(t, err)
	time.AfterFunc(1*time.Millisecond, func() {
		_ = fl.Unlock()
	})
	ok, err = fl1.RetryLockWithTimeout(cancel, 1*time.Microsecond, 2*time.Millisecond)
	assert.True(t, ok)
	assert.NoError(t, nil)
	err = fl1.Unlock()
	assert.NoError(t, err)

	err = fl.Lock()
	assert.NoError(t, err)
	time.AfterFunc(1*time.Millisecond, func() {
		close(cancel)
	})
	ok, err = fl1.RetryLockWithTimeout(cancel, 1*time.Microsecond, 2*time.Millisecond)
	assert.False(t, ok)
	assert.NoError(t, nil)
	err = fl1.Unlock()
	assert.NoError(t, err)

	// RLock
	cancel = make(chan struct{})

	err = fl.Lock()
	assert.NoError(t, err)
	ok, err = fl1.RetryRLockWithTimeout(cancel, 1*time.Microsecond, 1*time.Millisecond)
	assert.False(t, ok)
	assert.NoError(t, nil)
	err = fl.Unlock()
	assert.NoError(t, err)

	err = fl.Lock()
	assert.NoError(t, err)
	time.AfterFunc(1*time.Millisecond, func() {
		_ = fl.Unlock()
	})
	ok, err = fl1.RetryRLockWithTimeout(cancel, 1*time.Microsecond, 2*time.Millisecond)
	assert.True(t, ok)
	assert.NoError(t, nil)
	err = fl1.RUnlock()
	assert.NoError(t, err)

	err = fl.Lock()
	assert.NoError(t, err)
	time.AfterFunc(1*time.Millisecond, func() {
		close(cancel)
	})
	ok, err = fl1.RetryRLockWithTimeout(cancel, 1*time.Microsecond, 2*time.Millisecond)
	assert.False(t, ok)
	assert.NoError(t, nil)
}
