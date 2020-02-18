package utils

import (
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
