package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "server:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Escolha um número de 1 a 5:")
	reader := bufio.NewReader(os.Stdin)

	for {
		text, _ := reader.ReadString('\n')
		if text == "exit\n" {
			break
		}
		conn.Write([]byte(text)) // Envia a mensagem para o servidor

		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Resposta do servidor: ", message)
	}
}
