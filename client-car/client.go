package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

type Car struct {
	Type         string `json:"type"`
	ID           int    `json:"id"`
	BatteryLevel int    `json:"batteryLevel"`
	Location     [2]int `json:"location"`
}

type CarState struct {
	ID           int
	Location     [2]int
	BatteryLevel int
	Status       string // "normal", "critico"
}

func main() {
	rand.Seed(time.Now().UnixNano())

	carID := getCarID()

	car := Car{
		ID:           carID,
		BatteryLevel: rand.Intn(51) + 50,
		Location:     [2]int{rand.Intn(100), rand.Intn(100)},
	}

	carState := &CarState{
		ID:           car.ID,
		Location:     car.Location,
		BatteryLevel: car.BatteryLevel,
		Status:       "normal",
	}

	conn, err := net.Dial("tcp", "server:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Printf("\n🚗 Carro %d conectado ao servidor!\n", car.ID)

	ident := map[string]interface{}{
		"type": "car",
		"id":   car.ID,
	}
	jsonData, _ := json.Marshal(ident)
	conn.Write(jsonData)

	criticalChan := make(chan CarState)

	go handleRequests(car, conn)
	go carMovement(carState, criticalChan)
	go handleCriticalData(conn, criticalChan)

	select {} // Mantém o programa rodando
}


// ===================== GOROUTINES =========================

/*
	 Essas goroutines funcionam simultaneamente, ou seja, o carro está:
	 - Movendo-se
	 - Ouvindo requisições do servidor
	 - Enviando dados críticos para o servidor
*/

func carMovement(car *CarState, criticalChan chan CarState) {
	for {
		time.Sleep(time.Second)

		if car.Status == "normal" {
			car.Location[0] += rand.Intn(11)
			car.Location[1] += rand.Intn(11)
			car.BatteryLevel = batteryLevel(car.BatteryLevel)

			displayBattery(*car)

			if car.BatteryLevel <= 20 {
				car.Status = "critico"
				fmt.Println("⚠️  ALERTA! 🚨 Bateria crítica! ")
				criticalChan <- *car
				continue
			}

			fmt.Printf("📍 Coordenadas: %v\n", car.Location)
		}
	}
}

func handleCriticalData(conn net.Conn, criticalChan chan CarState) {
	for {
		carCritical := <-criticalChan
		data := fmt.Sprintf("%d, %d, %d\n", carCritical.Location[0], carCritical.Location[1], carCritical.BatteryLevel)

		_, err := conn.Write([]byte(data))
		if err != nil {
			fmt.Println("Erro ao enviar dados críticos:", err)
			return
		}
		fmt.Println("\n✅ Dados críticos enviados ao servidor:", data)
	}
}

func handleRequests(car Car, conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Erro ao ler requisição do servidor:", err)
			return
		}

		var request map[string]string
		err = json.Unmarshal(buf[:n], &request)
		if err != nil {
			fmt.Println("Erro ao decodificar requisição:", err)
			return
		}

		if request["action"] == "request_car_data" {
			sendCarData(car, conn)
			fmt.Println("Dados do carro enviados ao servidor.")
		}
	}
}

// ===================== AUXILIARES =========================

func batteryLevel(batteryLevel int) int {
	batteryLevel -= 5
	if batteryLevel <= 20 {
		batteryLevel = 20
	}
	return batteryLevel
}

func sendCarData(car Car, conn net.Conn) {
	jsonData, err := json.Marshal(car)
	if err != nil {
		fmt.Println("Erro ao converter para JSON:", err)
		return
	}
	conn.Write(jsonData)
}

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

func getCarID() int {
	return int(time.Now().UnixNano() % 10000)
}
