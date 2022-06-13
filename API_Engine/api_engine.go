package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

/**
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
func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", Home)
	r.HandleFunc("/cambiarCiudad/{lat}/{lan}", cambiarCiudad).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
