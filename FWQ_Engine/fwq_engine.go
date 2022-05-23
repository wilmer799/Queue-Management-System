package main

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime/pprof"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/segmentio/kafka-go"
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
	Parque       string `json:"parqueAtracciones"`
}

/*
* Estructura de las atracciones
 */
type atraccion struct {
	ID           string `json:"id"`
	TCiclo       int    `json:"tciclo"`
	NVisitantes  int    `json:"nvisitantes"`
	Posicionx    int    `json:"posicionx"`
	Posiciony    int    `json:"posiciony"`
	TiempoEspera int    `json:"tiempoEspera"`
	Parque       string `json:"parqueAtracciones"`
}

//var visitantesDelEngine []string

/*
 * @Description : Función main de fwq_engine
 * @Author : Wilmer Fabricio Bravo Shuira
 */
func main() {

	IpKafka := os.Args[1]
	PuertoKafka := os.Args[2]
	numeroVisitantes := os.Args[3]
	IpFWQWaiting := os.Args[4]
	PuertoWaiting := os.Args[5]

	fmt.Println("Arrancando un engine que atiende peticiones por " + IpKafka + ":" + PuertoKafka + ", limita el parque a " + numeroVisitantes + " visitantes y manda peticiones a un servidor de tiempos de espera situado en " + IpFWQWaiting + ":" + PuertoWaiting + ".\n")

	//Creamos el topic...Cambiar la Ipkafka en la función principal
	//Si no se ejecuta el programa, se cierra el kafka?
	crearTopics(IpKafka, PuertoKafka, "peticiones")
	crearTopics(IpKafka, PuertoKafka, "respuesta-login")
	crearTopics(IpKafka, PuertoKafka, "movimiento-mapa")

	// Visitantes, atracciones que se encuentran en la BD
	var visitantesRegistrados []visitante
	var conn = conexionBD()
	maxVisitantes, _ := strconv.Atoi(numeroVisitantes)
	establecerMaxVisitantes(conn, maxVisitantes)

	//Para empezar con el kafka
	ctx := context.Background()
	go consumidorEngine(IpKafka, PuertoKafka, ctx, maxVisitantes)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf("captured %v, stopping profiler and exiting..", sig)
			mensaje := "Engine no disponible"
			mensajeJson, err := json.Marshal(mensaje)
			if err != nil {
				fmt.Println("Error a la hora de codificar el mensaje: %v", err)
			}
			productorMapa(IpKafka, PuertoKafka, ctx, mensajeJson)

			/*for i := 0; i < len(visitantesDelEngine); i++ {

				// Al cerrar el parque tenemos que sacar a los visitantes de este
				sentenciaPreparada, err := conn.Prepare("UPDATE visitante SET dentroParque = 0, posicionx = 0, posiciony = 0, destinox = -1, destinoy = -1 WHERE id = ?")
				if err != nil {
					panic("Error al preparar la sentencia de modificación: " + err.Error())
				}

				// Ejecutar sentencia, un valor por cada '?'
				_, err = sentenciaPreparada.Exec(visitantesDelEngine[i])
				if err != nil {
					panic("Error al expulsar a los visitantes del parque: " + err.Error())
				}

				sentenciaPreparada.Close()
			}*/

			fmt.Println()
			fmt.Println("Engine apagado manualmente")
			pprof.StopCPUProfile()
			os.Exit(1)
		}
	}()

	for {
		visitantesRegistrados, _ = obtenerVisitantesBD(conn) // Obtenemos los visitantes registrados actualmente
		fmt.Println("*********** FUN WITH QUEUES RESORT ACTIVITY MAP *********")
		fmt.Println("ID   	" + "		Nombre      " + "	Pos.      " + "	Destino      " + "	DentroParque")
		//Hay que usar la función TrimSpace porque al parecer tras la obtención de valores de BD se agrega un retorno de carro a cada variable
		//Mostramos los visitantes registrados en la aplicación actualmente
		for i := 0; i < len(visitantesRegistrados); i++ {
			fmt.Println(strings.TrimSpace(visitantesRegistrados[i].ID) + " 		" + strings.TrimSpace(visitantesRegistrados[i].Nombre) +
				"    " + "	(" + strings.TrimSpace(strconv.Itoa(visitantesRegistrados[i].Posicionx)) + "," + strings.TrimSpace(strconv.Itoa(visitantesRegistrados[i].Posiciony)) +
				")" + "	    " + "    (" + strings.TrimSpace(strconv.Itoa(visitantesRegistrados[i].Destinox)) + "," + strings.TrimSpace(strconv.Itoa(visitantesRegistrados[i].Destinoy)) +
				")" + "	             " + strings.TrimSpace(strconv.Itoa(visitantesRegistrados[i].DentroParque)))
		}

		fmt.Println() // Para mejorar la visualización

		// Cada X segundos se conectará al servidor de tiempos para actualizar los tiempos de espera de las atracciones
		time.Sleep(time.Duration(5 * time.Second))
		atracciones, _ := obtenerAtraccionesBD(conn) // Obtenemos las atracciones actualizadas
		conexionTiempoEspera(conn, IpFWQWaiting, PuertoWaiting, atracciones)

		fmt.Println() // Para mejorar la visualización

	}

}

