// HE ESTADO INVESTIGANDO Y SE PUEDE DAR EL CASO DE QUE TE DE UN ERROR EL VSCODE AL IMPORTAR UN MÓDULO Y EN REALIDAD FUNCIONE ADECUADAMENTE.
// ESTO SE DEBE A QUE HAS ABIERTO EL VSCODE EN UNA RUTA QUE NO ES LA DEL PROYECTO ACTUAL SOBRE EL QUE ESTÁS TRABAJANDO.

// Un detalle a tener en cuenta es que si evitamos usar guiones en las rutas la importación de paquetes se hará de forma más automatizada
// y sino lo que tenemos que hacer es lo que hemos hecho aquí, que básicamente consiste en ejecutar el comando go mod init "nombre del módulo gestor de paquetes"
// y tras esto sólamente tendremos que hacer el import.
package main

//import "CursoGO/07-modularizacion/figuras"

//import "github.com/donvito/hellomod"
import "github.com/JoseHurtadoBaeza/figuras"

func main() {

	cua1 := figuras.Cuadrado{Lado: 10}
	figuras.Medidas(&cua1)

	//hellomod.SayHello()

	/*mensajes.Hola()
	mensajes.Imprimir()*/
	/*cua1 := figuras.Cuadrado{Lado: 8}
	cir1 := figuras.Circulo{Radio: 5}

	figuras.Medidas(&cua1)
	figuras.Medidas(&cir1)*/

	/*p1 := models.Persona{}
	p1.Constructor("Jose", 25)
	fmt.Println(p1)
	fmt.Println(p1.GetNombre())
	p1.SetNombre("Carlos")

	fmt.Println(p1)*/
}
