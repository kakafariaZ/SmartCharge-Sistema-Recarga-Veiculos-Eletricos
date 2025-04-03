package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	//	"strconv"
	//	"strings"
	"math"
)

// Função que calcula a distância entre o carro e os postos e retorna o melhor ponto de recarga
func calculateStationDistances(carLocation [2]int, stations map[int][2]int) int {
	var bestStation int
	var bestDistance int

	for station, location := range stations {
		distance := int(math.Sqrt(math.Pow(float64(carLocation[0]-location[0]), 2) +
			math.Pow(float64(carLocation[1]-location[1]), 2)))

		if bestDistance == 0 || distance < bestDistance {
			bestDistance = distance
			bestStation = station
		}
	}

	return bestStation
}

// Função que lê o JSON e retorna um mapa com ID da estação e sua localização
func LoadStationsFromJSON(filename string) (map[int][2]int, error) {
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

	// Criar o dicionário com ID -> Coordenadas
	stationsMap := make(map[int][2]int)

	// Processar os dados manualmente sem struct
	if stations, ok := rawData["charge_stations"].([]interface{}); ok {
		for _, station := range stations {
			stationMap := station.(map[string]interface{})
			id := int(stationMap["id"].(float64))
			locationData := stationMap["location"].([]interface{})
			location := [2]int{int(locationData[0].(float64)), int(locationData[1].(float64))}
			stationsMap[id] = location
		}
	}

	return stationsMap, nil
}

func requestStationData(conn net.Conn) map[string]interface{} {

	// Cria a estrutura da mensagem de requisição
	request := map[string]string{"action": "request_station_data"}
	jsonData, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Erro ao criar requisição JSON:", err)
		return nil
	}

	// Envia o pedido para o posto
	fmt.Println("Enviando requisição para o posto...")
	_, err = conn.Write(jsonData)
	if err != nil {
		fmt.Println("Erro ao requisitar dados do posto:", err)
		return nil
	}

	// Aguarda resposta do posto
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Erro ao receber resposta do posto:", err)
		return nil
	}

	// Decodifica a mensagem JSON recebida
	var message_station map[string]interface{}
	err = json.Unmarshal(buf[:n], &message_station)
	if err != nil {
		fmt.Println("Erro ao decodificar JSON:", err)
		return nil
	}

	// Exibe os dados recebidos
	//fmt.Println("Dados do posto recebidos:", message_station)

	return message_station
}

func requestCarData(conn net.Conn) map[string]interface{} {

	// Cria a estrutura da mensagem de requisição
	request := map[string]string{"action": "request_car_data"}
	jsonData, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Erro ao criar requisição JSON:", err)
		return nil
	}

	// Envia o pedido para o posto
	fmt.Println("Enviando requisição para o carro...")
	_, err = conn.Write(jsonData)
	if err != nil {
		fmt.Println("Erro ao requisitar dados do carro:", err)
		return nil
	}

	// Aguarda resposta do posto
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Erro ao receber resposta do carro:", err)
		return nil
	}

	// Decodifica a mensagem JSON recebida
	var message_car map[string]interface{}
	err = json.Unmarshal(buf[:n], &message_car)
	if err != nil {
		fmt.Println("Erro ao decodificar JSON:", err)
		return nil
	}

	// Exibe os dados recebidos
	//fmt.Println("Dados do posto recebidos:", message_station)

	return message_car
}

func messageIdentifier(message map[string]interface{}) string {
	// Identifica se a conexão é de um carro ou um posto de recarga
	clientType, ok := message["type"].(string)
	if !ok {
		fmt.Println("Mensagem inválida, sem tipo definido.")
		return clientType
	}

	// Identifica a origem da mensagem
	switch clientType {
	case "car":
		fmt.Println("Conexão de um carro detectada:", message)
		//handleCar(conn, message)
	case "station":
		fmt.Println("Conexão de um posto de recarga detectada:", message)
		//handleStation(conn, message)
	default:
		fmt.Println("Tipo desconhecido:", clientType)
	}

	return clientType
}

// Função para processar os dados enviados pelo cliente
func handleClient(conn net.Conn) {
	defer conn.Close()

	//buf := make([]byte, 1024) // Buffer para armazenar os dados recebidos

	for {
		// n, err := conn.Read(buf)
		// if err != nil {
		// 	if err == io.EOF {
		// 		fmt.Println("Cliente desconectado:", conn.RemoteAddr())
		// 	} else {
		// 		fmt.Println("Erro ao ler dados:", err)
		// 	}
		// 	break
		// }

		// // Decodifica a mensagem JSON recebida
		// var message_car map[string]interface{}
		// err = json.Unmarshal(buf[:n], &message_car)
		// if err != nil {
		// 	fmt.Println("Erro ao decodificar JSON:", err)
		// 	return
		// }

		message_car := requestCarData(conn) // Requisita os dados do carro
		messageIdentifier(message_car) // Identifica o tipo de mensagem recebida

		message_station := requestStationData(conn) // Requisita os dados do posto de recarga
		messageIdentifier(message_station) // Identifica o tipo de mensagem recebida

		//id := int(message["id"].(float64)) // Convertendo de JSON (float64) para int
		batteryLevel := int(message_car["batteryLevel"].(float64))
		//location := message["location"].([]interface{})
		locationInterface := message_car["location"].([]interface{})

		carLocation := [2]int{
			int(locationInterface[0].(float64)), // Converte cada elemento explicitamente para int
			int(locationInterface[1].(float64)),
		}

		// Armazena as localizações do posto na variável
		chargeStations, err := LoadStationsFromJSON("charge_stations_data.json")
		if err != nil {
			fmt.Println("Erro ao carregar estações de recarga:", err)
			return
		}

		// Verifica se o nível de bateria está crítico
		if batteryLevel <= 20 {
			// Se a bateria estiver crítica, chama a função para calcular a estação mais próxima
			bestStation := calculateStationDistances(carLocation, chargeStations)

			// Exibe o melhor posto de recarga
			fmt.Printf("Coordenadas recebidas: %d, %d\n", carLocation[0], carLocation[1])
			fmt.Printf("Nível de bateria crítico: %d%%\n", batteryLevel)
			fmt.Printf("Melhor Posto de Recarga: %d\n", bestStation)

			// Verifica se o posto selecionado está disponível
		}

		// Exibir os dados recebidos
		//fmt.Printf("Coordenadas recebidas: %d, %d\n", coord_x, coord_y)
		//fmt.Printf("Melhor Posto: %d\n", bestStation)
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
