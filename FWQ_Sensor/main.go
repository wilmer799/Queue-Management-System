package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const (
	topic             = "sd-events"
	brokerAddress     = "localhost:9092"
	timeServerAddress = "localhost:9094"
)

type sensor struct {
	IdAtraccion int
	Personas    int
}

func (this *sensor) enviaInformacion(ip int, puerto int, id int) {

}

func main() {

	// Nos guardamos los parámetros en formato entero
	ipBrokerGestorColas, err := strconv.Atoi(os.Args[1])

	if err != nil {
		fmt.Println("Error: Introduzca por parámetros IP, PUERTO e ID")
	}

	puertoBrokerGestorColas, err := strconv.Atoi(os.Args[2])

	if err != nil {
		fmt.Println("Error: Introduzca por parámetros IP, PUERTO e ID")
	}

	idAtraccion, err := strconv.Atoi(os.Args[3])

	if err != nil {
		fmt.Println("Error: Introduzca por parámetros IP, PUERTO e ID")
	}

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

	// Cada x segundos el sensor envía la información al servidor de tiempos
	intervalo := time.Tick(tiempoAleatorio * time.Second)

	s.enviaInformacion()

}
