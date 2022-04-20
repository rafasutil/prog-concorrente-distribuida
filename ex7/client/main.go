package main

import (
	"fmt"
	"os"
	"strconv"
	"github.com/streadway/amqp"
	"io/ioutil"
)

const BUFFERSIZE = 1024
const fileName = "file.txt"
const REQUESTS = 1

func main() {
	for i := 0; i < REQUESTS; i++ {
		conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		checkError(err)
		ch, err := conn.Channel()
		checkError(err)

		fileQueue, _ := ch.QueueDeclare(
			"file", // name
			false,  // durable
			false,  // delete when unused
			false,  // exclusive
			false,  // no-wait
			nil,    // arguments
		)

		requestQueue, err := ch.QueueDeclare(
			"request", // name
			false,     // durable
			false,     // delete when usused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)

		checkError(err)

		input, _ := ioutil.ReadFile(fileName)
		requestFile(input, ch, fileQueue, requestQueue)

		ch.Close()
		conn.Close()
	}

	os.Exit(0)
}

func requestFile(input []byte, ch *amqp.Channel, fileQ amqp.Queue, requestQ amqp.Queue) {
	// Getting File Size
	fileCh, err := ch.Consume(
		fileQ.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	checkError(err)

	// Sending File
	err = ch.Publish(
		"",            // exchange
		requestQ.Name, // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        input,
		})
	checkError(err)

	fileSize, _ := (strconv.ParseInt(string((<-fileCh).Body), 10, 64))
	file, err := os.Create(fileName + ".gz")
	checkError(err)
	defer file.Close()

	var recSize int64
	recSize = 0
	for {
		if (fileSize - recSize) < BUFFERSIZE {
			file.Write((<-fileCh).Body[:(fileSize - recSize)])
			recSize = fileSize
			break
		}
		file.Write((<-fileCh).Body)
		recSize += BUFFERSIZE
		if recSize == fileSize {
			break
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}