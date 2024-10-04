// semáforo contador/condicional :- un recurso compartido y un limite de procesos que acceden a esos recursos
package main

import (
	"fmt"
	"sync"
	"time"
)

type semaforo struct {
	contador int
	limite   int
	mu       sync.Mutex
}

func (s *semaforo) wait() {
	s.mu.Lock()
	defer s.mu.Unlock()

	//condición / límite
	for s.contador >= s.limite {
		//espera del proceso activo
		s.mu.Unlock()
		//tiempo por hacer algo
		time.Sleep(time.Millisecond * 50)
		s.mu.Lock()
	}
	s.contador++
}

func (s *semaforo) signal() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.contador--

}

func proceso(id int, s *semaforo, wg *sync.WaitGroup) {
	defer wg.Done()
	s.wait()
	//SC
	fmt.Printf("El proceso %d, adquirió el semáforo\n", id)
	//simular el trabajo del proceso con un time
	time.Sleep(time.Millisecond * 50)
	fmt.Printf("El proceso %d, liberó el semáforo\n", id)
	s.signal()
}

func main() {
	var wg sync.WaitGroup
	nroprocesos := 10
	limite := 3
	sem := &semaforo{limite: limite}

	//lanzar los procesos
	for i := 0; i < nroprocesos; i++ {
		wg.Add(1)
		go proceso(i, sem, &wg)
	}

	wg.Wait()
	fmt.Printf("Todos los procesos culminaron!!!")
}
