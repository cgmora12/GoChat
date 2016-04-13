/*
Servidor de una arquitectura cliente-servidor

uso:
go run servidorPedro.go

Driver MySQL:
http://stackoverflow.com/questions/11353679/whats-the-recommended-way-to-connect-to-mysql-from-go
https://github.com/ziutek/mymysql

Para saber el tipo de una variable:
importar: "reflect"
Y para imprimir el tipo por pantalla: fmt.Println(reflect.TypeOf(conn))
*/

package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/ziutek/mymysql/godrv"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native" // Native engine
	// _ "github.com/ziutek/mymysql/thrsafe" // Thread safe engine
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net"
	"strings"
)

// Función para comprobar errores (ahorra escritura)
func chk(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	fmt.Println("-- Servidor GoChat --")
	server()
}

// Gestiona las conexiones al servidor
func server() {
	// Se crea el map para almacenar los usuarios logueados.
	// Se utiliza un map porque se almacenará el port y el nombreUsuario
	// map[nombreUsuario]port
	usuariosLogueados := make(map[string]string)
	// map[port]conn
	connUsuariosLogueados := make(map[string]net.Conn)

	ln, err := net.Listen("tcp", "localhost:1337") // Se escucha en espera de conexión
	chk(err)
	defer ln.Close() // Se cierra la conexión al final

	for { // Bucle infinito, se sale con ctrl+c
		conn, err := ln.Accept() // Para cada nueva petición de conexión
		chk(err)
		go func() { // Función lambda (función anónima) en concurrencia. Con "go func()" se crea una goroutine para cada conexión.

			_, port, err := net.SplitHostPort(conn.RemoteAddr().String()) // Se obtiene el puerto remoto para identificar al cliente
			chk(err)

			fmt.Println("conexión: ", conn.LocalAddr(), " <--> ", conn.RemoteAddr())

			scanner := bufio.NewScanner(conn) // Con el scanner se trabaja con la entrada línea a línea (por defecto)

			for scanner.Scan() { // Se escanea la conexión
				textoRecibido := scanner.Text()

				fmt.Println("cliente[", port, "]: ", textoRecibido) // Se muestra el mensaje del cliente

				// Se comprueba si el mensaje recibido corresponde con algún método del servidor
				if strings.HasPrefix(textoRecibido, "Registro:") {
					procesarRegistro(conn, textoRecibido)

				} else if strings.HasPrefix(textoRecibido, "Login:") {
					procesarLogin(conn, textoRecibido, port, usuariosLogueados, connUsuariosLogueados)

				} else if strings.HasPrefix(textoRecibido, "SalirChat:") {
					// Se envia algo para que el scanner del cliente pueda reaccionar
					// (si no se envia nada el cliente se quedaría escuchando indefinidamente)
					fmt.Fprintln(conn, "")

				} else if strings.HasPrefix(textoRecibido, "GetLogueados:") {
					// Se envia algo para que el scanner del cliente pueda reaccionar
					// (si no se envia nada el cliente se quedaría escuchando indefinidamente)
					var textoAEnviar string = "GetLogueados:"
					for key, value := range usuariosLogueados {
						if value != port {
							textoAEnviar += (key + ":")
						}
					}
					fmt.Fprintln(conn, textoAEnviar)

				} else if strings.HasPrefix(textoRecibido, "SalaPrivada:") {
					enviarAlDestino(textoRecibido, usuariosLogueados, connUsuariosLogueados)

				} else if strings.HasPrefix(textoRecibido, "VerPerfiles:") {
					devolverPerfiles(conn)

				} else { // Si el mensaje recibido no se corresponde con ningún método del servidor
					//fmt.Fprintln(conn, "ack del servidor: ", textoRecibido) // Se envia el ack al cliente

					enviarATodos(textoRecibido, port, usuariosLogueados, connUsuariosLogueados)
				}
			}

			conn.Close() // Se cierra la conexión al finalizar el cliente (EOF se envía con ctrl+d o ctrl+z según el sistema)
			procesarLogout(port, usuariosLogueados, connUsuariosLogueados)
			fmt.Println("cierre[", port, "]")
		}()
	}
}

