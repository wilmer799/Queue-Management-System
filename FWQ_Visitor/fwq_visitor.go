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

/*
* Estructura del visitante
 */
type visitante struct {
	ID           string `json:"id"`
	Nombre       string `json:"nombre"`
	Password     string `json:"contraseña"`
	Posicionx    int    `json:"posicionx"`
	Posiciony    int    `json:"posiciony"`
	Destinox     int    `json:"destinox"`
	Destinoy     int    `json:"destinoy"`
	DentroParque int    `json:"dentroParque"`
	Parque       string `json:"parqueAtracciones"`
}

const (
	connType = "tcp"
)

var v visitante

/**
* Función main de los visitantes
**/
func main() {
	//Argumentos iniciales
	IpFWQ_Registry := os.Args[1]
	PuertoFWQ := os.Args[2]
	IpBroker := os.Args[3]
	PuertoBroker := os.Args[4]
	crearTopic(IpBroker, PuertoBroker, "inicio-sesion")          // Creamos un topic para el inicio de sesión de los visitantes en el parque
	crearTopic(IpBroker, PuertoBroker, "movimientos-visitantes") // Creamos un topic para el envío de los movimientos y la recepción del mapa
	crearTopic(IpBroker, PuertoBroker, "salir-parque")           // Creamos un topic para la solicitud de salida del parque
	fmt.Println("**Bienvenido al parque de atracciones**")
	fmt.Println("La IP del registro es la siguiente: " + IpFWQ_Registry + ":" + PuertoFWQ)
	fmt.Println("La IP del Broker es el siguiente: " + IpBroker + ":" + PuertoBroker)
	fmt.Println()
	MenuParque(IpFWQ_Registry, PuertoFWQ, IpBroker, PuertoBroker)
	ctx := context.Background()
	defer SalidaParque(v, IpBroker, PuertoBroker, ctx)

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
			ctx := context.Background()
			SalidaParque(v, IpBroker, PuertoBroker, ctx)
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

	fmt.Print("Por favor introduce tu alias:")
	alias, _ := reader.ReadString('\n')

	fmt.Print("y tu password:")
	password, _ := reader.ReadString('\n')

	ctx := context.Background()

	mensaje := string(alias) + ":" + string(password)

	//var mapa string // Variable donde almacenaremos el mapa pasado por el engine

	v := visitante{ // Guardamos la información del visitante que nos hace falta
		ID:           string(alias),
		Password:     string(password),
		Posicionx:    0,
		Posiciony:    0,
		Destinox:     -1,
		Destinoy:     -1,
		DentroParque: 0,
	}

	// Mandamos al engine las credenciales de inicio de sesión del visitante para entrar al parque
	productorLogin(IpBroker, PuertoBroker, mensaje, ctx)

	// Recibe del engine el mapa actualizado o un mensaje de parque cerrado
	respuestaEngine := consumidorLogin(IpBroker, PuertoBroker)

	if respuestaEngine == "Parque cerrado" {
		fmt.Println("Parque cerrado")
	} else {
		//mapa = formarMapa(respuestaEngine)
		v.DentroParque = 1 // El visitante está dentro del parque
	}

	for v.DentroParque == 1 { // Mientras el visitante esté dentro del parque vamos mandando los movimientos
		//go movimientoVisitante(v, mapa, IpBroker, PuertoBroker, ctx) // El visitante se desplaza una posición para alcanzar la atracción y le envía cada movimiento al engine
		time.Sleep(1 * time.Second) // Esperamos un segundo hasta volver a enviar el movimiento del visitante
	}

}

/* Función que permite a un visitante abandonar el parque */
func SalidaParque(v visitante, IpBroker string, PuertoBroker string, ctx context.Context) {

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Desea abandonar el parque (si/no): ")
	abandonar, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println("Perdone, no entiendo si quiere abandonar el parque o no")
	}

	if string(abandonar) == "s" || string(abandonar) == "si" || string(abandonar) == "SI" || string(abandonar) == "sI" || string(abandonar) == "Si" {
		v.DentroParque = 0
		mensaje := v.ID + ":" + "Salir"
		productorSalir(IpBroker, PuertoBroker, mensaje, ctx)
		fmt.Println("Gracias por venir al parque, espero que vuelvas cuanto antes")
	}

}

/* Función que se encarga de enviar las credenciales de inicio de sesión */
func productorLogin(IpBroker, PuertoBroker, credenciales string, ctx context.Context) {

	var brokerAddress string = IpBroker + ":" + PuertoBroker
	var topic string = "inicio-sesion"

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,
	})

	err := w.WriteMessages(ctx, kafka.Message{
		Key:   []byte("Key-Login"),
		Value: []byte(credenciales),
	})
	if err != nil {
		panic("No se pueden encolar las credenciales: " + err.Error())
	}

	fmt.Println("Enviando credenciales -> " + credenciales)

}

