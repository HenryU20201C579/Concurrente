package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func productor(id int, buffer chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 5; i++ {
		//simular la producción de datos
		dato := rand.Intn(100)
		fmt.Printf("Productor %d produce %d\n", id, dato)
		//enviar el dato al buffer
		buffer <- dato
		//retardo simulando el trabajo
		time.Sleep(time.Microsecond * 20)
	}
}

func consumidor(id int, buffer chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	//var d int
	for i := 0; i < 5; i++ {
		//simular el procesamiento de los datos del buffer
		//d=<-buffer
		fmt.Printf("Consumidor %d consume %d\n", id, <-buffer)
		//retardo simulando el trabajo
		time.Sleep(time.Microsecond * 20)
	}
}

func main() {
	var wg sync.WaitGroup
	//Buffer finito con capacidad de 5
	bufferSize := 5
	chDatos := make(chan int, bufferSize)

	//definir el nro de consumidores y productores
	nroProductores := 3
	nroConsumidores := 2

	//iniciar los procesos concurrentes
	//lanzar los productores
	for i := 0; i < nroProductores; i++ {
		wg.Add(1) //agregar el goroutine al grupo
		go productor(i, chDatos, &wg)
	}
	//lanzar los consumidores
	for i := 0; i < nroConsumidores; i++ {
		wg.Add(1)
		go consumidor(i, chDatos, &wg)
	}

	wg.Wait() //espera que los goroutines culminen
	fmt.Println("Todos los productores y consumidores culminaron su trabajo!!!!")
}
