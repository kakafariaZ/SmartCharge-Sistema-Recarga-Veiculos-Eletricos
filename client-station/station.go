package main

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

type Station struct {
	Type string `json:"type"`
	ID   int    `json:"id"`
}

type StationStatus struct {
	Type       string `json:"type"`
	StationID  int    `json:"station_id"`
	CarsInLine int    `json:"cars_in_line"`
	Available  bool   `json:"available"`
}

// QueueManager controla a fila de carros para o posto
type QueueManager struct {
	queue []int
	mutex sync.Mutex
}

func (q *QueueManager) AddCar(carID int) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.queue = append(q.queue, carID)
	fmt.Printf("[FILA] Carro %d adicionado à fila.\n", carID)
	fmt.Printf("[FILA] Fila atual: %v\n", q.queue)
}

func (q *QueueManager) RemoveCar() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if len(q.queue) == 0 {
		return -1
	}
	carID := q.queue[0]
	q.queue = q.queue[1:]
	fmt.Printf("[FILA] Carro %d removido da fila.\n", carID)
	fmt.Printf("[FILA] Fila atual: %v\n", q.queue)
	return carID
}

func (q *QueueManager) Count() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return len(q.queue)
}

func (q *QueueManager) IsAvailable() bool {
	return q.Count() == 0
}

func main() {

	// Conecta ao servidor na porta 8080
	conn, err := net.Dial("tcp", "server:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	queue := &QueueManager{}

	// Envia identificação como posto de recarga
	station := Station{
		Type: "station",
		ID:   2, // Pode vir de config/env
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

	fmt.Println("✅ Ponto de Recarga conectado ao servidor...")

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
			bestStation := int(request["best_station_id"].(float64))
			location := request["car_location"].([]interface{})
			x := int(location[0].(float64))
			y := int(location[1].(float64))

			fmt.Printf("\n📩 Requisição recebida para atender carro %d na localização [%d, %d] - POSTO: %d\n",
				carID, x, y, bestStation)

			// Adiciona o carro à fila
			queue.AddCar(carID)

			// Envia resposta ao servidor com status do posto
			status := StationStatus{
				Type:       "station_status",
				StationID:  station.ID,
				CarsInLine: queue.Count(),
				Available:  queue.IsAvailable(),
			}

			statusData, _ := json.Marshal(status)
			_, err := conn.Write(statusData)
			if err != nil {
				fmt.Println("Erro ao enviar status do posto:", err)
			} else {
				fmt.Printf("📤 Status enviado ao servidor: carros na fila = %d | disponível = %v\n\n",
					status.CarsInLine, status.Available)
			}
		}
	}
}
