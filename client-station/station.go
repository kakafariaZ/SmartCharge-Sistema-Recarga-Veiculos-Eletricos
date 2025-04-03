package main

import (
	"fmt"
	"net"
)

func main() {

	// Conecta ao servidor na porta 8080
	conn, err := net.Dial("tcp", "server:8080")
	if err != nil {
		panic(err)
	}

	defer conn.Close() // Fecha a conexão

	// Teste de conexão com o servidor
	fmt.Println("Ponto de Recarga conectando ao servidor...")

}
