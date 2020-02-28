package utils

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// Trace gets goroutine stack info
/*	use it in defer func with panic recovery:
	defer func() {
		if err := recover(); err != nil {
			message := fmt.Sprintf("%s", err)
			log.Println(Trace(message))
		}
	}
*/
func Trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

// PanicToErr extract panic from goroutine
// and wrap it to input error
//	see usage in `trace_test.go`
func PanicToErr(err *error) {
	// extract panic
	if e := recover(); e != nil {
		message := fmt.Sprintf("Panic Recovered from: %s", e)
		// log.Println(Trace(message))
		// wrap panic info to input error
		if err != nil {
			*err = errors.New(Trace(message))
		}
	}
}

// Fallback returns fallback if orig func panic
func Fallback(orig func() interface{}, fallback interface{}) (ret interface{}) {
	defer func() {
		if recover() != nil {
			ret = fallback
		}
	}()
	ret = orig()
	return
}
