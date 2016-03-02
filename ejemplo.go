/*
Ejemplo 1

Este programa copia de la entrada a la salida carácter a carácter,
restringiéndose a un alfabeto limitado y pasando a mayúsculas.
Permite leer de la entrada y salida estándar o usar ficheros

ejemplos de uso:

go run ejemplo1.go

go run ejemplo1.go fichentrada.txt fichsalida.txt


-lectura y escritura
-entrada y salida estándar
-ficheros
-parámetros en línea de comandos (os.Args)
*/

package main

import (
	"fmt"
	"net"
)

func handleConnection(conn net.Conn) {
	// try to read data from the connection
	data := make([]byte, 512)
	n, err := conn.Read(data)

	if err != nil {
		panic(err)
	}
	s := string(data[:n])

	// print the request data
	fmt.Println(s)

	// import "bytes"
	var str = []string{"Hi there!"}
	var buf bytes.Buffer
	for _, s := range str {
		buf.WriteString(s)
	}
	// write the data to the connection
	_, err = conn.Write(buf.Bytes())

	if err != nil {
		panic(err)
	}

	// close the connection
	conn.Close()
}

func main() {
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		// handle error
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}
		go handleConnection(conn)
	}

}
