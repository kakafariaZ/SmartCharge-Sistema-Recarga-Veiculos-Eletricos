package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor:", err)
		return
	}
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Erro ao ler mensagem:", err)
			return
		}

		msg := string(buffer[:n])
		fmt.Println("Posto recebeu:", msg)

		response := "Baia disponível para carregamento"
		_, _ = conn.Write([]byte(response))
	}
}
