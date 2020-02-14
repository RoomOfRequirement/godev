package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"strings"
	"time"
	"unsafe"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GenerateRandomInt returns a random int
func GenerateRandomInt() int {
	return rand.Int() - rand.Int()
}

// GenerateRandomIntInRange returns a random int in range [start, end)
func GenerateRandomIntInRange(start, end int) int {
	if start > end {
		start, end = end, start
	}
	return start + rand.Intn(end-start)
}

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

// GenerateRandomString returns a random string consists of alphabets with length of `length`
func GenerateRandomString(length int) string {
	// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
	s := make([]byte, length)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := length-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			s[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&s))
}

// PartialFunc based on reflect
func PartialFunc(funcMap map[string]interface{}, funcName string, funcArgs ...interface{}) (resSlice []reflect.Value, err error) {
	f := reflect.ValueOf(funcMap[funcName])
	if len(funcArgs) != f.Type().NumIn() {
		err = errors.New("invalid number of funcArgs")
		return
	}
	in := make([]reflect.Value, len(funcArgs))
	for idx, arg := range funcArgs {
		in[idx] = reflect.ValueOf(arg)
	}
	resSlice = f.Call(in)
	return
}

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
