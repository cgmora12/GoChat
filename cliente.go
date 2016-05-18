/*
Cliente de una estructura cliente-servidor

Uso:
go run clientePedro.go
*/

package main

import (
	"bufio"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/codahale/chacha20"
	//"github.com/dchest/chacha20"
	"github.com/howeyc/gopass"
	//"github.com/tang0th/go-chacha20/chacha"
	//"golang.org/x/crypto/bcrypt"
	"math/big"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var cifradorPublico cipher.Stream
var errorPublico error

// función para comprobar errores
func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// cifrado para la sala pública
	// (se utiliza una clave predefinida puesto que se trata de una sala pública en la que cualquier usuario puede entrar y ver los mensajes en claro)
	cifradorPublico, errorPublico = chacha20.New([]byte("12345678123456781234567812345678"), []byte("nonce123"))
	checkError(errorPublico)

	elegirOpcionMain()
}

func elegirOpcionMain() {

	salir := false
	for !salir {

		fmt.Println("\n\n-- Cliente GoChat --")
		fmt.Println("1.- Login")
		fmt.Println("2.- Registro")
		fmt.Println("3.- Salir")
		fmt.Print("Elija una opción: ")

		reader := bufio.NewReader(os.Stdin)
		opcionElegida, err := reader.ReadString('\n')
		checkError(err)
		opcionElegida = strings.TrimRight(opcionElegida, "\r\n")

		switch opcionElegida {
		case "1":
			login()

		case "2":
			registro()

		case "3":
			salir = true
			fmt.Println("\nHasta pronto!")

		default:
			fmt.Println("\nOpción '" + opcionElegida + "' desconocida. Introduzca una opción válida (1, 2 ó 3)")
		}
	}

}

