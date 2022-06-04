package main

import (
	"gorm/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	// Migra si no existe el modelo. Si ya existe no lo va a migrar.
	//models.MigrarUser() // Se conecta la BD y migra el modelo user en la BD

	// Rutas
	mux := mux.NewRouter()

	// Endpoint
	mux.HandleFunc("/api/user/", handlers.GetUsers).Methods("GET")
	mux.HandleFunc("/api/user/{id:[0-9]+}", handlers.GetUser).Methods("GET")
	mux.HandleFunc("/api/user/", handlers.CreateUser).Methods("POST")
	mux.HandleFunc("/api/user/{id:[0-9]+}", handlers.UpdateUser).Methods("PUT")
	mux.HandleFunc("/api/user/{id:[0-9]+}", handlers.DeleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":3000", mux))

}
