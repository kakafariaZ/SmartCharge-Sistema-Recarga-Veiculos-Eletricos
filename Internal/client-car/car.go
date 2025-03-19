package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor:", err)
		return
	}
	defer conn.Close()

	for {
		batteryLevel := 20 // Simulação de bateria em nível crítico
		if batteryLevel < 30 {
			msg := "Carro solicitando posto de carregamento"
			_, _ = conn.Write([]byte(msg))
			fmt.Println("Mensagem enviada:", msg)
		}
		time.Sleep(5 * time.Second)
	}
}
