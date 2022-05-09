package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/rpc"
	"strconv"
	"time"
)

const CLIENTS = 50

func client(i int) {
	var reply string

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

	for {
		time.Sleep(time.Duration(rand.Intn(15)) * time.Second)

		err = client.Call("SushiBar.P", i, &reply)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("O cliente " + strconv.Itoa(i) + " sentou " + reply)

		err = client.Call("SushiBar.V", i, &reply)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("O cliente " + strconv.Itoa(i) + " levantou " + reply)
	}
}

func main() {
	for i := 0; i < CLIENTS; i++ {
		go client(i)
	}
	fmt.Scanln()
}
