package main

import "fmt"

func main() {

	// Slices (son mutables, por lo que podomos aumentar y reducir su tamaño)
	numeros := []int{1, 2, 3}

	fmt.Println(numeros, len(numeros))

	// Agregar datos
	numeros = append(numeros, 4)
	numeros = append(numeros, 5)

	fmt.Println(numeros, len(numeros))

	// Sub Slices
	subSlice := numeros[:2]

	numeros[0] = 100 // También se modifica en el slice, ya que los slices derivan de los arrays, es decir, se generan de un array padre

	fmt.Println(subSlice)
	fmt.Println(numeros)

	// Punteros
	// Longitud
	// Capacidad

	meses := []string{"Enero", "Febrero", "Marzo"}

	fmt.Printf("Len: %v, Cap: %v, %p \n", len(meses), cap(meses), meses)

	meses = append(meses, "Abril") // Se genera un nuevo slice con una nueva longitud y capacidad y agregando el nuevo elemento

	fmt.Printf("Len: %v, Cap: %v, %p \n", len(meses), cap(meses), meses)

}
