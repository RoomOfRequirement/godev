package file

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"
)

func TestNewWatcher(t *testing.T) {
	// file not exist
	w, err := NewWatcher("notexist.txt", func() {
		log.Println("changed")
	}, time.Second)
	assert.Error(t, err)
	assert.Nil(t, w)

	// nil callback
	path := "test.txt"
	file, _ := os.Create(path)
	defer func() {
		_ = file.Close()
		_ = os.Remove(path)
	}()
	w, err = NewWatcher(path, nil, time.Second)
	assert.Error(t, err)
	assert.Nil(t, w)

	// zero checkInterval
	w, err = NewWatcher(path, func() {
		log.Println("changed")
	}, 0)
	assert.Error(t, err)
	assert.Nil(t, w)

	// normal
	w, err = NewWatcher(path, func() {
		log.Println("changed")
	}, time.Second)
	assert.NoError(t, err)
	assert.NotNil(t, w)
	defer w.Close()

	_, err = file.WriteString("test")
	assert.NoError(t, err)
	time.Sleep(2 * time.Second)
}
