package main

import (
	"fmt"
	"math/rand"
	"os"
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
	idAtraccion := os.Args[3]

	for respuesta != "no" {

		fmt.Println("Desea añadir un sensor (si/no): ")
		fmt.Scanln(&respuesta)

		if respuesta == "si" {

			creaSensor(ipGestorColas, puertoGestorColas, idAtraccion)

		} else if respuesta == "no" {

		} else {

		}

	}

}

func creaSensor(ip, puerto, id string) {

	s := new(sensor)
	s.IdAtraccion = 
	s.Personas = rand.Intn(10-0) + 0

}
