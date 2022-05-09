package main

import (
	"fmt"
	"net"
	"net/rpc"
	"os"
	"sync"
)

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

func (s *Semaphore) P(i int, reply *string) error {
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
	*reply = "com sucesso!"
	return nil
}

func (s *Semaphore) V(i int, reply *string) error {
	s.c.L.Lock()
	s.n++
	fmt.Printf("Cliente %d levantou | %d clientes sentados\n", i, 5-s.n)
	if s.n == 5 {
		fmt.Println("MESA VAZIA!")
		s.must_wait = false
		s.c.Broadcast()
	}
	s.c.L.Unlock()
	*reply = "com sucesso!!"
	return nil
}

/* --- Servidor --- */
func servidor() {
	s := NewSemaphore(5)

	server := rpc.NewServer()
	err := server.RegisterName("SushiBar", s)
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
