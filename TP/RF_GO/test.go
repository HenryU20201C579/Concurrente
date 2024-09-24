package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	//"strconv"
	"os"

	"tp-test/RF"

	"time"
	//"math"
)

func main() {

	start := time.Now()
	f, _ := os.Open("diabetesV3.csv")
	defer f.Close()
	content, _ := ioutil.ReadAll(f)
	s_content := string(content)
	lines := strings.Split(s_content, "\n")

	inputs := make([][]interface{}, 0)
	targets := make([]string, 0)
	for _, line := range lines {

		line = strings.TrimRight(line, "\r\n")

		if len(line) == 0 {
			continue
		}
		tup := strings.Split(line, ",")
		pattern := tup[:len(tup)-1]
		target := tup[len(tup)-1]
		X := make([]interface{}, 0)
		for _, x := range pattern {
			X = append(X, x)
		}
		inputs = append(inputs, X)

		targets = append(targets, target)
	}
	train_inputs := make([][]interface{}, 0)

	train_targets := make([]string, 0)

	test_inputs := make([][]interface{}, 0)
	test_targets := make([]string, 0)

	for i, x := range inputs {
		if i%2 == 1 {
			test_inputs = append(test_inputs, x)
		} else {
			train_inputs = append(train_inputs, x)
		}
	}

	for i, y := range targets {
		if i%2 == 1 {
			test_targets = append(test_targets, y)
		} else {
			train_targets = append(train_targets, y)
		}
	}

	forest := RF.BuildForest(inputs, targets, 10, 500, len(train_inputs[0])) //100 trees

	test_inputs = train_inputs
	test_targets = train_targets
	err_count := 0.0
	for i := 0; i < len(test_inputs); i++ {
		output := forest.Predicate(test_inputs[i])
		expect := test_targets[i]
		//fmt.Println(output, expect)
		if output != expect {
			err_count += 1
		}
	}
	fmt.Println("success rate:", 1.0-err_count/float64(len(test_inputs)))

	fmt.Println(time.Since(start))

	// Ejecutar el menú para ingresar datos manualmente
	for {
		fmt.Println("---- Menú de Predicción de Diabetes ----")
		fmt.Println("1. Ingresar datos manualmente")
		fmt.Println("2. Salir")
		fmt.Print("Seleccione una opción: ")

		reader := bufio.NewReader(os.Stdin)
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		if choice == "1" {
			input := inputData()
			fmt.Println("Datos ingresados:", input)
			prob := forest.Predicate(input)
			fmt.Printf("Predicción: %s\n", prob)
		} else if choice == "2" {
			break
		} else {
			fmt.Println("Opción inválida, intente nuevamente.")
		}
	}
	//imprimir test
	/*fmt.Print("Test Input (Not diabetic): ")
	fmt.Println("Male, 15.0, 0, 0, never, 30.36, 6.1, 200")
	fmt.Println("Prediction:", forest.Predicate([]interface{}{"Male", "15.0", "0", "0", "never", "30.36", "6.1", "200"}))

	fmt.Print("Test Input (Diabetic): ")
	fmt.Println("Female, 53.0, 0, 0, former, 27.32, 7.0, 159")
	fmt.Println("Prediction:", forest.Predicate([]interface{}{"Female", "53.0", "0", "0", "former", "27.32", "7.0", "159"}))
	//fmt.Println("Prediction:", forest.Predicate(test_inputs[2]))*/
}

