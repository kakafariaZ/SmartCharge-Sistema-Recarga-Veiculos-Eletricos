//atualizado 2
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


func main() {

	// Conecta ao servidor na porta 8080
	conn, err := net.Dial("tcp", "server:8080")
	if err != nil {
		panic(err)
	}

	defer conn.Close() // Fecha a conexão

	fmt.Println("Servidor esperando requisições na porta 8080...")

	data := loadStationData("charge_stations_data.json")
	fmt.Println("Lendo arquivo\n")
	
	sendStationData(data, conn)
	fmt.Println("Enviando arquivo\n")

	for {
		// Processa a requisição em uma goroutine para suportar múltiplos clientes
		
		//go handleAvailabilityRequest(conn)
		//go handleLocationRequest(conn)
	}
}