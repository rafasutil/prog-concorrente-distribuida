package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	// "time"
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
	// time1 := time.Now()

	r, err := net.ResolveTCPAddr("tcp", "localhost:1313")
	checkError(err)

	// for i := 0; i < 10; i++{
		conn, err := net.DialTCP("tcp", nil, r)
		checkError(err)
		// if err != nil {
		// 	continue
		// }
		HelloClientTCP(len(dates), conn)
	// }


	// HelloClientTCP(len(dates))
	// HelloClientUDP(len(dates))

	// time2 := time.Now()
	// elapsedTime := float64(time2.Sub(time1).Milliseconds())
	// fmt.Println(elapsedTime)
}

func HelloClientTCP(n int, conn net.Conn) {

	// retorna o endereço do endpoint TCP
	// r, err := net.ResolveTCPAddr("tcp", "localhost:1313")
	// checkError(err)

	// connecta ao servidor (sem definir uma porta local)
	// conn, err := net.DialTCP("tcp", nil, r)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(0)
	// }

	// fecha connexão
	// defer func(conn *net.TCPConn) {
	// 	err := conn.Close()
	// 	checkError(err)
	// }(conn)
	defer conn.Close()

	for i := 0; i < n; i++ {
		for j:= 0; j < 1000; j++ {
		// cria request
		req := dates[i]


		// envia mensage para o servidor
		_, err := fmt.Fprintf(conn, req+"\n")
		checkError(err)

		// recebe resposta do servidor
		rep, err := bufio.NewReader(conn).ReadString('\n')
		checkError(err)
		fmt.Println("Getting info about the date: " + req)
		fmt.Print(rep)
		}
	}
}

func HelloClientUDP(n int) {
	req := make([]byte, 10)
	rep := make([]byte, 1024)

	// retorna o endereço do endpoint UDP
	addr, err := net.ResolveUDPAddr("udp", "localhost:1313")
	checkError(err)

	// conecta ao servidor -- não cria uma conexão
	conn, err := net.DialUDP("udp", nil, addr)
	checkError(err)

	// desconecta do servidor
	defer func(conn *net.UDPConn) {
		err := conn.Close()
		checkError(err)
	}(conn)

	for i := 0; i < n; i++ {
		for j:= 0; j < 1000; j++ {
		// cria request
		// req = []byte("Mensagem " + strconv.Itoa(i))
		req = []byte(dates[i])

		// envia request ao servidor
		_, err = conn.Write(req)
		checkError(err)

		// recebe resposta do servidor
		_, _, err := conn.ReadFromUDP(rep)
		checkError(err)

		fmt.Println("Getting info about the date: " + string(req))
		fmt.Println(string(rep))
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}