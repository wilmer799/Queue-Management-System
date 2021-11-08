package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	connType = "tcp"
)

/**
* Función main de los visitantes
**/
func main() {
	//Argumentos iniciales
	IpFWQ_Registry := os.Args[1]
	PuertoFWQ := os.Args[2]
	IpBroker := os.Args[3]
	PuertoBroker := os.Args[4]
	//var opcion int

	fmt.Println("**Bienvenido al parque de atracciones**")
	fmt.Println("La IP del registro es la siguiente:" + IpFWQ_Registry + ":" + PuertoFWQ)
	fmt.Println("La IP del Broker es el siguiente:" + IpBroker + ":" + PuertoBroker)
	ConsumidorKafkaVisitante(IpBroker, PuertoBroker)
	/*
		fmt.Print("Elige la opción que quieras realizar:")

		fmt.Scanln(&opcion)
		switch os := opcion; os {
		case 1:
			CrearPerfil(IpFWQ_Registry, PuertoFWQ, IpBroker, PuertoBroker)
		case 2:
			EditarPerfil(IpFWQ_Registry, PuertoFWQ)
		case 3:
			EntradaParque(IpFWQ_Registry, PuertoFWQ)
		case 4:
			SalidaParque(IpFWQ_Registry, PuertoFWQ)

		default:
			fmt.Println("Opción invalida, elige otra opción")
		}
	*/
}

func CrearPerfil(ipRegistry, puertoRegistry, IpBroker, PuertoBroker string) {
	fmt.Println("**********Creación de perfil***********")
	var informacionVisitante string
	conn, err := net.Dial(connType, ipRegistry+":"+puertoRegistry)
	if err != nil {
		fmt.Println("Error al conectarse:", err.Error())
		os.Exit(1)
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Introduce tu ID:")
		//Leer entrada hasta nueva linea, introduciendo llave
		//input es el string que se ha escrito

		id, _ := reader.ReadString('\n')
		conn.Write([]byte(id))
		fmt.Print("Introduce tu nombre:")
		nombre, _ := reader.ReadString('\n')
		conn.Write([]byte(nombre))
		fmt.Print("Introduce tu contraseña:")
		password, _ := reader.ReadString('\n')
		conn.Write([]byte(password))
		//Con la función TrimSpace eliminamos los saltos de linea de input, nombre y contraseña
		informacionVisitante = strings.TrimSpace(id) + "|" + strings.TrimSpace(nombre) + "|" + strings.TrimSpace(password)
		//Para empezar con el kafka
		ctx := context.Background()
		ProductorKafkaVisitantes(IpBroker, PuertoBroker, informacionVisitante, ctx)
		//conn.Write([]byte(id))
		//Escuchando por el relay
		//message, _ := bufio.NewReader(conn).ReadString('\n')
		//Print server relay
		//log.Print("Server relay:", message)
	}
}

func EditarPerfil(ipRegistry, puertoRegistry string) {
	fmt.Println("Has entrado a editar perfil")
	conn, err := net.Dial(connType, ipRegistry+":"+puertoRegistry)
	if err != nil {
		fmt.Println("Error al conectarse:", err.Error())
		os.Exit(1)
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Información del cliente que se quiere modificar:")
		fmt.Print("Introduce el ID:")
		id, _ := reader.ReadString('\n')
		conn.Write([]byte(id))
		fmt.Print("Introduce el nombre:")
		nombre, _ := reader.ReadString('\n')
		conn.Write([]byte(nombre))
		fmt.Print("Introduce la contraseña:")
		password, _ := reader.ReadString('\n')
		conn.Write([]byte(password))
		//message, _ := bufio.NewReader(conn).ReadString('\n')
		//log.Print("Server relay:", message)
	}

}

func EntradaParque(ipRegistry, puertoRegistry string) {
	fmt.Println("*Bienvenido al parque de atracciones*")
	conn, err := net.Dial(connType, ipRegistry+":"+puertoRegistry)
	if err != nil {
		fmt.Println("Error al conectarse:", err.Error())
		os.Exit(1)
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Por favor introduce tu alias:")
		input, _ := reader.ReadString('\n')
		fmt.Print("y tu password:")
		salida, _ := reader.ReadString('\n')
		//Enviamos la conexion del socket
		conn.Write([]byte(input))
		conn.Write([]byte(salida))
		message, _ := bufio.NewReader(conn).ReadString('\n')
		log.Print("Server relay:", message)
	}

}

func SalidaParque(ipRegistry, puertoRegistry string) {
	fmt.Println("Gracias por venir al parque, espero que vuelvas cuanto antes")
}

func ProductorKafkaVisitantes(IpBroker, PuertoBroker, mensaje string, ctx context.Context) {
	var broker1Addres string = IpBroker + ":" + PuertoBroker
	var broker2Addres string = IpBroker + ":" + PuertoBroker
	var topic string = "sd-events"
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker1Addres, broker2Addres},
		Topic:   topic,
	})
	for {
		err := w.WriteMessages(ctx, kafka.Message{
			Key:   []byte("Key-A"),                                 //[]byte(strconv.Itoa(i)),
			Value: []byte("Información del visitante: " + mensaje), //strconv.Itoa(i)),
		})
		if err != nil {
			panic("No se puede escribir mensaje" + err.Error())
		}
		//Tenemos que enviar la información de los visitantes
		//Por lo que llamaremos a esta función desde crear perfil o editar perfil e ingresar en el parque
		fmt.Println("Escribiendo:", mensaje)
		//Descanso
		time.Sleep(time.Second)
	}

}

/*
* Consumidor de kafka para un visitante en un grupo
 */
func ConsumidorKafkaVisitante(IpBroker, PuertoBroker string) {

	broker := IpBroker + ":" + PuertoBroker
	r := kafka.ReaderConfig(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   "sd-events",
		//De esta forma solo cogera los ultimos mensajes despues de unirse al cluster
		StartOffset: kafka.LastOffset,
	})
	reader := kafka.NewReader(r)
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("Ha ocurrido algún error a la hora de conectarse con kafka", err)
			continue
		}
		fmt.Println("[", string(m.Value), "]")
	}

}
