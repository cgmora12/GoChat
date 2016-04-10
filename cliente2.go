package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

// función para comprobar errores (ahorra escritura)
func chk(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	// connect to this socket
	conn, err := net.Dial("tcp", "localhost:1337")
	chk(err)
	defer conn.Close() // nos aseguramos que cerramos las conexiones aunque el programa falle
	for {
		fmt.Print("0")

		// read in input from stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Pedro: ")
		text, _ := reader.ReadString('\n')

		fmt.Print("1")

		// send to socket
		fmt.Fprintf(conn, text+"\n")

		fmt.Print("2")

		// listen for reply
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("César: " + message)

		fmt.Print("3")
	}
}
