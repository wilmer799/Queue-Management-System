package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Función para responder con los datos
func sendData(rw http.ResponseWriter, data interface{}, status int) {
	rw.Header().Set("Content-Type", "aplication/json")
	rw.WriteHeader(status)

	output, _ := json.Marshal(&data)
	fmt.Fprintln(rw, string(output))
}

// Función para responder con un error
func sendError(rw http.ResponseWriter, status int) {
	rw.WriteHeader(status)
	fmt.Fprintln(rw, "Resource Not Found")
}
