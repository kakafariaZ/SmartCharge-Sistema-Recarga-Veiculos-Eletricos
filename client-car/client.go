package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

// Estrutura que será usada para comunicação com o servidor (mensagens JSON)
type Car struct {
	Type         string  `json:"type"`         // Tipo do cliente (carro)
	ID           int     `json:"id"`           // Identificador único do carro
	BatteryLevel int     `json:"batteryLevel"` // Nível de bateria do carro
	Location     [2]int  `json:"location"`     // Coordenadas (x, y)
	Credit       float64 `json:"credit"`
}

// Estrutura usada para o controle interno do cliente (carro)
type CarState struct {
	ID           int
	Location     [2]int
	BatteryLevel int
	Status       string // "normal", "crÍtico"
	Credit       float64
}

func main() {
	rand.Seed(time.Now().UnixNano())

	carID := getCarID() // Gera um ID dinâmico para o carro

	///* ****

	// Cria o objeto usado para comunicação com o servidor
	car := Car{
		ID:           carID,
		BatteryLevel: rand.Intn(51) + 50,                     // Bateria entre 50 e 100%
		Location:     [2]int{rand.Intn(250), rand.Intn(250)}, // Coordenadas aleatórias entre 0 e 250
	}

	// Cria a estrutura de estado interno do carro
	carState := &CarState{
		ID:           car.ID,
		Location:     car.Location,
		BatteryLevel: car.BatteryLevel,
		Status:       "normal",
		Credit:       500.00,
	}

	// Conecta ao servidor TCP (definido no docker-compose como "server:8080")
	conn, err := net.Dial("tcp", "server:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Printf("\n🚗 Carro %d conectado ao servidor!\n", car.ID)

	// Envia a identificação do carro para o servidor
	ident := map[string]interface{}{
		"type": "car",
		"id":   car.ID,
	}
	jsonData, _ := json.Marshal(ident)
	conn.Write(jsonData)

	// Canal para comunicação entre goroutines em caso de bateria crítica
	criticalChan := make(chan CarState)

	// Inicia as goroutines (concorrência)
	go handleRequests(carState, conn)         // Escuta requisições do servidor
	go carMovement(carState, criticalChan)    // Simula movimentação e monitora bateria
	go handleCriticalData(conn, criticalChan) // Envia dados críticos ao servidor

	select {} // Mantém o programa rodando
}

// ==========================
// Goroutines
// ==========================

// Simula o movimento do carro e monitora a bateria
func carMovement(car *CarState, criticalChan chan CarState) {
	for {
		time.Sleep(time.Second) // Espera 1 segundo entre os movimentos

		if car.Status == "normal" {
			// Move o carro aleatoriamente
			car.Location[0] += rand.Intn(11)
			car.Location[1] += rand.Intn(11)

			// Atualiza o nível da bateria
			car.BatteryLevel = batteryLevel(car.BatteryLevel)

			// Mostra visualmente o nível da bateria
			displayBattery(*car)

			// Verifica se entrou em estado crítico
			if car.BatteryLevel <= 20 {
				car.Status = "critico"
				fmt.Println("⚠️  ALERTA! 🚨 Bateria crítica! ")
				// Envia dados para o canal crítico
				criticalChan <- *car
				moveToStation(car, car.Location) // move o carro até o posto
			}

			// Exibe a posição atual do carro
			fmt.Printf("📍 Coordenadas: %v\n", car.Location)
		}
	}
}

// Lida com o envio de dados críticos ao servidor
func handleCriticalData(conn net.Conn, criticalChan chan CarState) {
	for {
		carCritical := <-criticalChan // Espera por dados no canal
		// Prepara os dados em formato "x, y, bateria"
		data := fmt.Sprintf("%d, %d, %d\n",
			carCritical.Location[0],
			carCritical.Location[1],
			carCritical.BatteryLevel,
		)

		// Envia para o servidor
		_, err := conn.Write([]byte(data))
		if err != nil {
			fmt.Println("❌ Erro ao enviar dados críticos:", err)
			return
		}
		fmt.Println("\n✅ Dados críticos enviados ao servidor:", data)
	}
}

// Função para mover o carro até o posto
func moveToStation(car *CarState, stationLocation [2]int) {
	// Move o carro até a coordenada do posto
	if car.Location != stationLocation {
		fmt.Printf("Carro na posição %v 📍, movendo-se para o posto localizado em %v 📍\n", car.Location, stationLocation)
		for car.Location != stationLocation {
			//fmt.Printf("Carro na posição %v 📍, movendo-se para o posto localizado em %v 📍\n", car.Location, stationLocation)
			if car.Location[0] < stationLocation[0] {
				car.Location[0]++
			} else if car.Location[0] > stationLocation[0] {
				car.Location[0]--
			}

			if car.Location[1] < stationLocation[1] {
				car.Location[1]++
			} else if car.Location[1] > stationLocation[1] {
				car.Location[1]--
			}

			if car.Location == stationLocation {
				fmt.Printf("Carro na posição %v 📍, chegou no posto localizado em %v 📍\n", car.Location, stationLocation)
				fmt.Printf("Carro chegou ao posto!!! 🚘📌")
				batteryUpLevel(car)
				fmt.Println("Carro carregado com sucesso!🔋✅")
				fmt.Println("Carro desconectado da baia de carregamento!🔋🔌")
				break
			}
		}
	}
}

// Escuta requisições do servidor e responde com os dados do carro
// Função que escuta requisições do servidor e responde com os dados do carro
func handleRequests(car *CarState, conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("❌ Erro ao ler requisição do servidor:", err)
			return
		}

		// Decodifica a requisição recebida (espera JSON)
		var request map[string]interface{}
		err = json.Unmarshal(buf[:n], &request)
		if err != nil {
			fmt.Println("❌ Erro ao decodificar requisição:", err)
			return
		}

		// Se a ação for de requisição de dados, envia os dados do carro
		if request["action"] == "request_car_data" {
			sendCarData(car, conn)
			fmt.Println("✅ Dados do carro enviados ao servidor.")
		}

		// Se o servidor envia as coordenadas do posto para o carro, move o carro
		if request["action"] == "request_station_data" {
			// Aqui estamos assumindo que station_location é uma fatia de 2 elementos [interface{}]
			if locationData, ok := request["station_location"].([]interface{}); ok && len(locationData) == 2 {
				// Converte os valores para inteiros
				x, okX := locationData[0].(float64) // Tentando converter para float64
				y, okY := locationData[1].(float64) // Tentando converter para float64

				if okX && okY {
					stationLocation := [2]int{int(x), int(y)} // Converte para [2]int
					moveToStation(car, stationLocation)
				} else {
					fmt.Println("❌ Erro ao converter coordenadas da estação para int")
				}
			} else {
				fmt.Println("❌ Erro: 'station_location' não é uma lista de 2 elementos")
			}
		}
	}
}

