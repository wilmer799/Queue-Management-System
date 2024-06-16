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
}

// GO nos ofrece que nosotros podemos cargar todos los archivos html y
// luego utilizarlos en cada handler
// Con esto tendremos todos los archivos html cargados en la variable templates y
// de esta forma nos evitamos tener que cargar todos los templates en cada handler
//var templates = template.Must(template.New("T").ParseFiles("index.html", "base.html"))

// Si usamos la librería Parse Glob nos ahorramos tener que escribir todos los templates.
// Parse Glob carga todos los templates de un directorio
var templates = template.Must(template.New("T").ParseGlob("templates/**/*.html")) // Con el ** nos referimos a los subdirectorios
var errorTemplate = template.Must(template.ParseFiles("templates/error/error.html"))

// Handler error
func handlerError(rw http.ResponseWriter, status int) {
	rw.WriteHeader(status)
	errorTemplate.Execute(rw, nil)
}

// Función que se va a encargar de renderizar los templates y
// la usaremos en todos los handlers que tengamos
func renderTemplate(rw http.ResponseWriter, name string, data interface{}) {

	err := templates.ExecuteTemplate(rw, name, data)

	if err != nil {
		//panic(err)
		//http.Error(rw, "No es posible devolver el template", http.StatusInternalServerError)
		handlerError(rw, http.StatusInternalServerError)
	}

}

// Handlers
func Index(rw http.ResponseWriter, r *http.Request) {

	// fmt.Fprintln(rw, "Hola Mundo")
	// Renderiza el primer archivo indicado y los demás los tiene ahí para reutilizar por ejemplo la herencia
	//template, err := template.ParseFiles("index.html", "base.html") // ParseFiles devuelve dos valores: el template en sí y también un error

	// Otra forma de hacerlo para evitar algún problema o algún error.
	// Indicamos explícitamente que lo que vamos a renderizar es el index.html con New
	//template, err := template.New("index.html").ParseFiles("index.html", "base.html")

	// La función Must maneja internamente los errores, por lo que nos ayudar a simplificar el código que maneja los errores
	//template := template.Must(template.New("index.html").ParseFiles("index.html", "base.html"))

	usuario := Usuarios{"Jose", 25}

	//template.Execute(rw, usuario)

	renderTemplate(rw, "inde.html", usuario)

}

func Registro(rw http.ResponseWriter, r *http.Request) {

	renderTemplate(rw, "registro.html", nil)

}

func main() {

	// Mux
	mux := http.NewServeMux()
	mux.HandleFunc("/", Index)
	mux.HandleFunc("/registro", Registro)

	// Servidor
	server := &http.Server{
		Addr:    "localhost:3000",
		Handler: mux,
	}

	fmt.Println("El servidor está corriendo en el puerto 3000")
	fmt.Println("Run server: http://localhost:3000/")
	log.Fatal(server.ListenAndServe())

}*/
