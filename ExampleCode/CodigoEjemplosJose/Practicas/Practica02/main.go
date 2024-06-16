package main

import "fmt"

func cociente(a, b int) int {
	return a / b
}

func resto(a, b int) int {
	return a % b
}

func main() {

	var num1 int
	var num2 int

	fmt.Println("Calculo del cociente y el resto de dos números enteros")
	fmt.Println("Introduzca el primer número:")
	fmt.Scanln(&num1)
	fmt.Println("Introduzca el segundo número:")
	fmt.Scanln(&num2)

	coci := cociente(num1, num2)
	rest := resto(num1, num2)

	fmt.Printf("El cociente y el resto de dividir los números %d y %d es %d y %d respectivamente \n", num1, num2, coci, rest)

}
