package main

import "fmt"

func main() {

	// arrays
	var numeros [5]int // Es inmutable, es decir, es estático, no podemos agregarle elementos ni tampoco quitárselos, pero sí modificarlos

	numeros[0] = 10
	numeros[1] = 20
	numeros[2] = 30

	fmt.Println(numeros)
	fmt.Println(numeros[1])

	// Arrays con valores
	nombres := [3]string{"Jose", "Carlos", "Pepe"}

	fmt.Println(nombres)

	// No le indicamos la longitud exacta, pero la toma en función de los elementos que pongamos
	// Tampoco podemos alterar la cantidad de elementos del array en este caso
	colores := [...]string{
		"Rojo",
		"Verde",
		"Negro",
		"Azul",
	}

	fmt.Println(colores, len(colores))

	// Índices definidos
	monedas := [...]string{ // La longitud es 6, se coloca un carácter vacío para rellenar las posiciones restantes
		0: "Dolares",
		2: "Soles",
		3: "Pesos",
		5: "Euros",
	}

	fmt.Println(monedas, len(monedas))

	// sub Array
	subMoneda := monedas[0:3] // El elemento de la posición 0 es incluído, pero el de la posición 3 no
	//[:3] // Coge desde el principio hasta la posición 3 no incluída
	//[0:] // COge desde la posición 0 incluída hasta el final
	fmt.Println(subMoneda)

}
