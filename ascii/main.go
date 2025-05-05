package main

import (
	"fmt"
	"strconv"
)

func main() {
	// Non-printable characters
	for i := byte(0); i < 32; i++ {
		fmt.Printf("%d = %v \n", i, strconv.QuoteToASCII(string(i)))
	}
	fmt.Printf("%d = %v \n", 127, strconv.QuoteToASCII(string(127)))

	// Printable characters
	for i := byte(32); i < 127; i++ {
		fmt.Printf("%d = %v \n", i, string(i))
	}

	// To print all ASCII symbols from 0 to 127 there is a strconv.QuoteToASCII function that quotes special characters
	// passing symbols as-is to formattable print function interprets/executes symbols affecting output
	// e.g. 0x10 is '\n' symbol so it will put Line Feed and (0x1B = 27) is ESC can cut further bytes
	// following loops produces quite different output for almost the same starting 'i'
	fmt.Println("Table (27-127):")
	for i := byte(27); i < 127; i++ {
		fmt.Printf("%v ", string(i))
	}
	fmt.Println()
	fmt.Println("Table (28-127):")
	for i := byte(28); i < 127; i++ {
		fmt.Printf("%v ", string(i))
	}
	fmt.Println()

	// Printing Unicode symbol consisting of 4 bytes
	fmt.Println()
	fmt.Println([]byte("🔥"))
	fmt.Println('🔥')
	fmt.Println(string([]byte{240, 159, 148, 165}))
	fmt.Printf("%s", []byte{240, 159, 148, 165, 10})
}