/* Función que comprueba si el aforo del parque está completo o no */
func parqueLleno(db *sql.DB, maxAforo int) bool {

	var lleno bool = false

	// Comprobamos si las credenciales de acceso son válidas
	results, err := db.Query("SELECT * FROM visitante WHERE dentroParque = 1")

	if err != nil {
		fmt.Println("Error al hacer la consulta sobre la BD para comprobar el aforo: " + err.Error())
	}

	visitantesDentroParque := 0 // Variable en la que vamos a almacenar el número de visitantes que se encuentran en el parque

	// Vamos recorriendo las filas devueltas para obtener el nº de visitanes dentro del parque
	for results.Next() {
		visitantesDentroParque++
	}

	results.Close() // Cerramos la conexión a la BD

	// Si el aforo está al completo
	if visitantesDentroParque >= maxAforo {
		lleno = true
	}

	return lleno

}

/* Función que permite eliminar un element de un slice */
/*func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	return s[:len(s)-1]
}*/

/* Función que recibe del gestor de colas las credenciales de los visitantes que quieren iniciar sesión para entrar en el parque */
func consumidorEngine(IpKafka, PuertoKafka string, ctx context.Context, maxVisitantes int) {

	//Accediendo a la base de datos
	//Abrimos la conexion con la base de datos
	db, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/parque_atracciones")
	//Si la conexión falla mostrara este error
	if err != nil {
		panic(err.Error())
	}
	//Cierra la conexion con la bd
	defer db.Close()

	direccionKafka := IpKafka + ":" + PuertoKafka
	//Configuración de lector de kafka
	conf := kafka.ReaderConfig{
		Brokers:     []string{direccionKafka},
		Topic:       "peticiones", //Topico que hemos creado
		GroupID:     "visitantes",
		StartOffset: kafka.LastOffset,
	}

	reader := kafka.NewReader(conf)

	for {

		m, err := reader.ReadMessage(context.Background())

		if err != nil {
			fmt.Println("Ha ocurrido algún error a la hora de conectarse con el kafka", err)
		}

		//fmt.Println("Petición recibida: " + string(m.Value))

		cadenaPeticion := strings.Split(string(m.Value), ":")

		alias := cadenaPeticion[0]
		peticion := cadenaPeticion[1]
		destino := strings.Split(cadenaPeticion[2], ",")
		destinoX, _ := strconv.Atoi(strings.TrimSpace(destino[0]))
		destinoY, _ := strconv.Atoi(strings.TrimSpace(destino[1]))

		v := visitante{
			ID:       strings.TrimSpace(alias),
			Password: strings.TrimSpace(peticion),
			Destinox: destinoX,
			Destinoy: destinoY,
		}

		// Comprobamos si lo enviado son credenciales de acceso en cuyo caso se trata de una petición de login
		results, err := db.Query("SELECT * FROM visitante WHERE id = ? and contraseña = ?", v.ID, v.Password)

		if err != nil {
			fmt.Println("Error al hacer la consulta sobre la BD para el login: " + err.Error())
		}

		var respuesta string = ""

		// Si las credenciales coinciden con las de un visitante registrado en la BD y el parque no está lleno
		if results.Next() && !parqueLleno(db, maxVisitantes) {

			// Actualizamos el estado del visitante en la BD
			sentenciaPreparada, err := db.Prepare("UPDATE visitante SET dentroParque = 1, destinox = ?, destinoy = ? WHERE id = ?")
			if err != nil {
				panic("Error al preparar la sentencia de modificación: " + err.Error())
			}

			// Ejecutar sentencia, un valor por cada '?'
			_, err = sentenciaPreparada.Exec(v.Destinox, v.Destinoy, v.ID)
			if err != nil {
				panic("Error al actualizar el estado del visitante respecto al parque: " + err.Error())
			}

			// Nos guardamos los visitantes del parque asociados a este engine
			//visitantesDelEngine = append(visitantesDelEngine, v.ID)

			respuesta += alias + ":" + "Acceso concedido"
			productorLogin(IpKafka, PuertoKafka, ctx, respuesta)

			sentenciaPreparada.Close()

		} else if peticion == "IN" || peticion == "N" || peticion == "S" || peticion == "W" || peticion == "E" || peticion == "NW" ||
			peticion == "NE" || peticion == "SW" || peticion == "SE" { // Si se nos ha mandado un movimiento

			// Comprobamos que el alias pertenezca a un visitante que se encuentra en el parque
			results, err := db.Query("SELECT * FROM visitante WHERE id = ?", v.ID)

			if err != nil {
				fmt.Println("Error al hacer la consulta sobre la BD para el login: " + err.Error())
			}

			// Actualizamos el estado el destino del visitante en la BD
			sentenciaPreparada, err := db.Prepare("UPDATE visitante SET destinox = ?, destinoy = ? WHERE id = ?")
			if err != nil {
				panic("Error al preparar la sentencia de modificación de destino: " + err.Error())
			}

			// Ejecutar sentencia, un valor por cada '?'
			_, err = sentenciaPreparada.Exec(v.Destinox, v.Destinoy, v.ID)
			if err != nil {
				panic("Error al actualizar el destino del visitante: " + err.Error())
			}

			sentenciaPreparada.Close()

			if results.Next() {

				var mapa [20][20]string
				visitantesParque, _ := obtenerVisitantesParque(db)             // Obtenemos los visitantes del parque actualizados
				mueveVisitante(db, alias, peticion, visitantesParque)          // Movemos al visitante en base al movimiento recibido
				visitantesParqueActualizados, _ := obtenerVisitantesParque(db) // Obtenemos los visitantes del parque actualizados
				// Preparamos el mapa a enviar a los visitantes que se encuentra en el parque
				atracciones, _ := obtenerAtraccionesBD(db) // Obtenemos las atracciones actualizadas
				mapaActualizado := asignacionPosiciones(visitantesParqueActualizados, atracciones, mapa)
				var representacion string
				for i := 0; i < len(mapaActualizado); i++ {
					for j := 0; j < len(mapaActualizado); j++ {
						if j == 19 {
							representacion = representacion + mapaActualizado[i][j] + "\n"
						} else {
							representacion = representacion + mapaActualizado[i][j]
						}
					}
				}
				//Convertimos el mapaActualizado a formato jSON
				//Esta función devuelve un array de byte
				mapaJson, err := json.Marshal(representacion)
				//En formato jSon tiene encuenta el salto de linea por lo que hay que ver si al decodificarlo se quita
				if err != nil {
					fmt.Println("Error a la hora de codificar el mapa: %v", err)
				}
				productorMapa(IpKafka, PuertoKafka, ctx, mapaJson) // Mandamos el mapa actualizado a los visitantes que se encuentran en el parque
				results.Close()

			} else { // Si el alias no pertenece a un visitante del parque
				respuesta += alias + ":" + "Parque cerrado"
				productorLogin(IpKafka, PuertoKafka, ctx, respuesta)
				results.Close()
			}

		} else if peticion == "OUT" { // Si se nos ha solicitado una salida del parque

			// Sacamos del parque al visitante y reinciamos tanto su posición actual como su destino
			sentenciaPreparada, err := db.Prepare("UPDATE visitante SET dentroParque = 0, posicionx = 0, posiciony = 0, destinox = -1, destinoy = -1 WHERE id = ?")
			if err != nil {
				panic("Error al preparar la sentencia de modificación: " + err.Error())
			}

			// Ejecutar sentencia, un valor por cada '?'
			_, err = sentenciaPreparada.Exec(v.ID)
			if err != nil {
				panic("Error al actualizar el estado del visitante respecto al parque: " + err.Error())
			}

			/*encontrado := false

			for i := 0; i < len(visitantesDelEngine) && !encontrado; i++ {

				if visitantesDelEngine[i] == v.ID {
					visitantesDelEngine = remove(visitantesDelEngine, i)
					encontrado = true
				}

			}*/

			sentenciaPreparada.Close()

		} else { // Si las credenciales enviadas para iniciar sesión no son válidas

			if parqueLleno(db, maxVisitantes) {
				respuesta += alias + ":" + "Aforo al completo"
				productorLogin(IpKafka, PuertoKafka, ctx, respuesta)
			} else {
				respuesta += alias + ":" + "Parque cerrado"
				productorLogin(IpKafka, PuertoKafka, ctx, respuesta)
			}
		}

		results.Close()

	}

}

