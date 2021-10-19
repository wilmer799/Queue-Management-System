package main

import (
	"fmt"
	"time"

	appkafka "github.com/wilmer799/SDPracticas/kafka"
)

func main() {
	fmt.Println("Okay")
	appkafka.StartKafka()

	fmt.Println("Kafka ha iniciado correctamente")

	time.Sleep(10 * time.Minute)
}
