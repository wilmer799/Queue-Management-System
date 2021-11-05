package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	connType = "tcp"
)

func main() {

	IpFWQ_Registry := os.Args[1]
	PuertoFWQ := os.Args[2]
	IpBroker := os.Args[3]
	PuertoBroker := os.Args[4]

	var opcion int
	fmt.Println("Bienvenido al parque de atracciones")
	fmt.Println("La IP del registro es la siguiente:" + IpFWQ_Registry + ":" + PuertoFWQ)
	fmt.Println("La IP del Broker es el siguiente:" + IpBroker + ":" + PuertoBroker)

	fmt.Print("Elige la opción que quieras realizar:")
	fmt.Scanln(&opcion)
	switch os := opcion; os {
	case 1:
		CrearPerfil(IpFWQ_Registry, PuertoFWQ)
	case 2:
		EditarPerfil(IpFWQ_Registry, PuertoFWQ)
	case 3:
		EntradaParque(IpFWQ_Registry, PuertoFWQ)
	case 4:
		SalidaParque(IpFWQ_Registry, PuertoFWQ)

	default:
		fmt.Println("Opción invalida, elegie otra opción")
	}
}

func CrearPerfil(ipRegistry, puertoRegistry string) {
	fmt.Println("**********Creación de perfil***********")
	conn, err := net.Dial(connType, ipRegistry+":"+puertoRegistry)
	if err != nil {
		fmt.Println("Error al conectarse:", err.Error())
		os.Exit(1)
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enviando mensajes por parte del visitante:")
		//Leer entrada hasta nueva linea, introduciendo llave
		input, _ := reader.ReadString('\n')
		//Enviamos la conexion del socket
		conn.Write([]byte(input))
		//Escuchando por el relay
		message, _ := bufio.NewReader(conn).ReadString('\n')
		//Print server relay
		log.Print("Server relay:", message)
	}
}

func EditarPerfil(ipRegistry, puertoRegistry string) {
	fmt.Println("Has entrado a editar perfil")
	conn, err := net.Dial(connType, ipRegistry+":"+puertoRegistry)
	if err != nil {
		fmt.Println("Error al conectarse:", err.Error())
		os.Exit(1)
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Información del cliente que se quiere modificar:")
		input, _ := reader.ReadString('\n')
		//Enviamos la conexion del socket
		conn.Write([]byte(input))
		message, _ := bufio.NewReader(conn).ReadString('\n')
		log.Print("Server relay:", message)
	}

}

func EntradaParque(ipRegistry, puertoRegistry string) {
	fmt.Println("*Bienvenido al parque de atracciones*")
	conn, err := net.Dial(connType, ipRegistry+":"+puertoRegistry)
	if err != nil {
		fmt.Println("Error al conectarse:", err.Error())
		os.Exit(1)
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Por favor introduce tu alias:")
		input, _ := reader.ReadString('\n')
		fmt.Print("y tu password:")
		salida, _ := reader.ReadString('\n')
		//Enviamos la conexion del socket
		conn.Write([]byte(input))
		conn.Write([]byte(salida))
		message, _ := bufio.NewReader(conn).ReadString('\n')
		log.Print("Server relay:", message)
	}

}

func SalidaParque(ipRegistry, puertoRegistry string) {
	fmt.Println("Gracias por venir al parque, espero que vuelvas cuanto antes")
}

func ConexionKafka(IpBroker, PuertoBroker string, ctx context.Context) {
	//Es para pruebas
	i := 0
	var broker1Addres string = IpBroker + ":" + PuertoBroker
	var broker2Addres string = IpBroker + ":" + PuertoBroker
	var topic string = "sd-events"
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker1Addres, broker2Addres},
		Topic:   topic,
	})
	for {
		err := w.WriteMessages(ctx, kafka.Message{
			Key:   []byte(strconv.Itoa(i)),
			Value: []byte("Esto es un mensaje por parte de los visitantes" + strconv.Itoa(i)),
		})
		if err != nil {
			panic("No se puede escribir mensaje" + err.Error())
		}
		//Tenemos que enviar la información de los visitantes
		//Por lo que llamaremos a esta función desde crear perfil o editar perfil e ingresar en el parque
		fmt.Println("Escribiendo:", i)
		i++
		//Descanso
		time.Sleep(time.Second)
	}

}
