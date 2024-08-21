package main

import "fmt"

var Exportada string //variable pública

func main() {
	var noExportada string //variable privada
	noExportada = "si"
	saludo := "Hola mundo"

	fmt.Println(saludo)
	fmt.Println("No esportada ", noExportada)

}
