package eventloop

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestNewLoop(t *testing.T) {
	f := func(n int) {
		defer func() {
			if r := recover(); r != nil {
				if r.(string) != "invalid buffered events size" {
					t.Fatal(r)
				}
			}
		}()

		_ = NewLoop(n, "debug")
	}

	f(-1)
	f(10)
}

func TestLoop_Simple(t *testing.T) {
	l := NewLoop(6, "debug")

	num := []int{0, 1, 2, 3, 4, 5}
	for i := range num[:3] {
		name := strconv.Itoa(i)
		l.Push(&Event{
			Name: name,
			Data: i,
		})
		err := l.On(name, func(ctx context.Context, args ...interface{}) error {
			_, err := fmt.Println(name, args)
			return err
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	for i := range num[3:] {
		name := strconv.Itoa(i)
		l.Emit(name, i)
		err := l.On(name, func(ctx context.Context, args ...interface{}) error {
			_, err := fmt.Println(name, args)
			return err
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	l.Start()
	l.Start()
	time.Sleep(10 * time.Millisecond)
	defer l.Stop()
}

func TestLoop_Timeout(t *testing.T) {
	l := NewLoop(6, "debug")

	num := []int{0, 1, 2, 3, 4, 5}
	for i := range num[:3] {
		name := strconv.Itoa(i)
		l.Push(&Event{
			Name: name,
			Data: i,
		})
		err := l.OnWithTimeout(name, func(ctx context.Context, args ...interface{}) error {
			_, err := fmt.Println(name, args)
			return err
		}, 5*time.Millisecond)
		time.Sleep(5 * time.Millisecond)
		if err != nil {
			t.Fatal(err)
		}
	}
	for i := range num[3:] {
		name := strconv.Itoa(i)
		l.Emit(name, i)
		err := l.On(name, func(ctx context.Context, args ...interface{}) error {
			_, err := fmt.Println(name, args)
			return err
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	l.Start()

	time.Sleep(10 * time.Millisecond)
	defer l.Stop()
}

func TestLoop_Error(t *testing.T) {
	l := NewLoop(6, "debug")

	num := []int{0, 1, 2, 3, 4, 5}
	for i := range num[:3] {
		name := strconv.Itoa(i)
		l.Push(&Event{
			Name: name,
			Data: i,
		})
		err := l.OnWithTimeout(name, func(ctx context.Context, args ...interface{}) error {
			return errors.New("test")
		}, 5*time.Millisecond)
		time.Sleep(5 * time.Millisecond)
		if err != nil {
			t.Fatal(err)
		}
	}
	for i := range num[3:] {
		name := strconv.Itoa(i)
		l.Emit(name, i)
		err := l.On(name, func(ctx context.Context, args ...interface{}) error {
			return errors.New("test")
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	l.Start()

	time.Sleep(10 * time.Millisecond)
	defer l.Stop()
}

func TestLoop_AddEventListenerError(t *testing.T) {
	l := NewLoop(6, "debug")

	num := []int{0, 1, 2, 3, 4, 5}
	for i := range num {
		name := strconv.Itoa(i)
		l.Push(&Event{
			Name: name,
			Data: i,
		})
		err := l.OnWithTimeout(name, nil, 5*time.Millisecond)
		time.Sleep(5 * time.Millisecond)
		if err == nil {
			t.Fatal("should throw error: invalid nil callback")
		}
	}

	l.Start()

	time.Sleep(10 * time.Millisecond)
	defer l.Stop()
}

func TestLoop_NewLogger(t *testing.T) {
	levels := []string{"debug", "info", "warn", "error", "fatal", "panic"}
	for _, l := range levels {
		_ = newLogger(l)
	}

	_ = newLogger("test")
}

/*
// https://stackoverflow.com/questions/52734529/testing-zap-logging-for-a-logger-built-from-a-custom-config
// https://github.com/uber-go/zap/blob/747abfb0b3b130c9cf699e451f5aebfda379b5d1/example_test.go
// https://medium.com/@KoheiMisu/validate-the-behavior-of-processing-using-zap-with-zaptest-df0f1f693800
func TestLoop_NewLogger(t *testing.T) {
	levels := []string{"debug", "info", "warn", "error"}
	calls := []string{"Debug", "Info", "Warn", "Error"}
	for i, l := range levels {
		testLogger(calls, l, i, t)
	}
	time.Sleep(1 * time.Millisecond)
}

func testLogger(calls []string, l string, i int, t *testing.T) {
	logger := newLogger(l)
	funcsMap := map[string]interface{} {
		"Debug": logger.Debug,
		"Info": logger.Info,
		"Warn": logger.Warn,
		"Error": logger.Error,
	}
	done := capture()
	for j := 0; j < i; j++ {
		if _, err := call(funcsMap, calls[j], "test"); err != nil {
			t.Fatal(err)
		}
	}
	captured, err := done()
	if err != nil {
		t.Fatal(err)
	}
	if captured != "" {
		t.Fatal(captured)
	}

	done = capture()
	fmt.Println("asdasd")
	if _, err := call(funcsMap, calls[i], "test"); err != nil {
		t.Fatal(err)
	}
	captured, err = done()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(captured, "test") || !strings.Contains(captured, l) {
		t.Fatal(captured, strings.Contains(captured, "test"), strings.Contains(captured, l))
	}
}

func capture() func() (string, error) {
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	done := make(chan error, 1)
	save := os.Stdout
	os.Stdout = w
	var buf strings.Builder
	go func() {
		_, err := io.Copy(&buf, r)
		_ = r.Close()
		done <- err
	}()

	return func() (string, error) {
		os.Stdout = save
		_ = w.Close()
		err := <-done
		return buf.String(), err
	}
}

func call(m map[string]interface{}, name string, args ... interface{}) (result []reflect.Value, err error) {
	f := reflect.ValueOf(m[name])
	if len(args) > f.Type().NumIn() {
		err = errors.New("invalid number of args")
		return
	}
	in := make([]reflect.Value, len(args))
	for k, param := range args {
		in[k] = reflect.ValueOf(param)
	}
	result = f.Call(in)
	return
}
*/
