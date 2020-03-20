package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"math/rand"
	"time"
)

func main() {

	//Make a connection
	conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")
	defer conn.Close()

	//Create a channel
	ch, _ := conn.Channel()
	defer ch.Close()

	//Declare a queue
	q, err := ch.QueueDeclare(
		"data", // name of the queue
		false,  // should the message be persistent? also queue will survive if the cluster gets reset
		false,  // autodelete if there's no consumers (like queues that have anonymous names, often used with fanout exchange)
		false,  // exclusive means I should get an error if any other consumer subsribes to this queue
		false,  // no-wait means I don't want RabbitMQ to wait if there's a queue successfully setup
		nil,    // arguments for more advanced configuration
	)
	if err != nil {
		fmt.Println(err.Error())
	}

	go func() {
		for {
			time.Sleep(5 * time.Second)
			// We consume data from the queue named data using the channel we created in go.
			msgs, err := ch.Consume("data", "", false, false, false, false, nil)

			if err != nil {
				fmt.Println("error consuming the queue: " + err.Error())
			}

			// We loop through the messages in the queue and print them in the console.
			// The msgs will be a go channel, not an amqp channel
			fmt.Println("------------------------------")
			for msg := range msgs {
				fmt.Println("message received: " + string(msg.Body))
				err := msg.Ack(false)
				if err != nil {
					fmt.Print(err.Error())
				}
			}
			fmt.Println("------------------------------")

			// We close the connection after the operation has completed.
			defer conn.Close()
		}
	}()

	for {
		//Publish a message
		body := fmt.Sprintf("%f", generate())
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		fmt.Println("Message:", body)
		time.Sleep(1 * time.Second)
	}
}

// генерируем рандомные параметры
func generate() float64 {
	min := 0.0
	max := 100.0
	return min + rand.Float64()*(max-min)
}
