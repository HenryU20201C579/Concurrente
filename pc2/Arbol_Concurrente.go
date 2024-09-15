package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const numRecords = 1000000 // Número de registros

// Datos simulados
type Record struct {
	feature1 float64
	feature2 float64
	label    int
}

// Entrenar Decision Tree (simulación)
func entrenarDecisionTree(records []Record, id int, wg *sync.WaitGroup) {
	defer wg.Done() // Asegura que Done se llame cuando la goroutine termine

	start := time.Now()
	fmt.Printf("Entrenando modelo %d...\n", id)

	// Simulación de procesamiento
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

	elapsed := time.Since(start)
	fmt.Printf("Modelo %d entrenado en %s\n", id, elapsed)
}

// Función principal
func main() {
	// Crear dataset simulado
	fmt.Println("Creando dataset...")

	dataset := make([]Record, numRecords)
	for i := 0; i < numRecords; i++ {
		dataset[i] = Record{
			feature1: rand.Float64(),
			feature2: rand.Float64(),
			label:    rand.Intn(2),
		}
	}

	fmt.Println("Dataset creado con éxito!")

	// Dividir el dataset en partes
	partes := 4
	tamañoParte := len(dataset) / partes

	var wg sync.WaitGroup
	start := time.Now()

	// Entrenar los modelos en paralelo
	for i := 0; i < partes; i++ {
		inicio := i * tamañoParte
		fin := inicio + tamañoParte
		if i == partes-1 { // Asegurarse de incluir el resto de los registros en la última parte
			fin = len(dataset)
		}

		wg.Add(1)
		go entrenarDecisionTree(dataset[inicio:fin], i+1, &wg)
	}

	wg.Wait() // Esperar a que todas las goroutines terminen
	elapsed := time.Since(start)
	fmt.Printf("Entrenamiento completado en %s\n", elapsed)
}
