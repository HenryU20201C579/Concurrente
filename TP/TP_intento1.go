package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

type DecisionNode struct {
	feature     int
	threshold   float64
	trueBranch  *DecisionNode
	falseBranch *DecisionNode
	prediction  float64
}

type Data struct {
	Pregnancies              float64
	Glucose                  float64
	BloodPressure            float64
	SkinThickness            float64
	Insulin                  float64
	BMI                      float64
	DiabetesPedigreeFunction float64
	Age                      float64
	Outcome                  float64
}

// Función leer el CSV y cargar los datos
func loadCSV(filename string) ([]Data, error) {
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

	var data []Data
	for _, record := range records[1:] {
		pregnancies, _ := strconv.ParseFloat(record[0], 64)
		glucose, _ := strconv.ParseFloat(record[1], 64)
		bloodPressure, _ := strconv.ParseFloat(record[2], 64)
		skinThickness, _ := strconv.ParseFloat(record[3], 64)
		insulin, _ := strconv.ParseFloat(record[4], 64)
		bmi, _ := strconv.ParseFloat(record[5], 64)
		diabetesPedigree, _ := strconv.ParseFloat(record[6], 64)
		age, _ := strconv.ParseFloat(record[7], 64)
		outcome, _ := strconv.ParseFloat(record[8], 64)

		data = append(data, Data{
			Pregnancies:              pregnancies,
			Glucose:                  glucose,
			BloodPressure:            bloodPressure,
			SkinThickness:            skinThickness,
			Insulin:                  insulin,
			BMI:                      bmi,
			DiabetesPedigreeFunction: diabetesPedigree,
			Age:                      age,
			Outcome:                  outcome,
		})
	}
	return data, nil
}

// Función para crear el árbol
func buildTree(data []Data, depth int) *DecisionNode {
	if len(data) == 0 || depth == 0 {
		return nil
	}

	bestFeature := rand.Intn(8)
	bestThreshold := rand.Float64()

	trueBranch := make([]Data, 0)
	falseBranch := make([]Data, 0)

	for _, row := range data {
		value := row.Glucose
		if value > bestThreshold {
			trueBranch = append(trueBranch, row)
		} else {
			falseBranch = append(falseBranch, row)
		}
	}

	if len(trueBranch) == 0 || len(falseBranch) == 0 {
		return &DecisionNode{prediction: data[0].Outcome}
	}

	return &DecisionNode{
		feature:     bestFeature,
		threshold:   bestThreshold,
		trueBranch:  buildTree(trueBranch, depth-1),
		falseBranch: buildTree(falseBranch, depth-1),
	}
}

// Función predecir Random Forest
func predict(node *DecisionNode, data Data) float64 {
	if node.prediction != 0 {
		return node.prediction
	}

	value := data.Glucose
	if value > node.threshold {
		return predict(node.trueBranch, data)
	}
	return predict(node.falseBranch, data)
}

// Entrenar árboles concurrentemente
func trainForest(data []Data, numTrees int, depth int) []*DecisionNode {
	forest := make([]*DecisionNode, numTrees)
	var wg sync.WaitGroup
	wg.Add(numTrees)

	for i := 0; i < numTrees; i++ {
		go func(i int) {
			defer wg.Done()
			tree := buildTree(data, depth)
			forest[i] = tree
		}(i)
	}
	wg.Wait()
	return forest
}

// Funcion predicción
func forestPredict(forest []*DecisionNode, data Data) float64 {
	votes := 0.0
	for _, tree := range forest {
		votes += predict(tree, data)
	}
	return votes / float64(len(forest)) // Promediar los resultados
}

// Calcular precisión del modelo
func calculateAccuracy(forest []*DecisionNode, data []Data) float64 {
	correct := 0
	for _, row := range data {
		prediction := forestPredict(forest, row)
		if (prediction >= 0.5 && row.Outcome == 1) || (prediction < 0.5 && row.Outcome == 0) {
			correct++
		}
	}
	return float64(correct) / float64(len(data))
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Cargar el dataset
	data, err := loadCSV("diabetes.csv")
	if err != nil {
		log.Fatal(err)
	}

	// Dividir datos en entrenamiento (80%) y prueba (20%)
	trainSize := int(0.8 * float64(len(data)))
	trainData := data[:trainSize]
	testData := data[trainSize:]

	// Entrenar Random Forest
	numTrees := 100
	maxDepth := 10
	forest := trainForest(trainData, numTrees, maxDepth)

	// Calcular precisión en el conjunto de prueba
	accuracy := calculateAccuracy(forest, testData)
	fmt.Printf("Precisión del modelo: %.2f%%\n", accuracy*100)

	// Probar una predicción con un ejemplo
	example := testData[0]
	prob := forestPredict(forest, example)
	fmt.Printf("Probabilidad de que el paciente tenga diabetes: %.2f%%\n", prob*100)
}
