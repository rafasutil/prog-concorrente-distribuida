package main

import (
	"fmt"
	"os"
	"github.com/streadway/amqp"
	"compress/gzip"
	"bytes"
)

const BUFFERSIZE = 1024

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	checkError(err)
	defer conn.Close()

	ch, err := conn.Channel()
	checkError(err)
	defer ch.Close()

	fileQueue, err := ch.QueueDeclare(
		"file", // name
		false,  // durable
		false,  // delete when usused
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

	requestCh, err := ch.Consume(
		requestQueue.Name, // queue
		"",                // consumer
		true,              // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)

	for {
		handleClient(ch, requestCh, fileQueue)
	}
}

func publish(ch *amqp.Channel, queue amqp.Queue, output []byte) {
	ch.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        (output),
			ContentEncoding: "gzip",
		})
}

func handleClient(ch *amqp.Channel, requestCh <-chan amqp.Delivery, fileQ amqp.Queue) {
	file := (<-requestCh).Body
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	gz.Write(file);
	gz.Close();

	fileSize := len(file)
	compressedFileSize := len(b.Bytes())

	fmt.Println("Before: ", fileSize,". After: ", compressedFileSize)
	
	// Send compressed file to client
	publish(ch, fileQ, b.Bytes())
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s ", err.Error())
		os.Exit(1)
	}
}