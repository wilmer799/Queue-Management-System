package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const (
	host = "172.20.42.134"
	tipo = "tcp"
)

func main() {

	puerto := os.Args[1]

	fmt.Println("Arrancando el registrador atendiendo en " + host + ":" + puerto)

	l, err := net.Listen(tipo, host+":"+puerto)

	if err != nil {
		fmt.Println("Error escuchando", err.Error())
		os.Exit(1)
	}

	defer l.Close()

	for {

		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error conectando:", err.Error())
		}

		fmt.Println("Cliente " + c.RemoteAddr().String() + " conectado.")

		go manejoConexion(c)

	}

}

func manejoConexion(conexion net.Conn) {

	buffer, err := bufio.NewReader(conexion).ReadBytes('\n')

	if err != nil {
		fmt.Println("Cliente desconectado.")
		conexion.Close()
		return
	}

	visitante := strings.Split(string(buffer[:len(buffer)-1]), " ")

	log.Println("Visitante a registrar: " + visitante[0] + " | " + visitante[1] + " | " + visitante[2])

	conexion.Write(buffer)

	archivo, err := os.Create("./Visitante" + visitante[0] + ".txt")

	if err != nil {
		fmt.Println("Hubo un error")
		return
	}

	fmt.Fprintln(archivo, "ID: "+visitante[0])
	fmt.Fprintln(archivo, "Nombre: "+visitante[1])
	fmt.Fprintln(archivo, "Password: "+visitante[2])

	archivo.Close()

	manejoConexion(conexion)

}
