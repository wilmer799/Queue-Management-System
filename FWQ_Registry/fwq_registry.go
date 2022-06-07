package main

import (
	"bufio"
	"crypto/rand"
	"crypto/tls"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
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
	IdEnParque   string `json:"idEnParque"`
	UltimoEvento string `json:"ultimoEvento"`
	Parque       string `json:"parqueAtracciones"`
}

func main() {

	host := os.Args[1]
	puerto := os.Args[2]

	cert, err := tls.LoadX509KeyPair("cert/cert.pem", "cert/key.pem")
	if err != nil {
		log.Fatal(err)
	}

	config := tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAnyClientCert,
	}

	config.Rand = rand.Reader

	// Arrancamos el servidor y atendemos conexiones entrantes
	fmt.Println("Arrancando el Registry, atendiendo en " + host + ":" + puerto)

	//l, err := net.Listen("tcp", host+":"+puerto) // CONEXIONES INSEGURAS
	l, err := tls.Listen("tcp", host+":"+puerto, &config) // CONEXIONES SEGURAS
	if err != nil {
		log.Fatal("Error escuchando", err.Error())
	}

	// Cerramos el listener cuando se cierra la aplicación
	defer l.Close()

	// Bucle infinito hasta la salida del programa
	for {

		// Atendemos conexiones entrantes
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error conectando con un visitante: ", err.Error())
			continue
		}

		// Imprimimos la dirección de conexión del cliente
		log.Println("Visitante " + conn.RemoteAddr().String() + " conectado.")

		// Llamamos a la función de forma asíncrona y manejamos las conexiones de forma concurrente
		go manejoConexion(conn)

	}

}

/* Función que genera y devuelve el hash de la contraseña pasada por parámetro */
func HashPassword(password string) string {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	return string(hash)

}

/* Función que almacena los registros de auditoría en la tabla visitante */
func registroLog(db *sql.DB, conexion net.Conn, idVisitante, accion, descripcion string) {

	// Añadimos el evento de log de error al visitante
	sentenciaPreparada, err := db.Prepare("UPDATE visitante SET ultimoEvento = ? WHERE id = ?")
	if err != nil {
		panic("Error al preparar la sentencia de inserción: " + err.Error())
	}

	defer sentenciaPreparada.Close()

	var eventoLog string // Variable donde vamos a guardar la información de log que le vamos a pasar a la BD

	dateTime := time.Now()                        // Fecha y hora del evento
	ipVisitante := conexion.RemoteAddr().String() // IP del visitante
	accionRealizada := accion                     // Que acción se realiza
	descripcionEvento := descripcion              // Parámetros o descripción del evento

	eventoLog += dateTime.String() + " | "
	eventoLog += ipVisitante + " | "
	eventoLog += accionRealizada + " | "
	eventoLog += descripcionEvento

	_, err = sentenciaPreparada.Exec(eventoLog, idVisitante)
	if err != nil {
		panic("Error al registrar el evento de log: " + err.Error())
	}

}

/* Función que procesa concurrentemente los registros o actualizaciones de los visitantes */
func manejoConexion(conexion net.Conn) {

	// Lectura de la opción elegida por el visitante
	opcion, err := bufio.NewReader(conexion).ReadBytes('\n')

	// Cerramos la conexión de los clientes que se han desconectado
	if err != nil {
		log.Println("Visitante " + conexion.RemoteAddr().String() + " desconectado.")
		conexion.Close()
		return
	}

	opcionElegida := strings.TrimSpace(string(opcion))

	// Lectura del id del visistante hasta final de línea
	id, err := bufio.NewReader(conexion).ReadBytes('\n')

	// Cerramos la conexión de los clientes que se han desconectado
	if err != nil {
		log.Println("Visitante " + conexion.RemoteAddr().String() + " desconectado.")
		conexion.Close()
		return
	}

	// Lectura del nombre del visitante hasta final de línea
	nombre, err := bufio.NewReader(conexion).ReadBytes('\n')

	// Cerramos la conexión de los clientes que se han desconectado
	if err != nil {
		log.Println("Visitante " + conexion.RemoteAddr().String() + " desconectado.")
		conexion.Close()
		return
	}

	// Lectura del password del visitante hasta final de línea
	password, err := bufio.NewReader(conexion).ReadBytes('\n')

	// Cerramos la conexión de los clientes que se han desconectado
	if err != nil {
		log.Println("Visitante " + conexion.RemoteAddr().String() + " desconectado.")
		conexion.Close()
		return
	}

	// Si se ha solicitado un registro de usuario
	if opcionElegida == "1" {

		// Imprimimos la información del visitante a registrar
		fmt.Println("Visitante a registrar -> ID: " + strings.TrimSpace(string(id)) +
			" | Nombre: " + strings.TrimSpace(string(nombre)) + " | Password: " + strings.TrimSpace(string(password)))

		// Si se ha solicitado editar/actualizar el perfil de usuario existente
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
		ID:         strings.TrimSpace(string(id)),
		Nombre:     strings.TrimSpace(string(nombre)),
		Password:   strings.TrimSpace(string(password)),
		IdEnParque: strings.TrimSpace(string(id[0])),
	}

	// Si se ha solicitado un registro
	if opcionElegida == "1" {

		results, err := db.Query("SELECT * FROM visitante WHERE id = ?", v.ID) // Comprueba si el visitante está registrado en la aplicación

		// Comprobamos que no se produzcan errores al hacer la consulta
		if err != nil {
			panic("Error al hacer la consulta de la información del visitante indicado: " + err.Error())
		}

		defer results.Close() // Nos aseguramos de cerrar

		// Si el visitante ya se había registrado
		if results.Next() {

			conexion.Write([]byte("ERROR: El visitante ya estaba registrado en la aplicación"))
			conexion.Close()

			registroLog(db, conexion, v.ID, "Error", "El visitante ya está registrado en la aplicación") // Registramos el evento de log

		} else { // Si es un nuevo usuario

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
			sentenciaPreparada, err := db.Prepare("INSERT INTO visitante (id, nombre, contraseña, idEnParque) VALUES(?, ?, ?, ?)")
			if err != nil {
				panic("Error al preparar la sentencia de inserción: " + err.Error())
			}

			defer sentenciaPreparada.Close()

			// Ejecutar sentencia, un valor por cada '?'
			_, err = sentenciaPreparada.Exec(v.ID, v.Nombre, HashPassword(v.Password), v.IdEnParque)
			if err != nil {
				panic("Error al registrar el visitante: " + err.Error())
			}

			conexion.Write([]byte("Visitante registrado en el parque. Actualmente hay " + strconv.Itoa(visitantesActuales+1) + " visitantes registrados."))
			conexion.Close()

		}

	} else { // Si se ha solicitado una actualización

		results, err := db.Query("SELECT * FROM visitante WHERE id = ?", v.ID) // Comprueba si el visitante está registrado en la aplicación

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
			_, err = sentenciaPreparada.Exec(v.Nombre, HashPassword(v.Password), v.ID)
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
