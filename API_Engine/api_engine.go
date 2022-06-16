package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

/**
* Para obtener informacion de servidor
* https://parzibyte.me/blog/2018/12/03/servidor-web-go/
* Función principal para las rutas
* Api Key del servidor de tiempos
* sdpracticas
* c3d8572d0046f36f0c586caa0e2e1d23
* https://openweathermap.org/current
* https://api.openweathermap.org/data/2.5/weather?lat=35&lon=139&appid=c3d8572d0046f36f0c586caa0e2e1d23&lang=es
* La de abajo es una petición para conocer como se llama la cuidad y las coordenadas para las peticiones de arriba
* http://api.openweathermap.org/geo/1.0/direct?q=Orihuela&limit=5&appid=c3d8572d0046f36f0c586caa0e2e1d23
* Para obtener la temperatura en grados celsius, realizamos lo siguiente petición, basicamente es añadir un nuevo parametro
* https://api.openweathermap.org/data/2.5/weather?lat=40.4167047&lon=-3.7035825&appid=c3d8572d0046f36f0c586caa0e2e1d23&lang=es&units=metric
* https://riptutorial.com/go/example/6628/decoding-json-data-from-a-file
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
	IdEnParque   string `json:"idEnParque"`
	UltimoEvento string `json:"ultimoEvento"`
	Parque       string `json:"parqueAtracciones"`
}

type mapa struct {
	Fila       int    `json:"fila"`
	InfoParque string `json:"infoParque"`
}

type ciudad struct {
	Cuadrante   string  `json:"cuadrante"`
	Nombre      string  `json:"name"`
	Temperatura float32 `json:"temp"`
}

func main() {

	ip := os.Args[1]
	puerto := os.Args[2]

	// IMPLEMENTAMOS EL API REST
	// Rutas
	mux := mux.NewRouter()

	// Responder al cliente
	mux.HandleFunc("/visitantes", getVisitantes).Methods("GET")
	mux.HandleFunc("/mapa", getMapa).Methods("GET")
	mux.HandleFunc("/ciudades", getCiudades).Methods("GET")

	// SERVIDOR
	// Arrancamos el servidor https en una go routine
	//go http.ListenAndServeTLS(":8081", "cert.pem", "key.pem", mux)
	fmt.Println("Servidor API ENGINE corriendo en https://" + ip + ":" + puerto)
	//log.Fatal(http.ListenAndServe(":8080", http.HandlerFunc(redirectToHttps)))
	log.Fatal(http.ListenAndServe(":"+puerto, mux))

}

// INICIO BLOQUE RESPONSE

// Estructura con el formato de una respuesta http
type Response struct {
	Status        int         `json:"status"`
	Data          interface{} `json:"data"`
	Message       string      `json:"message"`
	contentType   string
	responseWrite http.ResponseWriter
}

/* Función que crea una respuesta por defecto para los clientes de la API REST */
func CreateDefaultResponse(rw http.ResponseWriter) Response {
	return Response{
		Status:        http.StatusOK,
		responseWrite: rw,
		contentType:   "application/json",
	}
}

/* Función que envía las respuestas a los clientes de la API REST */
func (resp *Response) Send() {
	resp.responseWrite.Header().Set("Content-Type", resp.contentType)
	resp.responseWrite.WriteHeader(resp.Status)

	// Marshall devuelve 2 valores: Los valores transformados en tipo byte y un error
	output, _ := json.Marshal(&resp) // Para responder con json
	//output, _ := xml.Marshal(&resp) // Para responder con xml
	//output, _ := yaml.Marshal(&resp) // Para responder con yaml
	fmt.Fprintln(resp.responseWrite, string(output))
}

/* Función que envía una respuesta a los clientes indicando que la consulta de los visitantes ha sido satisfactoria */
func SendDataGetVisitantes(rw http.ResponseWriter, data interface{}) {
	response := CreateDefaultResponse(rw)
	response.Data = data
	response.Message = "OK: Estado actual de los visitantes obtenido."
	response.Send()
}

/* Función que envía una respuesta a los clientes indicando que el registro ha sido satisfactorio */
func SendDataGetMapa(rw http.ResponseWriter, data interface{}) {
	response := CreateDefaultResponse(rw)
	response.Data = data
	response.Message = "OK: Estado actual del mapa obtenido."
	response.Send()
}

