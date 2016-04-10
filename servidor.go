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

// función para comprobar errores (ahorra escritura)
func chk(e error) {
	if e != nil {
		panic(e)
	}
}

var messageCliente = ""
var i = 0

func handleConnection(conn net.Conn) {

	i = i + 1
	numCliente := i

	for {
		// try to read data from the connection
		data := make([]byte, 512)
		n, err := conn.Read(data)
		chk(err)

		fmt.Println("Cliente: " + strconv.Itoa(numCliente))
		if data != nil {

			s := string(data[:n])

			fmt.Println("Datos de la conexion: " + s)

			// read in input from stdin
			//reader := bufio.NewReader(os.Stdin)
			fmt.Println("Mensaje cliente " + messageCliente)
			if !strings.Contains(messageCliente, strconv.Itoa(numCliente)) && messageCliente != "" {
				fmt.Println("Cliente num" + strconv.Itoa(numCliente))

				// send data
				fmt.Println("Mensaje enviado: " + s + " al otro cliente " + strconv.Itoa(numCliente))
				fmt.Fprintf(conn, s+strconv.Itoa(numCliente)+"\n")

				s = messageCliente

				// send previous
				fmt.Println("Mensaje enviado: " + s + " al cliente " + strconv.Itoa(numCliente))
				fmt.Fprintf(conn, s+"\n")

			} else {
				fmt.Println("Sin mensajes " + strconv.Itoa(numCliente))
				messageCliente = s + strconv.Itoa(numCliente)
				// send to socket
				fmt.Fprintf(conn, "Sin mensajes\n")
			}

			fmt.Println("Fin bucle cliente " + strconv.Itoa(numCliente))
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
	ln, err := net.Listen("tcp", "localhost:1337")
	chk(err)
	defer ln.Close() // nos aseguramos que cerramos las conexiones aunque el programa falle

	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}
		fmt.Println("Conexion con cliente")
		go handleConnection(conn)
	}

}
