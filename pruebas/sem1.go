package main

import (
	"fmt"
	"sync"
	"time"
)

const numPhilosophers = 5

type Chopstick struct {
	sync.Mutex
}

type Philosopher struct {
	id             int
	leftChopstick  *Chopstick
	rightChopstick *Chopstick
}

func (p *Philosopher) eat() {
	for i := 0; i < 3; i++ {
		p.leftChopstick.Lock()
		p.rightChopstick.Lock()

		fmt.Printf("Philosopher %d is eating.\n", p.id)
		time.Sleep(time.Second)

		p.rightChopstick.Unlock()
		p.leftChopstick.Unlock()

		fmt.Printf("Philosopher %d is thinking.\n", p.id)
		time.Sleep(time.Second)
	}
}

func main() {
	chopsticks := make([]*Chopstick, numPhilosophers)
	for i := 0; i < numPhilosophers; i++ {
		chopsticks[i] = &Chopstick{}
	}

	philosophers := make([]*Philosopher, numPhilosophers)
	for i := 0; i < numPhilosophers; i++ {
		philosophers[i] = &Philosopher{
			id:             i,
			leftChopstick:  chopsticks[i],
			rightChopstick: chopsticks[(i+1)%numPhilosophers],
		}
	}

	var wg sync.WaitGroup
	for _, philosopher := range philosophers {
		wg.Add(1)
		go func(p *Philosopher) {
			defer wg.Done()
			p.eat()
		}(philosopher)
	}

	wg.Wait()
}