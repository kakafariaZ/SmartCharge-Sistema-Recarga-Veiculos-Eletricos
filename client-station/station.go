package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

// ChargingStation representa um posto de abastecimento de carro elétrico.
type ChargingStation struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Location   [2]int    `json:"location"`
	Occupation bool      `json:"occupation"`
	Capacity   int       `json:"capacity"`
	Power      int       `json:"power"`
	Price      float64   `json:"price"`
	mu         sync.Mutex
}

// station é a instância do posto carregada do JSON.
var station ChargingStation

// Configuração do WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// loadStationData carrega os dados do posto a partir de um arquivo JSON.
func loadStationData(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&station)
	return err
}

// handleWebSocket gerencia a comunicação via WebSocket.
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		// Lê a mensagem recebida do cliente (carro)
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("WebSocket read error:", err)
			break
		}

		var request map[string]string
		err = json.Unmarshal(message, &request)
		if err != nil {
			fmt.Println("JSON parse error:", err)
			continue
		}

		// Processa a requisição recebida
		station.mu.Lock()
		var response map[string]interface{}

		switch request["action"] {
		case "check":
			response = map[string]interface{}{"available": !station.Occupation}
		case "connect":
			if station.Occupation {
				response = map[string]interface{}{"error": "Station occupied"}
			} else {
				station.Occupation = true
				response = map[string]interface{}{"message": "Car connected successfully"}
			}
		case "disconnect":
			if !station.Occupation {
				response = map[string]interface{}{"error": "No car connected"}
			} else {
				station.Occupation = false
				response = map[string]interface{}{"message": "Car disconnected successfully"}
			}
		default:
			response = map[string]interface{}{"error": "Invalid action"}
		}
		station.mu.Unlock()

		// Envia a resposta para o cliente
		respJSON, _ := json.Marshal(response)
		conn.WriteMessage(websocket.TextMessage, respJSON)
	}
}

func main() {
	// Carrega os dados do posto a partir do JSON
	err := loadStationData("station.json")
	if err != nil {
		fmt.Println("Error loading station data:", err)
		return
	}

	http.HandleFunc("/ws", handleWebSocket)

	fmt.Println("Charging station WebSocket server running on port 8080")
	http.ListenAndServe(":8080", nil)
}