package models

type ChargeStation struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Location     [2]int  `json:"location"`
	Occupation   bool    `json:"occupation"`
	Capacity     int     `json:"capacity"`
	Power        int     `json:"power"`
	Price        float64 `json:"price"`
}


