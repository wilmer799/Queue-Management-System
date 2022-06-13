package main

import (
	"fmt"
	"os"

	api "github.com/wilmer799/SDPracticas/API_Engine"
	visitante "github.com/wilmer799/SDPracticas/FWQ_Visitor"
)

var opcion int

func main() {

	//https://github.com/sanathkr/go-npm

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
	fmt.Println("Bienvenido al parque de atracciones")
	fmt.Println("La IP del registro es la siguiente:" + IpFWQ_Registry + ":" + PuertoFWQ)
	fmt.Println("La IP del Broker es el siguiente:" + IpBroker + ":" + PuertoBroker)

	fmt.Print("Elige la opción que quieras realizar:")
	fmt.Scanln(&opcion)
	switch os := opcion; os {
	case 1:
		visitante.CrearPerfil(IpFWQ_Registry, PuertoFWQ)
	case 2:
		visitante.EditarPerfil(IpFWQ_Registry, PuertoFWQ)
	case 3:
		visitante.EntradaParque(IpFWQ_Registry, PuertoFWQ)
	case 4:
		visitante.SalidaParque(IpFWQ_Registry, PuertoFWQ)

	default:
		fmt.Println("Opción invalida, elegie otra opción")
	}
}
