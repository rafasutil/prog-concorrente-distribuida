package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/rpc"
	"time"
	"math"
)

const CLIENTS = 9
const REQUESTS = 100
var responses = 0
var responsesTime [100]time.Duration

func client(id int) {
	var reply []byte
	name_of_file := "ex6.txt"

	client, err := rpc.Dial("tcp", "localhost:1313")
	if err != nil {
		log.Fatal(err)
	}

	defer func(client *rpc.Client) {
		var err = client.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(client)

	for i := 0; i < REQUESTS; i++ {
		t1 := time.Now()

		content, _ := ioutil.ReadFile("file/" + name_of_file)

		err = client.Call("File.Compress", content, &reply)

		if err != nil {
			log.Fatal(err)
		}

		if id == 0 {
			responsesTime[i] = time.Since(t1)
			responses++
		}

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

func main() {
	for i := 0; i < CLIENTS; i++ {
		go client(i)
	}

	calculate()
	fmt.Scanln()
}
