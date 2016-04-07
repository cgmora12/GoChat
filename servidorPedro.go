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

// función para comprobar errores (ahorra escritura)
func chk(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	fmt.Println("-- Servidor GoChat --")
	server()
}

// gestiona el modo servidor
func server() {
	ln, err := net.Listen("tcp", "localhost:1337") // escucha en espera de conexión
	chk(err)
	defer ln.Close() // nos aseguramos que cerramos las conexiones aunque el programa falle

	for { // búcle infinito, se sale con ctrl+c
		conn, err := ln.Accept() // para cada nueva petición de conexión
		chk(err)
		go func() { // lanzamos un cierre (lambda, función anónima) en concurrencia

			_, port, err := net.SplitHostPort(conn.RemoteAddr().String()) // obtenemos el puerto remoto para identificar al cliente (decorativo)
			chk(err)

			fmt.Println("conexión: ", conn.LocalAddr(), " <--> ", conn.RemoteAddr())

			scanner := bufio.NewScanner(conn) // el scanner nos permite trabajar con la entrada línea a línea (por defecto)

			for scanner.Scan() { // escaneamos la conexión
				textoRecibido := scanner.Text()

				fmt.Println("cliente[", port, "]: ", textoRecibido) // mostramos el mensaje del cliente
				fmt.Fprintln(conn, "ack: ", textoRecibido)          // enviamos ack al cliente

				// Se comprueba
				if strings.Contains(textoRecibido, "Registro:") {
					fmt.Println(textoRecibido)
					s := strings.Split(textoRecibido, ":")
					nombreUsuario, password := s[1], s[2]
					err := almacenarBD(nombreUsuario, password)

					if err != nil {
						fmt.Println("Ya existe el usuario ", nombreUsuario, " en la base de datos.")
					} else {
						fmt.Println("Usuario registrado: ", nombreUsuario, password)
					}
				}
			}

			conn.Close() // cerramos al finalizar el cliente (EOF se envía con ctrl+d o ctrl+z según el sistema)
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
