package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
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

func main() {

	var respuesta string

	ipGestorColas := os.Args[1]
	puertoGestorColas := os.Args[2]
	idAtraccion := strconv.Atoi(os.Args[3])

	for respuesta != "no" {

		fmt.Println("Desea añadir un sensor (si/no): ")
		fmt.Scanln(&respuesta)

		if respuesta == "si" {
			creaSensor(idAtraccion)
		} else if respuesta == "no" {

		} else {

		}

	}

}

func creaSensor(id int) {

	s := new(sensor)
	s.idAtraccion = id
	s.Personas = rand.Intn(10-0) + 0

}
