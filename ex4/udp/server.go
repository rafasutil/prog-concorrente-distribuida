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
		{"19/03/2022", 21, 34, false},
		{"20/03/2022", 22, 32, true},
		{"21/03/2022", 24, 30, true},
	}

	HandleServerUDP()

	_, _ = fmt.Scanln()
}

func HandleServerUDP() {
	req := make([]byte, 10)
	rep := make([]byte, 1024)

	addr, err := net.ResolveUDPAddr("udp", "localhost:1313")
	checkError(err)

	conn, err := net.ListenUDP("udp", addr)
	checkError(err)


	defer func(conn *net.UDPConn) {
		err := conn.Close()
		checkError(err)
	}(conn)

	fmt.Println("Servidor UDP aguardando requests...")

	for {
		_, addr, err := conn.ReadFromUDP(req)
		checkError(err)

		stringReq := string(req)
		stringReq = strings.ReplaceAll(stringReq, "\n", "")
		idx := indexOf(stringReq, datesWeather)
		rep = []byte(formatDateInfoText(idx))

		fmt.Println("Returning to client info about the date: " + string(req))
		_, err = conn.WriteTo(rep, addr)

		checkError(err)
	}
}

func indexOf(element string, data []DateWeather) (int) {
	for k, v := range data {
		if element == v.date {
			return k
		}
	}
	return -1
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

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(0)
	}
}