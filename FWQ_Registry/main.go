package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

/*
* Estructura del visitante
 */
type visitante struct {
	ID        string `json:"id"`
	Nombre    string `json:"nombre"`
	Password  string `json:"contraseña"`
	Posicionx int    `json:"posicionx"`
	Posiciony int    `json:"posiciony"`
	Destinox  int    `json:"destinox"`
	Destinoy  int    `json:"destinoy"`
	Parque    string `json:parqueAtracciones`
}

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

	vis := strings.Split(string(buffer[:len(buffer)-1]), " ")

	// Imprimimos la información del visitante a registrar o editar en la base de datos
	log.Println("Visitante a registrar/editar: " + vis[0] + " | " + vis[1] + " | " + vis[2])

	// Mandamos un mensaje de respuesta al cliente
	conexion.Write(buffer)

	// Accedemos a la base de datos, empezando por abrir la conexión
	db, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/parque_atracciones")

	// Comprobamos que no haya error al conectarse
	if err != nil {
		panic("Error al conectarse con la BD: " + err.Error())
	}

	defer db.Close() // Para que siempre se cierre la conexión con la BD al finalizar el programa

	results, err := db.Query("SELECT * FROM VISITANTE WHERE id == " + vis[0])

	// Comrpobamos que no se produzcan errores al hacer la consulta
	if err != nil {
		panic("Error al hacer la consulta a la BD: " + err.Error())
	}

	v := visitante{
		ID:       vis[0],
		Nombre:   vis[1],
		Password: vis[2],
	}

	// Comprobamos que la consulta haya devuelto alguna fila de la BD
	// Si el visitante existe en la BD
	if results.Next() {

		// MODIFICAMOS la información de dicho visitante en la BD
		// Preparamos para prevenir inyecciones SQL
		sentenciaPreparada, err := db.Prepare("UPDATE visitante (id, nombre, contraseña) VALUES(?, ?, ?)")
		if err != nil {
			panic("Error al preparar la sentencia de modificación: " + err.Error())
		}

		defer sentenciaPreparada.Close()

		// Ejecutar sentencia, un valor por cada '?'
		_, err = sentenciaPreparada.Exec(v.ID, v.Nombre, v.Password)
		if err != nil {
			panic("Error al modificar el visitante: " + err.Error())
		}

	} else { // Sino existe en la BD

		// INSERTAMOS el nuevo visitante en la BD
		// Preparamos para prevenir inyecciones SQL
		sentenciaPreparada, err := db.Prepare("INSERT INTO visitante (id, nombre, contraseña) VALUES(?, ?, ?)")
		if err != nil {
			panic("Error al preparar la sentencia de inserción: " + err.Error())
		}

		defer sentenciaPreparada.Close()

		// Ejecutar sentencia, un valor por cada '?'
		_, err = sentenciaPreparada.Exec(v.ID, v.Nombre, v.Password)
		if err != nil {
			panic("Error al modificar el visitante: " + err.Error())
		}
	}

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
