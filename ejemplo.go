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
	//"bufio"
	//"bytes"
	"fmt"
	//"github.com/R358/brace/latch"
	//"log"
	"net"
	//"os"
	"strconv"
	"strings"
)

var messageCliente = ""
var i = 0

func handleConnection(conn net.Conn) {

	i = i + 1
	numCliente := i

	for {
		// try to read data from the connection
		data := make([]byte, 512)
		n, err := conn.Read(data)

		if err != nil {
			panic(err)
		}

		if data != nil {

			s := string(data[:n])

			fmt.Print("-1")

			fmt.Print("0")

			// print the request data

			fmt.Print("1")

			// read in input from stdin
			//reader := bufio.NewReader(os.Stdin)

			fmt.Print("2")

			if !strings.Contains(messageCliente, strconv.Itoa(numCliente)) && messageCliente != "" {
				fmt.Print("if")
				s = messageCliente

				fmt.Print("Recibido: ")
				fmt.Println(s)

				// send to socket
				fmt.Fprintf(conn, s+strconv.Itoa(numCliente)+"\n")
			} else {
				fmt.Print("else" + strconv.Itoa(numCliente))
				messageCliente = s + strconv.Itoa(numCliente)
				// send to socket
				fmt.Fprintf(conn, "Sin mensajes\n")
			}

			fmt.Print("4")
		}
	}

	// write the data to the connection
	/*_, err = conn.Write(buf.Bytes())

	if err != nil {
		panic(err)
	}*/

	//c := latch.NewCountdownLatch(1)
	//c.Await()

	// close the connection
	//conn.Close()
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
