package main

import (
	"fmt"
	"impl"
)

var (
	n int
	s impl.Semaphore
)

func customer(i int) {
	for {
		// Sentar
		s.P(i)
		//time.Sleep(time.Second * 15)

		// Levantar
		s.V(i)
		//time.Sleep(time.Second * 30)
	}
}

func main() {
	s = *impl.NewSemaphore(5)

	for i := 1; i <= 8; i++ {
		go customer(i)
	}

	fmt.Scanln()
}
