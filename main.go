package main

import (
	"github.com/streadway/amqp"
)

func main() {

	url := "amqp://guest:guest@localhost:5672"

	// Connect to the rabbitMQ instance
	connection, err := amqp.Dial(url)
	// Create a channel from the connection. We'll use channels to access the data in the queue rather than the connection itself.
	channel, err := connection.Channel()

	// We create an exchange that will bind to the queue to send and receive messages
	err = channel.ExchangeDeclare("events", "topic", true, false, false, false, nil)

	if err != nil {
		panic("could not open RabbitMQ channel:" + err.Error())
	}
	if err != nil {
		panic("could not establish connection with RabbitMQ:" + err.Error())
	}
}
