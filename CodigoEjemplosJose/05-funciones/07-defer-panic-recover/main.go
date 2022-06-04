package main

import (
	"fmt"
	"os"
)

func main() {

	// Función con la que podemos evitar que nuestra aplicación se detenga de forma brusca por un panic
	defer func() {

		if error := recover(); error != nil { // Recover retorna un error para controlar los panics
			fmt.Println("Ups!, al parecer el programa no finalizó de forma correcta")
		}

	}()

	if file, error := os.Open("hoa.txt"); error != nil {
		panic("No fue posible obtener el archivo!!")
	} else {
		defer func() { // Se ejecutará al final
			fmt.Println("Cerrando el archivo!")
			file.Close()
		}()
		contenido := make([]byte, 254)
		long, _ := file.Read(contenido) // Devuelve la longitud y el error
		contenidoArchivo := string(contenido[:long])
		fmt.Println(contenidoArchivo)
	}

}
