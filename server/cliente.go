package main

import (
	"bufio"
	"fmt"
	"net"
	"string"
	"os"
)

func main() {
	conn, err := net.Dial("top", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Digite uma mensagem:")
	reader := bufio.NewReader(os.Stdin)

	for {
		text, _ := reader.ReaderString ('\n')
		if text == "exit\n" {
			break
		}
		conn.Write ([]byte(text))

		message, _ := bufio.NewReader(conn).ReaderString('\n')
		fmt.Print("Resposta do servidor: ", message)
	}
}