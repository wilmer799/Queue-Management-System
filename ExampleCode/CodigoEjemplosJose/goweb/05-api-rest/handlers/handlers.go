package handlers

import (
	"apirest/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetUsers(rw http.ResponseWriter, r *http.Request) {
	// EL CÓDIGO COMENTADO ES ANTERIOR A LA PRIMERA REFACTORIZACIÓN
	//fmt.Fprintln(rw, "Lista todos los usuarios")

	/*rw.Header().Set("content-type", "application/json") // Para responder con json
	rw.Header().Set("content-type", "text/xml") // Para responder con xml
	// No hay un tipo de dato específico para yaml

	db.Connect()
	users, _ := models.ListUsers()
	db.Close()
	// Marshall devuelve 2 valores: Los valores transformados en tipo byte y un error
	output, _ := json.Marshal(users) // Para responder con json
	//output, _ := xml.Marshal(users) // Para responder con xml
	//output, _ := yaml.Marshal(users) // Para responder con yaml
	fmt.Fprintln(rw, string(output))*/

	if users, err := models.ListUsers(); err != nil {
		models.SendNotFound(rw)
	} else {
		models.SendData(rw, users)
	}

}

func GetUser(rw http.ResponseWriter, r *http.Request) {
	// EL CÓDIGO COMENTADO ES ANTERIOR A LA PRIMERA REFACTORIZACIÓN
	//fmt.Fprintln(rw, "Obtiene un usuario")

	/*rw.Header().Set("content-type", "application/json") // Para responder con json
	//rw.Header().Set("content-type", "text/xml") // Para responder con xml
	// No hay un tipo de dato específico para yaml

	// Obtener ID
	vars := mux.Vars(r) // Mos devuelve un mapa de tipo string
	userId, _ := strconv.Atoi(vars["id"])

	db.Connect()
	user, _ := models.GetUser(userId)
	db.Close()
	// Marshall devuelve 2 valores: Los valores transformados en tipo byte y un error
	output, _ := json.Marshal(user) // Para responder con json
	//output, _ := xml.Marshal(users) // Para responder con xml
	//output, _ := yaml.Marshal(users) // Para responder con yaml
	fmt.Fprintln(rw, string(output))*/

	if user, err := getUserByRequest(r); err != nil {
		models.SendNotFound(rw)
	} else {
		models.SendData(rw, user)
	}

}

func CreateUser(rw http.ResponseWriter, r *http.Request) {
	// EL CÓDIGO COMENTADO ES ANTERIOR A LA PRIMERA REFACTORIZACIÓN
	//fmt.Fprintln(rw, "Crea un usuario")

	/*rw.Header().Set("content-type", "application/json") // Para responder con json
	//rw.Header().Set("content-type", "text/xml") // Para responder con xml
	// No hay un tipo de dato específico para yaml

	// Obtener usuario/registro
	user := models.User{}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&user); err != nil {
		fmt.Fprintln(rw, http.StatusUnprocessableEntity)
	} else {
		db.Connect()
		user.Save() // Método creado por nosotros que inserta o actualiza dependiendo de si el usuario existe ya o no
		db.Close()
	}

	// Marshall devuelve 2 valores: Los valores transformados en tipo byte y un error
	output, _ := json.Marshal(user) // Para responder con json
	//output, _ := xml.Marshal(users) // Para responder con xml
	//output, _ := yaml.Marshal(users) // Para responder con yaml
	fmt.Fprintln(rw, string(output))*/

	// Obtener usuario/registro
	user := models.User{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&user); err != nil {
		models.SendUnprocessableEntity(rw)
	} else {
		user.Save()
		models.SendData(rw, user)
	}

}

func UpdateUser(rw http.ResponseWriter, r *http.Request) {
	// EL CÓDIGO COMENTADO ES ANTERIOR A LA PRIMERA REFACTORIZACIÓN
	//fmt.Fprintln(rw, "Actualiza un usuario")

	/*rw.Header().Set("content-type", "application/json") // Para responder con json
	//rw.Header().Set("content-type", "text/xml") // Para responder con xml
	// No hay un tipo de dato específico para yaml

	// Obtener registro
	user := models.User{}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&user); err != nil {
		fmt.Fprintln(rw, http.StatusUnprocessableEntity)
	} else {
		db.Connect()
		user.Save() // Método creado por nosotros que inserta o actualiza dependiendo de si el usuario existe ya o no
		db.Close()
	}

	// Marshall devuelve 2 valores: Los valores transformados en tipo byte y un error
	output, _ := json.Marshal(user) // Para responder con json
	//output, _ := xml.Marshal(users) // Para responder con xml
	//output, _ := yaml.Marshal(users) // Para responder con yaml
	fmt.Fprintln(rw, string(output))*/

	// Obtener registro
	var userId int64

	if user, err := getUserByRequest(r); err != nil {
		models.SendNotFound(rw)
	} else {
		userId = user.Id
	}

	user := models.User{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&user); err != nil {
		models.SendUnprocessableEntity(rw)
	} else {
		user.Id = userId
		user.Save()
		models.SendData(rw, user)
	}

}

func DeleteUser(rw http.ResponseWriter, r *http.Request) {
	// EL CÓDIGO COMENTADO ES ANTERIOR A LA PRIMERA REFACTORIZACIÓN
	//fmt.Fprintln(rw, "Elimina un usuario")

	/*rw.Header().Set("content-type", "application/json") // Para responder con json
	//rw.Header().Set("content-type", "text/xml") // Para responder con xml
	// No hay un tipo de dato específico para yaml

	// Obtener ID
	vars := mux.Vars(r) // Mos devuelve un mapa de tipo string
	userId, _ := strconv.Atoi(vars["id"])

	db.Connect()
	user, _ := models.GetUser(userId)
	user.Delete()
	db.Close()
	// Marshall devuelve 2 valores: Los valores transformados en tipo byte y un error
	output, _ := json.Marshal(user) // Para responder con json
	//output, _ := xml.Marshal(users) // Para responder con xml
	//output, _ := yaml.Marshal(users) // Para responder con yaml
	fmt.Fprintln(rw, string(output))*/

	if user, err := getUserByRequest(r); err != nil {
		models.SendNotFound(rw)
	} else {
		user.Delete()
		models.SendData(rw, user)
	}

}

// Función reutilizable para simplificar el código en el uso
// de GetUser tanto en la parte de editar como en la eliminar
func getUserByRequest(r *http.Request) (models.User, error) {

	// Obtener ID
	vars := mux.Vars(r) // Mos devuelve un mapa de tipo string
	userId, _ := strconv.Atoi(vars["id"])

	if user, err := models.GetUser(userId); err != nil {
		return *user, err
	} else {
		return *user, nil
	}
}
