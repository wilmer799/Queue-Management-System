package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/segmentio/kafka-go"
)

/*
* Estructura de las atracciones
 */
type atraccion struct {
	ID           string `json:"id"`
	TCiclo       int    `json:"tciclo"`
	NVisitantes  int    `json:"nvisitantes"`
	Posicionx    int    `json:"posicionx"`
	Posiciony    int    `json:"posiciony"`
	TiempoEspera int    `json:"tiempoEspera"`
	Parque       string `json:"parqueAtracciones"`
}

func main() {

	puertoEscucha := os.Args[1]
	ipBrokerGestorColas := os.Args[2]
	puertoBrokerGestorColas := os.Args[3]

	var conexionBD = conexionBD()
	var atracciones []atraccion

	atracciones, _ = obtenerAtraccionesBD(conexionBD)

	go recibeInformacionSensor(ipBrokerGestorColas, puertoBrokerGestorColas, atracciones)

	go atiendeEngine(puertoEscucha, atracciones)

}

/*
* Consumidor de kafka para recibir la información de los sensores
 */
func recibeInformacionSensor(IpBroker, PuertoBroker string, atracciones []atraccion) {

	broker := IpBroker + ":" + PuertoBroker
	r := kafka.ReaderConfig(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   "sensor-servidorTiempos",
		// Para que empiece a leer desde el último registro
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

		infoSensor := strings.Split(string(m.Value), ":")

		idAtraccion := infoSensor[0]
		personasEnCola, _ := strconv.Atoi(infoSensor[1])

		encontrado := false

		// Buscamos la atracción indicada por el sensor para calcular su tiempo de espera actual
		for i := 0; i < len(atracciones) && !encontrado; i++ {

			if atracciones[i].ID == idAtraccion {
				encontrado = true
				atracciones[i].TiempoEspera = calculaTiempoEspera(atracciones[i], personasEnCola)
			}

		}

	}

}

/* Función que calcula el tiempo de espera de una atracción dada una cantidad de personas en la cola */
func calculaTiempoEspera(a atraccion, personasEnCola int) int {

	tiempoEspera := 0

	// Mientras queden personas en la cola de la atracción
	for personasEnCola > a.NVisitantes {

		tiempoEspera += a.TCiclo
		personasEnCola -= a.NVisitantes

	}

	return tiempoEspera

}

/*
* Función que abre una conexion con la bd
 */
func conexionBD() *sql.DB {
	//Accediendo a la base de datos
	//Abrimos la conexion con la base de datos
	db, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/parque_atracciones")
	//Si la conexión falla mostrara este error
	if err != nil {
		panic(err.Error())
	}
	//Cierra la conexion con la bd
	//defer db.Close()
	return db
}

/*
* Función que obtiene las atracciones del parque
* @return []atraccion : Array con las atracciones del parque
* @return error : Error en caso de que no se ha podido obtener las atracciones
 */
func obtenerAtraccionesBD(db *sql.DB) ([]atraccion, error) {

	//Ejecutamos la sentencia
	results, err := db.Query("SELECT * FROM atraccion")
	if err != nil {
		return nil, err //devolvera nil y error en caso de que no se pueda hacer la consulta
	}

	//Cerramos la base de datos
	defer results.Close()

	//Declaramos el array de visitantes
	var atraccionesParque []atraccion

	//Recorremos los resultados obtenidos por la consulta
	for results.Next() {
		//   var nombreVariable tipoVariable
		//Variable donde guardamos la información de cada una filas de la sentencia
		var fwq_atraccion atraccion
		if err := results.Scan(&fwq_atraccion.ID, &fwq_atraccion.TCiclo,
			&fwq_atraccion.NVisitantes, &fwq_atraccion.Posicionx,
			&fwq_atraccion.Posiciony, &fwq_atraccion.TiempoEspera,
			&fwq_atraccion.Parque); err != nil {
			return atraccionesParque, err
		}
		//Vamos añadiendo las atracciones al array
		atraccionesParque = append(atraccionesParque, fwq_atraccion)
	}

	if err = results.Err(); err != nil {
		return atraccionesParque, err
	}
	return atraccionesParque, nil

}

func atiendeEngine(puertoEscuha string, atracciones []atraccion) {

}
