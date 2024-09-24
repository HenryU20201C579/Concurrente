package main

import (
	"fmt"
	"sync"
)

//Uso del WaitGroup

func main() {
	//Crear el WaitGroup para esperar q todos los goroutines terminen su ejecución
	var wg sync.WaitGroup
	//Creamos el mutex para exclusion
	var mu sync.Mutex

	//compartir un recurso por memoria
	var counter int

	//Número de goroutines
	numGoroutines := 5

	//Agregar el nro de goroutines al WG
	wg.Add(numGoroutines)

	//iniciar multiples goroutines
	for i := 0; i < numGoroutines; i++ {
		//función anónima q se conviente en goroutine
		go func(id int) {
			//garantizar la finalización del goroutine y descontar del WG
			defer wg.Done()

			//bloqueo
			mu.Lock()
			//simular la SC
			fmt.Println(id, "SC")
			counter++
			defer mu.Unlock()

			fmt.Printf("Goroutine %d incrementa el contador a %d\n", id, counter)
		}(i)
	}

	//WG esperar a que todos los goroutines finalicen
	wg.Wait()

	fmt.Printf("El valor final del contador es %d\n", counter)
}
