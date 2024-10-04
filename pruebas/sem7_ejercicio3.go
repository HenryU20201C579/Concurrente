package main

import (
	"fmt"
	"sync"
	"time"
)

type Group struct {
	name   string
	size   int
	mutex  *sync.Mutex
	queue  chan *Person
}

type Person struct {
	name   string
	group  *Group
}

func (g *Group) enter() {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	fmt.Printf("%s entra al baño.\n", g.name)
	time.Sleep(time.Second * 5)

	g.queue <- &Person{name: g.name}
}

func (g *Group) leave() {
	person := <-g.queue
	fmt.Printf("%s sale del baño.\n", person.name)
	time.Sleep(time.Second * 5)
}

func main() {
	maleGroup := &Group{
		name:  "Grupo de Varones",
		size:  5,
		mutex: &sync.Mutex{},
		queue: make(chan *Person, 5),
	}

	femaleGroup := &Group{
		name:  "Grupo de Mujeres",
		size:  5,
		mutex: &sync.Mutex{},
		queue: make(chan *Person, 5),
	}

	childrenGroup := &Group{
		name:  "Grupo de Niños",
		size:  5,
		mutex: &sync.Mutex{},
		queue: make(chan *Person, 5),
	}

	var wg sync.WaitGroup
	wg.Add(15)

	for i := 0; i < 5; i++ {
		go func(g *Group) {
			defer wg.Done()
			g.enter()
		}(maleGroup)

		go func(g *Group) {
			defer wg.Done()
			g.leave()
		}(maleGroup)
	}

	for i := 0; i < 5; i++ {
		go func(g *Group) {
			defer wg.Done()
			g.enter()
		}(femaleGroup)

		go func(g *Group) {
			defer wg.Done()
			g.leave()
		}(femaleGroup)
	}

	for i := 0; i < 5; i++ {
		go func(g *Group) {
			defer wg.Done()
			g.enter()
		}(childrenGroup)

		go func(g *Group) {
			defer wg.Done()
			g.leave()
		}(childrenGroup)
	}

	wg.Wait()
}