func registro() {

	fmt.Println("\n- Registro -")

	// Se obtienen los datos de registro validados

	var reader *bufio.Reader
	var nombreUsuario string
	var nombreCompleto string
	var pais string
	var provincia string
	var localidad string
	var email string
	var err error
	re := regexp.MustCompile("^[a-zA-Z0-9_ ]*$")
	reEmail := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

	nombreOk := false
	for !nombreOk {
		fmt.Print("Nombre de usuario: ")
		reader = bufio.NewReader(os.Stdin)
		nombreUsuario, err = reader.ReadString('\n')
		checkError(err)
		nombreUsuario = strings.TrimRight(nombreUsuario, "\r\n")

		if len(nombreUsuario) >= 5 && len(nombreUsuario) <= 20 {
			if re.MatchString(nombreUsuario) {
				nombreOk = true
			} else {
				fmt.Println("El nombre de usuario debe ser alfa-numérico.\n")
			}
		} else {
			fmt.Println("El nombre de usuario debe contener entre 5 y 20 caracteres.\n")
		}
	}

	fmt.Print("Contraseña: ")
	//reader = bufio.NewReader(os.Stdin)
	//password, err := reader.ReadString('\n')
	//checkError(err)
	//password = strings.TrimRight(password, "\r\n")
	pass, err := gopass.GetPasswdMasked()
	checkError(err)
	password := string(pass[:])

	nombreCompletoOk := false
	for !nombreCompletoOk {
		fmt.Print("Nombre completo: ")
		reader = bufio.NewReader(os.Stdin)
		nombreCompleto, err = reader.ReadString('\n')
		checkError(err)
		nombreCompleto = strings.TrimRight(nombreCompleto, "\r\n")

		if len(nombreCompleto) >= 5 && len(nombreCompleto) <= 40 {
			if re.MatchString(nombreCompleto) {
				nombreCompletoOk = true
			} else {
				fmt.Println("El nombre debe ser alfa-numérico.\n")
			}
		} else {
			fmt.Println("El nombre debe contener entre 5 y 40 caracteres.\n")
		}
	}

	paisOk := false
	for !paisOk {
		fmt.Print("País: ")
		reader = bufio.NewReader(os.Stdin)
		pais, err = reader.ReadString('\n')
		checkError(err)
		pais = strings.TrimRight(pais, "\r\n")

		if len(pais) <= 40 {
			if re.MatchString(pais) {
				paisOk = true
			} else {
				fmt.Println("El país debe ser una cadena alfa-numérica.\n")
			}
		} else {
			fmt.Println("El país debe contener menos de 40 caracteres.\n")
		}
	}

	provinciaOk := false
	for !provinciaOk {
		fmt.Print("Província: ")
		reader = bufio.NewReader(os.Stdin)
		provincia, err = reader.ReadString('\n')
		checkError(err)
		provincia = strings.TrimRight(provincia, "\r\n")

		if len(provincia) <= 40 {
			if re.MatchString(provincia) {
				provinciaOk = true
			} else {
				fmt.Println("La provincia debe ser una cadena alfa-numérica.\n")
			}
		} else {
			fmt.Println("La provincia debe contener menos de 40 caracteres.\n")
		}
	}

	localidadOk := false
	for !localidadOk {
		fmt.Print("Localidad: ")
		reader = bufio.NewReader(os.Stdin)
		localidad, err = reader.ReadString('\n')
		checkError(err)
		localidad = strings.TrimRight(localidad, "\r\n")

		if len(localidad) <= 40 {
			if re.MatchString(localidad) {
				localidadOk = true
			} else {
				fmt.Println("La localidad debe ser una cadena alfa-numérica.\n")
			}
		} else {
			fmt.Println("La localidad debe contener menos de 40 caracteres.\n")
		}
	}

	emailOk := false
	for !emailOk {
		fmt.Print("Email: ")
		reader = bufio.NewReader(os.Stdin)
		email, err = reader.ReadString('\n')
		checkError(err)
		email = strings.TrimRight(email, "\r\n")

		if len(email) <= 40 {
			if len(email) == 0 || reEmail.MatchString(email) {
				emailOk = true
			} else {
				fmt.Println("El formato del email no es correcto (email@company.net).\n")
			}
		} else {
			fmt.Println("El email debe contener menos de 40 caracteres.\n")
		}
	}

	// TLS
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	checkError(err)
	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}

	// Se envian los datos al servidor
	conn, err := tls.Dial("tcp", "localhost:1337", &config) // Se llama al servidor
	checkError(err)
	defer conn.Close() // Se cierra la conexión al finalizar

	//fmt.Println("conectado a ", conn.RemoteAddr())

	var datos string = "Registro#&" + nombreUsuario + "#&" + password + "#&" + nombreCompleto + "#&" + pais + "#&" + provincia + "#&" + localidad + "#&" + email

	fmt.Fprintln(conn, datos) // Se envian los datos al servidor

	netscan := bufio.NewScanner(conn)  // Se crea un scanner para la conexión (datos desde el servidor)
	netscan.Scan()                     // Se escanea la conexión
	fmt.Println("\n" + netscan.Text()) // Se muestra el mensaje desde el servidor
}

func login() {
	fmt.Println("\n- Login -")

	// Se obtienen los datos de registro

	re := regexp.MustCompile("^[a-zA-Z0-9_ ]*$")
	var nombreUsuario string
	var err error

	nombreOk := false
	for !nombreOk {
		fmt.Print("Nombre de usuario: ")
		reader := bufio.NewReader(os.Stdin)
		nombreUsuario, err = reader.ReadString('\n')
		checkError(err)
		nombreUsuario = strings.TrimRight(nombreUsuario, "\r\n")

		if re.MatchString(nombreUsuario) {
			nombreOk = true
		} else {
			fmt.Println("El nombre de usuario debe ser alfa-numérico.\n")
		}
	}

	fmt.Print("Contraseña: ")
	//reader = bufio.NewReader(os.Stdin)
	//password, err := reader.ReadString('\n')
	//checkError(err)
	//password = strings.TrimRight(password, "\r\n")
	pass, err := gopass.GetPasswdMasked()
	checkError(err)
	password := string(pass[:])

	// TLS
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	checkError(err)
	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}

	// Se envian los datos al servidor
	conn, err := tls.Dial("tcp", "localhost:1337", &config) // Se llama al servidor
	checkError(err)
	defer conn.Close() // Se cierra la conexión al finalizar

	//fmt.Println("conectado a ", conn.RemoteAddr())

	var datos string = "Login#&" + nombreUsuario + "#&" + password

	fmt.Fprintln(conn, datos) // Se envian los datos al servidor

	netscan := bufio.NewScanner(conn) // Se crea un scanner para la conexión (datos desde el servidor)
	netscan.Scan()                    // Se escanea la conexión

	textoRecibido := netscan.Text()
	textoAMostrar := strings.Replace(textoRecibido, "#&", ":", -1) // -1 significa que no hay límite de coincidencias para reemplazar.
	fmt.Println("\n" + textoAMostrar)                              // Se muestra el mensaje desde el servidor

	s := strings.Split(textoRecibido, "#&")
	substringRespuesta := s[0]

	if substringRespuesta != "Error" {
		elegirOpcionChat(nombreUsuario, conn)
	}
}

