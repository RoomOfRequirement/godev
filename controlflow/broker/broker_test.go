package broker

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"sync"
	"testing"
)

func TestBroker(t *testing.T) {
	broker := newBroker()
	var _ Broker = broker

	assert.False(t, broker.isStopped())
	err := broker.Stop()
	assert.True(t, broker.isStopped())
	assert.NoError(t, err)

	bk := New()
	subA, err := bk.Subscribe("a")
	assert.NoError(t, err)
	assert.NoError(t, bk.Publish("a", []byte("test")))
	assert.Error(t, bk.Publish("x", []byte("test")), "no such topic")
	assert.Equal(t, string(<-subA), "test")
	assert.NoError(t, bk.Unsubscribe("test", subA))
	assert.NoError(t, bk.Stop())
	assert.Error(t, bk.Publish("a", []byte("test")), "broker stopped")
	subB, err := bk.Subscribe("b")
	assert.Nil(t, subB)
	assert.Error(t, err, "broker stopped")
	assert.Error(t, bk.Unsubscribe("a", subA), "broker stopped")
	assert.NoError(t, bk.Stop())

	bk = New()
	wg := &sync.WaitGroup{}
	for i := 0; i < 8; i++ {
		topic := strconv.Itoa(i)
		payload := []byte("test")
		sub, err := bk.Subscribe(topic)
		assert.NoError(t, err)
		wg.Add(1)

		go assert.NoError(t, bk.Publish(topic, payload))

		go func() {
			defer wg.Done()
			assert.Equal(t, string(<-sub), string(payload))
			assert.NoError(t, bk.Unsubscribe(topic, sub))
		}()
	}
	wg.Wait()
}

func BenchmarkBroker_Publish(b *testing.B) {
	b.ReportAllocs()
	broker := New()
	for i := 0; i < b.N; i++ {
		sub, _ := broker.Subscribe("test")
		_ = broker.Publish("test", []byte("test"))
		<-sub
		_ = broker.Unsubscribe("test", sub)
	}
}
