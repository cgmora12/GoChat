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

	ln, err := net.Listen("tcp", "localhost:1337") // Se escucha en espera de conexión
	chk(err)
	defer ln.Close() // Se cierra la conexión al final

	for { // Bucle infinito, se sale con ctrl+c
		conn, err := ln.Accept() // Para cada nueva petición de conexión
		chk(err)
		go func() { // Función lambda (función anónima) en concurrencia

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
					procesarLogin(conn, textoRecibido, usuariosLogueados, port)

				} else { // Si el mensaje recibido no se corresponde con ningún método del servidor
					fmt.Fprintln(conn, "ack: ", textoRecibido) // Se envia el ack al cliente
				}
			}

			conn.Close() // Se cierra la conexión al finalizar el cliente (EOF se envía con ctrl+d o ctrl+z según el sistema)
			fmt.Println("cierre[", port, "]")
			procesarLogout(usuariosLogueados, port)
		}()
	}
}

func procesarRegistro(conn net.Conn, textoRecibido string) {
	fmt.Println(textoRecibido)
	s := strings.Split(textoRecibido, ":")
	nombreUsuario, password := s[1], s[2]
	err := registrarBD(nombreUsuario, password)

	if err != nil {
		respuestaServidor := "Ya existe el usuario " + nombreUsuario + " en la base de datos."
		fmt.Println(respuestaServidor)
		fmt.Fprintln(conn, "Respuesta del servidor: ", respuestaServidor)
	} else {
		respuestaServidor := "Usuario registrado: " + nombreUsuario + " Contraseña: " + password
		fmt.Println(respuestaServidor)
		fmt.Fprintln(conn, "Respuesta del servidor: ", respuestaServidor)
	}
}

func registrarBD(nombreUsuario string, password string) error {
	database := "gochat"
	user := "usuarioGo"
	passwordBD := "usuarioGo"
	con, err := sql.Open("mymysql", database+"/"+user+"/"+passwordBD)
	defer con.Close()

	_, err = con.Exec("insert into usuarios(nombreUsuario, password) values (?, ?)", nombreUsuario, password)

	//fmt.Println("Almacenado ", nombreUsuario, " en la base de datos.")
	return err
}

func procesarLogin(conn net.Conn, textoRecibido string, usuariosLogueados map[string]string, port string) {
	fmt.Println(textoRecibido)
	s := strings.Split(textoRecibido, ":")
	nombreUsuario, password := s[1], s[2]

	if buscarUsuarioLogueado(nombreUsuario, usuariosLogueados) {
		respuestaServidor := "Nombre de usuario y/o contraseña incorrectos o el usuario ya está logueado"
		fmt.Println(respuestaServidor)
		fmt.Fprintln(conn, "Respuesta del servidor:Error: ", respuestaServidor)
	} else {
		loginBD(nombreUsuario, password)

		count := loginBD(nombreUsuario, password)
		if count == 1 {
			respuestaServidor := "Usuario correcto."
			fmt.Println(respuestaServidor)
			fmt.Fprintln(conn, "Respuesta del servidor: ", respuestaServidor)
			usuariosLogueados[nombreUsuario] = port
		} else {
			respuestaServidor := "Nombre de usuario y/o contraseña incorrectos o el usuario ya está logueado"
			fmt.Println(respuestaServidor)
			fmt.Fprintln(conn, "Respuesta del servidor:Error: ", respuestaServidor)
		}
	}
}

func loginBD(nombreUsuario string, password string) int {
	database := "gochat"
	user := "usuarioGo"
	passwordBD := "usuarioGo"

	db := mysql.New("tcp", "", "127.0.0.1:3306", user, passwordBD, database)
	err := db.Connect()
	chk(err)

	rows, res, err := db.Query("SELECT count(*) FROM usuarios WHERE nombreUsuario = '%s' AND password = '%s'", nombreUsuario, password)
	chk(err)

	// Obtener valores por el nombre de la columna devuelta.
	columna := res.Map("count(*)")
	valor := rows[0].Int(columna)
	//fmt.Println(valor)

	return valor
}

func buscarUsuarioLogueado(nombreUsuario string, usuariosLogueados map[string]string) bool {
	_, ok := usuariosLogueados[nombreUsuario]
	return ok
}

func procesarLogout(usuariosLogueados map[string]string, port string) {
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
			fmt.Println("Logout ", key, " correcto")
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
