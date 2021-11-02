package FWQ_Visitor

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func CrearPerfil(ipRegistry, puertoRegistry string) {
	fmt.Println("Creación de perfil")
	conn, err := net.Dial("tcp", ipRegistry+":"+puertoRegistry)
	if err != nil {
		fmt.Println("Error al conectarse:", err.Error())
		os.Exit(1)
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Información del cliente")
		input, _ := reader.ReadString('\n')
		message, _ := bufio.NewReader(conn).ReadString('\n')
		log.Print("Server relay:", message)
	}
}

func EditarPerfil(ipRegistry, puertoRegistry string) {
	fmt.Println("Has entrado a editar perfil")
	conn, err := net.Dial("tcp", ipRegistry+":"+puertoRegistry)
	if err != nil {
		fmt.Println("Error al conectarse:", err.Error())
		os.Exit(1)
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Información del cliente que se quiere modificar")
		input, _ := reader.ReadString('\n')
		message, _ := bufio.NewReader(conn).ReadString('\n')
		log.Print("Server relay:", message)
	}

}

func EntradaParque(ipRegistry, puertoRegistry string) {
	fmt.Println("*Bienvenido al parque de atracciones*")
	conn, err := net.Dial("tcp", ipRegistry+":"+puertoRegistry)
	if err != nil {
		fmt.Println("Error al conectarse:", err.Error())
		os.Exit(1)
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Por favor introduce tu alias y tu password")
		input, _ := reader.ReadString('\n')
		message, _ := bufio.NewReader(conn).ReadString('\n')
		log.Print("Server relay:", message)
	}

}

func SalidaParque(ipRegistry, puertoRegistry string) {

}
