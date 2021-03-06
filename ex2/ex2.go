package main

import (
	"fmt"
	"sync"
)

/*
	Semaphore Implementation
*/

type ISemaphore interface {
	P()
	V()
}

type Semaphore struct {
	n         int32
	c         *sync.Cond
	must_wait bool
}

func NewSemaphore(N int32) *Semaphore {
	return &Semaphore{
		n:         N,
		c:         sync.NewCond(new(sync.Mutex)),
		must_wait: false,
	}
}

func (s *Semaphore) P(i int) {
	s.c.L.Lock()
	for s.n <= 0 || s.must_wait {
		fmt.Printf("Cliente %d deve esperar os 5 levantarem\n", i)
		s.c.Wait()
	}
	s.n--
	if s.n == 0 {
		s.must_wait = true
	}
	fmt.Printf("Cliente %d sentou | %d clientes sentados\n", i, 5-s.n)
	s.c.L.Unlock()
}

func (s *Semaphore) V(i int) {
	s.c.L.Lock()
	s.n++
	fmt.Printf("Cliente %d levantou | %d clientes sentados\n", i, 5-s.n)
	if s.n == 5 {
		fmt.Println("MESA VAZIA!")
		s.must_wait = false
		s.c.Broadcast()
	}
	s.c.L.Unlock()
}

/*
	Main Program
*/

var (
	n int
	s Semaphore
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
	s = *NewSemaphore(5)

	for i := 1; i <= 8; i++ {
		go customer(i)
	}

	fmt.Scanln()
}