func elegirOpcionChat(nombreUsuario string, conn net.Conn) {
	salir := false
	for !salir {
		fmt.Println("\n\n-- GoChat --")
		fmt.Println("-- Usuario: " + nombreUsuario + " --\n")

		fmt.Println("1.- Sala pública")
		fmt.Println("2.- Salas privadas")
		fmt.Println("3.- Ver perfiles de usuarios")
		fmt.Println("4.- Logout")
		fmt.Print("Elija una opción: ")

		reader := bufio.NewReader(os.Stdin)
		opcionElegida, err := reader.ReadString('\n')
		checkError(err)
		opcionElegida = strings.TrimRight(opcionElegida, "\r\n")

		switch opcionElegida {
		case "1":
			salaPublica(conn, nombreUsuario)

		case "2":
			salasPrivadas(conn, nombreUsuario)

		case "3":
			verPerfiles(conn, nombreUsuario)

		case "4":
			salir = true
			fmt.Println("\nLogout correcto")

		default:
			fmt.Println("\nOpción '" + opcionElegida + "' desconocida. Introduzca una opción válida (1, 2, 3 o 4)")
		}
	}
}

func salaPublica(conn net.Conn, nombreUsuario string) {

	// Se crean dos canales (channels)
	done1 := make(chan bool)
	done2 := make(chan bool)
	quit := make(chan bool)

	fmt.Println("\n\n-- Sala pública --")
	fmt.Println("-- Usuario: " + nombreUsuario + " --")
	fmt.Println("Escriba 'Salir' para volver al menú de usuario")
	netscan := bufio.NewScanner(conn) // Se crea un scanner para la conexión (datos desde el servidor)

	// Goroutine para leer los mensajes
	go func() {
		// Para mantener abierta la conexión
		for netscan.Scan() {
			select {
			case _, ok := <-quit:
				if ok {
					done1 <- true
					return
				} else {
					fmt.Println("Error: Canal 'quit' cerrado")
				}
			default:
				textoRecibidoCodificado := netscan.Text()
				//fmt.Println("\nTexto recibido codificado: *" + strings.Split(textoRecibidoCodificado, ": ")[2])
				textoRecibido := []byte(strings.Split(textoRecibidoCodificado, ": ")[2])
				if strings.Split(textoRecibidoCodificado, ": ")[0] == "Sala pública" {
					cifradorPublico.XORKeyStream(textoRecibido, textoRecibido) //dst, src
				}
				fmt.Println("\n" + strings.Split(textoRecibidoCodificado, ": ")[0] + ": " + strings.Split(textoRecibidoCodificado, ": ")[1] +
					": " + string(textoRecibido))
				fmt.Print("Continúe su mensaje: ")
			}
		}
		// Para indicar a la función que la goroutine ya ha acabado.
		done1 <- true
	}()

	// Goroutine para escribir los mensajes
	go func() {
		//fmt.Println("conectado a ", conn.RemoteAddr())

		keyscan := bufio.NewScanner(os.Stdin) // Se crea un scanner para la entrada estándar (teclado)

		fmt.Print("Escriba su mensaje: ")
		for keyscan.Scan() { // Se escanea la entrada
			fmt.Print("Escriba su mensaje: ")
			textoAEnviar := keyscan.Text()

			// Se comprueba si el mensaje enviado corresponde con algún método del servidor
			if strings.Contains(textoAEnviar, "#&") {
				fmt.Println("Error: Secuencia '#&' inválida")
				fmt.Print("Escriba su mensaje: ")
			} else if textoAEnviar == "Salir" {
				fmt.Fprintln(conn, "SalirChat#&")
				quit <- true
				done2 <- true
				return
			} else { // Si el mensaje recibido no se corresponde con ningún método del servidor
				textoAEnviarCodificado := []byte(textoAEnviar)
				cifradorPublico.XORKeyStream(textoAEnviarCodificado, textoAEnviarCodificado) //dst, src
				textoPreparado := string(textoAEnviarCodificado)
				fmt.Fprintln(conn, textoPreparado) // Se envia la entrada al servidor
			}
		}
		// Para indicar a la función que la goroutine ya ha acabado.
		done2 <- true
	}()

	// Para que se espere a que las goroutines acaben.
	<-done1
	<-done2
	close(quit)
}

