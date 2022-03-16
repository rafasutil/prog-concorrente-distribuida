package main

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
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
	}

		addr, err := net.ResolveUDPAddr("udp", "localhost:1313")
		checkError(err)	

		var wg sync.WaitGroup

		// open the file named results.csv, and all new writes will be appended
		file, err := os.OpenFile("results.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		checkError(err)
	
		for i := 0; i < 10; i++{
			wg.Add(1)
			go HandleClientUDP(len(dates), addr, &wg, i, file)
		}
	
		wg.Wait()

}


func HandleClientUDP(n int, addr *net.UDPAddr, wg *sync.WaitGroup, clientIndex int, file *os.File) {
	req := make([]byte, 10)
	rep := make([]byte, 1024)

	conn, err := net.DialUDP("udp", nil, addr)
	checkError(err)

	defer conn.Close()
	defer wg.Done()

	for i := 0; i < 10000; i++ {
		time1 := time.Now()
		req = []byte(dates[i%n])

		_, err := conn.Write(req)
		checkError(errr)

		_, _, err = conn.ReadFromUDP(rep)
		checkError(err)

		fmt.Println("Getting info about the date: " + string(req))
		fmt.Println(string(rep))

		time2 := time.Now()
		elapsedTime := float64(time2.Sub(time1).Nanoseconds()) / 1000000 // time in ms
		
		// just print in file is is the first client
		if(clientIndex == 0){
			f.WriteString(fmt.Sprintf("%f", elapsedTime) + "\n")
		}
	}

}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(0)
	}
}