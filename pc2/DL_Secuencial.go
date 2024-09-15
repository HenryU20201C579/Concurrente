package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// Definición de la red neuronal profunda
type RedNeuronal struct {
	Entradas   int
	Ocultas    int
	Salidas    int
	PesosOcultos [][]float64
	PesosSalidas [][]float64
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

// Inicializar los pesos de la red neuronal con valores aleatorios
func (red *RedNeuronal) inicializarPesos() {
	red.PesosOcultos = make([][]float64, red.Ocultas)
	for i := range red.PesosOcultos {
		red.PesosOcultos[i] = make([]float64, red.Entradas)
		for j := range red.PesosOcultos[i] {
			red.PesosOcultos[i][j] = rand.Float64() * 0.1
		}
	}

	red.PesosSalidas = make([][]float64, red.Salidas)
	for i := range red.PesosSalidas {
		red.PesosSalidas[i] = make([]float64, red.Ocultas)
		for j := range red.PesosSalidas[i] {
			red.PesosSalidas[i][j] = rand.Float64() * 0.1
		}
	}
}

// Función de activación Sigmoid
func sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

// Derivada de la función de activación Sigmoid
func derivadaSigmoid(x float64) float64 {
	sig := sigmoid(x)
	return sig * (1 - sig)
}

// Propagación hacia adelante
func (red *RedNeuronal) propagacionHaciaAdelante(entrada []float64) ([]float64, []float64) {
	ocultas := make([]float64, red.Ocultas)
	salidas := make([]float64, red.Salidas)
	
	for i := 0; i < red.Ocultas; i++ {
		sum := 0.0
		for j := 0; j < red.Entradas; j++ {
			sum += entrada[j] * red.PesosOcultos[i][j]
		}
		ocultas[i] = sigmoid(sum)
	}
	
	for i := 0; i < red.Salidas; i++ {
		sum := 0.0
		for j := 0; j < red.Ocultas; j++ {
			sum += ocultas[j] * red.PesosSalidas[i][j]
		}
		salidas[i] = sigmoid(sum)
	}
	
	return salidas, ocultas
}

// Retropropagación para ajustar los pesos
func (red *RedNeuronal) retropropagacion(entrada []float64, ocultas []float64, salidaEsperada float64, salida []float64, tasaAprendizaje float64) {
	error := salidaEsperada - salida[0]
	deltaSalida := error * derivadaSigmoid(salida[0])

	// Ajuste de pesos de salida
	for i := 0; i < red.Salidas; i++ {
		for j := 0; j < red.Ocultas; j++ {
			red.PesosSalidas[i][j] += tasaAprendizaje * deltaSalida * ocultas[j]
		}
	}

	// Cálculo del error en la capa oculta
	deltaOculta := make([]float64, red.Ocultas)
	for i := 0; i < red.Ocultas; i++ {
		sum := 0.0
		for j := 0; j < red.Salidas; j++ {
			sum += deltaSalida * red.PesosSalidas[j][i]
		}
		deltaOculta[i] = sum * derivadaSigmoid(ocultas[i])
	}

	// Ajuste de pesos de entrada a capa oculta
	for i := 0; i < red.Ocultas; i++ {
		for j := 0; j < red.Entradas; j++ {
			red.PesosOcultos[i][j] += tasaAprendizaje * deltaOculta[i] * entrada[j]
		}
	}
}

// Entrenamiento de la red neuronal
func (red *RedNeuronal) entrenar(dataset [][]float64, etiquetas []int, epochs int, tasaAprendizaje float64) {
	for epoch := 0; epoch < epochs; epoch++ {
		for i, entrada := range dataset {
			salidaEsperada := float64(etiquetas[i])
			salida, ocultas := red.propagacionHaciaAdelante(entrada)
			red.retropropagacion(entrada, ocultas, salidaEsperada, salida, tasaAprendizaje)
		}
	}
}

// Clasificar una instancia usando la red neuronal
func (red *RedNeuronal) clasificar(instancia []float64) int {
	salida, _ := red.propagacionHaciaAdelante(instancia)
	if salida[0] > 0.5 {
		return 1
	}
	return 0
}

// Evaluar la precisión del modelo en el conjunto de prueba
func evaluarPrecision(red *RedNeuronal, datasetPrueba [][]float64, etiquetasPrueba []int) float64 {
	correctos := 0
	for i, instancia := range datasetPrueba {
		prediccion := red.clasificar(instancia)
		if prediccion == etiquetasPrueba[i] {
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

	// Inicializar la red neuronal
	entradas := numCaracteristicas
	ocultas := 20
	salidas := 1
	red := RedNeuronal{Entradas: entradas, Ocultas: ocultas, Salidas: salidas}
	red.inicializarPesos()
	fmt.Println("Red neuronal inicializada.")

	// Medir el tiempo de ejecución
	start := time.Now()

	// Entrenar la red neuronal
	epochs := 10
	tasaAprendizaje := 0.01
	red.entrenar(datasetEntrenamiento, etiquetasEntrenamiento, epochs, tasaAprendizaje)
	fmt.Println("Entrenamiento completado.")

	// Evaluar el modelo
	precision := evaluarPrecision(&red, datasetPrueba, etiquetasPrueba)
	fmt.Printf("Precisión del modelo: %.2f%%\n", precision*100)

	// Calcular el tiempo de ejecución
	elapsed := time.Since(start)
	fmt.Printf("Tiempo de ejecución: %s\n", elapsed)
}
