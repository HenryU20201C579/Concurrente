// Cena de los filósofos
package main

import (
	"fmt"
	"sync"
)

// cada proceso
func filosofo(id int, tenedorIzq, tenedorDer sync.Mutex) {
	for {
		fmt.Printf("Filósofo %d, está pensando\n", id)
		tenedorIzq.Lock()
		tenedorDer.Lock()
		fmt.Printf("Folósofo %d, está comiendo\n", id)
		tenedorIzq.Unlock()
		tenedorDer.Unlock()
	}
}

func main() {
	//crear el arreglo de tenedores que accede c/filósofo
	tenedores := make([]sync.Mutex, 5)
	//Iniciar los 4 procesos
	go filosofo(1, tenedores[0], tenedores[1])
	go filosofo(2, tenedores[1], tenedores[2])
	go filosofo(3, tenedores[2], tenedores[3])
	go filosofo(4, tenedores[3], tenedores[4])

	filosofo(5, tenedores[4], tenedores[0])
}
