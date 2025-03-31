//atualizado
package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
)

// ChargingStation representa um posto de abastecimento de carro elétrico.
type ChargingStation struct {
	Type       string  `json:"type"`
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Location   [2]int  `json:"location"`
	Occupation bool    `json:"occupation"`
	Power      int     `json:"power"`
	Price      float64 `json:"price"`
	mu         sync.Mutex
}

// station é a instância do posto carregada do JSON.
//var station ChargingStation

// loadStationData carrega os dados do posto a partir de um arquivo JSON.
func loadStationData(filename string) ChargingStation {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Erro ao ler JSON")
	}

	// Decodificar o JSON diretamente para a estrutura ChargingStation
	var station ChargingStation
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&station); err != nil {
		return ChargingStation{}
	}

	// Retornar a instância de ChargingStation carregada
	return station
}

func sendStationData(station ChargingStation, conn net.Conn) {
	defer conn.Close()

	// Informações
	message := ChargingStation{
		Type:       station.Type,
		ID:         station.ID,
		Occupation: station.Occupation,
		Location:   station.Location,
	}

	// Convertendo para JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Erro ao converter para JSON:", err)
		return
	}

	// Enviando JSON par//a o servidor
	_, err = conn.Write(jsonData)
	if err != nil {
		fmt.Println("Erro ao enviar dados:", err)
	}

}

// func handleLocationRequest(conn net.Conn) {
// 	defer conn.Close()

// 	// Lendo a requisição do servidor solicitante
// 	reader := bufio.NewReader(conn)
// 	request, err := reader.ReadString('\n')
// 	if err != nil {
// 		fmt.Println("Erro ao ler a requisição:", err)
// 		return
// 	}

// 	fmt.Println("Requisição recebida:", request)

// 	// Consultando a localização do posto
// 	station.mu.Lock()
// 	isLocation := station.Location // Retorna a localização com uma tupla
// 	station.mu.Unlock()

// 	//Criando resposta
// 	response := fmt.Sprintf("%t\n", isLocation) // Responde com a localização com uma tupla
// 	conn.Write([]byte(response))                // Enviando resposta

// }

// func handleAvailabilityRequest(conn net.Conn) {
// 	defer conn.Close()

// 	// Lendo a requisição do servidor solicitante
// 	reader := bufio.NewReader(conn)
// 	request, err := reader.ReadString('\n')
// 	if err != nil {
// 		fmt.Println("Erro ao ler a requisição:", err)
// 		return
// 	}

// 	fmt.Println("Requisição recebida:", request)

// 	// Consultando a disponibilidade do posto
// 	station.mu.Lock()
// 	isAvailable := !station.Occupation // Se ocupado, retorna false; se livre, true
// 	station.mu.Unlock()

// 	// Criando a resposta
// 	response := fmt.Sprintf("%t\n", isAvailable) // Responde apenas "true" ou "false"
// 	conn.Write([]byte(response))                 // Enviando resposta
// }

func main() {
	// Carrega os dados do posto a partir do JSON
	// err := loadStationData("charge_stations_data.json")
	// if err != nil {
	// 	fmt.Println("Erro ao carregar os dados do posto:", err)
	// 	return
	// }

	//fmt.Printf("Dados do posto carregados: %+v\n", station)

	// Criando um servidor TCP para responder requisições
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Servidor esperando requisições na porta 8080...")

	for {
		// Aceita conexões de outros servidores perguntando sobre disponibilidade
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
			continue
		}

		// Processa a requisição em uma goroutine para suportar múltiplos clientes
		data := loadStationData("charge_stations_data.json")
		go sendStationData(data, conn)
		//go handleAvailabilityRequest(conn)
		//go handleLocationRequest(conn)
	}
}
