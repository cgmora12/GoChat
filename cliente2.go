package main

import "net"
import "fmt"
import "bufio"
import "os"

func main() {

	// connect to this socket
	conn, _ := net.Dial("tcp", "127.0.0.1:8081")
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
