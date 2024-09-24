package RF

import (
	"encoding/json"    // Para serialización y deserialización JSON
	"fmt"              // Para formatear cadenas y salida en consola
	"math"             // Operaciones matemáticas básicas
	"math/rand"        // Generación de números aleatorios
	"os"               // Para manipulación de archivos (abrir, crear)
	"sync"             // Para concurrencia: mutex y sincronización
	"time"             // Para obtener la hora actual (usada para la semilla aleatoria)
)

// Estructura `Forest` que contiene un slice de punteros a `Tree` (árboles de decisión)
type Forest struct {
	Trees []*Tree
}

// `BuildForest` crea un bosque aleatorio con `treesAmount` cantidad de árboles.
// Recibe las entradas (`inputs`), etiquetas (`labels`), cantidad de árboles (`treesAmount`),
// cantidad de muestras (`samplesAmount`) y cantidad de características seleccionadas (`selectedFeatureAmount`).
func BuildForest(inputs [][]interface{}, labels []string, treesAmount, samplesAmount, selectedFeatureAmount int) *Forest {
	// Inicializa la semilla aleatoria usando el tiempo actual para evitar generar siempre el mismo conjunto de números aleatorios.
	rand.Seed(time.Now().UnixNano())

	// Crea una instancia del bosque.
	forest := &Forest{}
	// Reserva espacio en memoria para almacenar los punteros a los árboles.
	forest.Trees = make([]*Tree, treesAmount)

	// Canal para señalar cuándo una goroutine termina de entrenar un árbol.
	done_flag := make(chan bool)

	// Contador de progreso que lleva el número de árboles ya entrenados.
	prog_counter := 0

	// Mutex para proteger el acceso al contador de progreso, evitando condiciones de carrera.
	mutex := &sync.Mutex{}

	// Bucle para lanzar las goroutines que entrenarán cada árbol de decisión.
	for i := 0; i < treesAmount; i++ {
		go func(x int) {
			// Imprime en consola cuándo comienza a construirse un árbol.
			fmt.Printf(">> %v buiding %vth tree...\n", time.Now(), x)

			// Llama a `BuildTree` para construir un árbol de decisión y lo almacena en el bosque.
			forest.Trees[x] = BuildTree(inputs, labels, samplesAmount, selectedFeatureAmount)

			// Bloquea el acceso al contador de progreso para incrementarlo de manera segura.
			mutex.Lock()
			prog_counter += 1
			// Imprime el porcentaje de progreso actual.
			fmt.Printf("%v tranning progress %.0f%%\n", time.Now(), float64(prog_counter)/float64(treesAmount)*100)
			// Desbloquea el mutex.
			mutex.Unlock()

			// Señala a través del canal que el entrenamiento de un árbol ha terminado.
			done_flag <- true
		}(i) // Pasamos el índice `i` a la goroutine para identificar el árbol.
	}

	// Espera hasta que todos los árboles hayan terminado su entrenamiento.
	for i := 1; i <= treesAmount; i++ {
		<-done_flag // Bloquea hasta que se reciba una señal de cada árbol.
	}

	// Imprime un mensaje cuando todos los árboles han sido entrenados.
	fmt.Println("all done.")
	return forest // Devuelve el bosque entrenado.
}

// `DefaultForest` crea un bosque con parámetros por defecto. 
// Calcula la cantidad de muestras y características seleccionadas como la raíz cuadrada del total.
func DefaultForest(inputs [][]interface{}, labels []string, treesAmount int) *Forest {
	// Número de características seleccionadas es la raíz cuadrada del número total de características.
	m := int(math.Sqrt(float64(len(inputs[0]))))
	// Número de muestras seleccionadas es la raíz cuadrada del número total de entradas.
	n := int(math.Sqrt(float64(len(inputs))))
	// Llama a `BuildForest` con los valores calculados.
	return BuildForest(inputs, labels, treesAmount, n, m)
}

// `Predicate` predice la clase para un conjunto de datos de entrada (`input`).
// Recorre todos los árboles en el bosque y cuenta las predicciones de cada árbol.
func (self *Forest) Predicate(input []interface{}) string {
	// Mapa para contar la frecuencia de predicciones para cada clase.
	counter := make(map[string]float64)

	// Recorre cada árbol del bosque para obtener la predicción.
	for i := 0; i < len(self.Trees); i++ {
		// `PredicateTree` devuelve las predicciones del árbol actual.
		tree_counter := PredicateTree(self.Trees[i], input)
		total := 0.0
		// Calcula el total de votos de este árbol.
		for _, v := range tree_counter {
			total += float64(v)
		}
		// Normaliza los votos de este árbol y los agrega al contador global.
		for k, v := range tree_counter {
			counter[k] += float64(v) / total
		}
	}

	// Encuentra la clase con la mayor cantidad de votos.
	max_c := 0.0
	max_label := ""
	for k, v := range counter {
		// Si la frecuencia de esta clase es mayor o igual a la máxima registrada, la actualiza.
		if v >= max_c {
			max_c = v
			max_label = k
		}
	}
	// Devuelve la clase con más votos.
	return max_label
}

// `DumpForest` guarda el bosque en un archivo JSON.
// Esto permite guardar el modelo entrenado para reutilizarlo más tarde.
func DumpForest(forest *Forest, fileName string) {
	// Abre o crea el archivo donde se almacenará el bosque.
	out_f, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		panic("failed to create " + fileName) // Error si no puede crear el archivo.
	}
	defer out_f.Close() // Asegura que el archivo se cierre después de su uso.
	// Crea un nuevo codificador JSON y almacena el bosque en el archivo.
	encoder := json.NewEncoder(out_f)
	encoder.Encode(forest)
}

// `LoadForest` carga un bosque desde un archivo JSON.
// Esto permite cargar un modelo entrenado previamente.
func LoadForest(fileName string) *Forest {
	// Abre el archivo que contiene el bosque.
	in_f, err := os.Open(fileName)
	if err != nil {
		panic("failed to open " + fileName) // Error si no puede abrir el archivo.
	}
	defer in_f.Close() // Asegura que el archivo se cierre después de su uso.
	// Crea un nuevo decodificador JSON para leer el contenido del archivo.
	decoder := json.NewDecoder(in_f)
	forest := &Forest{}
	// Decodifica el bosque desde el archivo.
	decoder.Decode(forest)
	return forest // Devuelve el bosque cargado.
}
