package main

import (
	"fmt"
	"sync"
	"time"
	"math/rand"
)

type Direction int

const (
	Left Direction = iota // i.e. 0
	Right // 1
)

func drive(wg *sync.WaitGroup, direction Direction) {
	if direction == Right {
		fmt.Println("car moving right joining the bridge")
		time.Sleep(1 * time.Second)
		fmt.Println("car moving right leaving the bridge")
		time.Sleep(1 * time.Second)
	}

	if direction == Left {
		fmt.Println("car moving left joining the bridge")
		time.Sleep(1 * time.Second)
		fmt.Println("car moving left leaving the bridge")
		time.Sleep(1 * time.Second)
	}

	wg.Done()
}

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go drive(&wg, Direction(rand.Intn(2)))
	}

	wg.Wait()

	fmt.Println("Acabou a fila de carros")
}