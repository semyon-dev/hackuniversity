package main

import (
	"encoding/binary"
	"fmt"
	"github.com/streadway/amqp"
	"math"
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

	for {
		//Publish a message
		body := generate()
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        float64ToByte(body),
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

func float64ToByte(f float64) []byte {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], math.Float64bits(f))
	return buf[:]
}