/* Función que envía el mensaje de respuesta a la petición de login de un visitante */
func productorLogin(IpBroker, PuertoBroker string, ctx context.Context, respuesta string) {

	var brokerAddress string = IpBroker + ":" + PuertoBroker
	var topic string = "respuesta-login"

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:          []string{brokerAddress},
		Topic:            topic,
		CompressionCodec: kafka.Snappy.Codec(),
	})

	err := w.WriteMessages(ctx, kafka.Message{
		Key:   []byte("Key-Login"),
		Value: []byte(respuesta),
	})

	if err != nil {
		fmt.Println("No se puede mandar el mensaje de respuesta a la petición de login: " + err.Error())
	}

}

/*
* Función que abre una conexion con la bd
 */
func conexionBD() *sql.DB {
	//Accediendo a la base de datos
	/*****Flate blod **/
	//Abrimos la conexion con la base de datos
	db, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/parque_atracciones")
	//Si la conexión falla mostrara este error
	if err != nil {
		panic(err.Error())
	}
	//Cierra la conexion con la bd
	//defer db.Close()
	return db
}

/*
* Función que obtiene todos los visitantes que se encuentran la BD
* @return []visitante : Arrays de los visitantes obtenidos en la sentencia
* @return error : Error en caso de que no se haya podido obtener ninguno
 */
