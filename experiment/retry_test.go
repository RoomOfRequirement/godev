package experiment

import (
	"errors"
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestRunWithBackOff(t *testing.T) {
	// success
	_ = RunWithBackOff(func() (success bool, err error) {
		return true, nil
	}, func(trialNo int, wait time.Duration, err error) {
		log.Println(trialNo, wait, err)
	}, 10, time.Second, 10*time.Second, 5*time.Second, 10*time.Second)

	// fail
	cancel := RunWithBackOff(func() (success bool, err error) {
		r := rand.Intn(5)
		switch r {
		case 0, 1, 2:
			return false, nil
		default:
			return false, errors.New("test error here")
		}
	}, func(trialNo int, wait time.Duration, err error) {
		log.Println(trialNo, wait, err)
	}, 10, time.Second, 10*time.Second, 5*time.Second, 10*time.Second)

	var wg sync.WaitGroup
	wg.Add(1)
	time.AfterFunc(20*time.Second, func() {
		cancel()
		wg.Done()
	})
	wg.Wait()
	time.Sleep(10 * time.Second)
}
