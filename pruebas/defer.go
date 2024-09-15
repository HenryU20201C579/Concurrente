package main

import (
	"fmt"
)

func primero(){
	fmt.Println("primero")
}

func segundo(){
    fmt.Println("segundo")
}

func main() {
	defer fmt.Println("hola")
	defer fmt.Println("hola2")
	defer primero()
	segundo()
}