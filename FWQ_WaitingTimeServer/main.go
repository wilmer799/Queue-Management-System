package main

func main() {

	puertoEscucha := os.Args[1];
	ipBrokerGestorColas := os.Args[2];
	puertoBrokerGestorColas := os.Args[3];



}

func StartKafka() {
	conf := kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "sensor-servidorTiempos", //Topico que hemos creado
		GroupID:  "g1",
		MaxBytes: 10,
	}

	reader := kafka.NewReader(conf)
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("Ha ocurrido algún error bro", err)
			continue
		}
		fmt.Println("El mensaje es : ", string(m.Value))
	}
}