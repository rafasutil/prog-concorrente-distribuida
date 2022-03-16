package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
	"sync"
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
		"19/03/2022",
		"20/03/2022",
		"21/03/2022",
	}
	var wg sync.WaitGroup

	// retorna o endereço do endpoint TCP
	r, err := net.ResolveTCPAddr("tcp", "localhost:1313")
	checkError(err)

	for i := 0; i < 10; i++{
		wg.Add(1)
		// connecta ao servidor (sem definir uma porta local)
		conn, err := net.DialTCP("tcp", nil, r)
		checkError(err)
		if err != nil {
			continue
		}
		go HandleClientTCP(len(dates), conn, &wg)
	}

	wg.Wait()
}

func HandleClientTCP(n int, conn net.Conn, wg *sync.WaitGroup) {
	time1 := time.Now()

	defer conn.Close()
	defer wg.Done()

	for i:= 0; i < 10000; i++ {
		// cria request
		req := dates[i%n]


		// envia mensage para o servidor
		_, err := fmt.Fprintf(conn, req+"\n")
		checkError(err)

		// recebe resposta do servidor
		rep, err := bufio.NewReader(conn).ReadString('\n')
		checkError(err)
		fmt.Println("Getting info about the date: " + req)
		fmt.Print(rep)
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