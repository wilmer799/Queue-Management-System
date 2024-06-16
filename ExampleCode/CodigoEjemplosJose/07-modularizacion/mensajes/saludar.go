package mensajes

import "fmt"

// Al poner el nombre de la función con una mayúscula la hacemos pública
func Hola() {
	fmt.Println("Hola desde el paquete mensajes")
}

const mensaje = "Hola desde mi constante"

// Esta en cambio es privada porque empiezar por minúscula.
// Este mismo criterio se usa para las variables y constantes.
func funcionPrivada() {
	fmt.Println("Hola desde la función privada")
}

func Imprimir() {
	fmt.Println(mensaje)
	funcionPrivada()
}
