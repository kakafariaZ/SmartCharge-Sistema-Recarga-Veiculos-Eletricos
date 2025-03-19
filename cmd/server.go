package main

import (
	"fmt"
	"net"
	"sync"
)

var clients = make(map[net.Conn]bool)
var broadcast = make(chan string)
var mutex = &sync.Mutex{}

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

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
		return
	}
	defer listener.Close()

	go handleMessages()

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
