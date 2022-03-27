package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type DateWeather struct {
	date               string // DD/MM/YYYY
	mininumTemperature int
	maximumTemperature int
	mayRain            bool
}

var datesWeather []DateWeather

func main() {
	datesWeather = []DateWeather{
		{"12/03/2022", 20, 28, true},
		{"13/03/2022", 23, 30, false},
		{"14/03/2022", 24, 32, true},
		{"15/03/2022", 21, 30, false},
		{"16/03/2022", 18, 27, true},
		{"17/03/2022", 25, 33, false},
		{"18/03/2022", 23, 31, true},
		{"19/03/2022", 21, 34, false},
		{"20/03/2022", 22, 32, true},
		{"21/03/2022", 24, 30, true},
	}

	// define o endpoint do servidor TCP
	r, err := net.ResolveTCPAddr("tcp", "localhost:1313")
	checkError(err)

	// cria um listener TCP
	ln, err := net.ListenTCP("tcp", r)
	checkError(err)

	for {
		// aguarda/aceita conexão
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go HandleTCP(conn)
	}

	_, _ = fmt.Scanln()
}

func HandleTCP(conn net.Conn) {
	fmt.Println("Servidor TCP aguardando conexão...")
	defer conn.Close()

	// recebe e processa requests
	for {
		// recebe request terminado com '\n'
		req, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil && err.Error() == "EOF" {
			conn.Close()
			break
		}

		// processa request
		req = strings.ReplaceAll(req, "\n", "")
		idx := indexOf(req, datesWeather)
		rep := formatDateInfoText(idx)

		// envia resposta
		fmt.Println("Returning to client info about the date: " + req)
		_, err = conn.Write([]byte(rep + "\n"))
		checkError(err)
	}
}

func indexOf(element string, data []DateWeather) int {
	for k, v := range data {
		if element == v.date {
			return k
		}
	}
	return -1
}

func formatDateInfoText(dateIndex int) string {
	if dateIndex == -1 {
		return "No info about this date"
	}

	dateWeather := datesWeather[dateIndex]
	rep := "In " + dateWeather.date +
		" the minimum temperature will be " +
		strconv.Itoa(dateWeather.mininumTemperature) +
		", the maximum temperature will be " +
		strconv.Itoa(dateWeather.maximumTemperature)

	if dateWeather.mayRain {
		rep += " and it will probably rain."
	} else {
		rep += " and it probably won't rain."
	}

	return rep
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(0)
	}
}