func procesarRegistro(conn net.Conn, textoRecibido string) {
	//fmt.Println(textoRecibido)
	s := strings.Split(textoRecibido, ":")
	nombreUsuario, password, nombreCompleto, pais, provincia, localidad, email := s[1], s[2], s[3], s[4], s[5], s[6], s[7]
	hashedPassword, errorHash := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	chk(errorHash)

	// Comparing the password with the hash
	comprobacion := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	fmt.Println(comprobacion) // nil means it is a match
	err := registrarBD(nombreUsuario, string(hashedPassword), nombreCompleto, pais, provincia, localidad, email)

	if err != nil {
		respuestaServidor := "Ya existe el usuario " + nombreUsuario + " en la base de datos."
		fmt.Println(respuestaServidor)
		fmt.Fprintln(conn, "Respuesta del servidor: ", respuestaServidor)
	} else {
		respuestaServidor := "Usuario registrado: " + nombreUsuario + " Contraseña: " + string(hashedPassword)
		fmt.Println(respuestaServidor)
		fmt.Fprintln(conn, "Respuesta del servidor: ", respuestaServidor)
	}
}

func registrarBD(nombreUsuario string, password string, nombreCompleto string, pais string, provincia string, localidad string, email string) error {
	database := "gochat"
	user := "usuarioGo"
	passwordBD := "usuarioGo"
	con, err := sql.Open("mymysql", database+"/"+user+"/"+passwordBD)
	defer con.Close()

	_, err = con.Exec("insert into usuarios(nombreUsuario, password, nombreCompleto, pais, provincia, localidad, email) values (?, ?, ?, ?, ?, ?, ?)", nombreUsuario, password, nombreCompleto, pais, provincia, localidad, email)

	//fmt.Println("Almacenado ", nombreUsuario, " en la base de datos.")
	return err
}

func procesarLogin(conn net.Conn, textoRecibido string, port string, usuariosLogueados map[string]string, connUsuariosLogueados map[string]net.Conn) {
	//fmt.Println(textoRecibido)
	s := strings.Split(textoRecibido, ":")
	nombreUsuario, password := s[1], s[2]

	if buscarUsuarioLogueado(nombreUsuario, usuariosLogueados) {
		respuestaServidor := "Nombre de usuario y/o contraseña incorrectos o el usuario ya está logueado"
		fmt.Println(respuestaServidor)
		fmt.Fprintln(conn, "Respuesta del servidor:Error: ", respuestaServidor)
	} else {
		comprobacion := loginBD(nombreUsuario, password)

		if comprobacion == nil {
			respuestaServidor := "Usuario correcto."
			fmt.Println(respuestaServidor)
			fmt.Fprintln(conn, "Respuesta del servidor: ", respuestaServidor)
			usuariosLogueados[nombreUsuario] = port
			connUsuariosLogueados[port] = conn
		} else {
			respuestaServidor := "Nombre de usuario y/o contraseña incorrectos o el usuario ya está logueado"
			fmt.Println(respuestaServidor)
			fmt.Fprintln(conn, "Respuesta del servidor:Error: ", respuestaServidor)
		}
	}
}

func loginBD(nombreUsuario string, password string) error {
	database := "gochat"
	user := "usuarioGo"
	passwordBD := "usuarioGo"
	comprobacion := errors.New("")

	db := mysql.New("tcp", "", "127.0.0.1:3306", user, passwordBD, database)
	err := db.Connect()
	chk(err)

	rows, res, err := db.Query("SELECT password FROM usuarios WHERE nombreUsuario = '%s'", nombreUsuario)
	chk(err)

	// Obtener valores por el nombre de la columna devuelta.

	passwordBd := res.Map("password")

	if rows != nil {
		valor := rows[0].Str(passwordBd)
		comprobacion = bcrypt.CompareHashAndPassword([]byte(valor), []byte(password))
	}

	return comprobacion
}

func buscarUsuarioLogueado(nombreUsuario string, usuariosLogueados map[string]string) bool {
	_, ok := usuariosLogueados[nombreUsuario]
	return ok
}

