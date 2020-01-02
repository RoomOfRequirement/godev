package flow

import (
	"errors"
	"syscall"
	"testing"
	"time"
)

func TestWaitSig_WaitFor(t *testing.T) {
	t.Parallel()

	f := func() error {
		return errors.New("test")
	}

	w := NewWait(syscall.SIGINT)
	e := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	if e != nil {
		t.Fatal(e)
	}
	err := w.WaitFor(f)
	if err == nil || err.Error() != "test" {
		t.Fatal("wrong", err)
	}
}

func TestWaitSig_WaitForOrAfter(t *testing.T) {
	t.Parallel()

	t.Run("signal not exceed timeout", func(t *testing.T) {
		t.Parallel()
		f := func() error {
			return errors.New("test")
		}
		start := time.Now().Unix()
		w := NewWait(syscall.SIGINT)
		e := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		if e != nil {
			t.Fatal(e)
		}
		err := w.WaitForOrAfter(f, time.Second)
		end := time.Now().Unix()
		if diff := end - start; diff >= 1 {
			t.Fatal(diff)
		}
		if err == nil || err.Error() != "test" {
			t.Fatal("wrong", err)
		}
	})

	t.Run("signal exceed timeout", func(t *testing.T) {
		t.Parallel()
		f := func() error {
			return errors.New("test")
		}
		start := time.Now().Unix()
		time.AfterFunc(2*time.Second, func() {
			e := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
			if e != nil {
				t.Fatal(e)
			}
		})
		w := NewWait(syscall.SIGINT)
		err := w.WaitForOrAfter(f, time.Second)
		end := time.Now().Unix()
		if diff := end - start; diff < 1 || diff >= 2 {
			t.Fatal(diff)
		}
		if err == nil || err.Error() != "test" {
			t.Fatal("wrong", err)
		}
	})
}
