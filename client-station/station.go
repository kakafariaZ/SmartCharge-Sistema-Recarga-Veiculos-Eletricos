package main

import (
	//"bufio"
	"fmt"
	"net"

	//"time"
	"encoding/json"
)

type Station struct {
	Type string `json:"type"`
	ID   int    `json:"id"`
}

func main() {

	// Conecta ao servidor na porta 8080
	conn, err := net.Dial("tcp", "server:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Envia identificação como posto de recarga
	station := Station{
		Type: "station",
		ID:   2, // Ou o ID correto do posto
	}
	
	jsonData, err := json.Marshal(station)
	if err != nil {
		fmt.Println("Erro ao converter identificação para JSON:", err)
		return
	}
	
	_, err = conn.Write(jsonData)
	if err != nil {
		fmt.Println("Erro ao enviar identificação:", err)
		return
	}

	fmt.Println("Ponto de Recarga conectado ao servidor...")

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Erro ao ler requisição do servidor:", err)
			return
		}

		var request map[string]interface{}
		err = json.Unmarshal(buf[:n], &request)
		if err != nil {
			fmt.Println("Erro ao decodificar requisição:", err)
			continue
		}

		if request["action"] == "request_station_data" {
			carID := int(request["car_id"].(float64))
			location := request["car_location"].([]interface{})
			x := int(location[0].(float64))
			y := int(location[1].(float64))
			
			fmt.Printf("Requisição recebida para atender carro %d na localização [%d, %d]\n", 
				carID, x, y)
			
			// Aqui você pode adicionar a lógica para responder ao servidor
			// confirmando que o posto está disponível, etc.
		}
	}
}
