package FWQ_Engine

import (
	"fmt"
	"os"
)

func main() {
	IpKafka := os.Args[1]
	PuertoKafka := os.Args[2]
	numeroVisitantes := os.Args[3]
	IpFWQWating := os.Args[4]
	PuertoWaiting := os.Args[5]
	fmt.Println("**Bienvenido al engine de la aplicación**")
	fmt.Println("La ip del apache kafka es el siguiente:" + IpKafka)
	fmt.Println("El puerto del apache kafka es el siguiente:" + PuertoKafka)
	fmt.Println("El número máximo de visitantes es el siguiente:" + numeroVisitantes)
	fmt.Println("La ip del servidor de espera es el siguiente:" + IpFWQWating)
	fmt.Println("El puerto del servidor de tiempo es el siguiente:" + PuertoWaiting)
}
