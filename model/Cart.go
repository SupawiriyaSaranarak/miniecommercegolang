package model

import "time"

type Cart struct {
	ID        uint   `gorm:"primary_key"`
	ProductID string `json:"productId"`
	Amount    int64  `json:"amount"`
	UserID    string `json:"userId"`
	CreatedAt time.Time
}
