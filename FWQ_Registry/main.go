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

	// Arrancamos el servidor y atendemos conexiones entrantes
	fmt.Println("Arrancando el registrador atendiendo en " + host + ":" + puerto)

	l, err := net.Listen(tipo, host+":"+puerto)

	if err != nil {
		fmt.Println("Error escuchando", err.Error())
		os.Exit(1)
	}

	// Cerramos el listener cuando se cierra la aplicación
	defer l.Close()

	// Bucle infinito hasta la salida del programa
	for {

		// Atendemos conexiones entrantes
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error conectando:", err.Error())
		}

		// Imprimimos la dirección de conexión del cliente
		fmt.Println("Cliente " + c.RemoteAddr().String() + " conectado.")

		// Manejamos las conexiones de forma concurrente
		go manejoConexion(c)

	}

}

func manejoConexion(conexion net.Conn) {

	// Lectura del buffer hasta el final de línea
	buffer, err := bufio.NewReader(conexion).ReadBytes('\n')

	// Cerramos la conexión de los clientes que se han desconectado
	if err != nil {
		fmt.Println("Cliente desconectado.")
		conexion.Close()
		return
	}

	visitante := strings.Split(string(buffer[:len(buffer)-1]), " ")

	// Imprimimos la información del visitante a registrar o editar en la base de datos
	log.Println("Visitante a registrar/editar: " + visitante[0] + " | " + visitante[1] + " | " + visitante[2])

	// Mandamos un mensaje de respuesta al cliente
	conexion.Write(buffer)

	// Reiniciamos el proceso
	manejoConexion(conexion)

	/* POR SI HICIERA FALTA HACERLO CON FICHEROS */
	/*archivo, err := os.Create("./Visitante" + visitante[0] + ".txt")

	if err != nil {
		fmt.Println("Hubo un error")
		return
	}

	fmt.Fprintln(archivo, "ID: "+visitante[0])
	fmt.Fprintln(archivo, "Nombre: "+visitante[1])
	fmt.Fprintln(archivo, "Password: "+visitante[2])

	archivo.Close()*/

}