// Función para ingresar datos manualmente
func inputData() []interface{} {
	reader := bufio.NewReader(os.Stdin)

	// Validar género
	var gender string
	for {
		fmt.Print("Ingrese género (F/M): ")
		genderInput, _ := reader.ReadString('\n')
		genderInput = strings.TrimSpace(genderInput)

		if strings.ToUpper(genderInput) == "F" {
			gender = "Female"
			break
		} else if strings.ToUpper(genderInput) == "M" {
			gender = "Male"
			break
		} else {
			fmt.Println("Entrada inválida. Por favor, ingrese 'F' o 'M'.")
		}
	}

	// Validar edad
	var age string
	for {
		fmt.Print("Ingrese edad: ")
		ageStr, _ := reader.ReadString('\n')
		age = strings.TrimSpace(ageStr)

		if _, err := strconv.ParseFloat(age, 64); err == nil {
			break
		} else {
			fmt.Println("Entrada inválida. Por favor, ingrese una edad válida.")
		}
	}

	// Validar hipertensión
	var hypertension string
	for {
		fmt.Print("Ingrese hipertensión (0 o 1): ")
		hypertensionStr, _ := reader.ReadString('\n')
		hypertension = strings.TrimSpace(hypertensionStr)

		if hypertension == "0" || hypertension == "1" {
			break
		} else {
			fmt.Println("Entrada inválida. Por favor, ingrese '0' o '1'.")
		}
	}

	// Validar enfermedad cardíaca
	var heartDisease string
	for {
		fmt.Print("Ingrese enfermedad cardíaca (0 o 1): ")
		heartDiseaseStr, _ := reader.ReadString('\n')
		heartDisease = strings.TrimSpace(heartDiseaseStr)

		if heartDisease == "0" || heartDisease == "1" {
			break
		} else {
			fmt.Println("Entrada inválida. Por favor, ingrese '0' o '1'.")
		}
	}

	// Validar historial de fumar
	var smokingHistory int
	var smokingHistoryEng string

	for {
		fmt.Print("Ingrese historial de fumar (0: Nunca, 1: Sin información, 2: Actualmente, 3: No actualmente, 4: Anteriormente, 5: Alguna vez): ")
		smokingHistoryStr, _ := reader.ReadString('\n')

		// Intentar convertir la entrada a un entero
		var err error
		smokingHistory, err = strconv.Atoi(strings.TrimSpace(smokingHistoryStr))

		// Verificar si hubo un error en la conversión o si el número está fuera del rango permitido
		if err != nil || smokingHistory < 0 || smokingHistory > 5 {
			fmt.Println("Entrada inválida. Por favor, ingrese un número entre 0 y 5.")
			continue // Vuelve a solicitar la entrada
		}

		// Conversión del historial de fumar a inglés
		switch smokingHistory {
		case 0:
			smokingHistoryEng = "never"
		case 1:
			smokingHistoryEng = "No Info"
		case 2:
			smokingHistoryEng = "current"
		case 3:
			smokingHistoryEng = "not current"
		case 4:
			smokingHistoryEng = "former"
		case 5:
			smokingHistoryEng = "ever"
		}

		// Si la entrada es válida, puedes salir del bucle
		break
	}

	// Validar BMI
	var bmi string
	for {
		fmt.Print("Ingrese BMI: ")
		bmiStr, _ := reader.ReadString('\n')
		bmi = strings.TrimSpace(bmiStr)

		if _, err := strconv.ParseFloat(bmi, 64); err == nil {
			break
		} else {
			fmt.Println("Entrada inválida. Por favor, ingrese un BMI válido.")
		}
	}

	// Validar nivel HbA1c
	var hba1cLevel string
	for {
		fmt.Print("Ingrese nivel HbA1c: ")
		hba1cLevelStr, _ := reader.ReadString('\n')
		hba1cLevel = strings.TrimSpace(hba1cLevelStr)

		if _, err := strconv.ParseFloat(hba1cLevel, 64); err == nil {
			break
		} else {
			fmt.Println("Entrada inválida. Por favor, ingrese un nivel HbA1c válido.")
		}
	}

	// Validar nivel de glucosa en sangre
	var bloodGlucoseLevel string
	for {
		fmt.Print("Ingrese nivel de glucosa en sangre: ")
		bloodGlucoseLevelStr, _ := reader.ReadString('\n')
		bloodGlucoseLevel = strings.TrimSpace(bloodGlucoseLevelStr)

		if _, err := strconv.ParseFloat(bloodGlucoseLevel, 64); err == nil {
			break
		} else {
			fmt.Println("Entrada inválida. Por favor, ingrese un nivel de glucosa en sangre válido.")
		}
	}

	// Retorna los datos como un slice de interface{}
	return []interface{}{
		gender,            // Gender (Female/Male)
		age,               // Age
		hypertension,      // Hypertension (0 o 1)
		heartDisease,      // Heart Disease (0 o 1)
		smokingHistoryEng, // Smoking History (never/No Info/current/not current/former/ever)
		bmi,               // BMI
		hba1cLevel,        // HbA1c Level
		bloodGlucoseLevel, // Blood Glucose Level
	}
}