func salasPrivadas(conn net.Conn, nombreUsuario string) {

	// Primero se envia un mensaje al servidor para obtener los usuarios logueados actualmente.
	fmt.Fprintln(conn, "GetLogueados#&")
	netscan := bufio.NewScanner(conn) // Se crea un scanner para la conexión (datos desde el servidor)
	textoRecibido := ""

	for netscan.Scan() {
		textoRecibido = netscan.Text()
		//fmt.Println(textoRecibido)

		if strings.HasPrefix(textoRecibido, "GetLogueados#&") {
			break
		}
	}

	usuariosLogueados := strings.Split(textoRecibido, "#&")
	// Es -2 porque hay uno inicial con "GetLogueados:" y con los ":" al final del string, se obtiene una última posición vacía
	numUsuarios := len(usuariosLogueados) - 2

	// El siguiente paso es elegir el usuario con quien se quiere hablar.
	salir := false
	for !salir {
		fmt.Println("\n\n-- Salas privadas --")
		fmt.Println("-- Usuario: " + nombreUsuario + " --")

		fmt.Println("Elija el usuario con quien quiera hablar:\n")

		i := 1
		if numUsuarios == 0 {
			fmt.Println("No hay ningún usuario más logueado.")

		} else {
			for ; i <= numUsuarios; i++ {
				// usuariosLogueados[i] porque me hay que saltarse la posición 0 ("GetLogueados:")
				fmt.Println("Usuario ", i, ".- ", usuariosLogueados[i])
			}
		}
		fmt.Println(i, ".- Salir de las salas privadas")
		fmt.Print("Opción elegida (introduzca el número): ")

		reader := bufio.NewReader(os.Stdin)
		opcionElegida, err := reader.ReadString('\n')
		checkError(err)
		opcionElegida = strings.TrimRight(opcionElegida, "\r\n")

		opcionElegidaInt, err := strconv.Atoi(opcionElegida)
		if err != nil {
			fmt.Println("\nOpción '" + opcionElegida + "' desconocida. Introduzca una opción válida")

		} else if opcionElegidaInt > 0 && opcionElegidaInt <= numUsuarios {
			usuarioElegido := usuariosLogueados[opcionElegidaInt]
			//fmt.Println("Elegido: " + opcionElegidaInt + ".- " + usuarioElegido)
			entrarSalaPrivada(conn, nombreUsuario, usuarioElegido)

			// Al salir de la sala privada no se vuelven a pedir los usuarios logueados.
			// Por lo que, si mientras se está en una sala privada se conecta otro usuario,
			// sólo con salir de la sala privada actual no se mostrará dicho usuario.
			// Así que, al salir de una sala privada, se sale al menú de usuario (no al menú de las salas privadas).
			salir = true

		} else if opcionElegidaInt == i {
			// Es la opción para salir de las salas privadas
			salir = true

		} else {
			fmt.Println("\nOpción '" + opcionElegida + "' desconocida. Introduzca una opción válida")
		}
	}
}

