package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string  `json:"name"`
	Code        string  `json:"code"`
	Stock       int     `json:"stock"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	Price       float64 `json:"price"` // Harga dalam mata uang default, misal USD
}
