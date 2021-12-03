package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
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
	IdParque     string `json:"idParque"`
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

/*
* Estructura del parque
 */
type parque struct {
	ID          string `json:"id"`
	AforoMaximo int    `json:"aforoMaximo"`
	AforoActual int    `json:"aforoActual"`
}

/*
 * @Description : Función main de fwq_engine
 * @Author : Wilmer Fabricio Bravo Shuira
 */
func main() {

	IpKafka := os.Args[1]
	PuertoKafka := os.Args[2]
	numeroVisitantes := os.Args[3]
	IpFWQWating := os.Args[4]
	PuertoWaiting := os.Args[5]

	//Creamos el topic...Cambiar la Ipkafka en la función principal
	//Si no se ejecuta el programa, se cierra el kafka?
	crearTopics(IpKafka, PuertoKafka, "peticion-login")
	crearTopics(IpKafka, PuertoKafka, "respuesta-login")
	crearTopics(IpKafka, PuertoKafka, "mapa")
	crearTopics(IpKafka, PuertoKafka, "movimientos")
	crearTopics(IpKafka, PuertoKafka, "salir")
	/*fmt.Println("**Bienvenido al engine de la aplicación**")
	fmt.Println("La ip del apache kafka es el siguiente:" + IpKafka)
	fmt.Println("El puerto del apache kafka es el siguiente:" + PuertoKafka)
	fmt.Println("El número máximo de visitantes es el siguiente:" + numeroVisitantes)
	fmt.Println("La ip del servidor de espera es el siguiente:" + IpFWQWating)
	fmt.Println("El puerto del servidor de tiempo es el siguiente:" + PuertoWaiting)
	fmt.Println()*/

	//Reserva de memoria para el mapa
	//var mapa [20][20]string

	// Visitantes, atracciones que se encuentran en la BD
	var visitantesRegistrados []visitante
	//var visitantesParque []visitante
	var atracciones []atraccion
	var parqueTematico []parque

	var conn = conexionBD()
	maxVisitantes, _ := strconv.Atoi(numeroVisitantes)
	establecerMaxVisitantes(conn, maxVisitantes)

	visitantesRegistrados, _ = obtenerVisitantesBD(conn)
	atracciones, _ = obtenerAtraccionesBD(conn)
	parqueTematico, _ = obtenerParqueDB(conn)

	// Esta parte la podemos suprimir, simplemente es a modo de comprobación
	fmt.Println("Visitantes registrados: ")
	fmt.Println(visitantesRegistrados)
	fmt.Println() // Para mejorar la salida por pantalla
	fmt.Println("Atracciones del parque: ")
	fmt.Println(atracciones)
	fmt.Println() // Para mejorar la salida por pantalla
	fmt.Println("Información del parque: ")
	fmt.Println(parqueTematico)
	fmt.Println() // Para mejorar la salida por pantalla

	fmt.Println("*********** FUN WITH QUEUES RESORT ACTIVITY MAP *********")
	fmt.Println("ID   	" + "		Nombre      " + "	Pos.      " + "	Destino")

	// Hay que usar la función TrimSpace porque al parecer tras la obtención de valores de BD se agrega un retorno de carro a cada variable
	for i := 0; i < len(visitantesRegistrados); i++ {
		fmt.Println(strings.TrimSpace(visitantesRegistrados[i].ID) + " 		" + strings.TrimSpace(visitantesRegistrados[i].Nombre) +
			"    " + "	(" + strings.TrimSpace(strconv.Itoa(visitantesRegistrados[i].Posicionx)) + "," + strings.TrimSpace(strconv.Itoa(visitantesRegistrados[i].Posiciony)) +
			")" + "    " + "	(" + strings.TrimSpace(strconv.Itoa(visitantesRegistrados[i].Destinox)) + "," + strings.TrimSpace(strconv.Itoa(visitantesRegistrados[i].Destinoy)) +
			")")
	}

	for {

		//Para empezar con el kafka
		ctx := context.Background()

		visitantesRegistrados, _ = obtenerVisitantesBD(conn) // Obtenemos los visitantes registrados actualmente
		atracciones, _ = obtenerAtraccionesBD(conn)          // Obtenemos las atracciones actualizadass

		go consumidorLogin(conn, IpKafka, PuertoKafka, ctx, visitantesRegistrados, maxVisitantes)

		go consumidorMovimientos(conn, IpKafka, PuertoKafka, ctx, atracciones)

		go consumidorSalir(conn, IpKafka, PuertoKafka, ctx)

		// Cada X segundos se conectará al servidor de tiempos para actualizar los tiempos de espera de las atracciones
		time.Sleep(time.Duration(5 * time.Second))
		tiempo := "Mándame los tiempos de espera actualizados"
		conexionTiempoEspera(conn, IpFWQWating, PuertoWaiting, tiempo)

	}

}

