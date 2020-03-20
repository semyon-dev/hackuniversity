// by Semyon

package main

import (
	"encoding/json"
	"fmt"
	"github.com/semyon-dev/hackuniversity/rabbitMQ/model"
	"github.com/streadway/amqp"
	"math/rand"
	"time"
)

var conn *amqp.Connection
var ch *amqp.Channel

func main() {

	//Make a connection
	var err error
	conn, err = amqp.Dial("amqp://guest:guest@192.168.1.106:5672/")
	if err != nil {
		fmt.Print(err.Error())
	}
	defer conn.Close()

	//Create a channel
	ch, err = conn.Channel()
	if err != nil {
		fmt.Print(err.Error())
	}
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
				Body:        body,
			})
		fmt.Println("Message:", body)
		time.Sleep(1 * time.Second)
	}
}

// генерируем рандомные параметры
func generate() []byte {
	min := 0.0
	max := 100.0
	data := model.Data{
		PRESSURE: rand.Float64() * (max - min),
		HUMIDITY: rand.Float64() * (max - min),
		TEMP_HOM: rand.Float64() * (max - min),
		TEMP_WOR: rand.Float64() * (max - min),
		LEVEL:    rand.Float64() * (max - min),
		MASS:     rand.Float64() * (max - min),
		WATER:    rand.Float64() * (max - min),
		LEVELCO2: rand.Float64() * (max - min),
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Print(err)
	}
	return jsonData
}
