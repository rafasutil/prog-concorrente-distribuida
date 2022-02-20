package main

import (
	"fmt"
	"sync"
	"time"
)

var M = 10

func main() {
	var wg sync.WaitGroup
	m := sync.Mutex{}
	cond := sync.NewCond(&m)

	wg.Add(1)
	go pessoa("Rafael", cond)
	go pessoa("Daniel", cond)
	go pessoa("Ayrton", cond)
	go cozinheiro(cond)
	wg.Wait()

}

func pessoa(nome string, cond *sync.Cond) {
	for true {
		// Pegar Porção
		cond.L.Lock()
		fmt.Printf("%s vai comer M=%d\n", nome, M)
		func() {
			for M == 0 {
				fmt.Printf("%s aguardando encher M=%d\n", nome, M)
				cond.Wait()
			}
			M--
			fmt.Printf("%s comeu M=%d\n", nome, M)
		}()
		cond.L.Unlock()
		// Comer
		time.Sleep(time.Second * 5)
	}
}

func cozinheiro(cond *sync.Cond) {
	for true {
		time.Sleep(time.Second * 10)
		// Encher Panela
		cond.L.Lock()
		fmt.Printf("Cozinheiro verifica M=%d\n", M)
		func() {
			if M == 0 {
				fmt.Printf("Cozinheiro está enchendo...\n")
				time.Sleep(time.Second * 3)
				M = 10
				cond.Broadcast()
			}
		}()
		cond.L.Unlock()
	}
}
