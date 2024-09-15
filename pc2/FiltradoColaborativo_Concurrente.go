package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

// Definir estructura de los datos
type Usuario struct {
	ID             int
	Calificaciones map[int]float64 // ID del producto -> calificación
}

// Crear dataset simulado de 1000 usuarios y 100 productos
func crearDataset(numUsuarios, numProductos int) []Usuario {
	rand.Seed(time.Now().UnixNano())
	dataset := make([]Usuario, numUsuarios)
	for i := 0; i < numUsuarios; i++ {
		calificaciones := make(map[int]float64)
		for j := 0; j < numProductos; j++ {
			calificaciones[j] = rand.Float64() * 5 // Calificación aleatoria entre 0 y 5
		}
		dataset[i] = Usuario{ID: i, Calificaciones: calificaciones}
	}
	return dataset
}

// Calcular la similitud entre dos usuarios usando la similitud de Pearson
func similitudPearson(u1, u2 Usuario) float64 {
	// Calcular promedios de calificaciones
	var sum1, sum2, sum1Sq, sum2Sq, pSum float64
	n := 0

	for prod, cal1 := range u1.Calificaciones {
		if cal2, ok := u2.Calificaciones[prod]; ok {
			n++
			sum1 += cal1
			sum2 += cal2
			sum1Sq += cal1 * cal1
			sum2Sq += cal2 * cal2
			pSum += cal1 * cal2
		}
	}

	if n == 0 {
		return 0 // No hay calificaciones en común
	}

	num := pSum - (sum1 * sum2 / float64(n))
	den := math.Sqrt((sum1Sq - sum1*sum1/float64(n)) * (sum2Sq - sum2*sum2/float64(n)))

	if den == 0 {
		return 0 // Evitar división por cero
	}

	return num / den
}

// Predecir la calificación de un usuario para un producto basado en similitud
func predecirCalificacion(usuario Usuario, producto int, dataset []Usuario, wg *sync.WaitGroup, mu *sync.Mutex, resultados map[int]float64) {
	defer wg.Done()

	var sumaSimilitudes, sumaPonderaciones float64

	for _, u := range dataset {
		if u.ID != usuario.ID {
			sim := similitudPearson(usuario, u)
			if sim > 0 {
				if cal, ok := u.Calificaciones[producto]; ok {
					sumaSimilitudes += sim
					sumaPonderaciones += sim * cal
				}
			}
		}
	}

	if sumaSimilitudes != 0 {
		resultado := sumaPonderaciones / sumaSimilitudes
		mu.Lock()
		resultados[usuario.ID] = resultado
		mu.Unlock()
	}
}

func main() {
	// Crear el dataset
	numUsuarios := 1000
	numProductos := 100
	fmt.Printf("Creando dataset con %d usuarios y %d productos...\n", numUsuarios, numProductos)
	dataset := crearDataset(numUsuarios, numProductos)
	fmt.Println("Dataset creado con éxito.")

	// Medir el tiempo de ejecución
	start := time.Now()

	// Parámetros de la predicción
	productoID := 5

	var wg sync.WaitGroup
	var mu sync.Mutex
	resultados := make(map[int]float64)

	// Predecir calificación para un producto específico para todos los usuarios usando goroutines
	for _, usuario := range dataset {
		wg.Add(1)
		go predecirCalificacion(usuario, productoID, dataset, &wg, &mu, resultados)
	}

	// Esperar a que todas las goroutines terminen
	wg.Wait()

	// Calcular el tiempo de ejecución
	elapsed := time.Since(start)
	fmt.Println("Predicciones completadas.")
	defer fmt.Printf("Tiempo de ejecución: %s\n", elapsed)

	// Imprimir algunas predicciones
	for id, prediccion := range resultados {
		fmt.Printf("Usuario %d predice una calificación de %.2f para el producto %d\n", id, prediccion, productoID)
	}
}
