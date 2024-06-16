package main

import (
	"fmt"
	"net/http"
	"time"
)

// Esta función se va a encargar de comprobar si un servidor está funcionando o no
// Recibe la url del servidor por parámetro
func revisarServidor(servidor string, canal chan string) {
	_, err := http.Get(servidor)

	if err != nil {
		//fmt.Println("No está disponible")
		canal <- servidor + " No está disponible" // Actualización de uso con canal
	} else {
		//fmt.Println(servidor, "Está funcionando")
		canal <- servidor + " Está funcionando" // Actualización de uso con canal
	}

}

func main() {

	inicio := time.Now()

	// Para comunicarse/tener esa comunicación entre hilos y saber qué está ´sucediendo vamos a implementar los canales
	canal := make(chan string) // Creamos un canal

	servidores := []string{
		"https://oregoom.com/",
		"https://www.udemy.com/",
		"https://www.youtube.com/",
		"https://www.facebook.com/",
		"https://www.google.com/",
	}

	for _, servidor := range servidores {
		//revisarServidor(servidor) // Ejecución de forma secuencial
		go revisarServidor(servidor, canal) // Ejecución de forma concurrente
		// Con la palabra reservada go indicamos la creación de múltiples hilos que nos van
		// a permitir ejecutar esta función al mismo tiempo.
	}

	// Ahora para obtener todo lo que nos está obteniendo el canal
	for i := 0; i < len(servidores); i++ {
		fmt.Println(<-canal) // Imprimimos lo que nos está devolviendo el canal
	}

	tiempoPaso := time.Since(inicio)

	fmt.Println("Tiempo de ejecución: ", tiempoPaso)

}
