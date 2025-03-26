package main

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

type Car struct {
	ID int `json:"id"`
	//User         User  `json:"name"`
	BatteryLevel int    `json:"batteryLevel"`
	Location     [2]int `json:"location"`
}

func carMovement(car Car, conn net.Conn) int {

	for { // Loop infinito para atualizar as posições
		time.Sleep(time.Second) // Espera 1 segundo a cada atualização

		/* Atualizando as coordenadas */
		// Atualiza a posição do carro
		car.Location[0] += rand.Intn(11) // Movimento no eixo X
		car.Location[1] += rand.Intn(11) // Movimento no eixo Y

		// Atualiza o nível da bateria
		car.BatteryLevel = batteryLevel(car1.BatteryLevel)

		// Verifica se a bateria está em nível crítico
		checkCriticalLevel(car.BatteryLevel)

		// Formata os dados como string ("car: [x, y]")
		data := fmt.Sprintf("%d, %d\n",
			car.Location[0], car.Location[1],)

		// Envia os dados para o servidor
		_, err := conn.Write([]byte(data))
		if err != nil {
			fmt.Println("Erro ao enviar dados:", err)
			break
		}

		fmt.Println("Dados enviados:", data)

		// Verifica o nível da bateria
		// Se a bateria acabar, parar a movimentação
		if car.BatteryLevel == 0 {
			fmt.Println("🔋 O carro parou! Bateria esgotada! 🚨")
			break
		}
	}
	return car.BatteryLevel
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

	// Criando os carros; Nível inicial da bateria (100%)
	car1 := Car{
		ID: 1,
		//User: models.User{Name: "João"},
		BatteryLevel: 100,
		Location: [2]int{
			rand.Intn(100),
			rand.Intn(100),
		},
	}

	car2 := Car{
		ID: 1,
		//User: models.User{Name: "João"},
		BatteryLevel: 100,
		Location: [2]int{
			rand.Intn(100),
			rand.Intn(100),
		},
	}

	// carCoordinates := map[string][]int{
	// 	"car1": {rand.Intn(100), rand.Intn(100)},
	// 	"car2": {rand.Intn(100), rand.Intn(100)},
	// }

	// Nível inicial da bateria (100%)
	// batteryLevel := 100

	// Conecta ao servidor na porta 8080
	conn, err := net.Dial("tcp", "server:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close() // Fecha a conexão

	// Teste de conexão com o servidor
	fmt.Println("Conectando ao servidor...")

	// Inicia a movimentação dos carros. Atualiza e envia as coordenadas ao servidor
	carMovement(car1, car2, conn)
}