func obtenerVisitantesBD(db *sql.DB) ([]visitante, error) {

	//Ejecutamos la sentencia
	results, err := db.Query("SELECT * FROM visitante")

	if err != nil {
		return nil, err //devolvera nil y error en caso de que no se pueda hacer la consulta
	}

	//Cerramos la base de datos
	defer results.Close()

	//Declaramos el array de visitantes
	var visitantes []visitante

	//Recorremos los resultados obtenidos por la consulta
	for results.Next() {

		//Variable donde guardamos la información de cada una filas de la sentencia
		var fwq_visitante visitante

		if err := results.Scan(&fwq_visitante.ID, &fwq_visitante.Nombre,
			&fwq_visitante.Password, &fwq_visitante.Posicionx,
			&fwq_visitante.Posiciony, &fwq_visitante.Destinox, &fwq_visitante.Destinoy,
			&fwq_visitante.DentroParque, &fwq_visitante.IdEnParque, &fwq_visitante.Parque); err != nil {
			return visitantes, err
		}

		//Vamos añadiendo los visitantes al array
		visitantes = append(visitantes, fwq_visitante)
	}

	if err = results.Err(); err != nil {
		return visitantes, err
	}

	return visitantes, nil

}

/*
* Función que obtiene todos los visitantes que se encuentran en el parque de la BD
* @return []visitante : Arrays de los visitantes obtenidos en la sentencia
* @return error : Error en caso de que no se haya podido obtener ninguno
 */
