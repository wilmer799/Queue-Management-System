package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
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
	crearTopic(IpBroker, PuertoBroker)
	fmt.Println("**Bienvenido al parque de atracciones**")
	fmt.Println("La IP del registro es la siguiente:" + IpFWQ_Registry + ":" + PuertoFWQ)
	fmt.Println("La IP del Broker es el siguiente:" + IpBroker + ":" + PuertoBroker)
	fmt.Println()
	MenuParque(IpFWQ_Registry, PuertoFWQ, IpBroker, PuertoBroker)
	defer SalidaParque(idUsuario)
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
			SalidaParque(idUsuario)
		default:
			fmt.Println("Opción invalida, elige otra opción")
		}
	}
}

/* Función que se conecta al módulo FWQ_Registry para crear un nuevo usuario */
func CrearPerfil(ipRegistry, puertoRegistry string) {

	fmt.Println("**********Creación de perfil***********")
	conn, err := net.Dial(connType, ipRegistry+":"+puertoRegistry)

	if err != nil {
		fmt.Println("Error al conectarse al Registry:", err.Error())
	} else { // Si el visitante establece conexión con el Registry indicado por parámetro

		conn.Write([]byte("1" + "\n")) // Le pasamos al Registry la opción elegida por el visitante

		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Introduce tu ID:")
		id, _ := reader.ReadString('\n')
		conn.Write([]byte(id))

		fmt.Print("Introduce tu nombre:")
		nombre, _ := reader.ReadString('\n')
		conn.Write([]byte(nombre))

		fmt.Print("Introduce tu contraseña:")
		password, _ := reader.ReadString('\n')
		conn.Write([]byte(password))

		//Escuchando por el relay el mensaje de respuesta del Registry
		message, _ := bufio.NewReader(conn).ReadString('\n')

		// Comprobamos si el Registry nos devuelve un mensaje de respuesta
		if message != "" {
			log.Print("Respuesta del Registry: ", message)
		} else {
			log.Print("Lo siento, el Registry no está disponible en estos momentos.")
		}

	}

}

/* Función que se conecta al módulo FWQ_Registry para editar o actualizar el perfil de un usuario existente */
func EditarPerfil(ipRegistry, puertoRegistry string) {

	fmt.Println("Has entrado a editar perfil")
	conn, err := net.Dial(connType, ipRegistry+":"+puertoRegistry)

	if err != nil {
		fmt.Println("Error al conectarse al Registry:", err.Error())
	} else { // Si el visitante establece conexión con el Registry indicado por parámetro

		conn.Write([]byte("2" + "\n")) // Le pasamos al Registry la opción elegida por el visitante

		reader := bufio.NewReader(os.Stdin)

		fmt.Println("Información del visitante que se quiere modificar:")
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

		// Comprobamos si el Registry nos devuelve un mensaje de respuesta
		if message != "" {
			log.Print("Respuesta del Registry: ", message)
		} else {
			log.Print("Lo siento, el Registry no está disponible en estos momentos.")
		}

	}

}

