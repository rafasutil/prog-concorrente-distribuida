package main

import (
	"fmt"
	"os"
	"strconv"
	"github.com/streadway/amqp"
	"io/ioutil"
	"time"
	"math"
)

const BUFFERSIZE = 1024
const fileName = "file.txt"
const REQUESTS = 10
const CLIENTS = 1
var responses = 0
var responsesTime [REQUESTS]time.Duration

func main() {
	for i := 0; i < CLIENTS; i++ {
		go client(i)
	}

	calculate()
}

func client(id int){
	for i := 0; i < REQUESTS; i++ {
		fmt.Println("Client", id, "requesting file", i)
		t1 := time.Now()
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
		requestFile(input, ch, fileQueue, requestQueue, i)

		if id == 0 {
			responsesTime[i] = time.Since(t1)
			responses++
		}

		ch.Close()
		conn.Close()
	}
}

func requestFile(input []byte, ch *amqp.Channel, fileQ amqp.Queue, requestQ amqp.Queue, i int) {
	// Sending File
	err := ch.Publish(
		"",            // exchange
		requestQ.Name, // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        input,
			ContentEncoding: "gzip",
		})
	checkError(err)

	// Getting compressed file
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


	fileContent := (<-fileCh).Body
	file, err := os.Create(strconv.Itoa(i) + ".gz")
	checkError(err)
	defer file.Close()
	file.Write(fileContent)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func calculate(){
	for {
		time.Sleep(1 * time.Second)
		if responses >= REQUESTS {
			sum := float64(0)
			for _, element := range responsesTime {
				t := element / time.Millisecond
				sum += float64(t)
			}

			average := sum / float64(REQUESTS)

			sum = float64(0)
			for _, element := range responsesTime {
				t := element / time.Millisecond
				sum += ((float64(t) - average) * (float64(t) - average))
			}

			deviation := math.Sqrt(float64(sum) / float64(time.Duration(REQUESTS-1)))
			fmt.Print(average, deviation)
			break
		}
	}
}