func entrarSalaPrivada(conn net.Conn, esteUsuario string, usuarioElegido string) {

	// Se crean dos canales (channels)
	done1 := make(chan bool)
	done2 := make(chan bool)
	quit := make(chan bool)

	fmt.Println("\n\n-- Sala privada con " + usuarioElegido + " --")
	fmt.Println("-- Usuario: " + esteUsuario + " --")
	fmt.Println("Escriba 'Salir' para volver al menú de usuario")
	fmt.Println("Espere a que se conecte el otro usuario...")

	netscan := bufio.NewScanner(conn) // Se crea un scanner para la conexión (datos desde el servidor)

	session_key := intercambioDeClaves(conn, esteUsuario, usuarioElegido, netscan)

	//fmt.Println("Clave de sesion: " + string(session_key))
	cifrador, err := chacha20.New(session_key, []byte("nonce123"))
	checkError(err)

	/* Como codificar
	out := []byte("hola")
	cifrador.XORKeyStream(out, out) //dst, src
	*/

	// Goroutine para leer los mensajes
	go func() {
		// Para mantener abierta la conexión
		for netscan.Scan() {
			select {
			case _, ok := <-quit:
				if ok {
					done1 <- true
					return
				} else {
					fmt.Println("Error: Canal 'quit' cerrado")
				}
			default:
				textoRecibidoCodificado := netscan.Text()

				if textoRecibidoCodificado == "Salir" {

					fmt.Println("\nEl usuario ha abandonado el chat...")
					fmt.Println("Escriba: Salir")

				} else {

					//fmt.Println("\nTexto recibido codificado: " + strings.Split(textoRecibidoCodificado, ":")[2])
					textoRecibido := []byte(strings.Split(textoRecibidoCodificado, ": ")[2])
					if strings.Split(textoRecibidoCodificado, ": ")[0] == "Sala privada" {
						cifrador.XORKeyStream(textoRecibido, textoRecibido) //dst, src
					}
					fmt.Println("\n" + strings.Split(textoRecibidoCodificado, ": ")[0] + ": " + strings.Split(textoRecibidoCodificado, ": ")[1] +
						": " + string(textoRecibido))
					fmt.Print("Continúe su mensaje: ")

				}
			}
		}
		// Para indicar a la función que la goroutine ya ha acabado.
		done1 <- true
	}()

	// Goroutine para escribir los mensajes
	go func() {
		//fmt.Println("conectado a ", conn.RemoteAddr())

		keyscan := bufio.NewScanner(os.Stdin) // Se crea un scanner para la entrada estándar (teclado)

		fmt.Print("Escriba su mensaje: ")
		for keyscan.Scan() { // Se escanea la entrada
			fmt.Print("Escriba su mensaje: ")
			textoAEnviar := keyscan.Text()

			// Se comprueba si el mensaje enviado corresponde con algún método del servidor
			if strings.Contains(textoAEnviar, "#&") {
				fmt.Println("Error: Secuencia '#&' inválida")
				fmt.Print("Escriba su mensaje: ")
			} else if textoAEnviar == "Salir" {
				fmt.Fprintln(conn, "SalirChat#&"+usuarioElegido)
				quit <- true
				done2 <- true
				return
			} else { // Si el mensaje recibido no se corresponde con ningún método del servidor
				textoAEnviarCodificado := []byte(textoAEnviar)
				cifrador.XORKeyStream(textoAEnviarCodificado, textoAEnviarCodificado) //dst, src
				textoPreparado := "SalaPrivada#&" + esteUsuario + "#&" + usuarioElegido + "#&" + string(textoAEnviarCodificado)
				fmt.Fprintln(conn, textoPreparado) // Se envia la entrada al servidor
			}
		}
		// Para indicar a la función que la goroutine ya ha acabado.
		done2 <- true
	}()

	// Para que se espere a que las goroutines acaben.
	<-done1
	<-done2
	close(quit)
}

