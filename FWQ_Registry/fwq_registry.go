package main

import (
	"bufio"
	"crypto/rand"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
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

	go lanzarServidor(host) // El servidor funcionará de forma paralela y concurrente a los sockets

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

// INICIO BLOQUE RESPONSE

// Estructura con el formato de una respuesta http
type Response struct {
	Status        int         `json:"status"`
	Data          interface{} `json:"data"`
	Message       string      `json:"message"`
	contentType   string
	responseWrite http.ResponseWriter
}

/* Función que crea una respuesta por defecto para los clientes de la API REST */
func CreateDefaultResponse(rw http.ResponseWriter) Response {
	return Response{
		Status:        http.StatusOK,
		responseWrite: rw,
		contentType:   "application/json",
	}
}

/* Función que envía las respuestas a los clientes de la API REST */
func (resp *Response) Send() {
	resp.responseWrite.Header().Set("Content-Type", resp.contentType)
	resp.responseWrite.WriteHeader(resp.Status)

	// Marshall devuelve 2 valores: Los valores transformados en tipo byte y un error
	output, _ := json.Marshal(&resp) // Para responder con json
	//output, _ := xml.Marshal(&resp) // Para responder con xml
	//output, _ := yaml.Marshal(&resp) // Para responder con yaml
	fmt.Fprintln(resp.responseWrite, string(output))
}

/* Función que envía los datos solicitados al cliente API REST */
func SendData(rw http.ResponseWriter, data interface{}) {
	response := CreateDefaultResponse(rw)
	response.Data = data
	response.Send()
}

/* Función utilizada junto a la de abajo al momento de eliminar,
recuperar una fila o todas las filas de la BD y
que se produzcan errores para poder manejarlos. */
func (resp *Response) NotFound() {
	resp.Status = http.StatusNotFound
	resp.Message = "Resource Not Found"
}

/* Función que manda una respuesta indicando que el recurso solicitado no ha sido encontrado */
func SendNotFound(rw http.ResponseWriter) {
	response := CreateDefaultResponse(rw)
	response.NotFound()
	response.Send()
}

/* Función utilizada junto a la de abajo al momento de eliminar,
recuperar una fila o todas las filas de la BD y
que se produzcan errores para poder manejarlos. */
func (resp *Response) YaExiste() {
	resp.Status = http.StatusNoContent
	resp.Message = "El visitante ya estaba registrado"
}

/* Función que manda una respuesta indicando que el recurso solicitado no ha sido encontrado */
func SendYaExiste(rw http.ResponseWriter) {
	response := CreateDefaultResponse(rw)
	response.YaExiste()
	response.Send()
}

/* Función utilizada junto a la de abajo al momento de insertar o
actualizar una fila de la BD y que se produzcan errores para poder manejarlos. */
func (resp *Response) UnprocessableEntity() {
	resp.Status = http.StatusUnprocessableEntity
	resp.Message = "UnprocessableEntity Not Found"
}

/* Función que manda una respuesta indicando que la entidad recibida no es procesable */
func SendUnprocessableEntity(rw http.ResponseWriter) {
	response := CreateDefaultResponse(rw)
	response.UnprocessableEntity()
	response.Send()
}

// FIN BLOQUE RESPONSE

// INICIO BLOQUE FUNCIONES VISITANTES

//Construir visitante
/*func NewUser(username, password, email string) *visitante {
	user := &visitante{Username: username, Password: password, Email: email}
	return user
}

//Obtener un visitante
func GetUser(id int) (*visitante, error) {
	user := NewUser("", "", "")

	sql := "SELECT id, username, password, email FROM users WHERE id=?"
	if rows, err := db.Query(sql, id); err != nil {
		return nil, err
	} else {
		for rows.Next() {
			rows.Scan(&user.Id, &user.Username, &user.Password, &user.Email)
		}
		return user, nil
	}
}*/

// FIN BLOQUE FUNCIONES VISITANTES

// INICIO BLOQUE HANDLERS

/* Función que devuelve un visitante dependiendo del id pasado en el request */
/*func getVisitanteByRequest(r *http.Request) (visitante, error) {
	//Obtener ID
	vars := mux.Vars(r)
	userId, _ := strconv.Atoi(vars["id"])

	if user, err := GetVisitante(userId); err != nil {
		return *user, err
	} else {
		return *user, nil
	}
}*/

/* Función manejadora para la creación del perfil de un visitante */
func crearPerfil(rw http.ResponseWriter, r *http.Request) {

	//Obtener ID
	vars := mux.Vars(r)
	userId := vars["id"]

	v := visitante{}
	v.ID = userId
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&v); err != nil {
		SendUnprocessableEntity(rw)
	} else {

		// Insertamos el nuevo visitante en la BD

		// Accedemos a la base de datos, empezando por abrir la conexión
		db, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/parque_atracciones")

		// Comprobamos que no haya error al conectarse
		if err != nil {
			panic("Error al conectarse con la BD: " + err.Error())
		}

		defer db.Close() // Para que siempre se cierre la conexión con la BD al finalizar el programa

		results, err := db.Query("SELECT * FROM visitante WHERE id = ?", v.ID) // Comprueba si el visitante está registrado en la aplicación

		// Comprobamos que no se produzcan errores al hacer la consulta
		if err != nil {
			panic("Error al hacer la consulta de la información del visitante indicado: " + err.Error())
		}

		defer results.Close() // Nos aseguramos de cerrar

		// Si el visitante ya se había registrado
		if results.Next() {

			//conexion.Write([]byte("ERROR: El visitante ya estaba registrado en la aplicación"))
			//conexion.Close()
			//SendData(rw, v) //TODO:

			RegistroLog(db, r.RemoteAddr, v.ID, "Error", "El visitante "+v.ID+" ya estaba registrado en la aplicación") // Registramos el evento de log

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
			_, err = sentenciaPreparada.Exec(v.ID, v.Nombre, HashPassword(v.Password), string(v.ID[0]))
			if err != nil {
				panic("Error al registrar el visitante: " + err.Error())
			}

			//conexion.Write([]byte("Visitante registrado en el parque. Actualmente hay " + strconv.Itoa(visitantesActuales+1) + " visitantes registrados."))
			//conexion.Close()

			RegistroLog(db, r.RemoteAddr, v.ID, "Alta", "Visitante "+v.ID+" registrado correctamente") // Registramos el evento de log

			SendData(rw, v) //TODO:

		}

	}

}

