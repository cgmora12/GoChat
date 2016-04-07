/*
Servidor de una arquitectura cliente-servidor

uso:
go run servidorPedro.go

Driver MySQL:
http://stackoverflow.com/questions/11353679/whats-the-recommended-way-to-connect-to-mysql-from-go
https://github.com/ziutek/mymysql
*/

package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/ziutek/mymysql/godrv"
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
					fmt.Println(textoRecibido)
					s := strings.Split(textoRecibido, ":")
					nombreUsuario, password := s[1], s[2]
					err := almacenarBD(nombreUsuario, password)

					if err != nil {
						respuestaServidor := "Ya existe el usuario " + nombreUsuario + " en la base de datos."
						fmt.Println(respuestaServidor)
						fmt.Fprintln(conn, "Respuesta del servidor: ", respuestaServidor)
					} else {
						respuestaServidor := "Usuario registrado: " + nombreUsuario + " Contraseña: " + password
						fmt.Println(respuestaServidor)
						fmt.Fprintln(conn, "Respuesta del servidor: ", respuestaServidor)
					}
				} else { // Si el mensaje recibido no se corresponde con ningún método del servidor
					fmt.Fprintln(conn, "ack: ", textoRecibido) // Se envia el ack al cliente
				}
			}

			conn.Close() // Se cierra la conexión al finalizar el cliente (EOF se envía con ctrl+d o ctrl+z según el sistema)
			fmt.Println("cierre[", port, "]")
		}()
	}
}

func almacenarBD(nombreUsuario string, password string) error {
	database := "gochat"
	user := "usuarioGo"
	passwordBD := "usuarioGo"
	con, err := sql.Open("mymysql", database+"/"+user+"/"+passwordBD)
	defer con.Close()

	_, err = con.Exec("insert into usuarios(nombreUsuario, password) values (?, ?)", nombreUsuario, password)

	//fmt.Println("Almacenado ", nombreUsuario, " en la base de datos.")
	return err
}
