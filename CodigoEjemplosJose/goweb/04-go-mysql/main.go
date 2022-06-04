package main

import (
	"fmt"
	"gomysql/db"
	"gomysql/models"
)

func main() {
	db.Connect()
	//db.Ping() // Verificamos la conexión

	//fmt.Println(db.ExistsTable("users"))
	//db.CreateTable(models.UserSchema, "users")

	//user := models.CreateUser("jose", "jose123", "jose@gmail.com")
	user := models.CreateUser("carlos", "carlos123", "carlos@gmail.com")
	fmt.Println(user)

	// Recuperamos todos los usuarios de la tabla users
	//users := models.ListUsers()
	//fmt.Println(users)

	// Recuperamos sólo un usuario de la tabla users
	//user := models.GetUser(2)
	//fmt.Println(user)
	/*
		user.Username = "juan"
		user.Password = "juan789"
		user.Email = "juan@gmail.com"
		user.Save()*/

	//user.Delete()

	//db.TruncateTable("users") // Elimina todas las filas de la tabla indicada
	fmt.Println(models.ListUsers())

	db.Close()
	// Si la ejecutamos después de cerrar la BD nos dará un panic
	//db.Ping()
}
