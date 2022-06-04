/*package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

// Funciones
func Saludar(nombre string) string {
	return "Hola " + nombre + " desde una función"
}

// Handler
func Index(rw http.ResponseWriter, r *http.Request) {

	funciones := template.FuncMap{
		"saludar": Saludar,
	}

	template, err := template.New("index.html").Funcs(funciones).ParseFiles("index.html")

	if err != nil {
		panic(err)
	} else {
		template.Execute(rw, nil)
	}

}

func main() {

	// Mux
	mux := http.NewServeMux()
	mux.HandleFunc("/", Index)

	// Servidor
	server := &http.Server{
		Addr:    "localhost:3000",
		Handler: mux,
	}

	fmt.Println("El servidor está corriendo en el puerto 3000")
	fmt.Println("Run server: http://localhost:3000/")
	log.Fatal(server.ListenAndServe())

}*/
