package main

import (
	"fmt"
	"math"
)

type Figura interface {
	area() int
	perimetro() int
}

type Cuadrado struct {
	ancho  int
	altura int
}

type Circulo struct {
	radio int
}

func (c *Cuadrado) area() int {
	return c.ancho * c.altura
}

func (c *Cuadrado) perimetro() int {
	return 2*c.ancho + 2*c.altura
}

func (c *Circulo) area() float64 {
	return math.Pi * float64(c.radio*c.radio)
}

func (c *Circulo) perimetro() float64 {
	return 2 * math.Pi * float64(c.radio)
}

func main() {

	cuadrado1 := Cuadrado{
		ancho:  5,
		altura: 10,
	}

	circulo1 := Circulo{
		radio: 7,
	}

	fmt.Println("El área del cuadrado es: ", cuadrado1.area())
	fmt.Println("El perímetro del cuadrado es: ", cuadrado1.perimetro())

	fmt.Println("El área del círculo es: ", circulo1.area())
	fmt.Println("El perímetro del círculo es: ", circulo1.perimetro())

}
