package models

type Car struct {
	ID int `json:"id"`
	//User         User  `json:"name"`
	BatteryLevel int    `json:"id"`
	Location     [2]int `json:"location"`
}
