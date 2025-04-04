package main

import (
	//"bufio"
	"fmt"
	"net"

	//"time"
	"encoding/json"
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

	for {
		// Leitura de dados do servidor
		fmt.Println("Aguardando requisição do servidor...")
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Erro ao ler requisição do servidor:", err)
			return
		}

		var request map[string]string
		err = json.Unmarshal(buf[:n], &request)
		if err != nil {
			fmt.Println("Erro ao decodificar requisição:", err)
			return
		}

		if request["action"] == "request_station_data" {
			fmt.Printf("Requisição para o posto de recarga %d recebida.\n", request["station_id"])
		}
	}

}
