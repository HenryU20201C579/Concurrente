package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

// Definir estructura de los datos
type Data struct {
	features []float64
	label    int
}

// Crear dataset a partir del archivo CSV
func crearDatasetDesdeCSV(filename string) ([]Data, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	dataset := []Data{}
	for i, record := range records {
		if i == 0 { // Saltar la fila de encabezado
			continue
		}
		// Extraer características y etiqueta
		features := []float64{}
		for _, value := range record[1:5] { // Cambia según las columnas que desees usar
			floatValue, _ := strconv.ParseFloat(value, 64)
			features = append(features, floatValue)
		}
		label, _ := strconv.Atoi(record[7]) // Cambia el índice según la columna de la etiqueta
		dataset = append(dataset, Data{features, label})
	}

	return dataset, nil
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
	filename := "Afiliados_activos_DM_SIS.csv"
	fmt.Printf("Creando dataset desde %s...\n", filename)
	dataset, err := crearDatasetDesdeCSV(filename)
	if err != nil {
		log.Fatal(err)
	}
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

	// Mostrar algunas predicciones
	fmt.Println("Ejecutando algunas predicciones...")
	for i := 0; i < 5 && i < len(dataset); i++ {
		prediccion := predecir(dataset[i].features, pesos)
		fmt.Printf("Predicción para registro %d: %v, Etiqueta real: %d\n", i, prediccion, dataset[i].label)
	}

	// Calcular tiempo de ejecución
	elapsed := time.Since(start)
	fmt.Printf("Tiempo de ejecución: %s\n", elapsed)
}