func obtenerVisitantesParque(db *sql.DB) ([]visitante, error) {

	//Ejecutamos la sentencia
	results, err := db.Query("SELECT * FROM visitante WHERE dentroParque = 1")

	if err != nil {
		return nil, err //devolvera nil y error en caso de que no se pueda hacer la consulta
	}

	//Cerramos la base de datos
	defer results.Close()

	//Declaramos el array de visitantes
	var visitantes []visitante

	//Recorremos los resultados obtenidos por la consulta
	for results.Next() {

		//Variable donde guardamos la información de cada una filas de la sentencia
		var fwq_visitante visitante

		if err := results.Scan(&fwq_visitante.ID, &fwq_visitante.Nombre,
			&fwq_visitante.Password, &fwq_visitante.Posicionx,
			&fwq_visitante.Posiciony, &fwq_visitante.Destinox, &fwq_visitante.Destinoy,
			&fwq_visitante.DentroParque, &fwq_visitante.IdEnParque, &fwq_visitante.Parque); err != nil {
			return visitantes, err
		}

		//Vamos añadiendo los visitantes al array
		visitantes = append(visitantes, fwq_visitante)
	}

	if err = results.Err(); err != nil {
		return visitantes, err
	}

	return visitantes, nil

}

/*
* Función que obtiene las atracciones del parque
* @return []atraccion : Array con las atracciones del parque
* @return error : Error en caso de que no se ha podido obtener las atracciones
 */
func obtenerAtraccionesBD(db *sql.DB) ([]atraccion, error) {

	//Ejecutamos la sentencia
	results, err := db.Query("SELECT * FROM atraccion")

	if err != nil {
		return nil, err //devolvera nil y error en caso de que no se pueda hacer la consulta
	}

	// Nos aseguramos de que se cierre la base de datos
	defer results.Close()

	//Declaramos el array de atracciones
	var atraccionesParque []atraccion

	//Recorremos los resultados obtenidos por la consulta
	for results.Next() {

		//Variable donde guardamos la información de cada una filas de la sentencia
		var fwq_atraccion atraccion

		if err := results.Scan(&fwq_atraccion.ID, &fwq_atraccion.TCiclo,
			&fwq_atraccion.NVisitantes, &fwq_atraccion.Posicionx,
			&fwq_atraccion.Posiciony, &fwq_atraccion.TiempoEspera,
			&fwq_atraccion.Parque); err != nil {
			return atraccionesParque, err
		}

		//Vamos añadiendo las atracciones al array
		atraccionesParque = append(atraccionesParque, fwq_atraccion)

	}

	if err = results.Err(); err != nil {
		return atraccionesParque, err
	}

	return atraccionesParque, nil

}

/*
* Función que forma el mapa del parque conteniendo a los visitantes y las atracciones
* @return [20][20]string : Matriz bidimensional representando el mapa
 */
func asignacionPosiciones(visitantes []visitante, atracciones []atraccion, mapa [20][20]string) [20][20]string {

	//Asignamos los id de los visitantes
	for i := 0; i < len(mapa); i++ {
		for j := 0; j < len(mapa[i]); j++ {
			for k := 0; k < len(visitantes); k++ {
				if i == visitantes[k].Posicionx && j == visitantes[k].Posiciony && visitantes[k].DentroParque == 1 {
					mapa[i][j] = visitantes[k].IdEnParque + "|"
				}
			}
		}
	}

	//Asignamos los valores de tiempo de espera de las atracciones
	for i := 0; i < len(mapa); i++ {
		for j := 0; j < len(mapa[i]); j++ {
			for k := 0; k < len(atracciones); k++ {
				if i == atracciones[k].Posicionx && j == atracciones[k].Posiciony {
					mapa[i][j] = strconv.Itoa(atracciones[k].TiempoEspera) + "|"
				}
			}
		}
	}

	// Las casillas del mapa que no tengan ni visitantes ni atracciones las representamos con una guión
	for i := 0; i < len(mapa); i++ {
		for j := 0; j < len(mapa[i]); j++ {
			if len(mapa[i][j]) == 0 {
				mapa[i][j] = "-" + "|"
			}
		}
	}
	return mapa
}

/*
* Función que se conecta al servidor de tiempos para obtener los tiempos de espera actualizados
 */
