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

	// Lista com as cores disponiveis para o cliente escolher e requisitar ao servidor
	data := [5]string{"rosa", "vermelho", "azul", "verde", "amarelo"}

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Conexão encerrada:", conn.RemoteAddr())
			return
		}
		fmt.Printf("Recebido de %s: %s", conn.RemoteAddr(), message)

		// Recebe o numero da cor escolhida pelo cliente e envia a cor correspondente
		// Caso o cliente escolha uma cor que não existe, o servidor envia uma mensagem de erro
		switch message {
		case "1\n":
			conn.Write([]byte(data[0] + "\n"))
		case "2\n":
			conn.Write([]byte(data[1] + "\n"))
		case "3\n":
			conn.Write([]byte(data[2] + "\n"))
		case "4\n":
			conn.Write([]byte(data[3] + "\n"))
		case "5\n":
			conn.Write([]byte(data[4] + "\n"))
		default:
			conn.Write([]byte("Cor não encontrada\n"))
		}

	}
}

func chargeStationGenerator(stations_amount, x_limit, y_limit int) []string {

}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	fmt.Println("Servidor TCP rodando na porta 8080...")

	// Loop infinito para aceitar conexões
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
