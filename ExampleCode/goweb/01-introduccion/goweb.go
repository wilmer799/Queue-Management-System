package main

import (
	"fmt"
	"log"
	"net/http"
)

// Handlers
func Hola(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("El método es + " + r.Method)
	fmt.Fprintln(rw, "Hola Mundo")
}

func PaginaNF(rw http.ResponseWriter, r *http.Request) {
	http.NotFound(rw, r)
}

func Error(rw http.ResponseWriter, r *http.Request) {
	//http.Error(rw, "La página no funciona", 404)
	// En vez de utilizar los valores enteros, GO nos recomienda usar las constantes que podemos consultar en https://pkg.go.dev/net/http#pkg-constants
	http.Error(rw, "La página no funciona", http.StatusNotFound)
}

func Saludar(rw http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	fmt.Println(r.URL.RawQuery) // Con esto separamos la URL y los argumentos que estamos envíando en dicha URL
	fmt.Println(r.URL.Query())  // Con esto convertimos la URL en un mapa

	name := r.URL.Query().Get("name")
	age := r.URL.Query().Get("age")

	// Usamos Fprintf para imprimir por el navegador
	fmt.Fprintf(rw, "Hola, %s tu edad es %s !!", name, age)

	//http.Error(rw, "Este es un Error", http.StatusConflict)
}

func main() {

	// Un Mux es una ruta asociada a un handler
	// Mux
	mux := http.NewServeMux()
	mux.HandleFunc("/", Hola)
	mux.HandleFunc("/page", PaginaNF)
	mux.HandleFunc("/error", Error)
	mux.HandleFunc("/saludar", Saludar)

	// Router
	/*http.HandleFunc("/", Hola)
	http.HandleFunc("/page", PaginaNF)
	http.HandleFunc("/error", Error)
	http.HandleFunc("/saludar", Saludar)*/

	// Crear servidor

	// En vez de trabajar de la forma estándar para crear un servidor,
	// nosotros podemos estructurar y crear nuestro propio servidor
	server := &http.Server{
		Addr:    "localhost:3000",
		Handler: mux,
	}

	//http.ListenAndServe("localhost:3000", nil)
	fmt.Println("El servidor está corriendo en el puerto 3000")
	fmt.Println("Run server: http://localhost:3000/")

	// Con log.Fatal si se devuelve algún error vamos a poder logear dicho error.
	// Cuando aquí ponemos nil estamos trabajando con mux automáticos o mux por defecto.
	//log.Fatal(http.ListenAndServe("localhost:3000", nil))
	log.Fatal(server.ListenAndServe()) // Ahora utilizamos nuestro nuevo servidor estructurado a nuestro gusto

}
