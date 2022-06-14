package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
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
func main() {

	http.HandleFunc("/hola", func(w http.ResponseWriter, peticion *http.Request) {
		io.WriteString(w, "Solicitaste hola")
	})

	r := mux.NewRouter()

	r.HandleFunc("/hola", HomeHandler)
	r.HandleFunc("obtenerCiudad", ObtenerCiudad)

	direccion := ":8080"
	fmt.Println("Servidor listo escuchando en" + direccion)
	log.Fatal(http.ListenAndServe(direccion, nil))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Estamos probando la api rest con postman")
}

func ObtenerCiudad(w http.ResponseWriter, r *http.Request) {

}
