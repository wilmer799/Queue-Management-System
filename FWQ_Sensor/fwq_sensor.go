package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

type sensor struct {
	IdAtraccion string
	Personas    int
}

func main() {

	ipBrokerGestorColas := os.Args[1]

	puertoBrokerGestorColas := os.Args[2]

	idAtraccion := os.Args[3] // Convertimos a entero

	crearTopic(ipBrokerGestorColas, puertoBrokerGestorColas)

	// Comprobamos que el id de la atracción sea válido
	valido := false
	for i := 1; !valido && i < 17; i++ {

		if idAtraccion == "atraccion"+strconv.Itoa(i) {
			valido = true
		}

	}

	// Si el id pasado por parámetro no es válido
	if !valido {
		panic("Error: El id de la atracción no es válido. Introduzca atraccion1...16)")
	}

	brokerAddress := ipBrokerGestorColas + ":" + puertoBrokerGestorColas

	// Creamos un sensor
	s := new(sensor)
	s.IdAtraccion = idAtraccion

	// Generamos un número aleatorio de personas inicialmente en la cola de la atracción
	rand.Seed(time.Now().UnixNano()) // Utilizamos la función Seed(semilla) para inicializar la fuente predeterminada al requerir un comportamiento diferente para cada ejecución
	min := 0
	max := 60
	s.Personas = (rand.Intn(max-min+1) + min)
	fmt.Println("Sensor creado para la atracción (" + idAtraccion + ") en la que inicialmente hay " + strconv.Itoa(s.Personas) + " personas en cola")

	// Generamos un tiempo aleatorio entre 1 y 3 segundos
	rand.Seed(time.Now().UnixNano()) // Utilizamos la función Seed(semilla) para inicializar la fuente predeterminada al requerir un comportamiento diferente para cada ejecución
	min = 1
	max = 3
	tiempoAleatorio := (rand.Intn(max-min+1) + min)

	// Envíamos al servidor de tiempos el número de personas que se encuentra en la cola de la atracción
	enviaInformacion(s, brokerAddress, tiempoAleatorio)

	defer func() {
		fmt.Println("El sensor ha sido apagado.")
	}()

}

/* Función que envía mediante un productor de Kafka la información recogida por el sensor  */
func enviaInformacion(s *sensor, brokerAddress string, tiempoAleatorio int) {

	// Inicializamos el escritor
	escritor := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerAddress},
		Topic:   "sensor-servidorTiempos",
	})

	for {

		err := escritor.WriteMessages(context.Background(),
			kafka.Message{
				Key:   []byte("Atraccion " + s.IdAtraccion),
				Value: []byte(s.IdAtraccion + ":" + strconv.Itoa(s.Personas)),
			})

		if err != nil {
			panic("Error al conectarse al gestor de colas - No se puede mandar la información al servidor de tiempos de espera: " + err.Error())
		}

		// Generamos un número aleatorio de personas en cola
		rand.Seed(time.Now().UnixNano()) // Utilizamos la función Seed(semilla) para inicializar la fuente predeterminada al requerir un comportamiento diferente para cada ejecución
		min := 0
		max := 60
		s.Personas = (rand.Intn(max-min+1) + min)

		// Cada 1 a 3 segundos el sensor envía la información al servidor de tiempos
		time.Sleep(time.Duration(tiempoAleatorio) * time.Second)

		fmt.Println("En la atracción [" + s.IdAtraccion + "] hay " + strconv.Itoa(s.Personas) + " personas en cola")

	}

}

/*
* Función que crea el topic para el envio de los movimientos de los visitantes
 */
func crearTopic(IpBroker, PuertoBroker string) {
	topic := "sensor-servidorTiempos"
	//partition := 0
	//Broker1 se sustituira en localhost:9092
	//var broker1 string = IpBroker + ":" + PuertoBroker
	//el localhost:9092 cambiara y sera pasado por parametro
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
