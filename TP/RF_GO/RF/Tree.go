// a random forest implementation in GoLang
package RF

import (
	"math"      // Paquete utilizado para funciones matemáticas como logaritmos.
	"math/rand" // Paquete para generar números aleatorios.
)

// Declaramos dos constantes que representan tipos de columnas.
// 'CAT' es para columnas categóricas y 'NUMERIC' para columnas numéricas.
const CAT = "cat"
const NUMERIC = "numeric"

// Estructura que representa un nodo de un árbol de decisión.
type TreeNode struct {
	ColumnNo int          // Número de la columna por la que se divide en este nodo.
	Value    interface{}  // Valor específico de la columna por el que se hace la división.
	Left     *TreeNode    // Subárbol izquierdo (muestras que cumplen con la condición de división).
	Right    *TreeNode    // Subárbol derecho (muestras que no cumplen con la condición de división).
	Labels   map[string]int // Mapa que almacena las etiquetas de las muestras para nodos hoja (finales).
}

// Estructura que representa un árbol de decisión.
type Tree struct {
	Root *TreeNode // Nodo raíz del árbol de decisión.
}

// Función que genera un rango de enteros aleatorios entre 0 y N, seleccionando M elementos únicos.
// Esta función se utiliza para seleccionar un subconjunto aleatorio de características en los nodos del árbol (para Random Forest).
func getRandomRange(N int, M int) []int {
	tmp := make([]int, N) // Se crea un slice de tamaño N.
	for i := 0; i < N; i++ {
		tmp[i] = i // Inicializa el slice con valores secuenciales del 0 al N-1.
	}
	// Se seleccionan M valores aleatorios del slice.
	for i := 0; i < M; i++ {
		// Intercambia el valor en la posición i con otro valor aleatorio dentro de las posiciones restantes.
		j := i + int(rand.Float64()*float64(N-i))
		tmp[i], tmp[j] = tmp[j], tmp[i]
	}

	return tmp[:M] // Retorna los primeros M elementos aleatorios.
}

// Función que selecciona un subconjunto de muestras (filas) del conjunto de datos original usando los índices proporcionados.
func getSamples(ary [][]interface{}, index []int) [][]interface{} {
	// Se crea un nuevo slice para almacenar las muestras seleccionadas.
	result := make([][]interface{}, len(index))
	for i := 0; i < len(index); i++ {
		// Se agregan las filas seleccionadas a partir de los índices en el resultado.
		result[i] = ary[index[i]]
	}
	return result // Retorna el subconjunto de muestras.
}

// Función que selecciona un subconjunto de etiquetas de las muestras usando los índices proporcionados.
func getLabels(ary []string, index []int) []string {
	// Se crea un nuevo slice para almacenar las etiquetas seleccionadas.
	result := make([]string, len(index))
	for i := 0; i < len(index); i++ {
		// Se asignan las etiquetas correspondientes a los índices.
		result[i] = ary[index[i]]
	}
	return result // Retorna el subconjunto de etiquetas.
}

// Función para calcular la entropía de un conjunto de datos.
// La entropía mide la incertidumbre o impureza de las etiquetas en las muestras.
func getEntropy(ep_map map[string]float64, total int) float64 {
	// Normaliza los valores del mapa dividiendo las frecuencias de cada etiqueta por el total.
	for k := range ep_map {
		ep_map[k] = ep_map[k] / float64(total)
	}

	entropy := 0.0
	// Calcula la entropía utilizando la fórmula de entropía de Shannon.
	for _, v := range ep_map {
		entropy += v * math.Log(1.0/v)
	}

	return entropy // Retorna el valor de la entropía calculada.
}

