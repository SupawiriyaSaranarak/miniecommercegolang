package model

import (
	"time"
)

type Wallet struct {
	ID         uint    `gorm:"primary_key"`
	UserID     string  `json:"userId"`
	Value      float64 `json:"value"`
	PaymentImg string  `json:"paymentImg"`
	CreatedAt  time.Time
}
