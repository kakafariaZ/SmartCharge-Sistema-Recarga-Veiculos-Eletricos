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
	Type         string `json:"type"`           // Tipo do cliente (carro)
	ID           int    `json:"id"`			    // Identificador único do carro
	BatteryLevel int    `json:"batteryLevel"`	// Nível de bateria do carro
	Location     [2]int `json:"location"`		// Coordenadas (x, y)
}


// Estrutura usada para o controle interno do cliente (carro)
type CarState struct {
	ID           int
	Location     [2]int
	BatteryLevel int
	Status       string // "normal", "crÍtico"
}

func main() {
	rand.Seed(time.Now().UnixNano())

	carID := getCarID() // Gera um ID dinâmico para o carro

	// Cria o objeto usado para comunicação com o servidor
	car := Car{
		ID:           carID,
		BatteryLevel: rand.Intn(51) + 50,                      // Bateria entre 50 e 100%
		Location:     [2]int{rand.Intn(250), rand.Intn(250)},  // Coordenadas aleatórias entre 0 e 250
	}

	// Cria a estrutura de estado interno do carro
	carState := &CarState{
		ID:           car.ID,
		Location:     car.Location,
		BatteryLevel: car.BatteryLevel,
		Status:       "normal",
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
	go handleRequests(car, conn) // Escuta requisições do servidor
	go carMovement(carState, criticalChan) // Simula movimentação e monitora bateria
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
				continue
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
			fmt.Println("Erro ao enviar dados críticos:", err)
			return
		}
		fmt.Println("\n✅ Dados críticos enviados ao servidor:", data)
	}
}

// Escuta requisições do servidor e responde com os dados do carro
func handleRequests(car Car, conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Erro ao ler requisição do servidor:", err)
			return
		}

		
		// Decodifica a requisição recebida (espera JSON)
		var request map[string]string
		err = json.Unmarshal(buf[:n], &request)
		if err != nil {
			fmt.Println("Erro ao decodificar requisição:", err)
			return
		}

		// Se a ação for de requisição de dados, envia os dados do carro
		if request["action"] == "request_car_data" {
			sendCarData(car, conn)
			fmt.Println("Dados do carro enviados ao servidor.")
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

// Envia os dados do carro para o servidor em formato JSON
func sendCarData(car Car, conn net.Conn) {
	jsonData, err := json.Marshal(car)
	if err != nil {
		fmt.Println("Erro ao converter para JSON:", err)
		return
	}
	conn.Write(jsonData)
}

// Exibe o nível da bateria de forma visual
func displayBattery(car CarState) {
	totalBars := 20
	numHashMarks := (car.BatteryLevel * totalBars) / 100

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
