package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/segmentio/kafka-go"
)

/*
* Estructura del visitante
 */
type visitante struct {
	ID        string `json:"id"`
	Nombre    string `json:"nombre"`
	Password  string `json:"contraseña"`
	Posicionx int    `json:"posicionx"`
	Posiciony int    `json:"posiciony"`
	Destinox  int    `json:"destinox"`
	Destinoy  int    `json:"destinoy"`
	Parque    string `json:parqueAtracciones`
}

/**
 * @Description : Función main de fwq_engine
 * @Author : Wilmer Fabricio Bravo Shuira
**/
func main() {
	IpKafka := os.Args[1]
	PuertoKafka := os.Args[2]
	numeroVisitantes := os.Args[3]
	IpFWQWating := os.Args[4]
	PuertoWaiting := os.Args[5]
	fmt.Println("**Bienvenido al engine de la aplicación**")
	fmt.Println("La ip del apache kafka es el siguiente:" + IpKafka)
	fmt.Println("El puerto del apache kafka es el siguiente:" + PuertoKafka)
	fmt.Println("El número máximo de visitantes es el siguiente:" + numeroVisitantes)
	fmt.Println("La ip del servidor de espera es el siguiente:" + IpFWQWating)
	fmt.Println("El puerto del servidor de tiempo es el siguiente:" + PuertoWaiting)
	//Reserva de memoria para el mapa
	var mapa [20][20]byte
	//Asignamos unos valores para comprobar que esta bien construida la matriz
	mapa[1][2] = 1
	mapa[2][1] = 2
	mapa[3][1] = 3
	//Accediendo a la base de datos
	//Abrimos la conexion con la base de datos

	db, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/parque_atracciones")
	//Si la conexión falla mostrara este error
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	//Ejecutamos la sentencia
	results, err := db.Query("SELECT * FROM visitante")
	if err != nil {
		panic(err.Error()) //En caso de error a la hora de hacer la select
	}
	//Recorremos los resultados obtenidos por la consulta
	for results.Next() {
		//   var nombreVariable tipoVariable
		var fwq_visitante visitante
		err = results.Scan(&fwq_visitante.ID, &fwq_visitante.Nombre,
			&fwq_visitante.Password, &fwq_visitante.Posicionx,
			&fwq_visitante.Posiciony, &fwq_visitante.Destinox, &fwq_visitante.Destinoy,
			&fwq_visitante.Parque)
		if err != nil {
			panic(err.Error())
		}
		log.Println(fwq_visitante.ID, fwq_visitante.Nombre, fwq_visitante.Password,
			fwq_visitante.Posicionx, fwq_visitante.Posiciony, fwq_visitante.Destinox,
			fwq_visitante.Destinoy, fwq_visitante.Parque)
	}
	//Cada una de las casillas, su valor entero representa el tiempo en minutos de una atracción
	//Cada uno de los personajes tenemos que representarlo por algo
	//Esto se le asignara cuando entre al parque
	//El mapa se carga de la base de datos al arrancar la aplicación
	fmt.Println("*********** FUN WITH QUEUES RESORT ACTIVITY MAP *********")
	fmt.Println("ID   " + "Nombre      " + "Pos.      " + "Destino")
	//Matriz transversal bidimensional
	for i := 0; i < len(mapa); i++ {
		for j := 0; j < len(mapa[i]); j++ {
			fmt.Print(mapa[i][j], " ")
		}
		fmt.Println()
	}
	conexionKafka()

	//Tiene que ser cada X tiempo para que actualize la matriz
	//tiempoEspera(IpFWQWating, PuertoKafka)
}

/**
* Función que se conecta al servidor de tiempo de espera
**/
func tiempoEspera(IpFWQWating, PuertoWaiting string) {
	fmt.Println("***Tiempo de espera***")
	var connType string = "tcp"
	conn, err := net.Dial(connType, IpFWQWating+":"+PuertoWaiting)
	if err != nil {
		fmt.Println("Error al conectarse:", err.Error())
		os.Exit(1)
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("***Actualizando los tiempos de espera***")
		//Leer entrada hasta nueva linea, introduciendo llave
		input, _ := reader.ReadString('\n')
		//Enviamos la conexion del socket
		conn.Write([]byte(input))
		//Escuchando por el relay
		message, _ := bufio.NewReader(conn).ReadString('\n')
		//Print server relay
		log.Print("Server relay:", message)
	}
}

/**
* Función que conecta el engine con el kafka
**/
func conexionKafka() {
	//Configuración de lector de kafka
	conf := kafka.ReaderConfig{
		//El broker habra que cambiarlo por otro
		Brokers:  []string{"localhost:9092"},
		Topic:    "sd-events", //Topico que hemos creado
		GroupID:  "g1",
		MaxBytes: 10,
	}
	reader := kafka.NewReader(conf)
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("Ha ocurrido algún error a la hora de conectarse con kafka", err)
			continue
		}
		fmt.Println("El mensaje es desde el terminal wilmer : ", string(m.Value))
	}
}
