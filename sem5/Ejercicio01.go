package main

//Estructura semáforo usando Sync.Mutex
import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var m sync.Mutex
	var in string
	//un conjunto de procesos concurrentes
	//q se van a sincronizar
	for i := 0; i < 10; i++ {
		//goroutines
		go func(j int) {
			m.Lock() //wait
			fmt.Println(j, " SC line01: Inicio")
			time.Sleep(time.Millisecond * 20)
			fmt.Println(j, " SC line02: Fin")
			m.Unlock() //signal
		}(i)
	}
	fmt.Scanln(&in)

}
