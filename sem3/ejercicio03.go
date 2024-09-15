// Primer Intento
package main

import (
	"fmt"
	"time"
)

var turno int = 1

func p() {
	for { //bucle infinito
		fmt.Println("Line01 SNC - P")
		for turno != 1 {
			//esperar el proceso P
		}
		fmt.Println("Line01 SC - P")
		turno = 2
	}
}

func q() {
	for { //bucle infinito
		fmt.Println("Line01 SNC - Q")
		for turno != 2 {
			//espera el proceso Q
		}
		fmt.Println("Line01 SC - Q")
		turno = 1
	}
}

func main() {
	go p()
	go q()

	time.Sleep(time.Minute * 120)
}
