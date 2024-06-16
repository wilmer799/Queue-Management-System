/*package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

// Estructuras
type Usuarios struct {
	UserName string
	Edad     int
	Activo   bool
	Admin    bool
	Cursos   []Curso
}

type Curso struct {
	Nombre string
}

// Handler
func Index(rw http.ResponseWriter, r *http.Request) {

	c1 := Curso{"Go"}
	c2 := Curso{"Python"}
	c3 := Curso{"Java"}
	c4 := Curso{"Javascript"}

	// fmt.Fprintln(rw, "Hola Mundo")
	template, err := template.ParseFiles("index.html") // ParseFiles devuelve dos valores: el template en sí y también un error

	cursos := []Curso{c1, c2, c3, c4}
	usuario := Usuarios{"Jose", 25, true, false, cursos}

	if err != nil {
		panic(err)
	} else {
		//template.Execute(rw, nil)     // No envíamos datos al html
		template.Execute(rw, usuario) // Envíamos datos al html
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