// ==========================
// Funções Auxiliares
// ==========================

// Reduz o nível da bateria gradualmente
func batteryLevel(batteryLevel int) int {
	batteryLevel -= 5
	if batteryLevel <= 20 {
		batteryLevel = 20 // não deixa ir abaixo de 20
	}
	return batteryLevel
}

// Função para carregar a bateria do carro
func batteryUpLevel(car *CarState) int {
	for i := 0; i < 10; i++ { // repete 10 vezes
		time.Sleep(1 * time.Second)
		fmt.Println("🔋⚡ Carro carregando...")

		// Exibe o estado atual da bateria
		displayBattery(*car)

		// Aumenta o nível da bateria
		car.BatteryLevel += 10

		// Se atingir ou ultrapassar 100, ajusta e faz o débito
		if car.BatteryLevel >= 100 {
			// Exibe a bateria após o aumento
			displayBattery(*car)

			// Garante que o nível da bateria não ultrapasse 100
			car.BatteryLevel = 100

			// Altera o status do carro
			car.Status = "normal"

			// Realiza o débito de 79.99 da conta do usuário
			car.Credit -= 79.99

			time.Sleep(1 * time.Second)
			fmt.Println("Foi debitado 79.99 reais da sua conta!💸💳")

			time.Sleep(1 * time.Second)

			// Exibe o saldo restante na conta do usuário
			fmt.Printf("Você tem %.2f reais na sua conta!🏦💵\n", car.Credit)

			// Sai do loop após a bateria atingir 100
			break
		}
	}

	// Retorna o nível final da bateria
	return car.BatteryLevel
}

// Envia os dados do carro para o servidor em formato JSON
func sendCarData(car *CarState, conn net.Conn) {
	// Cria um mapa com os dados do carro para enviar como JSON
	data := map[string]interface{}{
		"id":            car.ID,
		"location":      car.Location,
		"battery_level": car.BatteryLevel,
	}

	// Serializa os dados para JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Erro ao criar JSON para o carro:", err)
		return
	}

	// Envia os dados para o servidor
	_, err = conn.Write(jsonData)
	if err != nil {
		fmt.Println("Erro ao enviar dados do carro:", err)
	}
}

// Exibe o nível da bateria de forma visual
func displayBattery(car CarState) {
	totalBars := 20
	numHashMarks := (car.BatteryLevel * totalBars) / 100

	//time.Sleep(1 * time.Second)
	//fmt.Print("\033[H\033[2J") //Print limpa o terminal com esse comando

	fmt.Printf("\n     🚗 ID: %d\n", car.ID)
	fmt.Println("┌──────────────────────┐")
	fmt.Printf("│   Bateria: %3d%%      │\n", car.BatteryLevel)
	fmt.Println("├──────────────────────┤")
	fmt.Printf("││%s%s││\n", strings.Repeat("█", numHashMarks), strings.Repeat(" ", totalBars-numHashMarks))
	fmt.Println("└──────────────────────┘")
}

// Gera um ID único baseado no timestamp
func getCarID() int {
	return int(time.Now().UnixNano() % 10000)
}
