package shared

import "sync"

// ConsumerControl is an enumeration used by the Producer.Control() channel
type ConsumerControl int

const (
	// ConsumerControlStop will cause the consumer to halt and shutdown.
	ConsumerControlStop = ConsumerControl(1)

	// ConsumerControlRoll notifies the consumer about a reconnect or reopen request
	ConsumerControlRoll = ConsumerControl(2)
)

// Consumer is an interface for plugins that recieve data from outside sources
// and generate Message objects from this data.
type Consumer interface {
	// Consume should implement to main loop that fetches messages from a given
	// source and pushes it to the Message channel.
	Consume(*sync.WaitGroup)

	// IsActive returns true if the consumer is ready to generate messages.
	IsActive() bool

	// Control returns write access to this consumer's control channel.
	// See ConsumerControl* constants.
	Control() chan<- ConsumerControl

	// Messages returns a read only access to the messages generated by the
	// consumer.
	Messages() <-chan Message
}
