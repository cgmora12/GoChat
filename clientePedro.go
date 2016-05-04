/*
Cliente de una estructura cliente-servidor

Uso:
go run clientePedro.go
*/

package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/howeyc/gopass"
	"math/big"
	"net"
	"os"
	"strconv"
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
		//fmt.Println("3.- Entrar modo cliente")
		fmt.Println("3.- Salir")
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
			/*
				case "3":
					fmt.Println("\n- Entrar modo cliente -")
					client()
			*/
		case "3":
			salir = true

		default:
			fmt.Println("\nOpción '" + opcionElegida + "' desconocida. Introduzca una opción válida (1, 2 o 3)")
		}
	}
}

/*
func client() {
	conn, err := net.Dial("tcp", "localhost:1337") // Se llama al servidor
	chk(err)
	defer conn.Close() // Se cierra la conexión al final

	fmt.Println("conectado a " + conn.RemoteAddr())

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
*/

func registro() {
	fmt.Println("\n- Registro -")

	// Se obtienen los datos de registro
	fmt.Print("Nombre de usuario: ")
	reader := bufio.NewReader(os.Stdin)
	nombreUsuario, err := reader.ReadString('\n')
	chk(err)
	nombreUsuario = strings.TrimRight(nombreUsuario, "\r\n")

	fmt.Print("Contraseña: ")
	//reader = bufio.NewReader(os.Stdin)
	//password, err := reader.ReadString('\n')
	//chk(err)
	//password = strings.TrimRight(password, "\r\n")
	pass, err := gopass.GetPasswdMasked()
	chk(err)

	fmt.Print("Nombre completo: ")
	reader = bufio.NewReader(os.Stdin)
	nombreCompleto, err := reader.ReadString('\n')
	chk(err)
	nombreCompleto = strings.TrimRight(nombreCompleto, "\r\n")

	fmt.Print("País: ")
	reader = bufio.NewReader(os.Stdin)
	pais, err := reader.ReadString('\n')
	chk(err)
	pais = strings.TrimRight(pais, "\r\n")

	fmt.Print("Província: ")
	reader = bufio.NewReader(os.Stdin)
	provincia, err := reader.ReadString('\n')
	chk(err)
	provincia = strings.TrimRight(provincia, "\r\n")

	fmt.Print("Localidad: ")
	reader = bufio.NewReader(os.Stdin)
	localidad, err := reader.ReadString('\n')
	chk(err)
	localidad = strings.TrimRight(localidad, "\r\n")

	fmt.Print("Email: ")
	reader = bufio.NewReader(os.Stdin)
	email, err := reader.ReadString('\n')
	chk(err)
	email = strings.TrimRight(email, "\r\n")

	// Se envian los datos al servidor
	conn, err := net.Dial("tcp", "localhost:1337") // Se llama al servidor
	chk(err)
	defer conn.Close() // Se cierra la conexión al finalizar

	fmt.Println("conectado a ", conn.RemoteAddr())

	var datos string = "Registro#&" + nombreUsuario + "#&" + string(pass[:]) + "#&" + nombreCompleto + "#&" + pais + "#&" + provincia + "#&" + localidad + "#&" + email

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
	//reader = bufio.NewReader(os.Stdin)
	//password, err := reader.ReadString('\n')
	//chk(err)
	//password = strings.TrimRight(password, "\r\n")
	pass, err := gopass.GetPasswdMasked()
	chk(err)

	// Se envian los datos al servidor
	conn, err := net.Dial("tcp", "localhost:1337") // Se llama al servidor
	chk(err)
	defer conn.Close() // Se cierra la conexión al finalizar

	fmt.Println("conectado a ", conn.RemoteAddr())

	var datos string = "Login#&" + nombreUsuario + "#&" + string(pass[:])

	fmt.Fprintln(conn, datos) // Se envian los datos al servidor

	netscan := bufio.NewScanner(conn) // Se crea un scanner para la conexión (datos desde el servidor)
	netscan.Scan()                    // Se escanea la conexión

	textoRecibido := netscan.Text()
	textoAMostrar := strings.Replace(textoRecibido, "#&", ":", -1) // -1 significa que no hay límite de coincidencias para reemplazar.
	fmt.Println(textoAMostrar)                                     // Se muestra el mensaje desde el servidor

	s := strings.Split(textoRecibido, "#&")
	substringRespuesta := s[1]

	if substringRespuesta != "Error" {
		elegirOpcionChat(nombreUsuario, conn)
	}
}

