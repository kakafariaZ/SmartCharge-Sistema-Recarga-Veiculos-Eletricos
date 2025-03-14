package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Conexão encerrada:", conn.RemoteAddr())
			return
		}
		fmt.Printf("Recebido de %s: %s", conn.RemoteAddr(), message)
		response := strings.ToUpper(message)
		conn.Write([]byte(response))
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	fmt.Println("Servidor TCP rodando na porta 8080...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
			continue
		}
		fmt.Println("Nova conexão de:", conn.RemoteAddr())
		go handleConnection(conn)
	}
}
