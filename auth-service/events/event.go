package events

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

/**

Producer -------> Exchange --------> Queue -------> Consumer

*/

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs",  // name
		"topic", // type
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
}