func intercambioDeClaves(conn net.Conn, esteUsuario string, usuarioElegido string, netscan *bufio.Scanner) []byte {

	// Intercambio de claves
	cli_keys, err := rsa.GenerateKey(rand.Reader, 1024) // generamos un par de claves (privada, pública) para el cliente
	checkError(err)
	cli_keys.Precompute() // aceleramos su uso con un precálculo

	cli_keys_json, err := json.Marshal(cli_keys.PublicKey)
	checkError(err)
	//fmt.Println(string(cli_keys_json))

	clavePub := "Clave#&" + esteUsuario + "#&" + usuarioElegido + "#&" + string(cli_keys_json)
	fmt.Fprintln(conn, clavePub) // Se envia la entrada al servidor

	var claveRecibida string
	for netscan.Scan() {

		claveRecibida = netscan.Text()
		if strings.HasPrefix(claveRecibida, "Clave:") {

			clavePub := "Clave#&" + esteUsuario + "#&" + usuarioElegido + "#&" + string(cli_keys_json)
			fmt.Fprintln(conn, clavePub) // Se envia la entrada al servidor
			break
		}

	}

	//fmt.Println("Mensaje recibido:" + claveRecibida)
	cli_pub := strings.Split(string(claveRecibida), "Clave:")[1]
	//fmt.Println("Clave recibida:" + cli_pub)

	var cli_pub_key rsa.PublicKey
	cli_pub_trozo1 := strings.Split(cli_pub, ":")[1]
	//fmt.Println("Trozo1:" + cli_pub_trozo1)
	cli_pub_trozo2 := strings.Split(cli_pub_trozo1, ",")[0]
	//fmt.Println("Trozo2:" + cli_pub_trozo2)
	//fmt.Println("E:" + strings.Split(strings.Split(cli_pub, ":")[2], "}")[0])
	cli_pub_key.E, err = strconv.Atoi(strings.Split(strings.Split(cli_pub, ":")[2], "}")[0])
	bigint := new(big.Int)
	bigint.SetString(cli_pub_trozo2, 10)
	cli_pub_key.N = bigint

	//cli_token := make([]byte, 48) // 384 bits (256 bits de clave + 128 bits para el IV)
	buff := make([]byte, 256)   // contendrá el token cifrado con clave pública (puede ocupar más que el texto en claro)
	cli_token := randString(32) // generación del token aleatorio para el cliente

	//fmt.Println("token creado ", cli_token)

	// ciframos el token del cliente con la clave pública del otro cliente
	enctoken, err := rsa.EncryptPKCS1v15(rand.Reader, &cli_pub_key, []byte(cli_token))
	checkError(err)

	//fmt.Println("token cifrado ", string(enctoken)) // Se envia la entrada al servidor
	enctokenString := strings.Replace(string(enctoken), "\n", "-----n", -1)

	//fmt.Println("token cifrado sin endlines ", enctokenString) // Se envia la entrada al servidor
	// Falla al encriptar y poner saltos de linea
	tokenCliente := "Token#&" + esteUsuario + "#&" + usuarioElegido + "#&" + enctokenString
	fmt.Fprintln(conn, tokenCliente) // Se envia la entrada al servidor

	var tokenRecibido string
	for netscan.Scan() {

		tokenRecibido = netscan.Text()
		if strings.HasPrefix(tokenRecibido, "Token:") {
			//fmt.Println("Token recibido:" + tokenRecibido)
			break
		}

	}

	//fmt.Println("Token recibido sin prefijo: " + strings.Split(tokenRecibido, "Token:")[1])
	//fmt.Println("Token recibido sin prefijo ni endlines: " + strings.Replace(strings.Split(tokenRecibido, "Token:")[1], "-----n", "\n", -1))

	buff = []byte(strings.Replace(strings.Split(tokenRecibido, "Token:")[1], "-----n", "\n", -1))

	// desciframos el token del otro cliente con nuestra clave privada
	session_key, err := rsa.DecryptPKCS1v15(rand.Reader, cli_keys, buff)
	checkError(err)

	// realizamos el XOR entre ambos tokens (cliente y servidor acaban con la misma clave de sesión)
	var i int
	for i = 0; i < len(cli_token); i++ {
		session_key[i] ^= cli_token[i]
	}

	return session_key
}

