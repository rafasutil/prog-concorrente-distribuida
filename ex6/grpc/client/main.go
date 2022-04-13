package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"math"
	"time"
	"google.golang.org/grpc"
	pb "ex6/fileservice"
)

const CLIENTS = 10
const REQUESTS = 100
var responses = 0

var responsesTime [100]time.Duration

func compressFile(client pb.FileServiceClient, id int) {
	name_of_file := "ex6.txt"
	for i := 0; i < REQUESTS; i++ {
		content, err := ioutil.ReadFile("../file/" + name_of_file)
		if err != nil {
			fmt.Println(err)
		}
		timeBeforeRequest := time.Now()
		_, err = client.Compress(context.Background(), &pb.FileRequest{File: content})
		if err != nil {
			fmt.Println(err)
		}
		if id == 0 {
			responsesTime[i] = time.Since(timeBeforeRequest)
			responses++
		}
	}
}

func client(i int){
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	defer conn.Close()
	c := pb.NewFileServiceClient(conn)
	if err != nil {
		fmt.Printf("did not connect: %v", err)
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	compressFile(c, i)
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
	
}
