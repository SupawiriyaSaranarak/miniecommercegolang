package db

import (
	"main/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

//GetDB - call this method to get db
func GetDB() *gorm.DB {
	return db
}

//SetupDB - setup database for sharing to all api
func SetupDB() {
	dsn := "user=postgres password=Aa123456 dbname=miniecommerce port=5432 sslmode=disable TimeZone=Asia/Bangkok"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	database.AutoMigrate(&model.User{})
	database.AutoMigrate(&model.Product{})
	database.AutoMigrate(&model.Order{})
	database.AutoMigrate(&model.Cart{})
	database.AutoMigrate(&model.Wallet{})
	

	db = database
}