/* Función que se encarga de enviar los movimientos de los visitantes al engine */
func productorMovimientos(IpBroker, PuertoBroker, movimiento string, ctx context.Context) {

	var brokerAddress string = IpBroker + ":" + PuertoBroker
	var topic string = "movimientos-visitantes"

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,
	})

	err := w.WriteMessages(ctx, kafka.Message{
		Key:   []byte("Key-Moves"),
		Value: []byte(movimiento),
	})
	if err != nil {
		panic("No se puede encolar el movimiento: " + err.Error())
	}

	fmt.Println("Enviando movimiento -> " + movimiento)

}

/* Función que se encarga de mandar la solicitud de salida del parque al engine */
func productorSalir(IpBroker, PuertoBroker, peticion string, ctx context.Context) {

	var brokerAddress string = IpBroker + ":" + PuertoBroker
	var topic string = "salir-parque"

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,
	})

	err := w.WriteMessages(ctx, kafka.Message{
		Key:   []byte("Key-Moves"),
		Value: []byte(peticion),
	})
	if err != nil {
		panic("No se puede encolar la solicitud: " + err.Error())
	}

	fmt.Println("Enviando solicitud -> " + peticion)

}

/* Función que recibe el mensaje de parque cerrado por parte del engine o no */
func consumidorLogin(IpBroker, PuertoBroker string) string {

	respuesta := ""

	broker := IpBroker + ":" + PuertoBroker
	r := kafka.ReaderConfig(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   "inicio-sesion",
		//De esta forma solo cogera los ultimos mensajes despues de unirse al cluster
		StartOffset: kafka.LastOffset,
	})

	reader := kafka.NewReader(r)

	m, err := reader.ReadMessage(context.Background())

	if err != nil {
		panic("Ha ocurrido algún error a la hora de conectarse con kafka: " + err.Error())
	}

	if m.Value != nil {
		respuesta = string(m.Value)
	}

	return respuesta

}

/* Función que recibe el mapa del engine y lo devuelve formateado */
func consumidorMapa(IpBroker, PuertoBroker string) [][]string {

	var mapaFormateado [][]string

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
		panic("Ha ocurrido algún error a la hora de conectarse con kafka: " + err.Error())
	}

	// Procesamos el mapa enviado y lo convertimos a un array bidimensional de strings
	// PARA PROBAR A VER LO QUE MUESTRA EN CADA CASO
	fmt.Println(string(m.Value))
	fmt.Println(m.Value)
	fmt.Println(m.Value[0])
	fmt.Println(m.Value[1])
	fmt.Println(string(m.Value[0]))
	fmt.Println(string(m.Value[1]))

	//mapaFormateado = formateaMapa(m)

	return mapaFormateado

}

/* Función que formatea el mapa y lo convierte en un array bidimensional de strings */
/*func formateaMapa(mapa kafka.Message) [][]string {

	var mapaFormateado [][]string

	for i := 0; i < ; i++ {

		for j := 0; j < ; j++ {

			for k := 0; k < ; k++ {

			}

		}

	}

	return mapaFormateado

}*/

/* Función que se encarga de ir moviendo al visitante hasta alcanzar el destino */
func movimientoVisitante(v visitante, mapa [][]string, IpBroker string, PuertoBroker string, ctx context.Context) {

	for v.DentroParque == 1 { // Mientras el visitante esté dentro del parque vamos mandando los movimientos

		// Si el visitante no sabe a qué atracción dirigirse
		if v.Destinox == -1 && v.Destinoy == -1 {

			//Elegimos una atracción al azar del mapa entre las que el tiempo de espera sea menor de 60 minutos

			// Actualizamos la coordenadas de destino del visitante

		}

		// El visitante realiza un movimiento para acercarse a su destino

		movimiento := "N"

		mensaje := v.ID + ":" + movimiento

		// Enviamos el movimiento al engine
		productorMovimientos(IpBroker, PuertoBroker, mensaje, ctx)

		// Si el visitante se encuentra en la atracción
		if (v.Posicionx == v.Destinox) && (v.Posiciony == v.Destinoy) {

			time.Sleep(10 * time.Second) // Esperamos un tiempo para simular el tiempo de ciclo de la atracción
			// Ahora el visitante vuelve a desconocer su destino
			v.Destinox = -1
			v.Destinoy = -1

		}

		time.Sleep(1 * time.Second) // Esperamos un segundo hasta volver a enviar el movimiento del visitante

	}

}

/*
* Función que crea el topic para el envio de los movimientos de los visitantes
 */
func crearTopic(IpBroker, PuertoBroker, topic string) {

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