func conexionTiempoEspera(db *sql.DB, IpFWQWating, PuertoWaiting string, atracciones []atraccion) {

	fmt.Println() // Por limpieza
	fmt.Println("***Conexión con el servidor de tiempo de espera***")
	//fmt.Println("Arrancando el engine para atender los tiempos en el puerto: " + IpFWQWating + ":" + PuertoWaiting)
	var connType string = "tcp"
	conn, err := net.Dial(connType, IpFWQWating+":"+PuertoWaiting)

	if err != nil {
		fmt.Println("ERROR: El servidor de tiempos de espera no está disponible", err.Error())
	} else {

		fmt.Println("***Actualizando los tiempos de espera***")
		fmt.Println() // Por limpieza

		var infoAtracciones string = ""

		for i := 0; i < len(atracciones); i++ {
			infoAtracciones += atracciones[i].ID + ":"
			infoAtracciones += strconv.Itoa(atracciones[i].TCiclo) + ":"
			infoAtracciones += strconv.Itoa(atracciones[i].NVisitantes) + ":"
			infoAtracciones += strconv.Itoa(atracciones[i].TiempoEspera) + "|"
		}

		infoAtracciones += "\n" // Le añadimos el salto de línea porque los sockets los estamos leyendo hasta final de línea

		fmt.Println("Enviando información de las atracciones...")

		conn.Write([]byte(infoAtracciones))                        // Mandamos el id:tiempoCiclo:nºvisitantes de cada atracción en un string
		tiemposEspera, _ := bufio.NewReader(conn).ReadString('\n') // Obtenemos los tiempos de espera actualizados

		if tiemposEspera != "" {

			log.Println("Tiempos de espera actualizados: " + tiemposEspera)

			arrayTiemposEspera := strings.Split(tiemposEspera, "|")

			// Actualizamos los tiempos de espera de las atracciones en la BD
			actualizaTiemposEsperaBD(db, arrayTiemposEspera)

		} else {
			log.Println("Servidor de tiempos no disponible.")
		}

	}

}

/*
* Función que establece el aforo maximo permitido en el parque de atracciones
 */
func establecerMaxVisitantes(db *sql.DB, numero int) {

	//Ejecutamos la sentencia
	results, err := db.Query("SELECT * FROM parque")

	if err != nil {
		panic("Error al hacer la consulta del parque" + err.Error()) //devolvera nil y error en caso de que no se pueda hacer la consulta
	}

	//Cerramos la base de datos
	defer results.Close()

	//Recorremos los resultados obtenidos por la consulta
	if results.Next() {

		//Variable donde guardamos la información de cada una filas de la sentencia
		sentenciaPreparada, err := db.Prepare("UPDATE parque SET aforoMaximo = ? WHERE id = ?")

		if err != nil {
			panic("Error al preparar la sentencia" + err.Error()) //devolvera nil y error en caso de que no se pueda hacer la consulta
		}

		defer sentenciaPreparada.Close()

		_, err = sentenciaPreparada.Exec(numero, "SDpark")

		if err != nil {
			panic("Error al establecer el número máximo de visitantes" + err.Error())
		}

	}

}

/* Función que modifica las posiciones de los visitantes en el parque en base a sus movimientos */
func mueveVisitante(db *sql.DB, id, movimiento string, visitantes []visitante) {

	var nuevaPosicionX int
	var nuevaPosicionY int

	for i := 0; i < len(visitantes); i++ {

		if id == visitantes[i].ID { // Modificamos la posición del visitante recibido por kafka

			switch movimiento {
			case "N":
				visitantes[i].Posicionx--
			case "S":
				visitantes[i].Posicionx++
			case "W":
				visitantes[i].Posiciony--
			case "E":
				visitantes[i].Posiciony++
			case "NW":
				visitantes[i].Posicionx--
				visitantes[i].Posiciony--
			case "NE":
				visitantes[i].Posicionx--
				visitantes[i].Posiciony++
			case "SW":
				visitantes[i].Posicionx++
				visitantes[i].Posiciony--
			case "SE":
				visitantes[i].Posicionx++
				visitantes[i].Posiciony++
			}

			if visitantes[i].Posicionx == -1 {
				visitantes[i].Posicionx = 19
			} else if visitantes[i].Posicionx == 20 {
				visitantes[i].Posicionx = 0
			}

			if visitantes[i].Posiciony == -1 {
				visitantes[i].Posiciony = 19
			} else if visitantes[i].Posiciony == 20 {
				visitantes[i].Posiciony = 0
			}

			nuevaPosicionX = visitantes[i].Posicionx
			nuevaPosicionY = visitantes[i].Posiciony

		}

	}

	// MODIFICAMOS la posición de dicho visitante en la BD
	// Preparamos para prevenir inyecciones SQL
	sentenciaPreparada, err := db.Prepare("UPDATE visitante SET posicionx = ?, posiciony = ? WHERE id = ?")
	if err != nil {
		panic("Error al preparar la sentencia de modificación: " + err.Error())
	}

	defer sentenciaPreparada.Close()

	// Ejecutar sentencia, un valor por cada '?'
	_, err = sentenciaPreparada.Exec(nuevaPosicionX, nuevaPosicionY, id)
	if err != nil {
		panic("Error al modificar la posición del visitante en la BD: " + err.Error())
	}

}

