package main

import "fmt"

func main() {
	hola := "Hola"
	mundo := "Mundo"

	fmt.Println(hola)
	fmt.Println(mundo)

	nombre := "Jose"
	edad := 25

	// Si quisieramos dejar el espacio de un tabulador podemos usar \t
	fmt.Printf("Hola, %s y su edad es %d \n", nombre, edad)

	// %v se usa cuando no sabemos qué tipo de dato se va a imprimir en esta parte
	fmt.Printf("Hola, %v y su edad es %v \n", nombre, edad)

	// Sprintf formatea la información y la retorna formateada
	mensaje := fmt.Sprintf("Hola, %v y su edad es %v", nombre, edad)

	fmt.Println(mensaje)

	// Esto lo podemos usar para saber de qué tipo es una variable
	fmt.Printf("nombre: %T \n", nombre)
	fmt.Printf("nombre: %T \n", edad)

	fmt.Print("Ingrese otro nombre: ")
	fmt.Scanln(&nombre)

	fmt.Println("Otro Nombre: ", nombre)

}
