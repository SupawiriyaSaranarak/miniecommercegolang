package main

import (
	"main/api"
	"time"

	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

func main() {
	router := gin.Default()
	//setup cors middleware options
	config := cors.Config{
		Origins:         "*",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		Methods:         "GET, POST, PUT, DELETE",
		Credentials:     true,
		ValidateHeaders: false,
		MaxAge:          1 * time.Minute,
	}
	//Apply middleware to router
	router.Use(cors.Middleware(config))

	router.Static("images", "./uploaded/images")
	//โลกภายนอกเข้าถึงรูปภาพ locolhost:8081/images/<ชื่อไฟล์ที่อยู๋ใน./uploaded/images>

	api.Setup(router)
	router.Run(":8081")
}
