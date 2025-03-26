package main

import (
	"fmt"
	"io"
	"net"
	"encoding/json"
	"os"
	"strconv"
	"strings"
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

// Função para processar os dados enviados pelo cliente
func handleClient(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024) // Buffer para armazenar os dados recebidos

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Cliente desconectado:", conn.RemoteAddr())
			} else {
				fmt.Println("Erro ao ler dados:", err)
			}
			break
		}


		/* ====== LÊ OS DADOS DO BUFFER E OS INTERPRETA COMO COORDENADAS ====== */
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
			fmt.Printf("Coordenadas recebidas: %d, %d\n", coord_x, coord_y)
			fmt.Printf("Nível de bateria crítico: %d%%\n", batteryLevel)
			fmt.Printf("Melhor Posto de Recarga: %d\n", bestStation)
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
