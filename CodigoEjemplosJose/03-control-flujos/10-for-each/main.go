package main

import "fmt"

func main() {

	nombres := [...]string{"Jose", "Carlos", "Pepe"}

	/*for i := 0; i < len(nombres); i++ {
		fmt.Println(nombres[i])
	}*/

	// Cuando trabajamos con for each tenemos que definir 2 variables
	// para recuperar el índice y también para recuperar el elemento
	for indice, elemento := range nombres {
		fmt.Println(indice, elemento)
	}

	// Pero si sólo queremos simplemente el elemento
	for _, elemento := range nombres {
		fmt.Println(elemento)
	}

	// O si sólo queremos los índices
	for indice, _ := range nombres {
		fmt.Println(indice)
	}

}
