package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestTrace(t *testing.T) {
	panicFunc := func() interface{} {
		panic("hi, panic here, need help")
		return nil
	}
	call := func() interface{} {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Println(Trace(message))
			}
		}()
		return panicFunc()
	}
	assert.NotPanics(t, func() {
		call()
	})
}