func procesarLogout(port string, usuariosLogueados map[string]string, connUsuariosLogueados map[string]net.Conn) {
	/*
		// Map antes del logout
		fmt.Println("-- Map antes del logout --")
		for key, value := range usuariosLogueados {
			fmt.Println("Key:", key, "Value:", value)
		}
	*/
	for key, value := range usuariosLogueados {
		if value == port {
			delete(usuariosLogueados, key)
			delete(connUsuariosLogueados, port)
			fmt.Println("Logout usuario", key, "correcto")
			break
		}
	}
	/*
		// Map después del logout
		fmt.Println("-- Map después del logout --")
		for key, value := range usuariosLogueados {
			fmt.Println("Key:", key, "Value:", value)
		}
	*/
}

func enviarATodos(textoRecibido string, portOrigen string, usuariosLogueados map[string]string, connUsuariosLogueados map[string]net.Conn) {
	for key, value := range usuariosLogueados {
		if value != portOrigen { // Para no enviarlo al origen
			usuarioOrigen := buscarUsuarioOrigen(portOrigen, usuariosLogueados)
			textoAEnviar := "Sala pública: " + usuarioOrigen + ": " + textoRecibido
			fmt.Fprintln(connUsuariosLogueados[value], textoAEnviar) // Se envia la entrada al cliente
			fmt.Println("Enviado '", textoRecibido, "' al usuario", key, "mediante el puerto", value)
		}
	}
}

func buscarUsuarioOrigen(portOrigen string, usuariosLogueados map[string]string) string {
	var usuarioOrigen string

	for key, value := range usuariosLogueados {
		if value == portOrigen {
			usuarioOrigen = key
		}
	}

	return usuarioOrigen
}

func enviarAlDestino(textoRecibido string, usuariosLogueados map[string]string, connUsuariosLogueados map[string]net.Conn) {
	s := strings.Split(textoRecibido, ":")
	usuarioOrigen, usuarioDestino, mensajeAEnviar := s[1], s[2], s[3]

	portOrigen := usuariosLogueados[usuarioOrigen]
	connOrigen := connUsuariosLogueados[portOrigen]

	portDestino := usuariosLogueados[usuarioDestino]
	connDestino := connUsuariosLogueados[portDestino]

	if portDestino == "" {
		envioOrigen := "El usuario " + usuarioDestino + " ya no está logueado."
		fmt.Fprintln(connOrigen, envioOrigen) // Se envia el mensaje de error al origen
		fmt.Println("Enviado '", envioOrigen, "' al usuario", usuarioOrigen, "mediante el puerto", portOrigen)

	} else {
		envioDestino := "Sala privada: " + usuarioOrigen + ": " + mensajeAEnviar
		fmt.Fprintln(connDestino, envioDestino) // Se envia el mensaje al destino
		fmt.Println("Enviado '", envioDestino, "' al usuario", usuarioDestino, "mediante el puerto", portDestino)
	}
}

func devolverPerfiles(conn net.Conn) {
	database := "gochat"
	user := "usuarioGo"
	passwordBD := "usuarioGo"

	db := mysql.New("tcp", "", "127.0.0.1:3306", user, passwordBD, database)
	err := db.Connect()
	chk(err)

	rows, res, err := db.Query("SELECT * FROM usuarios")
	chk(err)

	// Obtener valores por el nombre de la columna devuelta.
	nombreUsuario := res.Map("nombreUsuario")
	nombreCompleto := res.Map("nombreCompleto")
	pais := res.Map("pais")
	provincia := res.Map("provincia")
	localidad := res.Map("localidad")
	email := res.Map("email")
	tamano := len(rows)
	textoAEnviar := "VerPerfiles:"

	for i := 0; i < tamano; i++ {
		valorNombreUsuario := rows[i].Str(nombreUsuario)
		valorNombreCompleto := rows[i].Str(nombreCompleto)
		valorPais := rows[i].Str(pais)
		valorProvincia := rows[i].Str(provincia)
		valorLocalidad := rows[i].Str(localidad)
		valorEmail := rows[i].Str(email)

		textoAEnviar += "Nombre usuario = " + valorNombreUsuario + "-Nombre completo = " + valorNombreCompleto + "-Pais = " + valorPais + "-Provincia = " + valorProvincia + "-Localidad = " + valorLocalidad + "-Email = " + valorEmail + ":"
	}
	//fmt.Println(textoAEnviar)
	fmt.Fprintln(conn, textoAEnviar)
}
