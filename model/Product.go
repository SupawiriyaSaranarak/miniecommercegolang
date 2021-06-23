package model

import "time"

type Product struct {
	ID uint `gorm:"primary_key"`
	Name string 
	Stock int64
	Price float64
	ProductImg string
	UserID string
	CreatedAt time.Time
}