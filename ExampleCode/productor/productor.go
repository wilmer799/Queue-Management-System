package productor

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	topic         = "sd-events"
	broker1Addres = "localhost:9092"
	broker2Addres = "localhost:9093"
)

func StartProduce(ctx context.Context) {
	//Inicializamos el contador
	i := 0

	//Inicializamos el escritor
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker1Addres, broker2Addres},
		Topic:   topic,
	})

	for {
		err := w.WriteMessages(ctx, kafka.Message{
			Key:   []byte(strconv.Itoa(i)),
			Value: []byte("Esto es un mensaje" + strconv.Itoa(i)),
		})
		if err != nil {
			panic("No se puede escribir mensaje" + err.Error())
		}
		fmt.Println("Escribiendo:", i)
		i++
		//Descanso
		time.Sleep(time.Second)
	}
}
