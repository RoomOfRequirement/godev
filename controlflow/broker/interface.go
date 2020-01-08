package broker

// Broker interface
type Broker interface {
	Publish(topic string, payload []byte) error
	Subscribe(topic string) (<-chan []byte, error)
	Unsubscribe(topic string, subChan <-chan []byte) error
	Stop() error
}
