package main

import (
	"fmt"
	"strconv"
	"sync"
)

var ponte *Ponte

type Carro struct {
	Nome string
}

type Ponte struct {
	cond       *sync.Cond
	buffer     []Carro
	capacidade int
}

func NovaPonte() *Ponte {
	return &Ponte{
		cond:       sync.NewCond(new(sync.Mutex)),
		capacidade: 1,
	}
}

func (p *Ponte) Entrar(e Carro) {
	p.cond.L.Lock()
	for len(p.buffer) == p.capacidade {
		p.cond.Wait()
	}

	p.buffer = append(p.buffer, e)
	fmt.Println(e.Nome, " entrou na ponte!", len(p.buffer), "carros na ponte!")

	p.cond.Broadcast()
	p.cond.L.Unlock()
}

func (p *Ponte) Sair() {
	p.cond.L.Lock()
	for len(p.buffer) == 0 {
		p.cond.Wait()
	}

	e := p.buffer[0]
	p.buffer = p.buffer[1:]
	fmt.Println(e.Nome, " saiu da ponte!", len(p.buffer), "carros na ponte!")

	p.cond.Broadcast()
	p.cond.L.Unlock()
}

func main() {
	wg := sync.WaitGroup{}
	n := 100

	ponte = NovaPonte()
	/* Carro saindo da ponte independente do lado */
	for a := 0; a < n; a++ {
		wg.Add(1)
		go consumer(&wg)
	}
	/* Simular carro indo da esquerda para a direita */
	for b := 1; b <= n/2; b++ {
		wg.Add(1)
		go producer(b, &wg)
	}
	/* Simular carro indo da direita para a esquerda */
	for c := 1 + n/2; c <= n; c++ {
		wg.Add(1)
		go producer(c, &wg)
	}

	wg.Wait()
}

func consumer(wg *sync.WaitGroup) {
	defer wg.Done()
	ponte.Sair()
}
func producer(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	e := Carro{Nome: "Carro " + strconv.Itoa(id)}
	ponte.Entrar(e)
}