func elegirOpcionChat(nombreUsuario string, conn net.Conn) {
	salir := false
	for !salir {
		fmt.Println("\n\n-- GoChat --")
		fmt.Println("-- Usuario: " + nombreUsuario + " --")

		fmt.Println("Eliga una opción: ")
		fmt.Println("1.- Sala pública")
		fmt.Println("2.- Salas privadas")
		fmt.Println("3.- Ver perfiles de usuarios")
		fmt.Println("4.- Logout")
		fmt.Print("Opción elegida (introduzca el número): ")

		reader := bufio.NewReader(os.Stdin)
		opcionElegida, err := reader.ReadString('\n')
		chk(err)
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
				textoRecibido := netscan.Text()
				fmt.Println("\n" + textoRecibido)
				fmt.Print("Continúe su mensaje: ")
			}
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
			if strings.Contains(textoAEnviar, "#&") {
				fmt.Println("Error: Secuencia '#&' inválida")
				fmt.Print("Escriba su mensaje: ")
			} else if textoAEnviar == "Salir" {
				fmt.Fprintln(conn, "SalirChat#&")
				quit <- true
				done2 <- true
				return
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

		fmt.Println("Eliga el usuario con quien quiera hablar:\n")

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
		chk(err)
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

	netscan := bufio.NewScanner(conn) // Se crea un scanner para la conexión (datos desde el servidor)

	// Todo: Intercambio de claves
	cli_keys, err := rsa.GenerateKey(rand.Reader, 1024) // generamos un par de claves (privada, pública) para el cliente
	chk(err)
	cli_keys.Precompute() // aceleramos su uso con un precálculo

	cli_keys_json, err := json.Marshal(cli_keys.PublicKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(cli_keys_json))

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

	fmt.Println("Mensaje recibido:" + claveRecibida)
	cli_pub := strings.Split(string(claveRecibida), "Clave:")[1]
	fmt.Println("Clave recibida:" + cli_pub)

	var cli_pub_key rsa.PublicKey
	cli_pub_trozo1 := strings.Split(cli_pub, ":")[1]
	fmt.Println("Trozo1:" + cli_pub_trozo1)
	cli_pub_trozo2 := strings.Split(cli_pub_trozo1, ",")[0]
	fmt.Println("Trozo2:" + cli_pub_trozo2)
	fmt.Println("E:" + strings.Split(strings.Split(cli_pub, ":")[2], "}")[0])
	cli_pub_key.E, err = strconv.Atoi(strings.Split(strings.Split(cli_pub, ":")[2], "}")[0])
	bigint := new(big.Int)
	bigint.SetString(cli_pub_trozo2, 10)
	cli_pub_key.N = bigint

	//cli_token := make([]byte, 48) // 384 bits (256 bits de clave + 128 bits para el IV)
	buff := make([]byte, 256)   // contendrá el token cifrado con clave pública (puede ocupar más que el texto en claro)
	cli_token := randString(48) // generación del token aleatorio para el cliente

	fmt.Println("token creado ", cli_token)

	// ciframos el token del cliente con la clave pública del otro cliente
	enctoken, err := rsa.EncryptPKCS1v15(rand.Reader, &cli_pub_key, []byte(cli_token))
	chk(err)

	fmt.Println("token cifrado ", string(enctoken)) // Se envia la entrada al servidor
	enctokenString := strings.Replace(string(enctoken), "\n", "-----n", -1)

	fmt.Println("token cifrado sin endlines ", enctokenString) // Se envia la entrada al servidor
	// Falla al encriptar y poner saltos de linea
	tokenCliente := "Token#&" + esteUsuario + "#&" + usuarioElegido + "#&" + enctokenString
	fmt.Fprintln(conn, tokenCliente) // Se envia la entrada al servidor

	var tokenRecibido string
	for netscan.Scan() {

		tokenRecibido = netscan.Text()
		if strings.HasPrefix(tokenRecibido, "Token:") {
			fmt.Println("Token recibido:" + tokenRecibido)
			break
		}

	}

	fmt.Println("Token recibido sin prefijo: " + strings.Split(tokenRecibido, "Token:")[1])
	fmt.Println("Token recibido sin prefijo ni endlines: " + strings.Replace(strings.Split(tokenRecibido, "Token:")[1], "-----n", "\n", -1))

	buff = []byte(strings.Replace(strings.Split(tokenRecibido, "Token:")[1], "-----n", "\n", -1))

	// desciframos el token del otro cliente con nuestra clave privada
	session_key, err := rsa.DecryptPKCS1v15(rand.Reader, cli_keys, buff)
	chk(err)

	// realizamos el XOR entre ambos tokens (cliente y servidor acaban con la misma clave de sesión)
	for i := 0; i < len(cli_token); i++ {
		session_key[i] ^= cli_token[i]
	}

	fmt.Println("Clave de sesion: " + string(session_key))

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
				textoRecibido := netscan.Text()
				fmt.Println("\n" + textoRecibido)
				fmt.Print("Continúe su mensaje: ")
			}
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
			if strings.Contains(textoAEnviar, "#&") {
				fmt.Println("Error: Secuencia '#&' inválida")
				fmt.Print("Escriba su mensaje: ")
			} else if textoAEnviar == "Salir" {
				fmt.Fprintln(conn, "SalirChat#&")
				quit <- true
				done2 <- true
				return
			} else { // Si el mensaje recibido no se corresponde con ningún método del servidor
				textoPreparado := "SalaPrivada#&" + esteUsuario + "#&" + usuarioElegido + "#&" + textoAEnviar
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

func verPerfiles(conn net.Conn, nombreUsuario string) {
	salir := false
	for !salir {
		fmt.Println("\n\n-- Perfiles de usuarios --")
		fmt.Println("-- Usuario: " + nombreUsuario + " --")
		fmt.Println("Eliga una opción: ")
		fmt.Println("1.- Ver todos los usuarios")
		fmt.Println("2.- Buscar usuarios")
		fmt.Println("3.- Volver atrás")
		fmt.Print("Opción elegida (introduzca el número): ")

		reader := bufio.NewReader(os.Stdin)
		opcionElegida, err := reader.ReadString('\n')
		chk(err)
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
	fmt.Println("Usuarios:")

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
	fmt.Print("Texto a buscar: ")
	reader := bufio.NewReader(os.Stdin)
	buscado, err := reader.ReadString('\n')
	chk(err)
	buscado = strings.TrimRight(buscado, "\r\n")

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
	fmt.Println("Usuarios:")

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
