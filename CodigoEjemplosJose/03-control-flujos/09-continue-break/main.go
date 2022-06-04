package main

import "fmt"

func main() {

	for i := 0; i <= 10; i++ {

		if i == 5 {
			fmt.Println("Salta la iteración")
			continue // Salta a la siguiente iteración
		}

		if i == 8 {
			fmt.Println("Romper con bucle ")
			break // Romper el bucle
		}

		fmt.Println(i)

	}

}
