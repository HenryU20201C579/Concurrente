package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

// Definición de un modelo basado en factores latentes
type FactorizationMachine struct {
	numFeatures   int
	numFactors    int
	learningRate  float64
	regularization float64
	weights       []float64
	factors       [][]float64
	mu            sync.Mutex
}

// Crear dataset simulado con características y etiquetas
func crearDataset(numEjemplos, numCaracteristicas int) ([][]float64, []float64) {
	rand.Seed(time.Now().UnixNano())
	dataset := make([][]float64, numEjemplos)
	etiquetas := make([]float64, numEjemplos)
	for i := 0; i < numEjemplos; i++ {
		ejemplo := make([]float64, numCaracteristicas)
		for j := 0; j < numCaracteristicas; j++ {
			ejemplo[j] = rand.Float64() * 10 // Valor aleatorio entre 0 y 10
		}
		dataset[i] = ejemplo
		etiquetas[i] = rand.Float64() * 10 // Etiqueta aleatoria entre 0 y 10
	}
	return dataset, etiquetas
}

// Dividir el dataset en subconjuntos para entrenamiento y prueba
func dividirDataset(dataset [][]float64, etiquetas []float64, ratio float64) ([][]float64, []float64, [][]float64, []float64) {
	numEjemplos := len(dataset)
	indices := rand.Perm(numEjemplos)
	numEntrenamiento := int(float64(numEjemplos) * ratio)
	
	datasetEntrenamiento := make([][]float64, numEntrenamiento)
	etiquetasEntrenamiento := make([]float64, numEntrenamiento)
	datasetPrueba := make([][]float64, numEjemplos-numEntrenamiento)
	etiquetasPrueba := make([]float64, numEjemplos-numEntrenamiento)
	
	for i := 0; i < numEntrenamiento; i++ {
		datasetEntrenamiento[i] = dataset[indices[i]]
		etiquetasEntrenamiento[i] = etiquetas[indices[i]]
	}
	for i := numEntrenamiento; i < numEjemplos; i++ {
		datasetPrueba[i-numEntrenamiento] = dataset[indices[i]]
		etiquetasPrueba[i-numEntrenamiento] = etiquetas[indices[i]]
	}
	
	return datasetEntrenamiento, etiquetasEntrenamiento, datasetPrueba, etiquetasPrueba
}

// Inicializar el modelo con pesos y factores aleatorios
func (fm *FactorizationMachine) inicializar() {
	fm.weights = make([]float64, fm.numFeatures)
	for i := range fm.weights {
		fm.weights[i] = rand.Float64() * 0.1
	}

	fm.factors = make([][]float64, fm.numFeatures)
	for i := range fm.factors {
		fm.factors[i] = make([]float64, fm.numFactors)
		for j := range fm.factors[i] {
			fm.factors[i][j] = rand.Float64() * 0.1
		}
	}
}

// Predicción usando el modelo basado en factores latentes
func (fm *FactorizationMachine) predecir(entrada []float64) float64 {
	numFeatures := len(entrada)
	
	// Cálculo de la predicción
	prediccion := 0.0
	for i := 0; i < numFeatures; i++ {
		prediccion += fm.weights[i] * entrada[i]
	}

	interaction := 0.0
	for i := 0; i < numFeatures; i++ {
		for j := i + 1; j < numFeatures; j++ {
			interaction += entrada[i] * entrada[j] * (fm.factors[i][0] * fm.factors[j][0])
		}
	}
	prediccion += 0.5 * interaction

	return prediccion
}

// Entrenamiento del modelo usando descenso de gradiente
func (fm *FactorizationMachine) entrenar(dataset [][]float64, etiquetas []float64, epochs int) {
	var wg sync.WaitGroup

	for epoch := 0; epoch < epochs; epoch++ {
		wg.Add(len(dataset))
		for i, entrada := range dataset {
			go func(i int, entrada []float64) {
				defer wg.Done()
				salidaEsperada := etiquetas[i]
				prediccion := fm.predecir(entrada)
				error := salidaEsperada - prediccion

				fm.mu.Lock()
				// Actualización de los pesos
				for j := range fm.weights {
					fm.weights[j] += fm.learningRate * (error * entrada[j] - fm.regularization * fm.weights[j])
				}

				// Actualización de los factores
				for k := 0; k < fm.numFeatures; k++ {
					for l := 0; l < fm.numFactors; l++ {
						fm.factors[k][l] += fm.learningRate * (error * entrada[k] - fm.regularization * fm.factors[k][l])
					}
				}
				fm.mu.Unlock()
			}(i, entrada)
		}
		wg.Wait()
	}
}

// Evaluar la precisión del modelo en el conjunto de prueba
func evaluarPrecision(fm *FactorizationMachine, datasetPrueba [][]float64, etiquetasPrueba []float64) float64 {
	numCorrectos := 0
	var wg sync.WaitGroup
	wg.Add(len(datasetPrueba))

	for i, instancia := range datasetPrueba {
		go func(i int, instancia []float64) {
			defer wg.Done()
			prediccion := fm.predecir(instancia)
			if math.Abs(prediccion-etiquetasPrueba[i]) < 1.0 { // Tolerancia para la precisión
				fm.mu.Lock()
				numCorrectos++
				fm.mu.Unlock()
			}
		}(i, instancia)
	}

	wg.Wait()
	return float64(numCorrectos) / float64(len(etiquetasPrueba))
}

func main() {
	// Crear el dataset
	numEjemplos := 1000
	numCaracteristicas := 10
	fmt.Printf("Creando dataset con %d ejemplos y %d características...\n", numEjemplos, numCaracteristicas)
	dataset, etiquetas := crearDataset(numEjemplos, numCaracteristicas)
	fmt.Println("Dataset creado con éxito.")

	// Dividir el dataset
	ratioEntrenamiento := 0.8
	datasetEntrenamiento, etiquetasEntrenamiento, datasetPrueba, etiquetasPrueba := dividirDataset(dataset, etiquetas, ratioEntrenamiento)
	fmt.Printf("Dataset dividido en %d ejemplos de entrenamiento y %d ejemplos de prueba.\n",
		len(datasetEntrenamiento), len(datasetPrueba))

	// Crear y entrenar el modelo
	fm := FactorizationMachine{
		numFeatures:   numCaracteristicas,
		numFactors:    5,
		learningRate:  0.01,
		regularization: 0.01,
	}
	fm.inicializar()

	start := time.Now()
	fm.entrenar(datasetEntrenamiento, etiquetasEntrenamiento, 10)
	fmt.Println("Entrenamiento completado.")

	// Evaluar el modelo
	precision := evaluarPrecision(&fm, datasetPrueba, etiquetasPrueba)
	fmt.Printf("Precisión del modelo: %.2f%%\n", precision*100)

	// Calcular el tiempo de ejecución
	elapsed := time.Since(start)
	fmt.Printf("Tiempo de ejecución: %s\n", elapsed)
}
