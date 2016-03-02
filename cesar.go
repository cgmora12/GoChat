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
	"math"
	"os"
	"strconv"
	"unicode"
)

func main() {

	var fin *os.File  // fichero de entrada
	var fout *os.File // fichero de salida
	var err error     // receptor de error
	var despl string

	// alfabeto con el que vamos a trabajar
	alfabeto := map[rune]int{'A': 0, 'B': 1, 'C': 2, 'D': 3, 'E': 4, 'F': 5, 'G': 6, 'H': 7, 'I': 8, 'J': 9,
		'K': 10, 'L': 11, 'M': 12, 'N': 13, 'Ñ': 14, 'O': 15, 'P': 16, 'Q': 17, 'R': 18, 'S': 19,
		'T': 20, 'U': 21, 'V': 22, 'W': 23, 'X': 24, 'Y': 25, 'Z': 26}

	if len(os.Args) == 1 { // no hay parámetros, usamos entrada (teclado) y salida estándar (pantalla)

		fin = os.Stdin
		fout = os.Stdout

	} else if len(os.Args) == 2 { // tenemos el desplazamiento parámetros, usamos entrada (teclado) y salida estándar (pantalla)
		despl = os.Args[1]

		fin = os.Stdin
		fout = os.Stdout
	} else if len(os.Args) == 4 { // tenemos los nombres de los ficheros de entrada y salida en los parámetros
		despl = os.Args[1]

		fin, err = os.Open(os.Args[2]) // abrimos el primer fichero (entrada)
		if err != nil {
			panic(err)
		}
		defer fin.Close()

		fout, err = os.Create(os.Args[3]) // abrimos el segundo fichero (salida)
		if err != nil {
			panic(err)
		}
		defer fout.Close()
	} else { // error de parámetros
		fmt.Println("Número de parámetros incorrecto: se espera los ficheros de origen y destino (1 y 2, opcionales)")
		os.Exit(1)
	}

	for { // bucle infinito
		var c rune // carácter a leer

		_, err = fmt.Fscanf(fin, "%c", &c) // lectura de la entrada

		if err != nil { // si hay error (fin de fichero)
			break //salimos del bucle
		}

		C := unicode.ToUpper(c) // pasamos a mayúsculas

		var indice int
		indice, ok := alfabeto[C] // comprobamos que está en el alfabeto
		if ok {                   // en tal caso imprimimos en la salida
			var desplazamiento int
			desplazamiento, err := strconv.Atoi(despl)
			if err != nil {
				desplazamiento = 0
			}

			for r := range alfabeto {
				indiceNuevo, ok := alfabeto[r]
				if ok {
					if indiceNuevo == int(math.Mod(float64(desplazamiento+indice), 27)) {
						fmt.Fprintf(fout, "%c", r)
					}
				}

			}

		}

	}

}
