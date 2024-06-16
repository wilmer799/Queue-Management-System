package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// username:password@tcp(localhost:3306)/database
const url = "root:1234@tcp(localhost:3306)/goweb_db"

// Guarda la conexión
var db *sql.DB

// Realiza la conexión
func Connect() {

	connection, err := sql.Open("mysql", url)
	if err != nil {
		panic(err)
	}

	fmt.Println("Conexión exitosa con la BD")
	db = connection

}

// Cierra la conexión
func Close() {
	db.Close()
}

// Verificar la conexión
func Ping() {
	if err := db.Ping(); err != nil {
		panic(err)
	}
}

// Verifica si una tabla existe o no
func ExistsTable(tableName string) bool {

	sql := fmt.Sprintf("SHOW TABLES LIKE '%s'", tableName)
	// rows, err := db.Query(sql) // Forma estándar
	rows, err := Query(sql) // Utilizando nuestra nueva función
	if err != nil {
		fmt.Println("Error:", err)
	}

	return rows.Next() // Next devuelve un valor booleano

}

// Crea una tabla
func CreateTable(schema string, name string) {

	if !ExistsTable(name) {
		_, err := Exec(schema)
		if err != nil {
			fmt.Println(err)
		}
	}

}

// Función que elimina todas las filas de una tabla
func TruncateTable(tablename string) {
	sql := fmt.Sprintf("TRUNCATE %s", tablename)
	Exec(sql)
}

// Polimorfismo de Exec
func Exec(query string, args ...interface{}) (sql.Result, error) {

	Connect() // Refactorización del código
	result, err := db.Exec(query, args...)
	Close() // Refactorización del código
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

// Polimorfismo de Query
func Query(query string, args ...interface{}) (*sql.Rows, error) {
	Connect() // Refactorización del código
	rows, err := db.Query(query, args...)
	Close() // Refactorización del código
	if err != nil {
		fmt.Println(err)
	}

	return rows, err
}
