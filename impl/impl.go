package sema

import (
	"fmt"
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

func (s *Semaphore) P(i int) {
	s.c.L.Lock()
	if s.n == 0 {
		fmt.Printf("Cliente %d deve esperar os 5 levantarem\n", i)
		s.must_wait = true
	}
	for s.n <= 0 || s.must_wait {
		s.c.Wait()
	}
	s.n--
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
	}
	s.c.L.Unlock()
	s.c.Signal()
}
