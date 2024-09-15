package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Definir estructura de los datos
type Data struct {
	features []float64
	label    int
}

// Crear dataset simulado de 1 millón de registros
func crearDataset(n int) []Data {
	rand.Seed(time.Now().UnixNano())
	dataset := make([]Data, n)
	for i := 0; i < n; i++ {
		features := []float64{
			rand.Float64(),
			rand.Float64(),
			rand.Float64(),
			rand.Float64(),
		}
		label := rand.Intn(2) // Etiqueta aleatoria (0 o 1)
		dataset[i] = Data{features, label}
	}
	return dataset
}

// Función para entrenar el modelo SVM de forma secuencial
func entrenarSVM(dataset []Data, epochs int, learningRate float64) []float64 {
	pesos := make([]float64, len(dataset[0].features)) // Inicialización de pesos
	for epoch := 0; epoch < epochs; epoch++ {
		for _, data := range dataset {
			prediccion := predecir(data.features, pesos)
			error := float64(data.label) - prediccion
			for i := 0; i < len(pesos); i++ {
				pesos[i] += learningRate * error * data.features[i]
			}
		}
	}
	return pesos
}

// Función para predecir con los pesos actuales
func predecir(features []float64, pesos []float64) float64 {
	var suma float64
	for i := 0; i < len(features); i++ {
		suma += features[i] * pesos[i]
	}
	if suma >= 0 {
		return 1.0
	}
	return 0.0
}

func main() {
	// Crear el dataset
	numDatos := 1000000
	fmt.Printf("Creando dataset de %d registros...\n", numDatos)
	dataset := crearDataset(numDatos)
	fmt.Println("Dataset creado con éxito.")

	// Medir el tiempo de ejecución
	start := time.Now()

	// Entrenar el modelo SVM secuencialmente
	epochs := 10
	learningRate := 0.01
	fmt.Println("Entrenando el modelo SVM secuencialmente...")
	pesosFinales := entrenarSVM(dataset, epochs, learningRate)

	// Mostrar resultados
	fmt.Println("Entrenamiento completado.")
	fmt.Printf("Pesos finales: %v\n", pesosFinales)

	// Calcular tiempo de ejecución
	elapsed := time.Since(start)
	fmt.Printf("Tiempo de ejecución: %s\n", elapsed)
}
