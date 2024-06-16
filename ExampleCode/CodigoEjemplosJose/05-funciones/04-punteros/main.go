package main

import "fmt"

func modificarValor(cadena *string) {

	fmt.Printf("%p\n", cadena)
	*cadena = "Hola desde la función"

}

func main() {

	cadena := "Hola Mundo de Go"
	fmt.Printf("%p\n", &cadena)
	fmt.Println("Antes de la función:", cadena)

	//modificarValor(cadena) // Se pasa una copia, no la referencia en memoria
	modificarValor(&cadena)

	fmt.Println("Después de la función:", cadena)

}
