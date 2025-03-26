package main

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

// Função que calcula a distânicia entre o carro e os postos e retorna o melhor ponto de recarga
func calculateStationDistances(carLocation [2]int, stations map[string][2]int) string {
	var bestStation string
	var bestDistance int

	for station, location := range stations {
		distance := abs(carLocation[0]-location[0]) + abs(carLocation[1]-location[1])
		if bestDistance == 0 || distance < bestDistance {
			bestDistance = distance
			bestStation = station
		}
	}

	return bestStation
}

func main() {
	listener, err := net.Listen("tcp", ":8080") // Cria um servidor TCP escutando na porta 8080
	if err != nil {
		panic(err)
	}
	defer listener.Close() // Garante que o socket será fechado quando o servidor for interrompido
	fmt.Println("Servidor TCP rodando na porta 8080...")

	// Loop infinito para aceitar conexões
	for {
		conn, err := listener.Accept() // aguarda por novas conexões
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
			continue
		}
		fmt.Println("Nova conexão de:", conn.RemoteAddr())
		go handleClient(conn) // Inicia uma goroutine para processar os dados recebidos
	}
}

// Função para processar os dados enviados pelo cliente
func handleClient(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024) // Buffer para armazenar os dados recebidos

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Cliente desconectado:", conn.RemoteAddr())
			} else {
				fmt.Println("Erro ao ler dados:", err)
			}
			break
		}

		// Exibir os dados recebidos
		fmt.Println("Coordenadas recebidas:", string(buf[:n]))

		// Separando as coordenadas x e y
		coordinates := strings.Split(string(buf[:n]), ",")

		// Convertendo para números inteiros
		coord_x, err := strconv.Atoi(coordinates[0])
		if err != nil {
			fmt.Println("Erro ao converter coordenada x:", err)
			break
		}
		coord_y, err := strconv.Atoi(coordinates[1])
		if err != nil {
			fmt.Println("Erro ao converter coordenada y:", err)
			break
		}

	}
}