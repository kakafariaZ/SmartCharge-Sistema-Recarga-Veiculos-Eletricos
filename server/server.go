package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
)

// ChargingStation representa um posto de abastecimento de carro elétrico.
type ChargingStation struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Location   [2]int  `json:"location"`
	Occupation bool    `json:"occupation"`
	Power      float64 `json:"power"`
	Price      float64 `json:"price"`
}

// Estrutura auxiliar para capturar a lista de estações no JSON
type ChargingStationsData struct {
	ChargeStations []ChargingStation `json:"charge_stations"`
}

// Função que calcula a distância entre o carro e os postos e retorna o melhor ponto de recarga
func calculateStationDistances(carLocation [2]int, stations []ChargingStation) ChargingStation {
	var bestStation ChargingStation
	bestDistance := 2000000 // Inicializa com um valor alto para garantir que qualquer distância encontrada seja menor

	// Percorre todas as estações de recarga e calcula a distância
	for _, station := range stations {
		distance := int(math.Sqrt(math.Pow(float64(carLocation[0]-station.Location[0]), 2) +
			math.Pow(float64(carLocation[1]-station.Location[1]), 2)))

		// Se a distância for menor que a melhor distância encontrada, atualiza
		if distance < bestDistance {
			bestDistance = distance
			bestStation = station
			fmt.Print("Melhor posto de recarga encontrado: ", bestStation.Name, "\n")
			fmt.Print("Distância: ", bestDistance, "\n")
		}
	}

	return bestStation
}

// Função que lê o JSON e retorna um mapa com ID da estação e sua localização
func LoadStationsFromJSON(filename string) ([]ChargingStation, error) {
	// Abrir o arquivo JSON
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decodificar o JSON para um mapa genérico
	var rawData map[string]interface{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&rawData); err != nil {
		return nil, err
	}

	// Criar a lista de ChargingStation
	var stations []ChargingStation

	// Processar os dados manualmente sem struct
	if stationsData, ok := rawData["charge_stations"].([]interface{}); ok {
		for _, s := range stationsData {
			stationMap := s.(map[string]interface{})

			id := int(stationMap["id"].(float64))
			name := stationMap["name"].(string)
			locationData := stationMap["location"].([]interface{})
			location := [2]int{int(locationData[0].(float64)), int(locationData[1].(float64))}
			occupation := stationMap["occupation"].(bool)
			// power := stationMap["power"].(float64)
			// price := stationMap["price"].(float64)

			// Criar um objeto ChargingStation e adicioná-lo à lista
			stations = append(stations, ChargingStation{
				ID:         id,
				Name:       name,
				Location:   location,
				Occupation: occupation,
				// Power:      power,
				// Price:      price,
			})
		}
	}

	return stations, nil
}

func requestStation(conn net.Conn, stations []ChargingStation, bestStation int) {

	// Cria a estrutura da mensagem de requisição
	request := map[string]string{"action": "request_station_data"}
	for _, station := range stations {
		if station.ID == bestStation {
			request["station_id"] = strconv.Itoa(station.ID)
			request["station_name"] = station.Name
		}
	}
	jsonData, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Erro ao criar requisição JSON:", err)
		return
	}

	// Envia o pedido para o posto
	fmt.Println("Enviando requisição para o posto...")
	_, err = conn.Write(jsonData)
	if err != nil {
		fmt.Println("Erro ao requisitar:", err)
		return
	}

	// // Aguarda resposta do posto
	// buf := make([]byte, 1024)
	// n, err := conn.Read(buf)
	// if err != nil {
	// 	fmt.Println("Erro ao receber resposta do carro:", err)
	// 	return
	// }

	// // Decodifica a mensagem JSON recebida
	// var message_car map[string]interface{}
	// err = json.Unmarshal(buf[:n], &message_car)
	// if err != nil {
	// 	fmt.Println("Erro ao decodificar JSON:", err)
	// 	return
	// }

	return
}

// Função para processar os dados enviados pelo cliente
func handleClient(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024) // Buffer para armazenar os dados recebidos

	// Armazena as localizações do posto na variável
	chargeStations, err := LoadStationsFromJSON("charge_stations_data.json")
	if err != nil {
		fmt.Println("Erro ao carregar estações de recarga:", err)
		return
	}
	//fmt.Print(chargeStations, "\n")

	for {

		/* ====== LÊ OS DADOS DO BUFFER E OS INTERPRETA COMO COORDENADAS ====== */
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Cliente desconectado:", conn.RemoteAddr())
			} else {
				fmt.Println("Erro ao ler dados:", err)
			}
			break
		}

		// Separando as coordenadas x e y
		coordinates := strings.Split(string(buf[:n]), ",")

		// Convertendo para números inteiros
		// coordenada x
		coord_x, err := strconv.Atoi(strings.TrimSpace(coordinates[0]))
		if err != nil {
			fmt.Println("Erro ao converter coordenada x:", err)
			break
		}

		// coordenada x
		coord_y, err := strconv.Atoi(strings.TrimSpace(coordinates[1]))
		if err != nil {
			fmt.Println("Erro ao converter coordenada y:", err)
			break
		}

		batteryLevel, err := strconv.Atoi(strings.TrimSpace(coordinates[2]))
		if err != nil {
			fmt.Println("Erro ao converter nível de bateria:", err)
			break
		}

		// Armazena as coordenadas do carro na variável
		carLocation := [2]int{coord_x, coord_y}

		// Verifica se o nível de bateria está crítico
		if batteryLevel <= 20 {

			// Se a bateria estiver crítica, chama a função para calcular a estação mais próxima
			bestStation := calculateStationDistances(carLocation, chargeStations)

			// Exibe o melhor posto de recarga
			fmt.Printf("Coordenadas recebidas: %d, %d\n", carLocation[0], carLocation[1])
			fmt.Printf("Nível de bateria crítico: %d%%\n", batteryLevel)
			fmt.Printf("Melhor Posto de Recarga: %d\n", bestStation.ID)

			// Verifica se o posto selecionado está disponível
			requestStation(conn, chargeStations, bestStation.ID) // Envia a requisição para o posto
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080") // Cria um servidor TCP escutando na porta 8080
	if err != nil {
		panic(err)
	}
	defer listener.Close() // Garante que o socket será fechado quando o servidor for interrompido
	fmt.Println("Servidor TCP rodando na porta 8080...")

	// Loop infinito para aceitar conexões
	for {
		conn, err := listener.Accept() // aguarda por novas conexões
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
			continue
		}
		fmt.Println("Nova conexão de:", conn.RemoteAddr())
		go handleClient(conn) // Inicia uma goroutine para processar os dados recebidos
	}
}
