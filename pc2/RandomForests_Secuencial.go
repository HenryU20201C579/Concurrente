package main

import (
	"fmt"
	//"math"
	"math/rand"
	"time"
)

// Estructura para un árbol de decisión
type Nodo struct {
	Caracteristica int
	Valor          float64
	Izquierda      *Nodo
	Derecha        *Nodo
	Clasificacion  int
}

// Estructura para un árbol de decisión
type Arbol struct {
	Raiz *Nodo
}

// Crear dataset simulado de 1000 ejemplos con 10 características
func crearDataset(numEjemplos, numCaracteristicas int) ([][]float64, []int) {
	rand.Seed(time.Now().UnixNano())
	dataset := make([][]float64, numEjemplos)
	etiquetas := make([]int, numEjemplos)
	for i := 0; i < numEjemplos; i++ {
		ejemplo := make([]float64, numCaracteristicas)
		for j := 0; j < numCaracteristicas; j++ {
			ejemplo[j] = rand.Float64() * 10 // Valor aleatorio entre 0 y 10
		}
		dataset[i] = ejemplo
		etiquetas[i] = rand.Intn(2) // Etiqueta aleatoria: 0 o 1
	}
	return dataset, etiquetas
}

// Dividir el dataset en subconjuntos para entrenamiento y prueba
func dividirDataset(dataset [][]float64, etiquetas []int, ratio float64) ([][]float64, []int, [][]float64, []int) {
	numEjemplos := len(dataset)
	indices := rand.Perm(numEjemplos)
	numEntrenamiento := int(float64(numEjemplos) * ratio)
	
	datasetEntrenamiento := make([][]float64, numEntrenamiento)
	etiquetasEntrenamiento := make([]int, numEntrenamiento)
	datasetPrueba := make([][]float64, numEjemplos-numEntrenamiento)
	etiquetasPrueba := make([]int, numEjemplos-numEntrenamiento)
	
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

// Crear un árbol de decisión simple (placeholder)
func crearArbolDecision(dataset [][]float64, etiquetas []int) *Arbol {
	// Implementar la lógica real de creación del árbol aquí
	return &Arbol{}
}

// Clasificar una instancia usando el árbol de decisión
func (arbol *Arbol) clasificar(instancia []float64) int {
	// Implementar la lógica real de clasificación aquí
	return rand.Intn(2)
}

// Evaluar la precisión del modelo en el conjunto de prueba
func evaluarPrecision(arbol *Arbol, datasetPrueba [][]float64, etiquetasPrueba []int) float64 {
	correctos := 0
	for i, instancia := range datasetPrueba {
		if arbol.clasificar(instancia) == etiquetasPrueba[i] {
			correctos++
		}
	}
	return float64(correctos) / float64(len(etiquetasPrueba))
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
	fmt.Printf("Dataset dividido en entrenamiento (%.2f%%) y prueba (%.2f%%).\n", ratioEntrenamiento*100, (1-ratioEntrenamiento)*100)

	// Medir el tiempo de ejecución
	start := time.Now()

	// Entrenar el modelo
	arbol := crearArbolDecision(datasetEntrenamiento, etiquetasEntrenamiento)
	fmt.Println("Entrenamiento completado.")

	// Evaluar el modelo
	precision := evaluarPrecision(arbol, datasetPrueba, etiquetasPrueba)
	fmt.Printf("Precisión del modelo: %.2f%%\n", precision*100)

	// Calcular el tiempo de ejecución
	elapsed := time.Since(start)
	fmt.Printf("Tiempo de ejecución: %s\n", elapsed)
}
