package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"sync"
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

func TestPanicToErr(t *testing.T) {
	var err error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer PanicToErr(&err)
		panic("hi, panic here, need help")
	}()
	wg.Wait()
	assert.Error(t, err)
	t.Log(err)
}