/* Función que envía las credenciales de acceso del visitante para entrar en el parque */
func EntradaParque(ipRegistry, puertoRegistry, IpBroker, PuertoBroker string) {

	fmt.Println("*Bienvenido al parque de atracciones*")

	reader := bufio.NewReader(os.Stdin)

	for {

		fmt.Print("Por favor introduce tu alias:")
		alias, _ := reader.ReadString('\n')
		idUsuario = string(alias) // Nos guardamos el id del visitante para cuando quiera salir del parque
		fmt.Print("y tu password:")
		password, _ := reader.ReadString('\n')

		ctx := context.Background()

		mensaje := string(alias) + ":" + string(password)

		// Mandamos al engine las credenciales de inicio de sesión del visitante para entrar al parque
		ProductorKafkaVisitantes(IpBroker, PuertoBroker, mensaje, ctx)

		// Recibe del engine la posición actual y la posición de destino
		posiciones := ConsumidorKafkaVisitantes(IpBroker, PuertoBroker)

		posicionesVisitante := strings.Split(posiciones, "|")

		posicionActual := strings.Split(posicionesVisitante[0], ",")
		posicionDestino := strings.Split(posicionesVisitante[1], ",")

		posicionx, _ := strconv.Atoi(posicionActual[0])
		posiciony, _ := strconv.Atoi(posicionActual[1])
		destinox, _ := strconv.Atoi(posicionDestino[0])
		destinoy, _ := strconv.Atoi(posicionDestino[1])

		salir := false
		for !salir {

			movimientoVisitante(posicionx, posiciony, destinox, destinoy) // El visitante se desplaza una posición para alcanzar la atracción

			// Recibe del engine el mapa actualizado o un mensaje de parque cerrado
			respuestaEngine := ConsumidorKafkaVisitantes(IpBroker, PuertoBroker)

			// Mostramos el mapa o el mensaje de error
			fmt.Println(respuestaEngine)

			fmt.Println("Desea abandonar el parque (si/no): ")
			abandonar, _ := reader.ReadString('\n')

			if string(abandonar) == "s" || string(abandonar) == "si" || string(abandonar) == "SI" || string(abandonar) == "sI" || string(abandonar) == "Si" {
				SalidaParque(idUsuario)
				salir = true
			}

			time.Sleep(1 * time.Second) // Esperamos un segundo hasta volver a enviar el movimiento del visitante

		}

	}

}

/* Función que permite a un visitante abandonar el parque */
func SalidaParque(idUsuario string) {
	//aqui le pasamos el id del usuario
	//El cual se buscara en la bd y se eliminara al usuario
	fmt.Println("Gracias por venir al parque, espero que vuelvas cuanto antes")
}

/*
* Función que envian la información de los movimientos de los visitantes
 */
func ProductorKafkaVisitantes(IpBroker, PuertoBroker, mensaje string, ctx context.Context) {
	var brokerAddres string = IpBroker + ":" + PuertoBroker
	var topic string = "movimientos-visitantes"
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerAddres},
		Topic:   topic,
	})
	sigue := true
	for sigue {
		err := w.WriteMessages(ctx, kafka.Message{
			Key:   []byte("Key-A"), //[]byte(strconv.Itoa(i)),
			Value: []byte(mensaje), //strconv.Itoa(i)),
		})
		if err != nil {
			panic("No se puede escribir mensaje" + err.Error())
		}
		//Tenemos que enviar la información de los visitantes
		//Por lo que llamaremos a esta función desde crear perfil o editar perfil e ingresar en el parque
		fmt.Println("Escribiendo:", mensaje)
		//Descanso
		time.Sleep(time.Second)
		sigue = false
	}

}

/*
* Consumidor de kafka para un visitante en un grupo
 */
func ConsumidorKafkaVisitantes(IpBroker, PuertoBroker string) string {

	broker := IpBroker + ":" + PuertoBroker
	r := kafka.ReaderConfig(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   "mapa-visitantes",
		//De esta forma solo cogera los ultimos mensajes despues de unirse al cluster
		StartOffset: kafka.LastOffset,
	})

	reader := kafka.NewReader(r)

	m, err := reader.ReadMessage(context.Background())
	if err != nil {
		fmt.Println("Ha ocurrido algún error a la hora de conectarse con kafka", err)
	}

	fmt.Println("[", string(m.Value), "]")

	return string(m.Value)

}

/* Función que se encarga de ir moviendo al visitante hasta alcanzar el destino */
func movimientoVisitante(posicionx, posiciony, destinox, destinoy int) {
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
	var broker1 string = IpBroker + ":" + PuertoBroker
	//el localhost:9092 cambiara y sera pasado por parametro
	conn, err := kafka.Dial("tcp", broker1)
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
