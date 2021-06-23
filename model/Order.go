package model

import "time"

type Order struct {
	ID        uint    `gorm:"primary_key"`
	ProductID string  `json:"productId"`
	Amount    int64   `json:"amount"`
	Status    string  `gorm:"default:ordered"`
	Price     float64 `json:"price"`
	UserID    string  `json:"userId"`
	CreatedAt time.Time
}