// Función para calcular la impureza del índice de Gini.
// El índice de Gini mide la probabilidad de que una instancia sea clasificada incorrectamente.
func getGini(ep_map map[string]float64) float64 {
	total := 0.0
	// Suma las frecuencias de todas las etiquetas.
	for _, v := range ep_map {
		total += v
	}

	// Normaliza las frecuencias dividiendo por el total.
	for k := range ep_map {
		ep_map[k] = ep_map[k] / total
	}

	impure := 0.0
	// Calcula la impureza del índice de Gini sumando los productos cruzados de probabilidades diferentes.
	for k1, v1 := range ep_map {
		for k2, v2 := range ep_map {
			if k1 != k2 {
				impure += v1 * v2
			}
		}
	}
	return impure // Retorna el valor de la impureza calculada.
}
// Función que encuentra la mejor ganancia de información para una columna específica.
// Evalúa las divisiones posibles y determina el valor y la columna que ofrecen la mayor ganancia.
func getBestGain(samples [][]interface{}, c int, samples_labels []string, column_type string, current_entropy float64) (float64, interface{}, int, int) {
	var best_value interface{} // Mejor valor para dividir
	best_gain := 0.0           // Mejora máxima encontrada
	best_total_r := 0          // Número de elementos en la rama derecha
	best_total_l := 0          // Número de elementos en la rama izquierda

	// Almacena los valores únicos de la columna c
	uniq_values := make(map[interface{}]int)
	for i := 0; i < len(samples); i++ {
		uniq_values[samples[i][c]] = 1
	}

	// Itera sobre cada valor único para evaluar la ganancia de información
	for value, _ := range uniq_values {
		map_l := make(map[string]float64) // Almacena la distribución de etiquetas a la izquierda
		map_r := make(map[string]float64) // Almacena la distribución de etiquetas a la derecha
		total_l := 0                      // Número de elementos en la rama izquierda
		total_r := 0                      // Número de elementos en la rama derecha

		// Evaluación según el tipo de columna (categórica o numérica)
		if column_type == CAT {
			// Si la columna es categórica
			for j := 0; j < len(samples); j++ {
				if samples[j][c] == value {
					total_l += 1
					map_l[samples_labels[j]] += 1.0
				} else {
					total_r += 1
					map_r[samples_labels[j]] += 1.0
				}
			}
		}
		if column_type == NUMERIC {
			// Si la columna es numérica
			for j := 0; j < len(samples); j++ {
				if samples[j][c].(float64) <= value.(float64) {
					total_l += 1
					map_l[samples_labels[j]] += 1.0
				} else {
					total_r += 1
					map_r[samples_labels[j]] += 1.0
				}
			}
		}

		// Calcula las probabilidades de pertenecer a cada rama
		p1 := float64(total_r) / float64(len(samples))
		p2 := float64(total_l) / float64(len(samples))

		// Calcula la nueva entropía después de la división
		new_entropy := p1*getEntropy(map_r, total_r) + p2*getEntropy(map_l, total_l)

		// Ganancia de información
		entropy_gain := current_entropy - new_entropy

		// Si la ganancia de entropía es mayor a la mejor registrada, la actualiza
		if entropy_gain >= best_gain {
			best_gain = entropy_gain
			best_value = value
			best_total_l = total_l
			best_total_r = total_r
		}
	}

	// Retorna la mejor ganancia, el mejor valor para dividir y el tamaño de las ramas
	return best_gain, best_value, best_total_l, best_total_r
}

// Función para dividir el conjunto de muestras basado en un valor específico y una columna.
// Los índices de las muestras que cumplen con la condición se almacenan en part_l y part_r.
func splitSamples(samples [][]interface{}, column_type string, c int, value interface{}, part_l *[]int, part_r *[]int) {
	if column_type == CAT {
		// Si la columna es categórica
		for j := 0; j < len(samples); j++ {
			if samples[j][c] == value {
				*part_l = append(*part_l, j)
			} else {
				*part_r = append(*part_r, j)
			}
		}
	}
	if column_type == NUMERIC {
		// Si la columna es numérica
		for j := 0; j < len(samples); j++ {
			if samples[j][c].(float64) <= value.(float64) {
				*part_l = append(*part_l, j)
			} else {
				*part_r = append(*part_r, j)
			}
		}
	}
}

