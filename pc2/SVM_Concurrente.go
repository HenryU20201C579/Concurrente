package main

import (
	"fmt"
	"math/rand"
	"sync"
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

// Función para entrenar el modelo SVM en una parte del dataset
func entrenarParteSVM(dataset []Data, epochs int, learningRate float64, pesos []float64, wg *sync.WaitGroup, m *sync.Mutex) {
	defer wg.Done()

	for epoch := 0; epoch < epochs; epoch++ {
		for _, data := range dataset {
			prediccion := predecir(data.features, pesos)
			error := float64(data.label) - prediccion
			m.Lock()
			for i := 0; i < len(pesos); i++ {
				pesos[i] += learningRate * error * data.features[i]
			}
			m.Unlock()
		}
	}
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

	// Dividir el dataset en partes para entrenamiento concurrente
	numGoroutines := 4
	partSize := len(dataset) / numGoroutines
	epochs := 10
	learningRate := 0.01
	pesos := make([]float64, len(dataset[0].features))

	// Crear un grupo de espera y un mutex para sincronización
	var wg sync.WaitGroup
	var m sync.Mutex

	fmt.Println("Entrenando el modelo SVM concurrentemente...")

	// Lanzar goroutines para cada parte del dataset
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		startIndex := i * partSize
		endIndex := startIndex + partSize
		if i == numGoroutines-1 {
			endIndex = len(dataset) // Último trozo incluye los datos restantes
		}
		go entrenarParteSVM(dataset[startIndex:endIndex], epochs, learningRate, pesos, &wg, &m)
	}

	// Esperar a que todas las goroutines terminen
	wg.Wait()

	// Mostrar resultados
	fmt.Println("Entrenamiento completado.")
	fmt.Printf("Pesos finales: %v\n", pesos)

	// Calcular tiempo de ejecución
	elapsed := time.Since(start)
	fmt.Printf("Tiempo de ejecución: %s\n", elapsed)
}
