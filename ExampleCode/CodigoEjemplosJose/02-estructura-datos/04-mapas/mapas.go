package main

import "fmt"

// Los mapas son otra colección/estructura de datos, pero en este caso de elementos desordenados
// Están formados por clave-valor
func main() {

	dias := make(map[int]string) // Se indica de esta forma -> map[tipo de dato de la clave de este mapa]tipo de dato del valor de estas claves

	fmt.Println(dias)

	// Agregar datos
	dias[1] = "Lunes"
	dias[2] = "Martes"
	dias[3] = "Miércoles"
	dias[4] = "Jueves"
	dias[5] = "Viernes"
	dias[6] = "Sábado"
	dias[7] = "Domingo"

	fmt.Println(dias)

	dias[7] = "DOMINGO"

	fmt.Println(dias)

	delete(dias, 1)

	fmt.Println(dias, len(dias))

	// Los mapas se puede volver colecciones más complejas en las cuales pueden ser arrays, objetos, estructuras, interface o cualquier tipo de dato
	estudiantes := make(map[string][]int)

	estudiantes["Jose"] = []int{13, 15, 16}
	estudiantes["Carlos"] = []int{14, 13, 17}

	fmt.Println(estudiantes)

	fmt.Println(estudiantes["Jose"])

	fmt.Println(estudiantes["Jose"][1])
}
