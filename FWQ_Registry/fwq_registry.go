package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

/*
* Estructura del visitante
 */
type visitante struct {
	ID           string `json:"id"`
	Nombre       string `json:"nombre"`
	Password     string `json:"contraseña"`
	Posicionx    int    `json:"posicionx"`
	Posiciony    int    `json:"posiciony"`
	Destinox     int    `json:"destinox"`
	Destinoy     int    `json:"destinoy"`
	DentroParque int    `json:"dentroParque"`
	IdParque     string `json:"idParque"`
	Parque       string `json:"parqueAtracciones"`
}

func main() {

	host := os.Args[1]
	puerto := os.Args[2]

	// Arrancamos el servidor y atendemos conexiones entrantes
	fmt.Println("Arrancando el Registry, atendiendo en " + host + ":" + puerto)

	l, err := net.Listen("tcp", host+":"+puerto)

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
			fmt.Println("Error conectando con un visitante: ", err.Error())
		}

		// Imprimimos la dirección de conexión del cliente
		log.Println("Visitante " + c.RemoteAddr().String() + " conectado.")

		// Llamamos a la función de forma asíncrona y manejamos las conexiones de forma concurrente
		go manejoConexion(c)

	}

}

/* Función que procesa concurrentemente los registros o actualizaciones de los visitantes */
func manejoConexion(conexion net.Conn) {

	// Lectura de la opción elegida por el visitante
	opcion, err := bufio.NewReader(conexion).ReadBytes('\n')

	// Cerramos la conexión de los clientes que se han desconectado
	if err != nil {
		fmt.Println("Visitante " + conexion.RemoteAddr().String() + " desconectado.")
		conexion.Close()
		return
	}

	opcionElegida := strings.TrimSpace(string(opcion))

	// Lectura del id del visistante hasta final de línea
	id, err := bufio.NewReader(conexion).ReadBytes('\n')

	// Cerramos la conexión de los clientes que se han desconectado
	if err != nil {
		fmt.Println("Visitante " + conexion.RemoteAddr().String() + " desconectado.")
		conexion.Close()
		return
	}

	// Lectura del nombre del visitante hasta final de línea
	nombre, err := bufio.NewReader(conexion).ReadBytes('\n')

	// Cerramos la conexión de los clientes que se han desconectado
	if err != nil {
		fmt.Println("Visitante " + conexion.RemoteAddr().String() + " desconectado.")
		conexion.Close()
		return
	}

	// Lectura del password del visitante hasta final de línea
	password, err := bufio.NewReader(conexion).ReadBytes('\n')

	// Cerramos la conexión de los clientes que se han desconectado
	if err != nil {
		fmt.Println("Visitante " + conexion.RemoteAddr().String() + " desconectado.")
		conexion.Close()
		return
	}

	// Si se ha solicitado un registro de usuario
	if opcionElegida == "1" {

		// Imprimimos la información del visitante a registrar
		fmt.Println("Visitante a registrar -> ID: " + strings.TrimSpace(string(id)) +
			" | Nombre: " + strings.TrimSpace(string(nombre)) + " | Password: " + strings.TrimSpace(string(password)))

		// Si se ha solicitado editar/actualizar un perfil de usuario existente
	} else if opcionElegida == "2" {

		// Imprimimos la información del visitante a editar
		fmt.Println("Visitante a editar -> ID: " + strings.TrimSpace(string(id)) +
			" | Nombre: " + strings.TrimSpace(string(nombre)) + " | Password: " + strings.TrimSpace(string(password)))

	} else {
		conexion.Write([]byte("La opción elegida no es válida"))
		conexion.Close()
	}

	// Accedemos a la base de datos, empezando por abrir la conexión
	db, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/parque_atracciones")

	// Comprobamos que no haya error al conectarse
	if err != nil {
		panic("Error al conectarse con la BD: " + err.Error())
	}

	defer db.Close() // Para que siempre se cierre la conexión con la BD al finalizar el programa

	v := visitante{
		ID:       strings.TrimSpace(string(id)),
		Nombre:   strings.TrimSpace(string(nombre)),
		Password: strings.TrimSpace(string(password)),
		IdParque: strings.TrimSpace(string(id[0])),
	}

	// Si se ha solicitado un registro
	if opcionElegida == "1" {

		results, err := db.Query("SELECT * FROM visitante") // Devuelve los visitantes que hay registrados en la aplicación

		// Comrpobamos que no se produzcan errores al hacer la consulta
		if err != nil {
			panic("Error al hacer la consulta a la BD: " + err.Error())
		}

		defer results.Close() // Nos aseguramos de cerrar

		// Nos guardamos el nº actual de visitantes registrados en la aplicación
		visitantesActuales := 0
		for results.Next() {
			visitantesActuales += 1
		}

		// INSERTAMOS el nuevo visitante en la BD
		// Preparamos para prevenir inyecciones SQL
		sentenciaPreparada, err := db.Prepare("INSERT INTO visitante (id, nombre, contraseña, idParque) VALUES(?, ?, ?, ?)")
		if err != nil {
			panic("Error al preparar la sentencia de inserción: " + err.Error())
		}

		defer sentenciaPreparada.Close()

		// Ejecutar sentencia, un valor por cada '?'
		_, err = sentenciaPreparada.Exec(v.ID, v.Nombre, v.Password, v.IdParque)
		if err != nil {
			panic("Error al registrar el visitante: " + err.Error())
		}

		conexion.Write([]byte("Visitante registrado en el parque. Actualmente hay " + strconv.Itoa(visitantesActuales+1) + " visitantes registrados."))
		conexion.Close()

		// Actualizamos el número de visitantes que se encuentran en el parque
		// Preparamos para prevenir inyecciones SQL
		/*sentenciaPreparada, err = db.Prepare("UPDATE parque SET aforoActual = ? + 1 WHERE id = ?")
		if err != nil {
			panic("Error al preparar la sentencia de actualización del aforo del parque: " + err.Error())
		}

		defer sentenciaPreparada.Close()

		// Ejecutar sentencia, un valor por cada '?'
		_, err = sentenciaPreparada.Exec(visitantesActuales, "SDPark")
		if err != nil {
			panic("Error al modificar el aforo actual del parque: " + err.Error())
		}*/

	} else { // Si se ha solicitado una actualización

		results, err := db.Query("SELECT * FROM visitante WHERE id = ?", v.ID) // Devuelve los visitantes que hay registrados en la aplicación

		// Comprobamos que no se produzcan errores al hacer la consulta
		if err != nil {
			panic("Error al hacer la consulta de la información del visitante indicado: " + err.Error())
		}

		defer results.Close() // Nos aseguramos de cerrar

		// Si el ID del visitante indicado por el cliente existe
		if results.Next() {

			// MODIFICAMOS la información de dicho visitante en la BD
			// Preparamos para prevenir inyecciones SQL
			sentenciaPreparada, err := db.Prepare("UPDATE visitante SET nombre = ?, contraseña = ? WHERE id = ?")
			if err != nil {
				panic("Error al preparar la sentencia de actualización: " + err.Error())
			}

			defer sentenciaPreparada.Close()

			// Ejecutar sentencia, un valor por cada '?'
			_, err = sentenciaPreparada.Exec(v.Nombre, v.Password, v.ID)
			if err != nil {
				panic("Error al modificar el visitante: " + err.Error())
			}

			conexion.Write([]byte("Visitante actualizado correctamente"))
			conexion.Close()

		} else {
			conexion.Write([]byte("El id del visitante no existe"))
			conexion.Close()
		}

	}

	// Reiniciamos el proceso
	manejoConexion(conexion)

}
