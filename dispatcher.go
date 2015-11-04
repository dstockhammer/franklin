package franklin

import (
	"encoding/json"
	"log"
	"reflect"

	"github.com/streadway/amqp"
)

// Dispatcher is going to be a go port of the C# ServiceActivator.
// Currently it's just an arbitrary implementation that connects to AMQ
// and consumes messages.
type Dispatcher interface {
	Receive()
	Close()
}

type amqpDispatcher struct {
	subscriberRegistry    *SubscriberRegistry
	messageMapperRegistry *MessageMapperRegistry
	connection            *amqp.Connection
	channel               *amqp.Channel
	exchange              string
}

// Receive begins listening for messages on channels,
// and dispatching them to request handlers.
func (d *amqpDispatcher) Receive() {
	for key, consumer := range d.subscriberRegistry.Consumers() {
		queueName := d.subscriberRegistry.QueueForKey(key)
		q, err := d.channel.QueueDeclare(
			queueName, // name
			true,      // durable
			false,     // delete when usused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)
		failOnError(err, "Failed to declare queue")

		d.channel.QueueBind(
			queueName,  // name of the queue
			key,        // bindingKey
			d.exchange, // sourceExchange
			false,      // noWait
			nil,        // arguments
		)
		failOnError(err, "Failed to bind queue to exchange")

		deliveries, err := d.channel.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		failOnError(err, "Failed to register a consumer")

		go d.consume(key, consumer, deliveries)
	}

	forever := make(chan bool)
	<-forever
}

func (d *amqpDispatcher) Close() {
	d.channel.Close()
	d.connection.Close()
}

func (d *amqpDispatcher) consume(key string, handler MessageHandler, deliveries <-chan amqp.Delivery) {
	log.Printf("Consuming %s...", key)

	for delivery := range deliveries {
		log.Printf("%s: Received %d bytes: [%v] %q",
			key, len(delivery.Body), delivery.DeliveryTag, delivery.Body)

		messageType := handler.MessageType()
		var message Message
		var err error

		messageMapper := d.messageMapperRegistry.MapperForType(messageType)
		if messageMapper != nil {
			message, err = messageMapper.MapToMessage(delivery.Body)
		} else {
			// todo: investigate why this doesn't work
			message = reflect.Zero(messageType)
			err = json.Unmarshal([]byte(delivery.Body), &message)
		}

		if err != nil {
			log.Printf("Error: %s", err.Error())
		} else {
			err = handler.Handle(message)
			if err != nil {
				log.Printf("Error: %s", err.Error())
			}
		}
	}
}

// InitialiseDispatcher initialises an amqpDispatcher and
// establishes an AMQ connection.
func InitialiseDispatcher(url string, exchange string, subscriberRegistry *SubscriberRegistry, messageMapperRegistry *MessageMapperRegistry) Dispatcher {
	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	return &amqpDispatcher{
		subscriberRegistry:    subscriberRegistry,
		messageMapperRegistry: messageMapperRegistry,
		connection:            conn,
		channel:               ch,
		exchange:              exchange}
}