/* Función que envía una respuesta a los clientes indicando que el registro ha sido satisfactorio */
func SendDataGetCiudades(rw http.ResponseWriter, data interface{}) {
	response := CreateDefaultResponse(rw)
	response.Data = data
	response.Message = "OK: Información de las ciudades obtenida."
	response.Send()
}

func (resp *Response) NotFound() {
	resp.Status = http.StatusNotFound
	resp.Message = "ERROR: Resource Not Found"
}

/* Función que manda una respuesta indicando que el recurso solicitado no ha sido encontrado */
func SendNotFound(rw http.ResponseWriter) {
	response := CreateDefaultResponse(rw)
	response.NotFound()
	response.Send()
}

// FIN BLOQUE RESPONSE

// DEFINIMOS LOS HANDLERS
/* Función manejadora para la obtención del estado de los visitantes */
func getVisitantes(rw http.ResponseWriter, r *http.Request) {

	log.Println("Petición de consulta de estado de los visitantes -> " + r.URL.Path)

	visitantes := []visitante{}

	// Accedemos a la base de datos, empezando por abrir la conexión
	db, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/parque_atracciones")

	// Comprobamos que no haya error al conectarse
	if err != nil {
		panic("Error al conectarse con la BD: " + err.Error())
	}

	defer db.Close() // Para que siempre se cierre la conexión con la BD al finalizar el programa

	rows, err := db.Query("SELECT * FROM visitante")
	// Comprobamos que no se produzcan errores al hacer la consulta
	if err != nil {
		panic("Error al consultar el estado de los visitantes en la BD: " + err.Error())
	}

	for rows.Next() {
		v := visitante{}
		rows.Scan(&v.ID, &v.Nombre, &v.Password, &v.Posicionx, &v.Posiciony, &v.Destinox, &v.Destinoy, &v.DentroParque, &v.IdEnParque, &v.IdEnParque, &v.UltimoEvento)
		visitantes = append(visitantes, v)
	}

	// CONTINUAR IMPLEMENTACIÓN
	SendDataGetVisitantes(rw, visitantes)

}

/* Función manejadora para la obtención del estado del mapa del parque */
func getMapa(rw http.ResponseWriter, r *http.Request) {

	log.Println("Petición de consulta de estado del mapa -> " + r.URL.Path)

	filasParque := []mapa{}

	// Accedemos a la base de datos, empezando por abrir la conexión
	db, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/parque_atracciones")

	// Comprobamos que no haya error al conectarse
	if err != nil {
		panic("Error al conectarse con la BD: " + err.Error())
	}

	defer db.Close() // Para que siempre se cierre la conexión con la BD al finalizar el programa

	rows, err := db.Query("SELECT * FROM mapa")
	// Comprobamos que no se produzcan errores al hacer la consulta
	if err != nil {
		panic("Error al consultar el estado del mapa en la BD: " + err.Error())
	}

	for rows.Next() {
		fila := mapa{}
		rows.Scan(&fila.Fila, &fila.InfoParque)
		filasParque = append(filasParque, fila)
	}

	// CONTINUAR IMPLEMENTACIÓN
	SendDataGetMapa(rw, filasParque)

}

/* Función manejadora para la obtención de la información de las ciudades */
func getCiudades(rw http.ResponseWriter, r *http.Request) {

	log.Println("Petición de consulta de las ciudades -> " + r.URL.Path)

	ciudades := []ciudad{}

	// Accedemos a la base de datos, empezando por abrir la conexión
	db, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/parque_atracciones")

	// Comprobamos que no haya error al conectarse
	if err != nil {
		panic("Error al conectarse con la BD: " + err.Error())
	}

	defer db.Close() // Para que siempre se cierre la conexión con la BD al finalizar el programa

	rows, err := db.Query("SELECT * FROM ciudades")
	// Comprobamos que no se produzcan errores al hacer la consulta
	if err != nil {
		panic("Error al consultar la información de las ciudades: " + err.Error())
	}

	for rows.Next() {
		ciudad := ciudad{}
		rows.Scan(&ciudad.Cuadrante, &ciudad.Nombre, &ciudad.Temperatura)
		ciudades = append(ciudades, ciudad)
	}

	// CONTINUAR IMPLEMENTACIÓN
	SendDataGetCiudades(rw, ciudades)

}
