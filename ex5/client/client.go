package main

import (
	"fmt"
	"log"
	"net/rpc"
	"time"
)

const CLIENTS = 1
const REQUESTS = 10000

var dates = []string{
	"12/03/2022",
	"13/03/2022",
	"14/03/2022",
	"15/03/2022",
	"16/03/2022",
	"17/03/2022",
	"18/03/2022",
	"19/03/2022",
	"20/03/2022",
	"21/03/2022",
}

func client() {
	var reply string
	times := [10000]time.Duration{}

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

		err = client.Call("Weather.Result", dates[i%10], &reply)
		if err != nil {
			log.Fatal(err)
		}

		times[i] = time.Now().Sub(t1)
	}

	totalTime := time.Duration(0)
	for i := range times {
		totalTime += times[i]
	}
	fmt.Printf("Total Duration: %v [%v]\n", totalTime, REQUESTS)
}

func main() {

	for i := 0; i < CLIENTS; i++ {
		go client()
	}

	fmt.Scanln()
}
