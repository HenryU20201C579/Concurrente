// Ejemplificar condición de carrera
package main

import (
	"fmt"
	"time"
)

func stingy(money *int) {
	for i := 0; i < 1000000; i++ {
		*money += 10
	}
	fmt.Println("Stingy Done!")
}

func spendy(money *int) {
	for i := 0; i < 1000000; i++ {
		*money -= 10
	}
	fmt.Println("Spendy Done!")
}

func main() {
	money := 100
	go stingy(&money)
	go spendy(&money)

	time.Sleep(time.Millisecond * 200)
	fmt.Println("El dinero final en el banco es : ", money)
}
