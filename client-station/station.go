// atualizado 2
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

func handleRequests(conn net.Conn) {
	defer conn.Close()
	for {
		fmt.Println("Aguardando requisição do servidor...")
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Erro ao ler requisição do servidor:", err)
			return
		}

		fmt.Print("Requisição recebida do servidor: ", string(buf[:n]))

		var request map[string]string
		err = json.Unmarshal(buf[:n], &request)
		if err != nil {
			fmt.Println("Erro ao decodificar requisição:", err)
			return
		}

		fmt.Println("Requisição recebida:", request)

		if request["action"] == "request_station_data" {
			// Carrega dados do posto
			station := loadStationData("charge_stations_data.json")
			fmt.Println("Lendo arquivo JSON do posto...")

			sendStationData(station, conn)
			fmt.Println("Dados do posto enviados ao servidor.")
		}
	}
}

func main() {

	// Conecta ao servidor na porta 8080
	conn, err := net.Dial("tcp", "server:8080")
	if err != nil {
		panic(err)
	}

	//defer conn.Close() // Fecha a conexão

	// Teste de conexão com o servidor
	fmt.Println("Ponto de Recarga conectando ao servidor...")

	// data := loadStationData("charge_stations_data.json")
	// fmt.Println("Lendo arquivo\n")

	// sendStationData(data, conn)
	// fmt.Println("Enviando arquivo\n")

	handleRequests(conn) // Aguarda e responde requisições do servidor
	// for {
	// 	// Processa a requisição em uma goroutine para suportar múltiplos clientes

	// 	//go handleAvailabilityRequest(conn)
	// 	//go handleLocationRequest(conn)
	// }
}
