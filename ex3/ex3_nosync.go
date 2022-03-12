package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var ponte *Ponte

type Carro struct {
	Nome string
}

type Ponte struct {
	buffer     []Carro
	capacidade int
}

func NovaPonte() *Ponte {
	return &Ponte{
		capacidade: 1,
	}
}

func (p *Ponte) Entrar(e Carro) {
	for len(p.buffer) == p.capacidade {
		//fmt.Println("Adicionando mais de 1 carro na ponte.")
		time.Sleep(time.Second * 3)
	}

	p.buffer = append(p.buffer, e)
	fmt.Println(e.Nome, " entrou na ponte!", len(p.buffer), "carros na ponte!")
}

func (p *Ponte) Sair() {
	for len(p.buffer) == 0 {
		//fmt.Println("Está tentando sair da da ponte sem entrar.")
		time.Sleep(time.Second * 3)
	}

	e := p.buffer[0]
	p.buffer = p.buffer[1:]
	fmt.Println(e.Nome, " saiu da ponte!", len(p.buffer), "carros na ponte!")
}

func main() {
	wg := sync.WaitGroup{}
	n := 20

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
