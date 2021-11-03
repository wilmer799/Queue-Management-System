package FWQ_Visitor

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const (
	connType = "tcp"
)

func CrearPerfil(ipRegistry, puertoRegistry string) {
	fmt.Println("Creación de perfil")
	conn, err := net.Dial(connType, ipRegistry+":"+puertoRegistry)
	if err != nil {
		fmt.Println("Error al conectarse:", err.Error())
		os.Exit(1)
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Enviando mensajes por parte del visitante:")
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
		fmt.Println("Información del cliente que se quiere modificar")
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

func ConexionKafka(IpBroker, PuertoBroker string) {

}
