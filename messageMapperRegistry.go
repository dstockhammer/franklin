package franklin

import "errors"

// MessageMapperRegistry blabla
type MessageMapperRegistry struct {
	mappers map[string]MessageMapper
}

// MessageMapper maps messages to and from their string representations
type MessageMapper interface {
	// MapToMessage creates a typed message from its JSON representation
	MapToMessage(json []byte) (Message, error)
}

// Register adds a mapper registration to the message mapper registry
func (r *MessageMapperRegistry) Register(key string, mapper MessageMapper) error {
	if r.mappers[key] != nil {
		return errors.New("Registration with key <" + key + "> already exists")
	}

	r.mappers[key] = mapper

	return nil
}

// MapperForKey returns the message mapper for a key
func (r *MessageMapperRegistry) MapperForKey(key string) MessageMapper {
	return r.mappers[key]
}

// NewMessageMapperRegistry creates an empty message mapper registry
func NewMessageMapperRegistry() *MessageMapperRegistry {
	return &MessageMapperRegistry{
		mappers: make(map[string]MessageMapper)}
}
