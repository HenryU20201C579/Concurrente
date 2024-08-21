// Ordenamiento de números
package main

import (
	"fmt"
	"sort"
)

func ordenarNumeros(numeros []int) []int {
	sort.Ints(numeros)
	return numeros
}

func main() {
	listaNumeros := []int{4, 2, 8, 1, 6}
	fmt.Println("Números ordenados:", ordenarNumeros(listaNumeros))
}
