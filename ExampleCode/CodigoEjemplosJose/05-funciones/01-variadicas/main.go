package main

import "fmt"

// Esto lo hacemos para indicar que no sabemos cuántos valores vamos a recibir en la función
func sumar(nombre string, numeros ...int) (string, int) {

	mensaje := fmt.Sprintf("La suma de %s es: ", nombre)
	var total int
	for _, num := range numeros {
		total += num
	}

	return mensaje, total

}

func main() {

	mensaje, result := sumar("Jose", 10, 20, 40, 70, 60)

	fmt.Println(mensaje, result)

}
