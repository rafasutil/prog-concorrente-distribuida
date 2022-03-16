package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
	"sync"
	// "strconv"
)

var dates []string

func main() {
	dates = []string{
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
	var wg sync.WaitGroup

	r, err := net.ResolveTCPAddr("tcp", "localhost:1313")
	checkError(err)
	
	// open the file named results.csv, and all new writes will be appended
	file, err := os.OpenFile("results.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		checkError(err)

	for i := 0; i < 10; i++{
		wg.Add(1)
		conn, err := net.DialTCP("tcp", nil, r)
		checkError(err)
		if err != nil {
			continue
		}
		go HandleClientTCP(len(dates), conn, &wg, i, file)
	}

	wg.Wait()
}

func HandleClientTCP(n int, conn net.Conn, wg *sync.WaitGroup, clientIndex int, file *os.File) {

	defer conn.Close()
	defer wg.Done()

	for i:= 0; i < 10000; i++ {
		time1 := time.Now()

		req := dates[i%n]

		_, err := fmt.Fprintf(conn, req+"\n")
		checkError(err)

		rep, err := bufio.NewReader(conn).ReadString('\n')
		checkError(err)

		fmt.Println("Getting info about the date: " + req)
		fmt.Print(rep)
		
		time2 := time.Now()
		elapsedTime := float64(time2.Sub(time1).Nanoseconds()) / 1000000 // time in ms
		
		// just print in file is is the first client
		if(clientIndex == 0){
			file.WriteString(fmt.Sprintf("%f", elapsedTime) + "\n")
		}
	}

}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(0)
	}
}