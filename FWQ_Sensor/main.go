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

const (
	topic             = "sensor-servidorTiempos"
	timeServerAddress = "localhost:9094"
)

type sensor struct {
	IdAtraccion int
	Personas    int
}

func main() {

	ipBrokerGestorColas := os.Args[1]

	puertoBrokerGestorColas := os.Args[2]

	idAtraccion, err := strconv.Atoi(os.Args[3])

	if err != nil {
		panic("Error: Introduzca por parámetros IP, PUERTO e ID " + err.Error())
	}

	brokerAddress := ipBrokerGestorColas + ":" + puertoBrokerGestorColas

	// Creamos un sensor
	s := new(sensor)
	s.IdAtraccion = idAtraccion
	// Generamos un número aleatorio de personas
	rand.Seed(time.Now().UnixNano()) // Utilizamos la función Seed(semilla) para inicializar la fuente predeterminada al requerir un comportamiento diferente para cada ejecución
	min := 0
	max := 10
	s.Personas = (rand.Intn(max-min+1) + min)
	fmt.Printf("Sensor creado para la atracción %d en la que inicialmente hay %d personas\n", s.IdAtraccion, s.Personas)

	// Generamos un tiempo aleatorio entre 1 y 3 segundos
	rand.Seed(time.Now().UnixNano()) // Utilizamos la función Seed(semilla) para inicializar la fuente predeterminada al requerir un comportamiento diferente para cada ejecución
	min = 1
	max = 3
	tiempoAleatorio := (rand.Intn(max-min+1) + min)

	// Envíamos al servidor de tiempos el número de personas que se encuentra en la cola de la atracción
	ctx := context.Background()
	enviaInformacion(s, ctx, brokerAddress, tiempoAleatorio)

}

/* Función que envía mediante un productor de Kafka la información recogida por el sensor  */
func enviaInformacion(s *sensor, ctx context.Context, brokerAddress string, tiempoAleatorio int) {

	// Inicializamos el escritor
	escritor := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerAddress, timeServerAddress},
		Topic:   topic,
	})

	for {

		err := escritor.WriteMessages(ctx, kafka.Message{
			Value: []byte("En la atracción " + strconv.Itoa(s.IdAtraccion) + " hay actualmente " + strconv.Itoa(s.Personas) + " personas"),
		})
		if err != nil {
			panic("Error: No se puede escribir el mensaje: " + err.Error())
		}

		// Generamos un número aleatorio de personas
		rand.Seed(time.Now().UnixNano()) // Utilizamos la función Seed(semilla) para inicializar la fuente predeterminada al requerir un comportamiento diferente para cada ejecución
		min := 0
		max := 10
		s.Personas = (rand.Intn(max-min+1) + min)

		// Cada x segundos el sensor envía la información al servidor de tiempos
		time.Sleep(time.Duration(tiempoAleatorio) * time.Second)
	}

}
