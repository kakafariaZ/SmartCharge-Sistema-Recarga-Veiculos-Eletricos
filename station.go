package main

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"os"
	"math/rand"
	"time"
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

// Função para obter o ID do ponto de recarga, garantindo que não se repita
func getStationID() int {
	// Lê o arquivo JSON com os IDs já usados
	file, err := os.Open("used_ids.json")
	if err != nil {
		// Se o arquivo não existir, cria um novo arquivo com uma lista vazia de IDs
		if os.IsNotExist(err) {
			file, err = os.Create("used_ids.json")
			if err != nil {
				fmt.Println("Erro ao criar o arquivo:", err)
				return -1
			}
			// Inicializa com uma lista vazia de IDs usados
			json.NewEncoder(file).Encode(map[string][]int{"used_ids": []int{}})
			file.Close()
			file, err = os.Open("used_ids.json")
			if err != nil {
				fmt.Println("Erro ao ler o arquivo:", err)
				return -1
			}
		} else {
			fmt.Println("Erro ao abrir o arquivo:", err)
			return -1
		}
	}
	defer file.Close()

	// Lê os IDs usados a partir do arquivo
	var data map[string][]int
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("Erro ao decodificar o arquivo:", err)
		return -1
	}

	// Gera um ID aleatório entre 1 e 5
	rand.Seed(time.Now().UnixNano())
	var newID int
	usedIDs := data["used_ids"]

	// Garante que o novo ID não tenha sido utilizado antes
	for {
		newID = rand.Intn(5) + 1
		if !contains(usedIDs, newID) {
			break
		}
	}

	// Adiciona o novo ID à lista de usados
	usedIDs = append(usedIDs, newID)
	// Exibi a lista de IDs usados
	fmt.Printf("IDs usados: %v\n", usedIDs)

	// Atualiza o arquivo com a lista de IDs usados
	file, err = os.Create("used_ids.json")
	if err != nil {
		fmt.Println("Erro ao criar o arquivo:", err)
		return -1
	}
	defer file.Close()
	data["used_ids"] = usedIDs
	err = json.NewEncoder(file).Encode(data)
	if err != nil {
		fmt.Println("Erro ao codificar o arquivo:", err)
		return -1
	}
	fmt.Printf("Novo ID de estação gerado: %d\n", newID)

	return newID
}

// Função auxiliar para verificar se um ID já foi usado
func contains(ids []int, id int) bool {
	for _, v := range ids {
		if v == id {
			return true
		}
	}
	return false
}

func main() {

	// Conecta ao servidor na porta 8080
	conn, err := net.Dial("tcp", "server:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	queue := &QueueManager{}

	// Pega o ID do ponto de recarga a partir do ID do container	
	stationID := getStationID()

	station := Station{
		Type: "station",
		ID:   stationID, // Pode vir de config/env
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

	fmt.Printf("✅ Ponto de Recarga %d conectado ao servidor...", station.ID)

	// Loop para receber requisições do servidor
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
				StationID:  bestStation,
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
