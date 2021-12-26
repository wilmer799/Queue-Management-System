package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
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

func main() {

	host := os.Args[1]
	puertoEscucha := os.Args[2]
	ipBrokerGestorColas := os.Args[3]
	puertoBrokerGestorColas := os.Args[4]

	crearTopic(ipBrokerGestorColas, puertoBrokerGestorColas)

	// Arrancamos el servidor y atendemos conexiones entrantes
	fmt.Println("Servidor de tiempos atendiendo en " + host + ":" + puertoEscucha)

	l, err := net.Listen("tcp", host+":"+puertoEscucha)

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
			fmt.Println("Error conectando con el engine:", err.Error())
		}

		// Imprimimos la dirección de conexión del cliente
		log.Println("Cliente engine " + c.RemoteAddr().String() + " conectado.")

		// Manejamos las conexiones de forma concurrente
		go manejoConexion(ipBrokerGestorColas, puertoBrokerGestorColas, c)

	}

}

/* Función que calcula el tiempo de espera de una atracción dada una cantidad de personas en la cola */
func calculaTiempoEspera(a atraccion, personasEnCola int) int {

	tiempoEspera := 0

	// Mientras queden personas en la cola de la atracción
	for personasEnCola > a.NVisitantes {

		tiempoEspera += a.TCiclo
		personasEnCola -= a.NVisitantes

	}

	return tiempoEspera

}

// Función que maneja la lógica para una única petición de conexión
func manejoConexion(IpBroker, PuertoBroker string, conn net.Conn) {

	// Lectura del buffer de entrada hasta el final de línea
	buffer, err := bufio.NewReader(conn).ReadBytes('\n')

	fmt.Println("Petición del Engine: " + string(buffer))

	infoAtracciones := strings.Split(string(buffer), "|")

	var atracciones []atraccion

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

		atracciones = append(atracciones, a)

	}

	fmt.Println("Longitud atracciones: " + strconv.Itoa(len(atracciones)))

	// Cerrar las conexiones con engines desconectados
	if err != nil {
		log.Println("Engine" + conn.RemoteAddr().String() + " desconectado.")
		conn.Close()
		return
	}

	broker := IpBroker + ":" + PuertoBroker
	r := kafka.ReaderConfig(kafka.ReaderConfig{
		Brokers:     []string{broker},
		Topic:       "sensor-servidorTiempos",
		GroupID:     "sensores",
		StartOffset: kafka.LastOffset,
	})

	reader := kafka.NewReader(r)

	for {

		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("Ha ocurrido algún error a la hora de conectarse con kafka", err)
		}

		if strings.Contains(string(m.Value), "desconectado") {
			fmt.Println(string(m.Value))
			conn.Close()
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

			tiemposEspera := ""

			// Formamos la cadena con los tiempos de espera que le vamos a mandar al engine
			for i := 0; i < len(atracciones); i++ {
				tiemposEspera += atracciones[i].ID + ":" + strconv.Itoa(atracciones[i].TiempoEspera) + "|"
			}

			// Mandamos una cadena separada por barras con los tiempos de espera de cada atracción al engine
			conn.Write([]byte(tiemposEspera))
			conn.Close()

		}

	}

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
