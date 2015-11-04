package franklin

import (
	"fmt"
	"reflect"
)

// MessageMapperRegistry is currently not in use
type MessageMapperRegistry struct {
	mappers map[reflect.Type]MessageMapper
}

// Register adds a mapper registration to the message mapper registry
func (r *MessageMapperRegistry) Register(mapper MessageMapper) error {
	messageType := mapper.MessageType()

	if _, exists := r.mappers[messageType]; exists {
		return fmt.Errorf("Mapper for %s is already registered", messageType)
	}

	r.mappers[messageType] = mapper

	return nil
}

// MapperForType returns the message mapper for a message type
func (r *MessageMapperRegistry) MapperForType(messageType reflect.Type) MessageMapper {
	if mapper, exists := r.mappers[messageType]; exists {
		return mapper
	}

	return nil
}

// NewMessageMapperRegistry creates an empty message mapper registry
func NewMessageMapperRegistry() *MessageMapperRegistry {
	return &MessageMapperRegistry{
		mappers: make(map[reflect.Type]MessageMapper)}
}
