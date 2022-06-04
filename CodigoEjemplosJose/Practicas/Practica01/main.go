package main

import "fmt"

func suma(a, b int) int {
	return a + b
}

func main() {

	var num1 int
	var num2 int

	fmt.Println("Suma de dos números")
	fmt.Println("Introduzca el primer número:")
	fmt.Scanln(&num1)
	fmt.Println("Introduzca el segundo número:")
	fmt.Scanln(&num2)
	resultado := suma(num1, num2)
	fmt.Printf("El resultado de sumar %d y %d es igual a %d \n", num1, num2, resultado)

}
