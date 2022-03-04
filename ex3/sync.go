package main

import (
	"fmt"
	"sync"
	"time"
	"math/rand"
)

type Direction int

const (
	Left Direction = iota
	Right
)

func drive(id int, wg *sync.WaitGroup, m *sync.Mutex, direction Direction) {
	m.Lock()

	if direction == Right {
		fmt.Println("car moving right joining the bridge")
		time.Sleep(100 * time.Millisecond)
		fmt.Println("car moving right leaving the bridge")
		time.Sleep(100 * time.Millisecond)
	}

	if direction == Left {
		fmt.Println("car moving left joining the bridge")
		time.Sleep(100 * time.Millisecond)
		fmt.Println("car moving left leaving the bridge")
		time.Sleep(100 * time.Millisecond)
	}

	m.Unlock()
	wg.Done()
}

func main() {
	var wg sync.WaitGroup
	var m sync.Mutex

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go drive(i, &wg, &m, Direction(rand.Intn(2)))
	}

	wg.Wait()

	fmt.Println("Acabou a fila de carros")
}