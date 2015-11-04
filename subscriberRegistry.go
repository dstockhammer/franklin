package franklin

import (
	"fmt"
	"reflect"
)

// SubscriberRegistry blabla
type SubscriberRegistry struct {
	handlers map[string]MessageHandler
	queues   map[string]string
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

// QueueForKey returns the queue for a key
func (r *SubscriberRegistry) QueueForKey(key string) string {
	return r.queues[key]
}

// Register adds a handler registration to the subscriber registry
func (r *SubscriberRegistry) Register(key string, queue string, handler MessageHandler) error {
	if _, exists := r.keys[handler.MessageType()]; exists {
		return fmt.Errorf("Registration for key %s already exists", key)
	}

	r.handlers[key] = handler
	r.queues[key] = queue
	r.keys[handler.MessageType()] = key

	return nil
}

// RegisterExternal adds an external handler registration to the subscriber registry
func (r *SubscriberRegistry) RegisterExternal(key string, messageType reflect.Type) error {
	if _, exists := r.keys[messageType]; exists {
		return fmt.Errorf("Registration for key %s already exists", key)
	}

	r.keys[messageType] = key

	return nil
}

// NewSubscriberRegistry creates an empty subscriber registry
func NewSubscriberRegistry() *SubscriberRegistry {
	return &SubscriberRegistry{
		handlers: make(map[string]MessageHandler),
		queues:   make(map[string]string),
		keys:     make(map[reflect.Type]string)}
}
