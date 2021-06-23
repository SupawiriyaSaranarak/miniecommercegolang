package model

import "time"

// https://gorm.io/docs/model.html

type User struct {
	ID        uint   `gorm:"primary_key"`
	Email     string `gorm:"unique" form:"email" binding"required"`
	Password  string `form:"password" binding"required"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Address   string `json:"address"`
	Phone     string `json:"phone"`
	Status    string `gorm:"default:user"`
	CreatedAt time.Time
}
