package broker

import (
	"errors"
	"sync"
	"time"
)

// CacheMessagePerChan means only cache 10 messages in default
const CacheMessagePerChan = 10
// DefaultPubTimeout means default publish timeout in default
const DefaultPubTimeout = time.Millisecond * 10

type broker struct {
	topics map[string][]chan []byte

	sync.RWMutex

	stop chan struct{}
}

// New returns a new Broker
func New() Broker {
	return newBroker()
}

func newBroker() *broker {
	return &broker{
		topics:  make(map[string][]chan []byte),
		RWMutex: sync.RWMutex{},
		stop:    make(chan struct{}),
	}
}

func (b *broker) isStopped() bool {
	select {
	case <-b.stop:
		return true
	default:
		return false
	}
}

func (b *broker) Publish(topic string, payload []byte) error {
	if b.isStopped() {
		return errors.New("broker stopped")
	}

	b.RLock()
	subChans, found := b.topics[topic]
	b.RUnlock()
	if !found {
		return errors.New("no such topic")
	}

	// 100 subs per goroutine
	l := len(subChans)
	cc := l / 100 + 1

	for i := 0; i < cc; i++ {
		go func(i int) {
			for j := i; j < l; j += cc {
				select {
				case <-b.stop:
					// stopped
					return
				case subChans[j] <- payload:
					// pub
				case <-time.After(DefaultPubTimeout):
					// timeout
				}
			}
		}(i)
	}
	return nil
}

func (b *broker) Subscribe(topic string) (<-chan []byte, error) {
	if b.isStopped() {
		return nil, errors.New("broker stopped")
	}

	subChan := make(chan[]byte, CacheMessagePerChan)
	b.Lock()
	b.topics[topic] = append(b.topics[topic], subChan)
	b.Unlock()
	return subChan, nil
}

func (b *broker) Unsubscribe(topic string, subChan <-chan []byte) error {
	if b.isStopped() {
		return errors.New("broker stopped")
	}

	b.RLock()
	subChans, ok := b.topics[topic]
	b.RUnlock()
	if !ok {
		return nil
	}

	newSubChans := make([]chan []byte, 0, len(subChans))
	for _, sc := range subChans {
		if sc == subChan {
			continue
		}
		newSubChans = append(newSubChans, sc)
	}

	b.Lock()
	b.topics[topic] = newSubChans
	b.Unlock()
	return nil
}

func (b *broker) Stop() error {
	if b.isStopped() {
		return nil
	}
	close(b.stop)
	b.Lock()
	b.topics = nil
	b.Unlock()
	return nil
}
