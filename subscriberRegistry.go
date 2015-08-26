package franklin

import (
	"errors"
	"reflect"
)

// SubscriberRegistry blabla
type SubscriberRegistry struct {
	handlers map[string]MessageHandler
	keys     map[reflect.Type]string
}

// Consumers returns the subscriber registrations
func (r *SubscriberRegistry) Consumers() map[string]MessageHandler {
	return r.handlers
}

// KeyForMessage returns the key for a message
func (r *SubscriberRegistry) KeyForMessage(message Message) string {
	messageType := reflect.TypeOf(message)
	return r.keys[messageType]
}

// Register adds a handler registration to the subscriber registry
func (r *SubscriberRegistry) Register(key string, handler MessageHandler) error {
	if r.handlers[key] != nil {
		return errors.New("Registration with key <" + key + "> already exists")
	}

	r.handlers[key] = handler
	r.keys[handler.MessageType()] = key

	return nil
}

// NewSubscriberRegistry creates an empty subscriber registry
func NewSubscriberRegistry() *SubscriberRegistry {
	return &SubscriberRegistry{
		handlers: make(map[string]MessageHandler),
		keys:     make(map[reflect.Type]string)}
}
