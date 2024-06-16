package main

import "fmt"

// Struct persona
type Persona struct {
	nombre string
	edad   int
}

// Métodos
func (p *Persona) imprimir() {

	fmt.Printf("Nombre: %s\nEdad: %d\n", p.nombre, p.edad)

}

func (p *Persona) editNombre(nuevoNombre string) {

	p.nombre = nuevoNombre

}

// Herencia
type Empleado struct {
	Persona
	sueldo float64
}

func main() {

	p1 := Persona{"Jose", 25}

	//fmt.Println(p1)
	p1.imprimir()

	//p1.nombre = "Carlos"
	p1.editNombre("Carlos")

	//fmt.Println(p1)
	p1.imprimir()

	p2 := Persona{
		nombre: "Pepe",
		edad:   62,
	}

	//fmt.Println(p2)
	p2.imprimir()
	p2.editNombre("Juan")
	p2.imprimir()

	em1 := Empleado{
		sueldo: 500,
	}
	em1.nombre = "Pedro"
	em1.edad = 30
	em1.imprimir()
	fmt.Println(em1) // Se muestra con dos llaves por la herencia

}
