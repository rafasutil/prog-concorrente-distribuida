package main

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
	"strconv"
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

		// retorna o endereço do endpoint UDP
		addr, err := net.ResolveUDPAddr("udp", "localhost:1313")
		checkError(err)	

		var wg sync.WaitGroup
	
		for i := 0; i < 10; i++{
			wg.Add(1)
			go HandleClientUDP(len(dates), addr, &wg)
		}
	
		wg.Wait()

}


func HandleClientUDP(n int, addr *net.UDPAddr, wg *sync.WaitGroup) {
	time1 := time.Now()
	req := make([]byte, 10)
	rep := make([]byte, 1024)

	conn, err := net.DialUDP("udp", nil, addr)
	checkError(err)

	defer conn.Close()
	defer wg.Done()

	for i := 0; i < 10000; i++ {
		req = []byte(dates[i%n])

		_, err := conn.Write(req)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		_, _, err = conn.ReadFromUDP(rep)
		checkError(err)

		fmt.Println("Getting info about the date: " + string(req))
		fmt.Println(string(rep))
	}
	time2 := time.Now()
	elapsedTime := float64(time2.Sub(time1).Milliseconds())
	f, err := os.OpenFile("results.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	checkError(err)
	f.WriteString(strconv.Itoa(int(elapsedTime)) + "\n")
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}