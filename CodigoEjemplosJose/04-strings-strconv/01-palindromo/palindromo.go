package main

import (
	"fmt"
	"strings"
)

// Función que devuelve la cadena pasada por parámetro al revés
func reverse(cadena string) string {

	arrayCadena := strings.Split(cadena, "")
	arrayInvertida := make([]string, 0)

	for i := len(arrayCadena) - 1; i >= 0; i-- {
		arrayInvertida = append(arrayInvertida, arrayCadena[i])
	}
	//fmt.Println(arrayCadena)
	//fmt.Println(arrayInvertida)

	return strings.Join(arrayInvertida, "")
}

func esPalindromo(palabra string) bool {

	//fmt.Println(palabra)
	palabra = strings.ToLower(palabra)

	//fmt.Println(palabra)
	//palabra = strings.ToUpper(palabra)
	//fmt.Println(palabra)

	//palabra = strings.Replace(palabra, "Z", "S", 2)ç
	palabra = strings.ReplaceAll(palabra, " ", "")
	//fmt.Println(palabra)

	//fmt.Println(reverse(palabra))
	palabraInvertida := reverse(palabra)

	return palabra == palabraInvertida

}

func main() {

	// Una palabra palíndroma es que se lee de igual forma de izquierda a derecha y viceversa
	// Un ejemplo podría ser la palabra "oso"
	//esPalindromo("Oso")
	if esPalindromo("ANA") { // También es una palabra palíndroma si le quitamos ese espacio
		fmt.Println("Es palíndromo")
	} else {
		fmt.Println("No es palíndromo")
	}

	//reverse("Luz azul")

}
