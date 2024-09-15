// Características de CPU y procesos
package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println("Número de CPUs: ", runtime.NumCPU())
	fmt.Println("GOMAXPROCS: ", runtime.GOMAXPROCS(0))
}
