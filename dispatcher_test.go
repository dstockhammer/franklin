package franklin

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/streadway/amqp"
)

type FooMessage struct {
	MessageID string `json:"messageId"`
	Value     string `json:"value"`
}

type FooMessageMapper struct {
}

func (m FooMessageMapper) MapToMessage(body []byte) (Message, error) {
	var message FooMessage
	err := json.Unmarshal(body, &message)
	return message, err
}

type FooMessageHandler struct {
	received chan FooMessage
}

func (h FooMessageHandler) MessageType() reflect.Type {
	return reflect.TypeOf(FooMessage{})
}

func (h FooMessageHandler) Handle(message Message) error {
	fooMesssage, ok := message.(FooMessage)
	if !ok {
		return errors.New("Message is not a FooMessage")
	}

	h.received <- fooMesssage

	return nil
}

func TestDispatcher(t *testing.T) {
	Convey("When receiving a foo message", t, func() {
		Convey("It should dispatch the message correctly", func() {
			key := "Foo.Message"

			received := make(chan FooMessage)

			fooHandler := &FooMessageHandler{
				received: received}

			subscriberRegistry := NewSubscriberRegistry()
			subscriberRegistry.Register(key, fooHandler)

			messageMapperRegistry := NewMessageMapperRegistry()
			messageMapperRegistry.Register(key, &FooMessageMapper{})

			dispatcher := &amqpDispatcher{
				subscriberRegistry:    subscriberRegistry,
				messageMapperRegistry: messageMapperRegistry}

			deliveries := make(chan amqp.Delivery)

			go dispatcher.consume(key, fooHandler, deliveries)

			deliveries <- amqp.Delivery{
				Body:        []byte("{\"messageId\": \"1337\", \"value\": \"hello world\"}"),
				DeliveryTag: 0}

			close(deliveries)

			receivedMessage := <-received

			close(received)

			So(receivedMessage, ShouldNotBeNil)
			So(receivedMessage.MessageID, ShouldEqual, "1337")
			So(receivedMessage.Value, ShouldEqual, "hello world")
		})
	})
}
