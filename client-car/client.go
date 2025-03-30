package main

import (
	"fmt"
	"math/rand"
	"net"
	"time"
	"strings"
)

type Car struct {
	ID             int    `json:"id"`
	//User         User   `json:"name"`
	BatteryLevel   int    `json:"batteryLevel"`
	Location       [2]int `json:"location"`
}


func carMovement(car Car, conn net.Conn) int {
	for { // Loop infinito para atualizar as posições
		time.Sleep(time.Second) // Espera 1 segundo a cada atualização

		/* Atualizando as coordenadas */
		// Atualiza a posição do carro
		car.Location[0] += rand.Intn(11) // Movimento no eixo X
		car.Location[1] += rand.Intn(11) // Movimento no eixo Y

		// Atualiza o nível da bateria
		car.BatteryLevel = batteryLevel(car.BatteryLevel)

		// Exibe a bateria no terminal
		displayBattery(car)

		// Verifica se a bateria está em nível crítico
		checkCriticalLevel(car.BatteryLevel, car.ID)

		// Formata os dados como string ("car: [x, y]"). Envia as coordenadas e o nível de bateria
		data := fmt.Sprintf("%d, %d, %d\n",
			car.Location[0], car.Location[1], car.BatteryLevel)

		// Envia os dados para o servidor
		_, err := conn.Write([]byte(data))
		if err != nil {
			fmt.Println("Erro ao enviar dados:", err)
			break
		}

		//fmt.Println("Dados enviados:", data)
	}

	return car.BatteryLevel
}

// Atualiza o nível da bateria do carro
func batteryLevel(batteryLevel int) int {
	//batteryConsumption := rand.Intn(11) // Consumo de bateria aleatório de 0% a 10%
	batteryLevel -= 5 // Diminui a bateria

	if batteryLevel <= 20 {
		batteryLevel = 20 // Garante que não fique negativo
	}

	fmt.Println("Nível de bateria:", batteryLevel)

	return batteryLevel
}

// Verifica se a bateria está em nível crítico
func checkCriticalLevel(batteryLevel int, carID int) {
	if batteryLevel <= 20 {
		fmt.Printf("⚠️  ALERTA: --- CARRO %d --- Bateria em nível crítico! 🚨 Nível de Bateria: %d%%\n", carID, batteryLevel)
	}
}

func getCarID() int {
	// Pegamos um número aleatório baseado no timestamp atual
	carID := int(time.Now().UnixNano() % 10000) // Pegamos os últimos 4 dígitos
	fmt.Printf("🆔 ID do carro gerado: %d\n", carID)
	return carID
}


// Função para exibir a barra de bateria no terminal
func displayBattery(car Car) {
	totalBars := 20              // Total de "blocos" da barra
	batteryPercentage := car.BatteryLevel
	numHashMarks := (batteryPercentage * totalBars) / 100 // Quantos "#" mostrar

	// Exibe a interface de bateria
	fmt.Printf("\n      USER    ID: %d\n", car.ID)
	fmt.Println(" -----------------------")
	// Usando strings.Repeat para repetir os caracteres
	fmt.Printf("|%s   |  %d%%\n", strings.Repeat("#", numHashMarks) + strings.Repeat(" ", totalBars - numHashMarks), batteryPercentage)
	fmt.Println(" -----------------------")
}


func main() {
	rand.Seed(time.Now().UnixNano()) // Inicializa a semente aleatória

	// Obtém o ID do carro da variável de ambiente
	carID := getCarID()

	// Criando os carros; Nível inicial da bateria (100%)
	car := Car{
		ID: carID,
		//User: models.User{Name: "João"},
		BatteryLevel: 100,
		Location: [2]int{
			rand.Intn(100),
			rand.Intn(100),
		},
	}

	// Conecta ao servidor na porta 8080
	conn, err := net.Dial("tcp", "server:8080")
	if err != nil {
		panic(err)
	}

	defer conn.Close() // Fecha a conexão

	// Teste de conexão com o servidor
	fmt.Println("Conectando ao servidor...")
	
	fmt.Printf("🚗 Carro %d conectado ao servidor!\n", car.ID)

	// Inicia a movimentação dos carros. Atualiza e envia as coordenadas ao servidor
	carMovement(car, conn)
}
