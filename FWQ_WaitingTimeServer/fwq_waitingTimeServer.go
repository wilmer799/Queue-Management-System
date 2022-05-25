package main

import (
	"bufio"
	"context"
	hho "crypto/rand"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/oklog/ulid"
	"github.com/segmentio/kafka-go"
)

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

var atracciones []atraccion // Variable global que irá almacenando la información actual de las atracciones del parque

func main() {

	host := os.Args[1]
	puertoEscucha := os.Args[2]
	ipBrokerGestorColas := os.Args[3]
	puertoBrokerGestorColas := os.Args[4]

	crearTopic(ipBrokerGestorColas, puertoBrokerGestorColas)

	// Arrancamos el servidor y atendemos conexiones entrantes
	fmt.Println("Servidor de tiempos atendiendo en " + host + ":" + puertoEscucha)
	fmt.Println() // Por limpieza

	l, err := net.Listen("tcp", host+":"+puertoEscucha)

	if err != nil {
		fmt.Println("Error escuchando", err.Error())
		os.Exit(1)
	}

	// Cerramos el listener cuando se cierra la aplicación
	defer l.Close()

	go recibirInfoSensores(ipBrokerGestorColas, puertoBrokerGestorColas)

	// Bucle infinito hasta la salida del programa en el que se tratan las peticiones recibidas por el engine
	for {

		// Atendemos conexiones entrantes
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error conectando con el engine:", err.Error())
		}

		// Imprimimos la dirección de conexión del cliente
		fmt.Println() // Por limpieza
		log.Println("Cliente engine " + c.RemoteAddr().String() + " conectado.\n")

		// Manejamos las conexiones de forma concurrente
		go manejoConexion(c)

	}

}

/* Función que se encarga de recibir el número de personas que hay en las colas por parte de los sensores */
func recibirInfoSensores(IpBroker, PuertoBroker string) {

	broker := IpBroker + ":" + PuertoBroker
	r := kafka.ReaderConfig(kafka.ReaderConfig{
		Brokers:     []string{broker},
		Topic:       "sensor-servidorTiempos",
		GroupID:     Ulid(),
		StartOffset: kafka.LastOffset,
	})

	reader := kafka.NewReader(r)

	for {

		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("Ha ocurrido algún error a la hora de conectarse con kafka", err)
		}

		if strings.Contains(string(m.Value), "desconectado") {
			log.Println(string(m.Value))
		} else {
			fmt.Println("[", string(m.Value)+" personas en cola", "]")

			infoSensor := strings.Split(string(m.Value), ":")

			idAtraccion := infoSensor[0]
			personasEnCola, _ := strconv.Atoi(infoSensor[1])

			encontrado := false

			// Buscamos la atracción indicada por el sensor para calcular su tiempo de espera actual
			for i := 0; i < len(atracciones) && !encontrado; i++ {

				if atracciones[i].ID == idAtraccion {
					encontrado = true
					atracciones[i].TiempoEspera = calculaTiempoEspera(atracciones[i], personasEnCola)
				}

			}

		}
	}

}

/* Función que calcula el tiempo de espera de una atracción dada una cantidad de personas en la cola */
func calculaTiempoEspera(a atraccion, personasEnCola int) int {

	tiempoEspera := 0

	// Mientras el número de personas en la cola sea mayor o igual al número de personas que pueden subir a la atracción
	for personasEnCola >= a.NVisitantes {

		tiempoEspera += a.TCiclo
		personasEnCola -= a.NVisitantes

	}

	return tiempoEspera

}

// Función que maneja la lógica para una única petición de conexión
func manejoConexion(conn net.Conn) {

	// Lectura del buffer de entrada hasta el final de línea
	buffer, err := bufio.NewReader(conn).ReadBytes('\n')

	// Cerrar las conexiones con engines desconectados
	if err != nil {
		log.Println("Engine" + conn.RemoteAddr().String() + " desconectado.\n")
		conn.Close()
		return
	}

	//fmt.Println("Petición del Engine: " + string(buffer))
	log.Println("Petición de un engine recibida.")

	infoAtracciones := strings.Split(string(buffer), "|")

	fmt.Println() // Por limpieza
	fmt.Println("Estado actual de las atracciones:")
	fmt.Println("ID:TCiclo:NVisitantes:TiempoEspera")
	// El -1 es porque la longitud del array es 17, ya que al terminar la cadena en la barra vertical
	// el split realizado coge un último elemento para formar el array que contiene un espacio en blanco.
	for i := 0; i < (len(infoAtracciones) - 1); i++ {

		fmt.Println(infoAtracciones[i])
		infoAtraccion := strings.Split(infoAtracciones[i], ":")

		var a = atraccion{
			ID:           "",
			TCiclo:       -1,
			NVisitantes:  -1,
			TiempoEspera: -1,
		}

		a.ID = infoAtraccion[0]
		tciclo, _ := strconv.Atoi(infoAtraccion[1])
		a.TCiclo = tciclo
		nvisitantes, _ := strconv.Atoi(infoAtraccion[2])
		a.NVisitantes = nvisitantes
		tiempoEspera, _ := strconv.Atoi(infoAtraccion[3])
		a.TiempoEspera = tiempoEspera

		encontrada := false

		// Si la atracción ya la tenemos almacenada, actualizamos su información y sino la añadimos al slice
		for i := 0; i < len(atracciones) && !encontrada; i++ {
			if atracciones[i].ID == a.ID {
				encontrada = true
			}
		}

		if !encontrada {
			atracciones = append(atracciones, a)
		}

	}

	fmt.Println() // Por limpieza

	//fmt.Println("Longitud atracciones: " + strconv.Itoa(len(atracciones)))

	// Enviamos al engine los tiempos de espera actuales
	tiemposEspera := ""

	// Formamos la cadena con los tiempos de espera que le vamos a mandar al engine
	for i := 0; i < len(atracciones); i++ {
		tiemposEspera += atracciones[i].ID + ":" + strconv.Itoa(atracciones[i].TiempoEspera) + "|"
	}

	// Mandamos una cadena separada por barras con los tiempos de espera de cada atracción al engine
	conn.Write([]byte(tiemposEspera))
	fmt.Println("Enviando los tiempos de espera actualizados...")
	fmt.Println() // Por limpieza
	conn.Close()  // Cerramos la conexión con el engine

	manejoConexion(conn) // Reiniciamos el proceso

}

/*
* Función que crea el topic para el envio de los movimientos de los visitantes
 */
func crearTopic(IpBroker, PuertoBroker string) {

	topic := "sensor-servidorTiempos"
	conn, err := kafka.Dial("tcp", IpBroker+":"+PuertoBroker)
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
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}
	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		panic(err.Error())
	}

}

func Ulid() string {
	t := time.Now().UTC()
	id := ulid.MustNew(ulid.Timestamp(t), hho.Reader)

	return id.String()
}
