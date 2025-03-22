package main

import (
	"fmt"
	"net"
	"math/rand" 
    "time"
)


func carMovement(carCoordinates map[string][]int, battery int, conn net.Conn) int {
	for { // Loop infinito para atualizar as posições
        time.Sleep(time.Second) // Espera 1 segundo a cada atualização

		/* Atualizando as coordenadas */
        for car := range carCoordinates {
            // Atualiza a posição do carro
            carCoordinates[car][0] += rand.Intn(11) // Movimento no eixo X
            carCoordinates[car][1] += rand.Intn(11) // Movimento no eixo Y
        }

		// Atualiza o nível da bateria
		battery = batteryLevel(battery)
		// Verifica se a bateria está em nível crítico
		checkCriticalLevel(battery)

		// Formata os dados como string ("car1: [x, y], car2: [x, y]")
		data := fmt.Sprintf("car1: [%d, %d], car2: [%d, %d]\n",
			carCoordinates["car1"][0], carCoordinates["car1"][1],
			carCoordinates["car2"][0], carCoordinates["car2"][1])

		// Envia os dados para o servidor
		_, err := conn.Write([]byte(data))
		if err != nil {
			fmt.Println("Erro ao enviar dados:", err)
			break
		}

		fmt.Println("Dados enviados:", data)

		// Se a bateria acabar, parar a movimentação
		if battery == 0 {
			fmt.Println("🔋 O carro parou! Bateria esgotada! 🚨")
			break
		}
	}
	return battery
}


// Atualiza o nível da bateria do carro
func batteryLevel(batteryLevel int) int {
	//batteryConsumption := rand.Intn(11) // Consumo de bateria aleatório de 0% a 10%
	batteryLevel -= 5 // Diminui a bateria

	if batteryLevel < 0 {
		batteryLevel = 0 // Garante que não fique negativo
	}

	fmt.Println("Nível de bateria:", batteryLevel)

	return batteryLevel
}


// Verifica se a bateria está em nível crítico
func checkCriticalLevel(batteryLevel int) {
	if batteryLevel <= 20 {
		fmt.Println("⚠️  ALERTA: Bateria em nível crítico! 🚨", batteryLevel, "%")
	}
}


func main() {
	rand.Seed(time.Now().UnixNano()) // Inicializa a semente aleatória
	carCoordinates := map[string][]int{
		"car1": {rand.Intn(100), rand.Intn(100)},
		"car2": {rand.Intn(100), rand.Intn(100)},
	}

	// Nível inicial da bateria (100%)
	batteryLevel := 100

	// Conecta ao servidor na porta 8080
	conn, err := net.Dial("tcp", "server:8080") 
	if err != nil {
		panic(err)
	}
	defer conn.Close() // Fecha a conexão

	// Teste de conexão com o servidor
	fmt.Println("Conectando ao servidor...")

	// Inicia a movimentação dos carros. Atualiza e envia as coordenadas ao servidor
	carMovement(carCoordinates, batteryLevel, conn)
}