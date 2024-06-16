package main

import "fmt"

func main() {

	// Con esta función podemos generar un slice vacío para luego nosotros agregarle valores
	numeros := make([]int, 3, 3) // Indicamos el tipo de dato, la longitud y la capacidad del slice

	numeros[0] = 100
	numeros[1] = 200
	numeros[2] = 300

	numeros = append(numeros, 400)

	fmt.Println(numeros, len(numeros), cap(numeros))

}
