package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"strconv"
)

type DateWeather struct {
	date string // DD/MM/YYYY
	mininumTemperature int
	maximumTemperature int
	mayRain bool
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
	}

	HandleServerUDP()

	_, _ = fmt.Scanln()
}

func HandleServerUDP() {
	req := make([]byte, 10)
	rep := make([]byte, 1024)

	// define o endpoint do servidor UDP
	addr, err := net.ResolveUDPAddr("udp", "localhost:1313")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	// prepara o endpoint UDP para receber requests
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	// fecha conn
	defer func(conn *net.UDPConn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}(conn)

	fmt.Println("Servidor UDP aguardando requests...")

	for {
		// recebe request
		_, addr, err := conn.ReadFromUDP(req)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		// processa request
		stringReq := string(req)
		stringReq = strings.ReplaceAll(stringReq, "\n", "")
		idx := indexOf(stringReq, datesWeather)
		rep = []byte(formatDateInfoText(idx))

		// envia reposta
		fmt.Println("Returning to client info about the date: " + string(req))
		_, err = conn.WriteTo(rep, addr)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}
}

func indexOf(element string, data []DateWeather) (int) {
	for k, v := range data {
		if element == v.date {
			return k
		}
	}
	return -1 //not found.
}

func formatDateInfoText(dateIndex int ) (string) {
	if(dateIndex == -1){
		return "No info about this date"
	}

	dateWeather := datesWeather[dateIndex]
	rep:= "In " + dateWeather.date +
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