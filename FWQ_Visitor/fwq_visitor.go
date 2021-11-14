package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	connType = "tcp"
)

var idUsuario string

/**
* Función main de los visitantes
**/
func main() {
	//Argumentos iniciales
	IpFWQ_Registry := os.Args[1]
	PuertoFWQ := os.Args[2]
	IpBroker := os.Args[3]
	PuertoBroker := os.Args[4]
	//crearTopic(IpBroker, PuertoBroker)
	fmt.Println("**Bienvenido al parque de atracciones**")
	fmt.Println("La IP del registro es la siguiente:" + IpFWQ_Registry + ":" + PuertoFWQ)
	fmt.Println("La IP del Broker es el siguiente:" + IpBroker + ":" + PuertoBroker)
	MenuParque(IpFWQ_Registry, PuertoFWQ, IpBroker, PuertoBroker)
}

/*
* Función que pinta el menu del parque
 */
func MenuParque(IpFWQ_Registry, PuertoFWQ, IpBroker, PuertoBroker string) {
	var opcion int
	//Guardamos la opcion elegida
	for {
		fmt.Println("***Menu parque de atracciones***")
		fmt.Println("1.Crear perfil")
		fmt.Println("2.Editar perfil")
		fmt.Println("3.Moverse por el parque")
		fmt.Println("4.Salir del parque")
		fmt.Print("Elige la acción a realizar:")
		fmt.Scanln(&opcion)

		switch os := opcion; os {
		case 1:
			CrearPerfil(IpFWQ_Registry, PuertoFWQ)
		case 2:
			EditarPerfil(IpFWQ_Registry, PuertoFWQ)
		case 3:
			EntradaParque(IpFWQ_Registry, PuertoFWQ, IpBroker, PuertoBroker)
		case 4:
			SalidaParque(IpFWQ_Registry, PuertoFWQ, idUsuario)
		default:
			fmt.Println("Opción invalida, elige otra opción")
		}
	}
}

func CrearPerfil(ipRegistry, puertoRegistry string) {
	fmt.Println("**********Creación de perfil***********")
	//var informacionVisitante string
	conn, err := net.Dial(connType, ipRegistry+":"+puertoRegistry)
	if err != nil {
		fmt.Println("Error al conectarse:", err.Error())
		os.Exit(1)
	}
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Introduce tu ID:")
	//Leer entrada hasta nueva linea, introduciendo llave
	//input es el string que se ha escrito
	id, _ := reader.ReadString('\n')
	idUsuario = id
	conn.Write([]byte(id))
	fmt.Print("Introduce tu nombre:")
	nombre, _ := reader.ReadString('\n')
	conn.Write([]byte(nombre))
	fmt.Print("Introduce tu contraseña:")
	password, _ := reader.ReadString('\n')
	conn.Write([]byte(password))
	//Solo nos interesa que llegue la información y se pueda dar de alta
	//Con la función TrimSpace eliminamos los saltos de linea de input, nombre y contraseña
	//informacionVisitante = strings.TrimSpace(id) + "|" + strings.TrimSpace(nombre) + "|" + strings.TrimSpace(password)

	//Escuchando por el relay
	message, _ := bufio.NewReader(conn).ReadString('\n')
	//Print server relay
	log.Print("Respuesta del Registry: ", message)

}

func EditarPerfil(ipRegistry, puertoRegistry string) {
	fmt.Println("Has entrado a editar perfil")
	conn, err := net.Dial(connType, ipRegistry+":"+puertoRegistry)
	if err != nil {
		fmt.Println("Error al conectarse:", err.Error())
		os.Exit(1)
	}
	reader := bufio.NewReader(os.Stdin)

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
	message, _ := bufio.NewReader(conn).ReadString('\n')
	log.Print("Respuesta del Registry:", message)

}

func EntradaParque(ipRegistry, puertoRegistry, IpBroker, PuertoBroker string) {
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
		//Llama al kafka para dibujar el mapa y la información del visitantes
		//ConsumidorKafkaVisitante(IpBroker, PuertoBroker)
		//ProductorKafkaVisitantes(IpBroker,PuertoBroker,mensaje, ctx)
		/*
			message, _ := bufio.NewReader(conn).ReadString('\n')
			log.Print("Server relay:", message) */
	}

}

func SalidaParque(ipRegistry, puertoRegistry, idUsuario string) {
	//aqui le pasamos el id del usuario
	//El cual se buscara en la bd y se eliminara al usuario
	fmt.Println("Gracias por venir al parque, espero que vuelvas cuanto antes")
}

/*
* Función que envian la información de los movimientos de los visitantes
 */
func ProductorKafkaVisitantes(IpBroker, PuertoBroker, mensaje string, ctx context.Context) {
	var broker1Addres string = IpBroker + ":" + PuertoBroker
	var broker2Addres string = IpBroker + ":" + PuertoBroker
	var topic string = "movimientos-visitantes"
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
		Topic:   "mapa-visitantes",
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

func movimientoParque() {
	//Primero el engine enviara el mapa con los visitantes y las atracciones.
	//Con esa información los visitantes se empezaran a mover
	//while
	//1. enviar mapa por el topic a los visitantes
	//2. mover los visitantes
	//3. Enviar información de movimiento por el topic
}

/*
* Función que crea el topic para el envio de los movimientos de los visitantes
 */
func crearTopic(IpBroker, PuertoBroker string) {
	topic := "movimientos-visitantes"
	//partition := 0
	//Broker1 se sustituira en localhost:9092
	//var broker1 string = IpBroker + ":" + PuertoBroker
	//el localhost:9092 cambiara y sera pasado por parametro
	conn, err := kafka.Dial("tcp", "localhost:9092")
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	controller, err := conn.Controller()

	if err != nil {
		panic(err.Error())
	}
	//Creamos una variable del tipo kafka.Conn
	var controllerConn *kafka.Conn
	//Le damos los valores necesarios para crear el controllerConn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err.Error())
	}
	defer controllerConn.Close()
	//Configuración del topic mapa-visitantes
	topicConfigs := []kafka.TopicConfig{
		kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}
	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		panic(err.Error())
	}

}
