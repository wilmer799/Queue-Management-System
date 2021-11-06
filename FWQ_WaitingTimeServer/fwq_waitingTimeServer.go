package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/segmentio/kafka-go"
)

func main() {

	puertoEscucha := os.Args[1]
	ipBrokerGestorColas := os.Args[2]
	puertoBrokerGestorColas := os.Args[3]

}

func StartKafka() {

	conf := kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "sd-events", //Topico que hemos creado
		GroupID:  "g1",
		MaxBytes: 10,
	}

	reader := kafka.NewReader(conf)
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("Ha ocurrido algún error", err)
			continue
		}
		fmt.Println("El mensaje es : ", string(m.Value))

		infoAtraccion := strings.Split(string(buffer[:len(buffer)-1]), ":")

		idAtraccion := infoAtraccion[0]
		personasAtraccion := infoAtraccion[1]

	}
}
