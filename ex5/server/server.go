package main

import (
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strconv"
)

/* --- Implementação --- */
type DateWeather struct {
	date               string // DD/MM/YYYY
	mininumTemperature int
	maximumTemperature int
	mayRain            bool
}

var datesWeather = []DateWeather{
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

func (t *DateWeather) Result(req string, reply *string) error {
	idx := indexOf(req, datesWeather)
	*reply = formatDateInfoText(idx)
	return nil
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

/* --- Servidor --- */
func servidor() {

	dateweather := new(DateWeather)

	server := rpc.NewServer()
	err := server.RegisterName("Weather", dateweather)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	ln, err := net.Listen("tcp", "localhost:1313")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	defer func(ln net.Listener) {
		var err = ln.Close()
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}(ln)

	fmt.Println("Servidor está pronto...")

	server.Accept(ln)

}

func main() {
	go servidor()
	_, _ = fmt.Scanln()
}
