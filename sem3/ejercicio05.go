// Sincronización con Estructuras de canales
package main

import "fmt"

var ch chan bool //declaración de canal tipo bool

func main() {
	ch = make(chan bool) //construcción del canal tipo bool

	//ch <- true

	go func() { //función anónima
		fmt.Println("Enviando datos!!")
		ch <- true
	}()

	<-ch

	fmt.Println("Datos de main!")
}
