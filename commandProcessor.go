package franklin

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/streadway/amqp"
)

// CommandProcessor is a very early stage go port of the C# CommandProcessor,
// currently only posting messages over AMQ is supported.
type CommandProcessor interface {
	Post(message Message) error
	Close()
}

type amqpCommandProcessor struct {
	subscriberRegistry *SubscriberRegistry
	connection         *amqp.Connection
	channel            *amqp.Channel
	exchange           string
}

func (cp *amqpCommandProcessor) Post(message Message) error {
	body, err := json.Marshal(message)
	if err != nil {
		log.Print(err.Error())
		return errors.New("Failed to marshal message")
	}

	routingKey := cp.subscriberRegistry.KeyForMessage(message)

	err = cp.channel.Publish(
		cp.exchange, // exchange
		routingKey,  // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body)})

	if err != nil {
		log.Print(err.Error())
		return errors.New("Failed to publish message")
	}

	log.Printf("Sent %s", body)
	return nil
}

func (cp *amqpCommandProcessor) Close() {
	cp.channel.Close()
	cp.connection.Close()
}

// InitialiseCommandProcessor initialises an amqpCommandProcessor and
// establishes an AMQ connection.
func InitialiseCommandProcessor(url string, exchange string, subscriberRegistry *SubscriberRegistry) CommandProcessor {
	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	err = ch.ExchangeDeclare(
		exchange, // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare exchange")

	return &amqpCommandProcessor{
		subscriberRegistry: subscriberRegistry,
		connection:         conn,
		channel:            ch,
		exchange:           exchange}
}