/* Función que envia el mapa a los visitantes */
func productorMapa(IpBroker, PuertoBroker string, ctx context.Context, mapa []byte) {

	var brokerAddress string = IpBroker + ":" + PuertoBroker
	var topic string = "movimiento-mapa"

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:          []string{brokerAddress},
		Topic:            topic,
		CompressionCodec: kafka.Snappy.Codec(),
	})

	err := w.WriteMessages(ctx, kafka.Message{
		Key:   []byte("Key-Mapa"), //[]byte(strconv.Itoa(i)),
		Value: []byte(mapa),
	})
	if err != nil {
		panic("No se puede mandar el mapa: " + err.Error())
	}
}

/* Función que actualiza los tiempos de espera de las atracciones en la BD */
func actualizaTiemposEsperaBD(db *sql.DB, tiemposEspera []string) {

	results, err := db.Query("SELECT * FROM atraccion")

	// Comrpobamos que no se produzcan errores al hacer la consulta
	if err != nil {
		panic("Error al hacer la consulta a la BD: " + err.Error())
	}

	defer results.Close() // Nos aseguramos de cerrar*/

	i := 0

	// Recorremos todas las filas de la consulta
	for results.Next() {

		// Preparamos para prevenir inyecciones SQL
		sentenciaPreparada, err := db.Prepare("UPDATE atraccion SET tiempoEspera = ? WHERE id = ?")
		if err != nil {
			panic("Error al preparar la sentencia de modificación: " + err.Error())
		}

		defer sentenciaPreparada.Close()

		infoAtraccion := strings.Split(tiemposEspera[i], ":") // Extraemos el id y el tiempo de espera de la atracción

		idAtraccion := infoAtraccion[0]

		nuevoTiempo, err := strconv.Atoi(infoAtraccion[1])

		if err != nil {
			panic("Error al convertir la cadena con el nuevo tiempo de la atracción")
		}

		// Ejecutar sentencia, un valor por cada '?'
		_, err = sentenciaPreparada.Exec(nuevoTiempo, idAtraccion)
		if err != nil {
			panic("Error al modificar el tiempo de espera de la atracción: " + err.Error())
		}

		i++

	}
}

/*
* Función que crea un topic para el envio de los visitantes
 */
func crearTopics(IpBroker, PuertoBroker, nombre string) {
	/**** IMPORTANTE CAMBIAR*/
	//Broker1 se sustituira en localhost:9092
	var broker1 string = IpBroker + ":" + PuertoBroker
	//el localhost:9092 cambiara y sera pasado por parametro
	conn, err := kafka.Dial("tcp", broker1)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	controller, err := conn.Controller()

	if err != nil {
		panic(err.Error())
	}
	//Creamos una variable del tipo kafka.Conn
	var controllerConn *kafka.Conn
	//Le damos los valores necesarios para crear el controllerConn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err.Error())
	}
	defer controllerConn.Close()
	//Configuración del topic mapa-visitantes
	topicConfigs := []kafka.TopicConfig{
		kafka.TopicConfig{
			Topic:             nombre,
			NumPartitions:     10, //Cambiamos el número de particiones
			ReplicationFactor: 1,
		},
	}
	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		panic(err.Error())
	}
}
