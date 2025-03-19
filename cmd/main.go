package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
)

// Estrutura representando um carro elétrico
type Car struct {
	ID       string  `json:"id"`
	Battery  int     `json:"battery"`
	Position float64 `json:"position"`
}

// Estrutura representando um posto de carregamento
type Station struct {
	ID       string  `json:"id"`
	Slots    int     `json:"slots"`
	Position float64 `json:"position"`
}

var (
	clients  = make(map[net.Conn]bool) // Mapa para armazenar conexões de clientes
	broadcast = make(chan string)       // Canal de comunicação
	mutex     = &sync.Mutex{}           // Mutex para evitar condições de corrida
)

// Função para lidar com conexões de clientes
func handleClient(conn net.Conn) {
	defer conn.Close()
	clients[conn] = true

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			fmt.Println("Cliente desconectado")
			return
		}
		message := string(buffer[:n])
		broadcast <- message
	}
}

// Função para distribuir mensagens entre os clientes conectados
func handleMessages() {
	for {
		msg := <-broadcast
		mutex.Lock()
		for client := range clients {
			_, _ = client.Write([]byte(msg))
		}
		mutex.Unlock()
	}
}

// Função para salvar dados em um arquivo JSON
func saveData(filename string, data interface{}) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Erro ao criar arquivo:", err)
		return
	}
	defer file.Close()
	json.NewEncoder(file).Encode(data)
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
		return
	}
	defer listener.Close()

	go handleMessages()

	cars := []Car{{ID: "C1", Battery: 100, Position: 10.5}}
	stations := []Station{{ID: "S1", Slots: 2, Position: 15.0}}

	saveData("cars.json", cars)
	saveData("stations.json", stations)

	fmt.Println("Servidor rodando na porta 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
			continue
		}
		go handleClient(conn)
	}
}
