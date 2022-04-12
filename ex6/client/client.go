package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"time"
)

const CLIENTS = 1
const REQUESTS = 1

func client() {
	var reply []byte
	times := [10000]time.Duration{}
	name_of_file := "ex6.txt"

	name_of_compressed_file := "ex6.gz"

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

		f, _ := os.Open("../file/" + name_of_file)
		read := bufio.NewReader(f)
		data, _ := ioutil.ReadAll(read)

		err = client.Call("File.Compress", data, &reply)

		if err != nil {
			log.Fatal(err)
		}

		g, _ := os.Create("../file/" + name_of_compressed_file)

		g.Write(reply)

		g.Close()

		times[i] = time.Since(t1)
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
