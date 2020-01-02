package lock

import (
	"io"
	"syscall"
)

// FdLock struct
type FdLock struct {
	fd int
}

// NewFdLock creates a new FdLock of input fd
func NewFdLock(fd int) (*FdLock, error) {
	// check validation by trying get file lock
	if err := syscall.FcntlFlock(uintptr(fd), syscall.F_GETFL, &syscall.Flock_t{
		Type:      syscall.F_WRLCK,
		Whence:    io.SeekStart,
		Pad_cgo_0: [4]byte{},
		Start:     0,
		Len:       0,
		Pid:       0,
		Pad_cgo_1: [4]byte{},
	}); err != nil {
		return nil, err
	}
	return &FdLock{fd: fd}, nil
}

// Lock adds write lock
func (fl *FdLock) Lock() error {
	if err := syscall.Flock(fl.fd, syscall.LOCK_EX); err != nil {
		return err
	}
	return nil
}

// Unlock removes write lock
func (fl *FdLock) Unlock() error {
	if err := syscall.Flock(fl.fd, syscall.LOCK_UN); err != nil {
		return err
	}
	return nil
}

// RLock adds read lock
func (fl *FdLock) RLock() error {
	if err := syscall.Flock(fl.fd, syscall.LOCK_SH); err != nil {
		return err
	}
	return nil
}

// RUnlock removes a read lock
func (fl *FdLock) RUnlock() error {
	if err := syscall.Flock(fl.fd, syscall.LOCK_UN); err != nil {
		return err
	}
	return nil
}

// Close unlocks fd and closes it
func (fl *FdLock) Close() error {
	if err := syscall.Flock(fl.fd, syscall.LOCK_UN); err != nil {
		return err
	}
	return syscall.Close(fl.fd)
}
