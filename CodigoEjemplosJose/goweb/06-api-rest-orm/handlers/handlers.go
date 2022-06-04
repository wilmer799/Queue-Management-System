package handlers

import (
	"encoding/json"
	"gorm/db"
	"gorm/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func GetUsers(rw http.ResponseWriter, r *http.Request) {
	// EL CÓDIGO COMENTADO ES ANTERIOR A LA REFACTORIZACIÓN FINAL
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
	fmt.Fprintln(rw, string(output))

	if users, err := models.ListUsers(); err != nil {
		models.SendNotFound(rw)
	} else {
		models.SendData(rw, users)
	}*/

	users := models.Users{}
	db.Database.Find(&users) // Le indicamos a find dónde vamos a recuperar los datos
	sendData(rw, users, http.StatusOK)

}

func GetUser(rw http.ResponseWriter, r *http.Request) {
	// EL CÓDIGO COMENTADO ES ANTERIOR A LA REFACTORIZACIÓN FINAL
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

	/*if user, err := getUserByRequest(r); err != nil {
		models.SendNotFound(rw)
	} else {
		models.SendData(rw, user)
	}*/

	if user, err := getUserById(r); err != nil {
		sendError(rw, http.StatusNotFound)
	} else {
		sendData(rw, user, http.StatusOK)
	}

}

// Función reutilizable para simplificar el código en el uso
// de GetUser tanto en la parte de editar como en la eliminar
func getUserById(r *http.Request) (models.User, *gorm.DB) {

	vars := mux.Vars(r)
	userId, _ := strconv.Atoi(vars["id"])

	user := models.User{}
	if err := db.Database.First(&user, userId); err.Error != nil { // Nos va a devolver un dato en función del id que le pasemos
		return user, err
	} else {
		return user, nil
	}

}

func CreateUser(rw http.ResponseWriter, r *http.Request) {
	// EL CÓDIGO COMENTADO ES ANTERIOR A LA REFACTORIZACIÓN FINAL
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
	/*	user := models.User{}
		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&user); err != nil {
			models.SendUnprocessableEntity(rw)
		} else {	sendData(rw, user, http.StatusOK)

	*/

	// Obtener usuario/registro
	user := models.User{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&user); err != nil {
		sendError(rw, http.StatusUnprocessableEntity)
	} else {
		db.Database.Save(&user)
		sendData(rw, user, http.StatusCreated)
	}

}

func UpdateUser(rw http.ResponseWriter, r *http.Request) {
	// EL CÓDIGO COMENTADO ES ANTERIOR A LA REFACTORIZACIÓN FINAL
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
	  fmt.Fprintln(rw, string(output	sendData(rw, user, http.StatusOK)

	/*	var userId int64

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
	*/

	var userId int64

	if user_ant, err := getUserById(r); err != nil {
		sendError(rw, http.StatusNotFound)
	} else {
		userId = user_ant.Id

		// Obtener usuario/registro
		user := models.User{}
		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&user); err != nil {
			sendError(rw, http.StatusUnprocessableEntity)
		} else {
			user.Id = userId
			db.Database.Save(&user)
			sendData(rw, user, http.StatusOK)
		}
	}

}

func DeleteUser(rw http.ResponseWriter, r *http.Request) {
	// EL CÓDIGO COMENTADO ES ANTERIOR A LA REFACTORIZACIÓN FINAL
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

	/*	if user, err := getUserByRequest(r); err != nil {
			models.SendNotFound(rw)
		} else {
			user.Delete()
			models.SendData(rw, user)
		}*/

	if user, err := getUserById(r); err != nil {
		sendError(rw, http.StatusNotFound)
	} else {
		db.Database.Delete(&user)
		sendData(rw, user, http.StatusOK)
	}

}

// DEPRECATED
// Función reutilizable para simplificar el código en el uso
// de GetUser tanto en la parte de editar como en la eliminar
/*func getUserByRequest(r *http.Request) (models.User, error) {

	// Obtener ID
	vars := mux.Vars(r) // Mos devuelve un mapa de tipo string
	userId, _ := strconv.Atoi(vars["id"])

	if user, err := models.GetUser(userId); err != nil {
		return *user, err
	} else {
		return *user, nil
	}
}*/