func verPerfiles(conn net.Conn, nombreUsuario string) {
	salir := false
	for !salir {
		fmt.Println("\n\n-- Perfiles de usuarios --")
		fmt.Println("-- Usuario: " + nombreUsuario + " --\n")
		fmt.Println("1.- Ver todos los usuarios")
		fmt.Println("2.- Buscar usuarios")
		fmt.Println("3.- Volver atrás")
		fmt.Print("Elija una opción: ")

		reader := bufio.NewReader(os.Stdin)
		opcionElegida, err := reader.ReadString('\n')
		checkError(err)
		opcionElegida = strings.TrimRight(opcionElegida, "\r\n")

		switch opcionElegida {
		case "1":
			verTodosPerfiles(conn)

		case "2":
			buscarUsuarios(conn)

		case "3":
			salir = true

		default:
			fmt.Println("\nOpción '" + opcionElegida + "' desconocida. Introduzca una opción válida (1, 2 o 3)")
		}
	}
}

func verTodosPerfiles(conn net.Conn) {
	// Primero se envia un mensaje al servidor para obtener los usuarios registrados actualmente.
	fmt.Fprintln(conn, "VerTodosPerfiles#&")
	netscan := bufio.NewScanner(conn) // Se crea un scanner para la conexión (datos desde el servidor)
	textoRecibido := ""

	for netscan.Scan() {
		textoRecibido = netscan.Text()
		//fmt.Println(textoRecibido)

		if strings.HasPrefix(textoRecibido, "VerTodosPerfiles#&") {
			break
		}
	}

	usuarios := strings.Split(textoRecibido, "#&")
	// Es -2 porque hay uno inicial con "GetLogueados:" y con los ":" al final del string, se obtiene una última posición vacía
	numUsuarios := len(usuarios) - 2
	//fmt.Println(usuarios)

	// El siguiente paso es elegir el usuario con quien se quiere hablar.
	fmt.Println("\nUsuarios:")

	i := 1
	if numUsuarios == 0 {
		fmt.Println("No hay ningún usuario registrado.")

	} else {
		for ; i <= numUsuarios; i++ {
			// usuarios[i] porque me hay que saltarse la posición 0 ("GetLogueados:")
			textoAMostrar := strings.Replace(usuarios[i], "-", "\n", -1) // -1 significa que no hay límite de coincidencias para reemplazar.
			fmt.Println("\nUsuario ", i, "\n", textoAMostrar)
		}
	}
}

func buscarUsuarios(conn net.Conn) {
	// Primero se le pide al usuario que escriba lo que quiere buscar.

	re := regexp.MustCompile("^[a-zA-Z0-9_ ]*$")
	var buscado string
	var err error

	buscarOk := false
	for !buscarOk {
		fmt.Print("Texto a buscar: ")
		reader := bufio.NewReader(os.Stdin)
		buscado, err = reader.ReadString('\n')
		checkError(err)
		buscado = strings.TrimRight(buscado, "\r\n")

		if re.MatchString(buscado) {
			buscarOk = true
		} else {
			fmt.Println("El usuario buscado debe ser una cadena alfa-numérica.")
		}
	}

	// Luego se envia el texto buscado al servidor para obtener los usuarios que contentan dicho texto.
	textoAEnviar := "BuscarUsuarios#&" + buscado
	fmt.Fprintln(conn, textoAEnviar)
	netscan := bufio.NewScanner(conn) // Se crea un scanner para la conexión (datos desde el servidor)
	textoRecibido := ""

	for netscan.Scan() {
		textoRecibido = netscan.Text()
		//fmt.Println(textoRecibido)

		if strings.HasPrefix(textoRecibido, "UsuariosEncontrados#&") {
			break
		}
	}

	usuarios := strings.Split(textoRecibido, "#&")
	// Es -2 porque hay uno inicial con "GetLogueados:" y con los ":" al final del string, se obtiene una última posición vacía
	numUsuarios := len(usuarios) - 2

	// El siguiente paso es elegir el usuario con quien se quiere hablar.
	fmt.Println("\nUsuarios:")

	i := 1
	if numUsuarios == 0 {
		fmt.Println("No se ha encontrado ningún usuario que contenga '" + buscado + "'.")

	} else {
		for ; i <= numUsuarios; i++ {
			// usuarios[i] porque me hay que saltarse la posición 0 ("GetLogueados:")
			textoAMostrar := strings.Replace(usuarios[i], "-", "\n", -1) // -1 significa que no hay límite de coincidencias para reemplazar.
			fmt.Println("\nUsuario ", i, "\n", textoAMostrar)
		}
	}
}

func randString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}
