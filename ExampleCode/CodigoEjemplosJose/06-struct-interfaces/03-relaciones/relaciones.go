package main

import "fmt"

type User struct {
	nombre string
	email  string
	activo bool
}

type Student struct {
	user   User
	codigo string
}

// Relación de uno a muchos
type Curso struct {
	titulo string
	videos []Video
}

type Video struct {
	titulo string
	curso  Curso
}

func main() {

	// Relación de uno a uno
	/*
		jose := User{
			nombre: "Jose",
			email:  "jose@gmail.com",
			activo: true,
		}

		carlos := User{
			nombre: "Carlos",
			email:  "carlos@gmail.com",
			activo: true,
		}

		joseStudent := Student{
			user:   jose,
			codigo: "001arsd",
		}

		fmt.Println(carlos)
		fmt.Println(joseStudent.user.nombre)*/

	// Relación de uno a muchos
	video1 := Video{titulo: "01-Introducción"}
	video2 := Video{titulo: "02-Instalación"}

	curso := Curso{
		titulo: "Curso de GO",
		videos: []Video{video1, video2},
	}

	video1.curso = curso
	video2.curso = curso

	fmt.Println(video1.curso.titulo)

	for _, video := range curso.videos {
		fmt.Println(video.titulo)
	}

}
