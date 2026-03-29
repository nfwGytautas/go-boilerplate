// Package broker provides a wrapper around a broker service for distributed messaging
// it provides a common interface that one can use and wraps around some popular providers
package broker

import "context"

// Producer defines a producer that is able to send a message to a broker, the interface
// is implemented by both async and sync producers
type Producer[T any] interface {
	Send(context.Context, T) error
}

// Consumer defines a consumer that is able to handle messages from a broker if it returns
// an error the message is deemed not consumed
type Consumer[T any] interface {
	Handle(context.Context, T) error
}
