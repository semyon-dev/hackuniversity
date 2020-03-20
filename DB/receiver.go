package DB

import (
	"fmt"
	"github.com/streadway/amqp"
	"time"
)

var message  chan string

func ReceiveMessages() {
	//Make a connection
	conn, err := amqp.Dial("amqp://guest:guest@192.168.1.106:5672/")
	if err != nil {
		fmt.Print(err.Error())
	}
	defer conn.Close()

	//Create a channel
	ch, err := conn.Channel()
	if err != nil {
		fmt.Print(err.Error())
	}
	defer ch.Close()
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
			message <- string(msg.Body)

		}
		fmt.Println("------------------------------")

		// We close the connection after the operation has completed.
		defer conn.Close()
	}
}

func insertInDB()  {
	for{
		select {
		case recieved := <- message:
			fmt.Println(recieved)
		}
	}
}