// Construcción recursiva del árbol de decisión.
func buildTree(samples [][]interface{}, samples_labels []string, selected_feature_count int) *TreeNode {
	column_count := len(samples[0])              // Número total de columnas
	split_count := selected_feature_count        // Número de características seleccionadas
	columns_choosen := getRandomRange(column_count, split_count) // Columnas seleccionadas al azar

	best_gain := 0.0
	var best_part_l []int = make([]int, 0, len(samples)) // Índices de la rama izquierda
	var best_part_r []int = make([]int, 0, len(samples)) // Índices de la rama derecha
	var best_total_l int = 0                             // Tamaño de la rama izquierda
	var best_total_r int = 0                             // Tamaño de la rama derecha
	var best_value interface{}                           // Mejor valor de división
	var best_column int                                  // Mejor columna para dividir
	var best_column_type string                          // Tipo de la mejor columna

	// Mapa que cuenta la frecuencia de cada etiqueta
	current_entropy_map := make(map[string]float64)
	for i := 0; i < len(samples_labels); i++ {
		current_entropy_map[samples_labels[i]] += 1
	}

	// Calcula la entropía actual del nodo
	current_entropy := getEntropy(current_entropy_map, len(samples_labels))

	// Itera sobre las columnas seleccionadas al azar para encontrar la mejor división
	for _, c := range columns_choosen {
		column_type := CAT
		if _, ok := samples[0][c].(float64); ok {
			column_type = NUMERIC
		}

		// Calcula la ganancia de información para la columna actual
		gain, value, total_l, total_r := getBestGain(samples, c, samples_labels, column_type, current_entropy)

		// Si la ganancia es mejor que la actual, actualiza los mejores valores
		if gain >= best_gain {
			best_gain = gain
			best_value = value
			best_column = c
			best_column_type = column_type
			best_total_l = total_l
			best_total_r = total_r
		}
	}

	// Si se encuentra una buena división, crea un nodo y divide el conjunto
	if best_gain > 0 && best_total_l > 0 && best_total_r > 0 {
		node := &TreeNode{}
		node.Value = best_value
		node.ColumnNo = best_column
		splitSamples(samples, best_column_type, best_column, best_value, &best_part_l, &best_part_r)
		node.Left = buildTree(getSamples(samples, best_part_l), getLabels(samples_labels, best_part_l), selected_feature_count)
		node.Right = buildTree(getSamples(samples, best_part_r), getLabels(samples_labels, best_part_r), selected_feature_count)
		return node
	}

	// Si no se encuentra una buena división, genera una hoja
	return genLeafNode(samples_labels)
}

// Genera un nodo hoja con las etiquetas correspondientes.
func genLeafNode(labels []string) *TreeNode {
	counter := make(map[string]int)
	for _, v := range labels {
		counter[v] += 1
	}

	node := &TreeNode{}
	node.Labels = counter
	return node
}

// Predice la clase de una entrada dada, recorriendo el árbol desde la raíz.
func predicate(node *TreeNode, input []interface{}) map[string]int {
	if node.Labels != nil { // Si es una hoja, retorna las etiquetas
		return node.Labels
	}

	c := node.ColumnNo
	value := input[c]

	// Según el tipo de dato de la columna, continúa la predicción
	switch value.(type) {
	case float64:
		if value.(float64) <= node.Value.(float64) && node.Left != nil {
			return predicate(node.Left, input)
		} else if node.Right != nil {
			return predicate(node.Right, input)
		}
	case string:
		if value == node.Value && node.Left != nil {
			return predicate(node.Left, input)
		} else if node.Right != nil {
			return predicate(node.Right, input)
		}
	}

	return nil
}

// Función que construye un árbol a partir de las entradas y etiquetas proporcionadas.
func BuildTree(inputs [][]interface{}, labels []string, samples_count, selected_feature_count int) *Tree {

	// Selecciona una muestra aleatoria del conjunto de datos
	samples := make([][]interface{}, samples_count)
	samples_labels := make([]string, samples_count)
	for i := 0; i < samples_count; i++ {
		j := int(rand.Float64() * float64(len(inputs)))
		samples[i] = inputs[j]
		samples_labels[i] = labels[j]
	}

	// Crea y construye el árbol
	tree := &Tree{}
	tree.Root = buildTree(samples, samples_labels, selected_feature_count)

	return tree
}

// Función para predecir la clase de una entrada utilizando el árbol construido.
func PredicateTree(tree *Tree, input []interface{}) map[string]int {
	return predicate(tree.Root, input)
}