/* Función manejadora para la modificación del perfil de un visitante */
func editarPerfil(rw http.ResponseWriter, r *http.Request) {

	// Accedemos a la base de datos, empezando por abrir la conexión
	db, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/parque_atracciones")

	// Comprobamos que no haya error al conectarse
	if err != nil {
		panic("Error al conectarse con la BD: " + err.Error())
	}

	defer db.Close() // Para que siempre se cierre la conexión con la BD al finalizar el programa

	//Obtener ID
	vars := mux.Vars(r)
	userId := vars["id"]

	/*if v, err := getUserByRequest(r); err != nil {
		SendNotFound(rw)
	} else {
		userId = v.ID
	}*/

	results, err := db.Query("SELECT * FROM visitante WHERE id = ?", userId) // Comprueba si el visitante está registrado en la aplicación

	// Comprobamos que no se produzcan errores al hacer la consulta
	if err != nil {
		panic("Error al hacer la consulta de la información del visitante indicado: " + err.Error())
	}

	defer results.Close() // Nos aseguramos de cerrar

	// Si el ID del visitante indicado por el cliente existe
	if results.Next() {

		v := visitante{}
		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&v); err != nil {
			SendUnprocessableEntity(rw)
		} else {

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

			//conexion.Write([]byte("Visitante actualizado correctamente"))
			//conexion.Close()

			RegistroLog(db, r.RemoteAddr, v.ID, "Modificación", "Visitante "+v.ID+" actualizado correctamente") // Registramos el evento de log

			//v.ID = userId
			SendData(rw, v)
		}

	} else {
		//conexion.Write([]byte("El id del visitante no existe"))
		//conexion.Close()
		SendNotFound(rw)
	}

}

// FIN BLOQUE HANDLERS

/* Función que se encarga de arrancar el servidor API REST */
func lanzarServidor(host string) {

	// IMPLEMENTAMOS EL API REST
	// Rutas
	mux := mux.NewRouter()

	// Responder al cliente
	mux.HandleFunc("/crear/{id:[A-Za-z0-9_]+}", crearPerfil).Methods("POST")
	mux.HandleFunc("/editar/{id:[A-Za-z0-9_]+}", editarPerfil).Methods("PUT")

	// Servidor
	fmt.Println("Servidor corriendo en https://" + host + ":3000") //TODO:
	log.Fatal(http.ListenAndServe(":3000", mux))

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
func RegistroLog(db *sql.DB, ipPuerto, idVisitante, accion, descripcion string) {

	// Añadimos el evento de log de error al visitante
	sentenciaPreparada, err := db.Prepare("UPDATE visitante SET ultimoEvento = ? WHERE id = ?")
	if err != nil {
		panic("Error al preparar la sentencia de inserción: " + err.Error())
	}

	defer sentenciaPreparada.Close()

	var eventoLog string // Variable donde vamos a guardar la información de log que le vamos a pasar a la BD

	dateTime := time.Now().Format("2006-01-02 15:04:05") // Fecha y hora del evento con formato personalizado
	//ipVisitante := conexion.RemoteAddr().String()
	ipVisitante := ipPuerto          // IP y puerto de quién ha provocado el evento
	accionRealizada := accion        // Que acción se realiza
	descripcionEvento := descripcion // Parámetros o descripción del evento

	eventoLog += dateTime + " | "
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

			RegistroLog(db, conexion.RemoteAddr().String(), v.ID, "Error", "El visitante "+v.ID+" ya estaba registrado en la aplicación") // Registramos el evento de log

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

			RegistroLog(db, conexion.RemoteAddr().String(), v.ID, "Alta", "Visitante "+v.ID+" registrado correctamente") // Registramos el evento de log

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

			RegistroLog(db, conexion.RemoteAddr().String(), v.ID, "Modificación", "Visitante "+v.ID+" actualizado correctamente") // Registramos el evento de log

		} else {
			conexion.Write([]byte("El id del visitante no existe"))
			conexion.Close()
		}

	}

	// Reiniciamos el proceso
	manejoConexion(conexion)

}
