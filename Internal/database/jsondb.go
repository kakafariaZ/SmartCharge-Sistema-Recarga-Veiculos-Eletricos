package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Car struct {
	ID      string `json:"id"`
	Battery int    `json:"battery"`
}

type Station struct {
	ID    string `json:"id"`
	Slots int    `json:"slots"`
}

func saveData(filename string, data interface{}) {
	file, _ := os.Create(filename)
	defer file.Close()
	json.NewEncoder(file).Encode(data)
}

func loadData(filename string, data interface{}) {
	file, _ := os.Open(filename)
	defer file.Close()
	json.NewDecoder(file).Decode(data)
}

func main() {
	cars := []Car{{ID: "C1", Battery: 100}}
	stations := []Station{{ID: "S1", Slots: 2}}

	saveData("cars.json", cars)
	saveData("stations.json", stations)

	var loadedCars []Car
	var loadedStations []Station

	loadData("cars.json", &loadedCars)
	loadData("stations.json", &loadedStations)

	fmt.Println("Carros:", loadedCars)
	fmt.Println("Postos:", loadedStations)
}
