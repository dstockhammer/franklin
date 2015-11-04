package franklin

import "reflect"

// MessageMapper maps messages to and from their string representations
type MessageMapper interface {
	// MessageType returns the message type that this mapper maps
	MessageType() reflect.Type

	// MapToMessage creates a typed message from its JSON representation
	MapToMessage(json []byte) (Message, error)
}
