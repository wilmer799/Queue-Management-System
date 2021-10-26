package main

import (
	"fmt"
	"os"
)

var Alias string
var usuario string
var password string
var opcion int

func main() {
	//fmt.Println("Okay")
	//appkafka.StartKafka()
	//appkafka "github.com/wilmer799/SDPracticas/kafka"
	//productorK "github.com/wilmer799/SDPracticas/productor"
	/*
		fmt.Println("Kafka ha iniciado correctamente")
		ctx := context.Background()
		fmt.Println("Envia los mensajes necesarios:")
		productorK.StartProduce(ctx)

		time.Sleep(10 * time.Minute)
	*/
	IpFWQ_Registry := os.Args[1]
	PuertoFWQ := os.Args[2]
	IpBroker := os.Args[3]
	PuertoBroker := os.Args[4]
	fmt.Println("Conectandose al registro de la aplicación")
	fmt.Println("La IP del registro es la siguiente:" + IpFWQ_Registry + ":" + PuertoFWQ)
	fmt.Println("La IP del Broker es el siguiente:" + IpBroker + ":" + PuertoBroker)
	fmt.Println("Bienvenido al parque de atracciones")
	fmt.Println("Elige la opción que quieras realizar:")
	fmt.Scanln(&opcion)
	switch os := opcion; os {
	case 1:
		//CrearPerfil(IpFWQ_Registry, PuertoFWQ)
		fmt.Println("Creación de perfil")
	case 2:
		fmt.Println("Edición de perfil")
	case 3:
		fmt.Println("Entrada al parque de atracciones")
	case 4:
		fmt.Println("Salida del parque de atracciones")
	default:
		fmt.Println("Opción invalida, elegie otra opción")
	}
}
