package main

import "fmt"

func main() {

	a := 20
	b := 10

	// Suma
	result := a + b
	fmt.Println("Suma:", result)

	// Resta
	result = a - b
	fmt.Println("Resta:", result)

	// Multiplicación
	result = a * b
	fmt.Println("Multi:", result)

	// División
	result = a / b
	fmt.Println("División:", result)

	var div float64 = 3.0 / 2.0
	fmt.Println("División:", div)

	// Modulo (resto de una división)
	result = 3 % 2
	fmt.Println("Módulo:", result)
}
