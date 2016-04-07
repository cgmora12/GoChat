/*
Cliente de una estructura cliente-servidor

Uso:
go run clientePedro.go
*/

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// función para comprobar errores (ahorra escritura)
func chk(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	elegirOpcion()
}

func elegirOpcion() {
	salir := false
	for !salir {
		fmt.Println("\n-- Cliente GoChat --")

		fmt.Println("Eliga una opción: ")
		fmt.Println("1.- Login")
		fmt.Println("2.- Registro")
		fmt.Println("3.- Entrar modo cliente")
		fmt.Println("4.- Salir")
		fmt.Print("Opción elegida (introduzca el número): ")

		reader := bufio.NewReader(os.Stdin)
		opcionElegida, err := reader.ReadString('\n')
		chk(err)
		opcionElegida = strings.TrimRight(opcionElegida, "\r\n")

		switch opcionElegida {
		case "1":
			fmt.Println("- Login -")

		case "2":
			registro()

		case "3":
			fmt.Println("- Entrar modo cliente -")
			client()

		case "4":
			salir = true

		default:
			fmt.Println("Opción '", opcionElegida, "' desconocida. Introduzca una opción válida (1, 2, 3 o 4)")
		}
	}
}

// gestiona el modo cliente
func client() {
	conn, err := net.Dial("tcp", "localhost:1337") // llamamos al servidor
	chk(err)
	defer conn.Close() // es importante cerrar la conexión al finalizar

	fmt.Println("conectado a ", conn.RemoteAddr())

	keyscan := bufio.NewScanner(os.Stdin) // scanner para la entrada estándar (teclado)
	netscan := bufio.NewScanner(conn)     // scanner para la conexión (datos desde el servidor)

	for keyscan.Scan() { // escaneamos la entrada
		fmt.Fprintln(conn, keyscan.Text())         // enviamos la entrada al servidor
		netscan.Scan()                             // escaneamos la conexión
		fmt.Println("servidor: " + netscan.Text()) // mostramos mensaje desde el servidor
	}
}

func registro() {
	fmt.Println("- Registro -")

	// Se obtienen los datos de registro
	fmt.Print("Nombre de usuario: ")
	reader := bufio.NewReader(os.Stdin)
	nombreUsuario, err := reader.ReadString('\n')
	chk(err)
	nombreUsuario = strings.TrimRight(nombreUsuario, "\r\n")

	fmt.Print("Contraseña: ")
	reader = bufio.NewReader(os.Stdin)
	password, err := reader.ReadString('\n')
	chk(err)
	password = strings.TrimRight(password, "\r\n")

	// Se envian los datos al servidor
	conn, err := net.Dial("tcp", "localhost:1337") // Se llama al servidor
	chk(err)
	defer conn.Close() // Se cierra la conexión al finalizar

	fmt.Println("conectado a ", conn.RemoteAddr())

	var datos string = "Registro:" + nombreUsuario + ":" + password

	fmt.Fprintln(conn, datos) // Se envian los datos al servidor
}