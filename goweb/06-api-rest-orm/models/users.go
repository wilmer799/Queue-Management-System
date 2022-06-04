package models

import "gorm/db"

// Responder con JSON
type User struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// Responder con XML
/*type User struct {
	Id       int64  `xml:"id"`
	Username string `xml:"username"`
	Password string `xml:"password"`
	Email    string `xml:"email"`
}*/

// Responder con YAML
/*type User struct {
	Id       int64  `yaml:"id"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Email    string `yaml:"email"`
}*/

type Users []User

// Función que se conecta la BD y migra el modelo user en la BD
// de forma que no tenemos que realizar/codificar la parte de SQL
func MigrarUser() {
	db.Database.AutoMigrate(User{}) // Migramos una estructura
}

/*const UserSchema string = `CREATE TABLE users (
	id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
	username VARCHAR(30) NOT NULL,
	password VARCHAR(100) NOT NULL,
	email VARCHAR(50),
	create_data TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)`

// Construir usuario
func NewUser(username, password, email string) *User {
	user := &User{Username: username, Password: password, Email: email}
	return user
}

// Crear usuario e inserta en la BD
func CreateUser(username, password, email string) *User {
	user := NewUser(username, password, email)
	//user.insert() // Deprecated
	user.Save()
	return user
}

// Insertar fila en la BD
func (user *User) insert() {
	sql := "INSERT INTO users SET username=?, password=?, email=?"
	result, _ := db.Exec(sql, user.Username, user.Password, user.Email)

	user.Id, _ = result.LastInsertId()

}

// Lista todas las filas de la tabla users
func ListUsers() (Users, error) {
	sql := "SELECT id, username, password, email FROM users"
	users := Users{}
	rows, err := db.Query(sql)

	for rows.Next() {
		user := User{}
		rows.Scan(&user.Id, &user.Username, &user.Password, &user.Email)
		users = append(users, user)
	}

	return users, err

}

// Obtener una sóla fila de la tabla users
func GetUser(id int) (*User, error) {
	user := NewUser("", "", "")

	sql := "SELECT id, username, password, email FROM users WHERE id=?"
	if rows, err := db.Query(sql, id); err != nil {
		return nil, err // Otra opción sería devolver un user vacío
	} else {

		for rows.Next() {
			rows.Scan(&user.Id, &user.Username, &user.Password, &user.Email)
		}

		return user, nil
	}
}

// Actualizar una fila de la tabla users
func (user *User) update() {
	sql := "UPDATE users SET username=?, password=?, email=? WHERE id=?"
	db.Exec(sql, user.Username, user.Password, user.Email, user.Id)
}

// Guardar o editar una fila de la tabla users
func (user *User) Save() {
	if user.Id == 0 {
		user.insert()
	} else {
		user.update()
	}
}

// Eliminar una fila de la tabla users
func (user *User) Delete() {
	sql := "DELETE FROM users WHERE id=?"
	db.Exec(sql, user.Id)
}*/
