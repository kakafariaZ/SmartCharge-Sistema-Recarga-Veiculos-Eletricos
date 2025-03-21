package main

import (
	"fmt"
	"net"
	"math/rand"
    "time"
)


func carMovement(carCoordinates map[string][]int, conn net.Conn) {	
	for { // Loop infinito para atualizar as posições
        time.Sleep(time.Second) // Espera 1 segundo a cada atualização

		/* Atualizando as coordenadas */
        for car := range carCoordinates {
            // Atualiza a posição do carro
            carCoordinates[car][0] += rand.Intn(11) // Movimento no eixo X
            carCoordinates[car][1] += rand.Intn(11) // Movimento no eixo Y
        }

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
	}
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Inicializa a semente aleatória
	carCoordinates := map[string][]int{
		"car1": {rand.Intn(100), rand.Intn(100)},
		"car2": {rand.Intn(100), rand.Intn(100)},
	}

	// Conecta ao servidor na porta 8080
	conn, err := net.Dial("tcp", "server:8080") 
	if err != nil {
		panic(err)
	}
	defer conn.Close() // Fecha a conexão

	// Teste de conexão com o servidor
	fmt.Println("Conectando ao servidor...")

	// Inicia a movimentação dos carros. Atualiza e envia as coordenadas ao servidor
	carMovement(carCoordinates, conn)
}