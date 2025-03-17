package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "servidor:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Digite uma mensagem:")
	reader := bufio.NewReader(os.Stdin)

	for {
		text, _ := reader.ReadString('\n')
		if text == "exit\n" {
			break
		}
		conn.Write([]byte(text))

		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Resposta do servidor: ", message)
	}
}
