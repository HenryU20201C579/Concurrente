package main

import (
	"fmt"
	"math"
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

// Función de activación sigmoide
func sigmoide(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

// Derivada de la función sigmoide
func derivadaSigmoide(x float64) float64 {
	return x * (1.0 - x)
}

// Función para entrenar una red neuronal en una parte del dataset
func entrenarRedNeuronal(dataset []Data, epochs int, learningRate float64, pesos [][]float64) {
	for epoch := 0; epoch < epochs; epoch++ {
		for _, data := range dataset {
			// Forward pass
			capaOculta := make([]float64, len(pesos[0]))
			for i := range capaOculta {
				for j, peso := range pesos[i] {
					capaOculta[i] += data.features[j] * peso
				}
				capaOculta[i] = sigmoide(capaOculta[i])
			}

			// Cálculo de la predicción final
			salida := 0.0
			for i := range capaOculta {
				salida += capaOculta[i] * pesos[len(pesos)-1][i]
			}
			salida = sigmoide(salida)

			// Cálculo del error
			errorSalida := float64(data.label) - salida

			// Backpropagation
			ajusteSalida := errorSalida * derivadaSigmoide(salida)

			// Ajustar pesos de la capa de salida
			for i := range capaOculta {
				pesos[len(pesos)-1][i] += learningRate * ajusteSalida * capaOculta[i]
			}

			// Ajustar pesos de la capa oculta
			for i := range capaOculta {
				ajusteOculta := ajusteSalida * pesos[len(pesos)-1][i] * derivadaSigmoide(capaOculta[i])
				for j := range data.features {
					pesos[i][j] += learningRate * ajusteOculta * data.features[j]
				}
			}
		}
	}
}

func main() {
	// Crear el dataset
	numDatos := 1000000
	fmt.Printf("Creando dataset de %d registros...\n", numDatos)
	dataset := crearDataset(numDatos)
	fmt.Println("Dataset creado con éxito.")

	// Medir el tiempo de ejecución
	start := time.Now()

	// Parámetros de la red neuronal
	numFeatures := 4           // Número de características de entrada
	numNeuronasOcultas := 5    // Número de neuronas en la capa oculta
	epochs := 10               // Número de épocas de entrenamiento
	learningRate := 0.01       // Tasa de aprendizaje

	// Inicializar los pesos de las capas de la red neuronal
	pesos := make([][]float64, numNeuronasOcultas+1) // Capas: oculta + salida
	for i := range pesos {
		if i < numNeuronasOcultas {
			pesos[i] = make([]float64, numFeatures) // Pesos de la capa oculta
		} else {
			pesos[i] = make([]float64, numNeuronasOcultas) // Pesos de la capa de salida
		}
		for j := range pesos[i] {
			pesos[i][j] = rand.Float64() // Inicialización aleatoria de pesos
		}
	}

	// Entrenar la red neuronal secuencialmente
	fmt.Println("Entrenando la red neuronal secuencialmente...")
	entrenarRedNeuronal(dataset, epochs, learningRate, pesos)

	// Calcular el tiempo de ejecución
	elapsed := time.Since(start)
	fmt.Println("Entrenamiento completado.")
	fmt.Printf("Tiempo de ejecución: %s\n", elapsed)
}