/* Función que comprueba si el aforo del parque está completo o no */
func parqueLleno(db *sql.DB, maxAforo int) bool {

	var lleno bool = false

	// Comprobamos si las credenciales de acceso son válidas
	results, err := db.Query("SELECT * FROM visitante WHERE dentroParque = 1")

	if err != nil {
		fmt.Println("Error al hacer la consulta sobre la BD para el login: " + err.Error())
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

/* Función que recibe del gestor de colas las credenciales de los visitantes que quieren iniciar sesión para entrar en el parque */
func consumidorLogin(db *sql.DB, IpKafka, PuertoKafka string, ctx context.Context, visitantesRegistrados []visitante, maxVisitantes int) {

	direccionKafka := IpKafka + ":" + PuertoKafka

	//Configuración de lector de kafka
	conf := kafka.ReaderConfig{
		//El broker habra que cambiarlo por otro
		Brokers:     []string{direccionKafka},
		Topic:       "peticion-login", //Topico que hemos creado
		StartOffset: kafka.LastOffset,
		//MaxBytes: 10,
	}

	reader := kafka.NewReader(conf)

	for {

		m, err := reader.ReadMessage(context.Background())

		if err != nil {
			fmt.Println("Ha ocurrido algún error a la hora de conectarse con el kafka", err)
		}

		//fmt.Println("Petición de inicio de sesión del visitante: ", string(m.Value))

		credenciales := strings.Split(string(m.Value), ":")

		alias := credenciales[0]
		password := credenciales[1]

		v := visitante{
			ID:       strings.TrimSpace(alias),
			Password: strings.TrimSpace(password),
		}

		//fmt.Println("El alias recibido es: " + alias)
		//fmt.Println("El password recibido es: " + password)

		// Comprobamos si las credenciales de acceso son válidas
		results, err := db.Query("SELECT * FROM visitante WHERE id = ? and contraseña = ?", v.ID, v.Password)

		if err != nil {
			fmt.Println("Error al hacer la consulta sobre la BD para el login: " + err.Error())
		}

		// Cerramos la base de datos
		//defer results.Close()

		var respuesta string = ""

		// Si las credenciales coinciden con las de un visitante registrado en la BD y el parque no está lleno
		if results.Next() && !parqueLleno(db, maxVisitantes) {

			// Actualizamos el estado del visitante en la BD
			sentenciaPreparada, err := db.Prepare("UPDATE visitante SET dentroParque = 1 WHERE id = ?")
			if err != nil {
				panic("Error al preparar la sentencia de modificación: " + err.Error())
			}

			//defer sentenciaPreparada.Close()

			// Ejecutar sentencia, un valor por cada '?'
			_, err = sentenciaPreparada.Exec(v.ID)
			if err != nil {
				panic("Error al actualizar el estado del visitante respecto al parque: " + err.Error())
			}

			respuesta += alias + ":" + "Acceso concedido"

			results.Close()
			sentenciaPreparada.Close()

		} else { // Sino entonces le informamos de que el parque está cerrado
			respuesta += alias + ":" + "Parque cerrado"
			results.Close()
		}

		productorLogin(IpKafka, PuertoKafka, ctx, respuesta)

	}

}

/* Función que envía el mensaje de respuesta a la petición de login de un visitante */
func productorLogin(IpBroker, PuertoBroker string, ctx context.Context, respuesta string) {

	var brokerAddress string = IpBroker + ":" + PuertoBroker
	var topic string = "respuesta-login"

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,
		//CompressionCodec: kafka.Snappy.Codec(),
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
* Función que convierte el mapa de tipo string en []byte
* @return []mapa : Retorna un array de bytes que va a ser enviado a los visitantes que se encuentren en el parque
 */
func convertirMapa(mapa [20][20]string) []byte {
	var mapaOneD []byte
	var cadenaMapa []byte
	for i := 0; i < len(mapa); i++ {
		for j := 0; j < len(mapa); j++ {
			cadenaMapa = []byte(mapa[i][j])
			mapaOneD = append(mapaOneD, cadenaMapa...)
		}
	}
	return mapaOneD
}

/*
* Función que abre una conexion con la bd
 */
func conexionBD() *sql.DB {
	//Accediendo a la base de datos
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
* Función que obtienen los parques
* @return []parque : Arrays de parque en la base de datos
* @return error : Error en caso de que no se pueda obtener parques
 */
func obtenerParqueDB(db *sql.DB) ([]parque, error) {

	//Cada parque sera un grupo // Idea
	results, err := db.Query("SELECT * FROM parque")

	if err != nil {
		return nil, err //devolvera nil y error en caso de que no se pueda hacer la consulta
	}

	//Cerramos la base de datos
	defer results.Close()

	//Declaramos el array de visitantes
	var parquesTematicos []parque

	//Recorremos los resultados obtenidos por la consulta
	for results.Next() {

		//Variable donde guardamos la información de cada una filas de la sentencia
		var parqueTematico parque

		if err := results.Scan(&parqueTematico.ID, &parqueTematico.AforoMaximo,
			&parqueTematico.AforoActual); err != nil {
			return parquesTematicos, err
		}

		//Vamos añadiendo los visitantes al array
		parquesTematicos = append(parquesTematicos, parqueTematico)

	}

	if err = results.Err(); err != nil {
		return parquesTematicos, err
	}

	return parquesTematicos, nil

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
			&fwq_visitante.DentroParque, &fwq_visitante.IdParque, &fwq_visitante.Parque); err != nil {
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
			&fwq_visitante.DentroParque, &fwq_visitante.IdParque, &fwq_visitante.Parque); err != nil {
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

	//Cerramos la base de datos
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

	//Asignamos valores a las posiciones del mapa
	for i := 0; i < len(mapa); i++ {
		for j := 0; j < len(mapa[i]); j++ {
			for k := 0; k < len(visitantes); k++ {
				if i == visitantes[k].Posicionx && j == visitantes[k].Posiciony && visitantes[k].DentroParque == 1 {
					mapa[i][j] = visitantes[k].IdParque + "|"
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
func conexionTiempoEspera(db *sql.DB, IpFWQWating, PuertoWaiting, tiempo string) {

	fmt.Println("***Conexión con el servidor de tiempo de espera***")
	fmt.Println("Arrancando el engine para atender los tiempos en el puerto: " + IpFWQWating + ":" + PuertoWaiting)
	var connType string = "tcp"
	conn, err := net.Dial(connType, IpFWQWating+":"+PuertoWaiting)

	if err != nil {
		fmt.Println("Error a la hora de escuchar", err.Error())
	} else {
		fmt.Println("***Actualizando los tiempos de espera***")

		conn.Write([]byte("Mándame los tiempos de espera actualizados" + "\n")) // Mandamos la petición
		tiemposEspera, _ := bufio.NewReader(conn).ReadString('\n')              // Obtenemos los tiempos de espera actualizados

		arrayTiemposEspera := strings.Split(tiemposEspera[:len(tiemposEspera)-1], "|")

		// Actualizamos los tiempos de espera de las atracciones en la BD
		actualizaTiemposEsperaBD(db, arrayTiemposEspera)

		if tiemposEspera != "" {
			log.Println("Tiempos de espera actualizados: " + tiemposEspera)
		} else {
			log.Println("Servidor de tiempos no disponible.")
		}

	}

}

/* Función que establece en la tabla parque de la BD el aforo máximo permitido */
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

/* Función que recibe los movimientos de los visitantes y actualiza el mapa en base a ellos */
func consumidorMovimientos(db *sql.DB, IpKafka, PuertoKafka string, ctx context.Context, atracciones []atraccion) {

	var mapa [20][20]string
	var visitantesParque []visitante
	direccionKafka := IpKafka + ":" + PuertoKafka

	//Configuración de lector de kafka
	conf := kafka.ReaderConfig{
		//El broker habra que cambiarlo por otro
		Brokers: []string{direccionKafka},
		Topic:   "movimientos", //Topico que hemos creado
		//MaxBytes: 10,
	}

	reader := kafka.NewReader(conf)

	for {

		m, err := reader.ReadMessage(context.Background())

		if err != nil {
			fmt.Println("Ha ocurrido algún error a la hora de conectarse con kafka", err)
			continue
		}

		fmt.Println("Mensaje desde el gestor de colas: ", string(m.Value))

		mensaje := strings.Split(string(m.Value), ":")

		idVisitante := mensaje[0]
		movimiento := mensaje[1]

		// Si un visitante nos ha indicado un movimiento
		if movimiento == " " || movimiento == "N" || movimiento == "S" || movimiento == "W" || movimiento == "E" ||
			movimiento == "NW" || movimiento == "NE" || movimiento == "SW" || movimiento == "SE" {
			visitantesParque, _ = obtenerVisitantesParque(db) // Obtenemos los visitantes del parque actualizados
			visitantesParque = mueveVisitante(db, idVisitante, movimiento, visitantesParque)
			// Preparamos el mapa a enviar a los visitantes que se encuentra en el parque
			mapa = asignacionPosiciones(visitantesParque, atracciones, mapa)
			mapaConvertido := convertirMapa(mapa)
			productorMapa(IpKafka, PuertoKafka, ctx, mapaConvertido) // Mandamos el mapa actualizado a los visitantes que se encuentran en el parque
		}

	}

}

/* Función que modifica las posiciones de los visitantes en el parque en base a sus movimientos */
func mueveVisitante(db *sql.DB, id, movimiento string, visitantes []visitante) []visitante {

	var nuevaPosicionX int
	var nuevaPosicionY int

	for i := 0; i < len(visitantes); i++ {
		if id == visitantes[i].ID {
			switch movimiento {
			case "N":
				if visitantes[i].Posicionx == 0 {
					visitantes[i].Posicionx = 19
				} else {
					visitantes[i].Posicionx--
				}
			case "S":
				if visitantes[i].Posicionx == 19 {
					visitantes[i].Posicionx = 0
				} else {
					visitantes[i].Posicionx++
				}
			case "W":
				if visitantes[i].Posiciony == 0 {
					visitantes[i].Posiciony = 19
				} else {
					visitantes[i].Posiciony--
				}
			case "E":
				if visitantes[i].Posiciony == 19 {
					visitantes[i].Posiciony = 0
				} else {
					visitantes[i].Posiciony++
				}
			case "NW":
				if visitantes[i].Posicionx == 0 && visitantes[i].Posiciony == 0 {
					visitantes[i].Posicionx = 19
					visitantes[i].Posiciony = 19
				} else if visitantes[i].Posicionx == 0 && visitantes[i].Posiciony > 0 {
					visitantes[i].Posicionx = 19
					visitantes[i].Posiciony--
				} else if visitantes[i].Posicionx > 0 && visitantes[i].Posiciony == 0 {
					visitantes[i].Posicionx--
					visitantes[i].Posiciony = 19
				} else {
					visitantes[i].Posicionx--
					visitantes[i].Posiciony--
				}
			case "NE":
				if visitantes[i].Posicionx == 0 && visitantes[i].Posiciony == 19 {
					visitantes[i].Posicionx = 19
					visitantes[i].Posiciony = 0
				} else if visitantes[i].Posicionx == 0 && visitantes[i].Posiciony > 0 {
					visitantes[i].Posicionx = 19
					visitantes[i].Posiciony++
				} else if visitantes[i].Posicionx > 0 && visitantes[i].Posiciony == 19 {
					visitantes[i].Posicionx--
					visitantes[i].Posiciony = 0
				} else {
					visitantes[i].Posicionx--
					visitantes[i].Posiciony++
				}
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

	return visitantes

}

/* Función que se encarga de que un visitante salga del parque */
func salirVisitanteParque(db *sql.DB, idVisitante string) {

	sentenciaPreparada, err := db.Prepare("UPDATE visitante SET dentroParque = 0 WHERE id = ?")
	if err != nil {
		panic("Error al preparar la sentencia de modificación: " + err.Error())
	}

	defer sentenciaPreparada.Close()

	// Ejecutar sentencia, un valor por cada '?'
	_, err = sentenciaPreparada.Exec(idVisitante)
	if err != nil {
		panic("Error al modificar el estado del visitante en el parque: " + err.Error())
	}

}

/* Función que envia el mapa a los visitantes */
func productorMapa(IpBroker, PuertoBroker string, ctx context.Context, mapa []byte) {

	var brokerAddress string = IpBroker + ":" + PuertoBroker
	var topic string = "mapa"

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:          []string{brokerAddress},
		Topic:            topic,
		CompressionCodec: kafka.Snappy.Codec(),
	})

	//https://cwiki.apache.org/confluence/display/KAFKA/A+Guide+To+The+Kafka+Protocol
	//https://docs.confluent.io/clients-confluent-kafka-go/current/overview.html
	//https://www.confluent.io/blog/5-things-every-kafka-developer-should-know/

	//https://en.m.wikipedia.org/wiki/SerDes
	//Compresion de mensajes
	//https://developer.ibm.com/articles/benefits-compression-kafka-messaging/
	err := w.WriteMessages(ctx, kafka.Message{
		Key:   []byte("Key-Mapa"), //[]byte(strconv.Itoa(i)),
		Value: []byte(mapa),
	})
	if err != nil {
		panic("No se puede mandar el mapa: " + err.Error())
	}
}

/* Función que recibe del gestor de colas las peticiones de salida del parque de los visitantes */
func consumidorSalir(db *sql.DB, IpKafka, PuertoKafka string, ctx context.Context) {

	direccionKafka := IpKafka + ":" + PuertoKafka

	//Configuración de lector de kafka
	conf := kafka.ReaderConfig{
		//El broker habra que cambiarlo por otro
		Brokers:     []string{direccionKafka},
		Topic:       "peticion-salir", //Topico que hemos creado
		StartOffset: kafka.LastOffset,
		//MaxBytes: 10,
	}

	reader := kafka.NewReader(conf)

	for {

		m, err := reader.ReadMessage(context.Background())

		if err != nil {
			fmt.Println("Ha ocurrido algún error a la hora de conectarse con el kafka", err)
		}

		mensaje := strings.Split(string(m.Value), ":")

		alias := mensaje[0]
		peticion := mensaje[1]

		v := visitante{
			ID: strings.TrimSpace(alias),
		}

		// Comprobamos si las credenciales de acceso son válidas
		results, err := db.Query("SELECT * FROM visitante WHERE id = ?", v.ID)

		if err != nil {
			fmt.Println("Error al hacer la consulta sobre la BD para el login: " + err.Error())
		}

		// Si las credenciales coinciden con las de un visitante registrado en la BD y el parque no está lleno
		if results.Next() && peticion == "Salir" {

			// Sacamos del parque al visitante
			sentenciaPreparada, err := db.Prepare("UPDATE visitante SET dentroParque = 0 WHERE id = ?")
			if err != nil {
				panic("Error al preparar la sentencia de modificación: " + err.Error())
			}

			// Ejecutar sentencia, un valor por cada '?'
			_, err = sentenciaPreparada.Exec(v.ID)
			if err != nil {
				panic("Error al actualizar el estado del visitante respecto al parque: " + err.Error())
			}

			results.Close()
			sentenciaPreparada.Close()

		} else { // Sino entonces le informamos de que el parque está cerrado
			results.Close()
		}

	}

}

/* Función que actualiza los tiempos de espera de las atracciones en la BD*/
func actualizaTiemposEsperaBD(db *sql.DB, tiemposEspera []string) {

	results, err := db.Query("SELECT * FROM atraccion")

	// Comrpobamos que no se produzcan errores al hacer la consulta
	if err != nil {
		panic("Error al hacer la consulta a la BD: " + err.Error())
	}

	defer results.Close() // Nos aseguramos de cerrar

	i := 0

	// Recorremos todas las filas de la consulta
	for results.Next() {

		// MODIFICAMOS la información de dicho visitante en la BD
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

		fmt.Println("Atracción modificada.")

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
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}
	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		panic(err.Error())
	}
}
