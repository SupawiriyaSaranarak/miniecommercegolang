package api

import (
	"github.com/gin-gonic/gin"
	"main/db"
)
 // Setup - call this method to setup routes
func Setup(router *gin.Engine) {
	db.SetupDB()
	SetupAuthenAPI(router)
	SetupProductAPI(router)
	SetupWalletAPI(router)
	SetupOrderAPI(router)
	SetupCartAPI(router)
}