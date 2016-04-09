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
	elegirOpcionMain()
}

func elegirOpcionMain() {
	salir := false
	for !salir {
		fmt.Println("\n\n-- Cliente GoChat --")

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
			login()

		case "2":
			registro()

		case "3":
			fmt.Println("\n- Entrar modo cliente -")
			client()

		case "4":
			salir = true

		default:
			fmt.Println("\nOpción '", opcionElegida, "' desconocida. Introduzca una opción válida (1, 2, 3 o 4)")
		}
	}
}

// gestiona el modo cliente
func client() {
	conn, err := net.Dial("tcp", "localhost:1337") // Se llama al servidor
	chk(err)
	defer conn.Close() // Se cierra la conexión al final

	fmt.Println("conectado a ", conn.RemoteAddr())

	keyscan := bufio.NewScanner(os.Stdin) // scanner para la entrada estándar (teclado)
	netscan := bufio.NewScanner(conn)     // scanner para la conexión (datos desde el servidor)

	for keyscan.Scan() { // Se escanea la entrada
		textoAEnviar := keyscan.Text()

		// Se comprueba si el mensaje enviado corresponde con algún método del servidor
		if strings.HasPrefix(textoAEnviar, "Registro:") || strings.HasPrefix(textoAEnviar, "Login:") {
			fmt.Println("No se puede enviar un mensaje con esa estructura")
		} else { // Si el mensaje recibido no se corresponde con ningún método del servidor
			fmt.Fprintln(conn, textoAEnviar)           // Se envia la entrada al servidor
			netscan.Scan()                             // Se escanea la conexión
			fmt.Println("servidor: " + netscan.Text()) // Se muestra el mensaje recibido desde el servidor
		}
	}
}

func registro() {
	fmt.Println("\n- Registro -")

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

	netscan := bufio.NewScanner(conn) // Se crea un scanner para la conexión (datos desde el servidor)
	netscan.Scan()                    // Se escanea la conexión
	fmt.Println(netscan.Text())       // Se muestra el mensaje desde el servidor
}

func login() {
	fmt.Println("\n- Login -")

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

	var datos string = "Login:" + nombreUsuario + ":" + password

	fmt.Fprintln(conn, datos) // Se envian los datos al servidor

	netscan := bufio.NewScanner(conn) // Se crea un scanner para la conexión (datos desde el servidor)
	netscan.Scan()                    // Se escanea la conexión

	textoRecibido := netscan.Text()
	fmt.Println(textoRecibido) // Se muestra el mensaje desde el servidor

	s := strings.Split(textoRecibido, ":")
	substringRespuesta := s[1]

	if substringRespuesta != "Error" {
		elegirOpcionChat(nombreUsuario, conn)
	}
}

func elegirOpcionChat(nombreUsuario string, conn net.Conn) {
	salir := false
	for !salir {
		fmt.Println("\n\n-- GoChat --")
		fmt.Println("-- Usuario:", nombreUsuario, "--")

		fmt.Println("Eliga una opción: ")
		fmt.Println("1.- Sala pública")
		fmt.Println("2.- Logout")
		fmt.Print("Opción elegida (introduzca el número): ")

		reader := bufio.NewReader(os.Stdin)
		opcionElegida, err := reader.ReadString('\n')
		chk(err)
		opcionElegida = strings.TrimRight(opcionElegida, "\r\n")

		switch opcionElegida {
		case "1":
			salaPublica(conn)

		case "2":
			salir = true

		default:
			fmt.Println("\nOpción '", opcionElegida, "' desconocida. Introduzca una opción válida (1 o 2)")
		}
	}
}

func salaPublica(conn net.Conn) {

	// Se crean dos canales (channels)
	done1 := make(chan bool)
	done2 := make(chan bool)

	// Goroutine para leer los mensajes
	go func() {
		netscan := bufio.NewScanner(conn) // Se crea un scanner para la conexión (datos desde el servidor)
		// Para mantener abierta la conexión
		for netscan.Scan() {
			textoRecibido := netscan.Text()
			fmt.Println(textoRecibido)
			fmt.Print("Escriba su mensaje: ")
		}
		// Para indicar a la función que la goroutine ya ha acabado.
		done1 <- true
	}()

	// Goroutine para escribir los mensajes
	go func() {
		fmt.Println("conectado a ", conn.RemoteAddr())

		keyscan := bufio.NewScanner(os.Stdin) // Se crea un scanner para la entrada estándar (teclado)

		fmt.Print("Escriba su mensaje: ")
		for keyscan.Scan() { // Se escanea la entrada
			fmt.Print("Escriba su mensaje: ")
			textoAEnviar := keyscan.Text()

			// Se comprueba si el mensaje enviado corresponde con algún método del servidor
			if strings.HasPrefix(textoAEnviar, "Registro:") || strings.HasPrefix(textoAEnviar, "Login:") {
				fmt.Println("No se puede enviar un mensaje con esa estructura")
			} else { // Si el mensaje recibido no se corresponde con ningún método del servidor
				fmt.Fprintln(conn, textoAEnviar) // Se envia la entrada al servidor
			}
		}
		// Para indicar a la función que la goroutine ya ha acabado.
		done2 <- true
	}()

	// Para que se espere a que las goroutines acaben.
	<-done1
	<-done2
}
