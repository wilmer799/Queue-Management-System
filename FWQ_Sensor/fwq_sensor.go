package main

import (
	"context"
	"fmt"
	"math/rand"
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

	brokerAddress := ipBrokerGestorColas + ":" + puertoBrokerGestorColas

	// Creamos un sensor
	s := new(sensor)
	s.IdAtraccion = idAtraccion

	// Generamos un número aleatorio de personas inicialmente en la cola de la atracción
	rand.Seed(time.Now().UnixNano()) // Utilizamos la función Seed(semilla) para inicializar la fuente predeterminada al requerir un comportamiento diferente para cada ejecución
	min := 0
	max := 40
	s.Personas = (rand.Intn(max-min+1) + min)
	fmt.Println("Sensor creado para la atracción (" + idAtraccion + ") en la que inicialmente hay " + strconv.Itoa(s.Personas) + " en cola")

	// Generamos un tiempo aleatorio entre 1 y 3 segundos
	rand.Seed(time.Now().UnixNano()) // Utilizamos la función Seed(semilla) para inicializar la fuente predeterminada al requerir un comportamiento diferente para cada ejecución
	min = 1
	max = 3
	tiempoAleatorio := (rand.Intn(max-min+1) + min)

	// Envíamos al servidor de tiempos el número de personas que se encuentra en la cola de la atracción
	enviaInformacion(s, brokerAddress, tiempoAleatorio)

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
			panic("Error: No se puede escribir el mensaje: " + err.Error())
		}

		// Generamos un número aleatorio de personas en cola
		rand.Seed(time.Now().UnixNano()) // Utilizamos la función Seed(semilla) para inicializar la fuente predeterminada al requerir un comportamiento diferente para cada ejecución
		min := 0
		max := 40
		s.Personas = (rand.Intn(max-min+1) + min)

		// Cada x segundos el sensor envía la información al servidor de tiempos
		time.Sleep(time.Duration(tiempoAleatorio) * time.Second)
	}

}
