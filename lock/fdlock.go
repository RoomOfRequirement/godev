package lock

import (
	"io"
	"syscall"
	"time"
)

// FdLock struct
//	https://linux.die.net/man/2/flock
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
	return syscall.Flock(fl.fd, syscall.LOCK_EX)
}

// TryLock used for some non-blocking case
func (fl *FdLock) TryLock() (bool, error) {
	return fl.try(syscall.LOCK_EX)
}

// TryLockWithTimeout keeps trying to get Lock until timeout
func (fl *FdLock) TryLockWithTimeout(timeout time.Duration) (bool, error) {
	return fl.tryWithTimeout(syscall.LOCK_EX, timeout)
}

// RetryLock keeps retrying to get Lock with delay until cancel
func (fl *FdLock) RetryLock(cancel chan struct{}, delay time.Duration) (bool, error) {
	return fl.retry(syscall.LOCK_EX, cancel, delay)
}

// RetryLockWithTimeout keeps retrying to get Lock with delay until cancel or timeout
func (fl *FdLock) RetryLockWithTimeout(cancel chan struct{}, delay, timeout time.Duration) (bool, error) {
	return fl.retryWithTimeout(syscall.LOCK_EX, cancel, delay, timeout)
}

// Unlock removes write lock
func (fl *FdLock) Unlock() error {
	return syscall.Flock(fl.fd, syscall.LOCK_UN)
}

// RLock adds read lock
func (fl *FdLock) RLock() error {
	return syscall.Flock(fl.fd, syscall.LOCK_SH)
}

// TryRLock used for some non-blocking case
func (fl *FdLock) TryRLock() (bool, error) {
	return fl.try(syscall.LOCK_SH)
}

// TryRLockWithTimeout keeps trying to get RLock until timeout
func (fl *FdLock) TryRLockWithTimeout(timeout time.Duration) (bool, error) {
	return fl.tryWithTimeout(syscall.LOCK_SH, timeout)
}

// RetryRLock keeps retrying to get RLock with delay until cancel
func (fl *FdLock) RetryRLock(cancel chan struct{}, delay time.Duration) (bool, error) {
	return fl.retry(syscall.LOCK_SH, cancel, delay)
}

// RetryRLockWithTimeout keeps retrying to get Lock with delay until cancel or timeout
func (fl *FdLock) RetryRLockWithTimeout(cancel chan struct{}, delay, timeout time.Duration) (bool, error) {
	return fl.retryWithTimeout(syscall.LOCK_SH, cancel, delay, timeout)
}

// RUnlock removes a read lock
func (fl *FdLock) RUnlock() error {
	return syscall.Flock(fl.fd, syscall.LOCK_UN)
}

// Close unlocks fd and closes it
func (fl *FdLock) Close() error {
	if err := syscall.Flock(fl.fd, syscall.LOCK_UN); err != nil {
		return err
	}
	return syscall.Close(fl.fd)
}

func (fl *FdLock) try(flag int) (bool, error) {
	// LOCK_NB -> non-blocking flag
	err := syscall.Flock(fl.fd, flag|syscall.LOCK_NB)

	switch err {
	case syscall.EWOULDBLOCK: // already locked
		return false, nil
	case nil:
		return true, nil
	}
	return false, err
}

func (fl *FdLock) tryWithTimeout(flag int, timeout time.Duration) (bool, error) {
	timer := time.NewTimer(timeout)
	for {
		select {
		case <-timer.C:
			timer.Stop()
			return false, nil
		default:
			if ok, err := fl.try(flag); ok || err != nil {
				return ok, err
			}
		}
	}
}

func (fl *FdLock) retry(flag int, cancel chan struct{}, delay time.Duration) (bool, error) {
	for {
		select {
		case <-cancel:
			return false, nil
		case <-time.After(delay):
			if ok, err := fl.try(flag); ok || err != nil {
				return ok, err
			}
		}
	}
}

func (fl *FdLock) retryWithTimeout(flag int, cancel chan struct{}, delay, timeout time.Duration) (bool, error) {
	timer := time.NewTimer(timeout)
	for {
		select {
		case <-cancel:
			return false, nil
		case <-time.After(delay):
			if ok, err := fl.try(flag); ok || err != nil {
				return ok, err
			}
		case <-timer.C:
			timer.Stop()
			return false, nil
		}
	}
}
