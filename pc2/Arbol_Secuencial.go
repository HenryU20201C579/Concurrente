package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Registro simulado con características y etiquetas
type Registro struct {
	feature1 float64
	feature2 float64
	label    int
}

// Función para crear un dataset simulado de 1 millón de registros
func crearDataset(n int) []Registro {
	rand.Seed(time.Now().UnixNano())
	dataset := make([]Registro, n)
	for i := 0; i < n; i++ {
		dataset[i] = Registro{
			feature1: rand.Float64(),
			feature2: rand.Float64(),
			label:    rand.Intn(2), // Etiquetas binarias 0 o 1
		}
	}
	return dataset
}

// Función para entrenar un Árbol de Decisión de manera secuencial
func entrenarDecisionTree(dataset []Registro, profundidad int) string {
	if profundidad == 0 || len(dataset) == 0 {
		return "Hoja"
	}
	// Dividir el dataset en dos ramas (ejemplo simple, basado en feature1)
	var izquierda, derecha []Registro
	for _, registro := range dataset {
		if registro.feature1 < 0.5 {
			izquierda = append(izquierda, registro)
		} else {
			derecha = append(derecha, registro)
		}
	}
	// Recursión para las ramas izquierda y derecha
	return fmt.Sprintf("Nodo(izquierda=%s, derecha=%s)", entrenarDecisionTree(izquierda, profundidad-1), entrenarDecisionTree(derecha, profundidad-1))
}

func main() {
	// Medir el tiempo de ejecución
	inicio := time.Now()

	// Crear dataset simulado de 1 millón de registros
	dataset := crearDataset(1000000)
	fmt.Println("Dataset creado con éxito!")

	// Entrenar Árbol de Decisión con profundidad máxima de 3
	arbol := entrenarDecisionTree(dataset, 3)
	fmt.Println("Árbol de Decisión entrenado:", arbol)

	// Calcular el tiempo transcurrido
	duracion := time.Since(inicio)
	fmt.Printf("Tiempo total de ejecución: %s\n", duracion)
}
