package main

import "fmt"

func precioVenta(valor float64) float64 {
	return valor + (valor * 0.18)
}

func hallarIGV(valor float64) float64 {
	return valor * 0.18
}

func main() {

	var valorVenta float64

	fmt.Println("Por favor, ingrese el valor de venta del producto:")
	fmt.Scanln(&valorVenta)

	igv := hallarIGV(valorVenta)
	precio := precioVenta(valorVenta)

	fmt.Printf("El IGV asociado a dicho precio de venta es: %f \n", igv)
	fmt.Printf("El precio de venta asociado a dicho precio de venta es: %f \n", precio)

}
