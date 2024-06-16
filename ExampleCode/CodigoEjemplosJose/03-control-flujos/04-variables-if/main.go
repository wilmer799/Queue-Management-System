package main

import "fmt"

func main() {

	if nombre, edad := "Carlos", 26; nombre == "Jose" {
		fmt.Println("Hola, ", nombre)
	} else {
		fmt.Printf("Nombre: %s - Edad: %d\n", nombre, edad)
	}

	// Aquí no estaría disponible la variable
	//(fmt.Println(nombre)

	// Obtener valor de mapa
	users := make(map[string]string)

	users["Jose"] = "jose@gmail.com"
	users["Carlos"] = "carlos@gmail.com"

	correo, ok := users["Jose"]
	fmt.Println(correo, ok)

	correo, ok = users["Juan"]
	fmt.Println(correo, ok) // Imprime un valor vacío y false

	if correo, ok := users["Jose"]; ok {
		fmt.Println(correo)
	} else {
		fmt.Println("No fue posible obtener el valor")
	}

	if correo, ok := users["Juan"]; ok {
		fmt.Println(correo)
	} else {
		fmt.Println("No fue posible obtener el valor")
	}

}
