package franklin

import "reflect"

// MessageHandler receives messages representation and handles them
type MessageHandler interface {
	// MessageType returns the message type that this handler subscribes to
	MessageType() reflect.Type

	// Handle handles a message of the type returned by MessageType()
	Handle(message Message) error
}
