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
	"sync"
)

// Adicione estas estruturas no início do arquivo
type ClientType int

const (
	CarType ClientType = iota
	StationType
)

type ClientConnection struct {
	Conn net.Conn
	Type ClientType
	ID   int
}

var (
	connections     []ClientConnection
	connectionsLock sync.Mutex
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

// Função para processar os dados enviados pelo cliente
func handleClient(conn net.Conn, chargeStations []ChargingStation) {
	defer conn.Close()

	// Primeiro, identifique se é um carro ou um posto se conectando
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Erro ao ler dados de identificação:", err)
		return
	}

	var clientData map[string]interface{}
	err = json.Unmarshal(buf[:n], &clientData)
	if err != nil {
		fmt.Println("Erro ao decodificar dados de identificação:", err)
		return
	}

	clientType, ok := clientData["type"].(string)
	if !ok {
		fmt.Println("Tipo de cliente não especificado")
		return
	}

	var clientID int
	if id, ok := clientData["id"].(float64); ok {
		clientID = int(id)
	}

	var cType ClientType
	if clientType == "car" {
		cType = CarType
		fmt.Printf("Carro %d conectado\n", clientID)
	} else if clientType == "station" {
		cType = StationType
		fmt.Printf("Posto de recarga %d conectado\n", clientID)
	} else {
		fmt.Println("Tipo de cliente desconhecido:", clientType)
		return
	}

	// Registre a conexão
	connectionsLock.Lock()
	connections = append(connections, ClientConnection{
		Conn: conn,
		Type: cType,
		ID:   clientID,
	})
	connectionsLock.Unlock()

	// Se for um carro, processe seus dados
	if cType == CarType {
		processCarData(conn, clientID)
	} else if cType == StationType {
		processStationData(conn, clientID)
	}
}

func processCarData(conn net.Conn, carID int) {
	buf := make([]byte, 1024)
	chargeStations, err := LoadStationsFromJSON("charge_stations_data.json")
	if err != nil {
		fmt.Println("Erro ao carregar estações de recarga:", err)
		return
	}

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Carro %d desconectado\n", carID)
			} else {
				fmt.Printf("Erro ao ler dados do carro %d: %v\n", carID, err)
			}
			removeConnection(conn)
			break
		}

		coordinates := strings.Split(string(buf[:n]), ",")
		if len(coordinates) < 3 {
			continue
		}

		coord_x, err := strconv.Atoi(strings.TrimSpace(coordinates[0]))
		if err != nil {
			fmt.Println("Erro ao converter coordenada x:", err)
			continue
		}

		coord_y, err := strconv.Atoi(strings.TrimSpace(coordinates[1]))
		if err != nil {
			fmt.Println("Erro ao converter coordenada y:", err)
			continue
		}

		batteryLevel, err := strconv.Atoi(strings.TrimSpace(coordinates[2]))
		if err != nil {
			fmt.Println("Erro ao converter nível de bateria:", err)
			continue
		}

		carLocation := [2]int{coord_x, coord_y}

		if batteryLevel <= 20 {
			bestStation := calculateStationDistances(carLocation, chargeStations)
			fmt.Printf("Carro %d - Bateria crítica: %d%%. Melhor Posto: %d\n",
				carID, batteryLevel, bestStation.ID)

				sendToCar(bestStation.ID, carID, carLocation, chargeStations)
				sendToStation(bestStation.ID, carID, carLocation, batteryLevel)
			

		}
	}
}

func processStationData(conn net.Conn, stationID int) {
	// Mantém a conexão aberta para receber requisições
	buf := make([]byte, 1024)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Posto %d desconectado\n", stationID)
			} else {
				fmt.Printf("Erro ao ler dados do posto %d: %v\n", stationID, err)
			}
			removeConnection(conn)
			break
		}
	}
}

func sendToCar(stationID int, carID int, carLocation [2]int, chargeStations []ChargingStation) {
	connectionsLock.Lock()
	defer connectionsLock.Unlock()

	for _, c := range connections {
		// Verifica se a conexão é do tipo "station" e se o ID corresponde
		// ao ID do posto de recarga
		// Se sim, envia a requisição
		// Se não, continua verificando as outras conexões

		if c.Type == CarType && c.ID == carID {
			// Envia coordenadas do posto para o carro
			bestStation := calculateStationDistances(carLocation, chargeStations)
			request := map[string]interface{}{
				"action":           "request_station_data",
				"best_station_id":  stationID,
				"car_id":           carID,
				"station_location": bestStation.Location, // Envia a localização da estação de recarga
			}

			jsonData, err := json.Marshal(request)
			if err != nil {
				fmt.Println("Erro ao criar requisição JSON:", err)
				return
			}

			_, err = c.Conn.Write(jsonData)
			if err != nil {
				fmt.Println("Erro ao enviar requisição para o carro:", err)
			} else {
				fmt.Printf("Requisição enviada para o carro %d sobre o posto %d\n", carID, stationID)
			}
			return
		}
	}

	fmt.Printf("Carro %d não encontrado entre as conexões ativas\n", carID)
}

func sendToStation(stationID int, carID int, carLocation [2]int, batteryLevel int) {
	connectionsLock.Lock()
	defer connectionsLock.Unlock()

	for _, c := range connections {

		// Verifica se a conexão é do tipo "station" e se o ID corresponde
		// ao ID do posto de recarga
		// Se sim, envia a requisição
		// Se não, continua verificando as outras conexões

		fmt.Printf("Verificando conexão: %d\n", c.ID)

		if c.Type == StationType && c.ID == stationID {
			fmt.Printf("Conexão encontrada para o posto %d\n", stationID)
			request := map[string]interface{}{
				"action":          "request_station_data",
				"best_station_id": stationID,
				"car_id":          carID,
				"car_location":    carLocation,
				"batteryLevel":    batteryLevel,
			}

			jsonData, err := json.Marshal(request)
			if err != nil {
				fmt.Println("Erro ao criar requisição JSON:", err)
				return
			}

			_, err = c.Conn.Write(jsonData)
			if err != nil {
				fmt.Println("Erro ao enviar requisição para o posto:", err)
			} else {
				fmt.Printf("Requisição enviada para o posto %d sobre o carro %d\n",
					stationID, carID)
			}
			return
		}
	}

	fmt.Printf("Posto %d não encontrado entre as conexões ativas\n", stationID)
}


func removeConnection(conn net.Conn) {
	connectionsLock.Lock()
	defer connectionsLock.Unlock()

	for i, c := range connections {
		if c.Conn == conn {
			connections = append(connections[:i], connections[i+1:]...)
			break
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
		go handleClient(conn, []ChargingStation{}) // Inicia uma goroutine para processar os dados recebidos
	}
}
