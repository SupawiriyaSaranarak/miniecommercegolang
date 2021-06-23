package api

import (
	"fmt"
	"main/db"
	"main/interceptor"
	"main/model"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupWalletAPI(router *gin.Engine) {
	walletAPI := router.Group("/api")
	{
		walletAPI.GET("/wallet", interceptor.JwtVerify, getBalance)
		walletAPI.POST("/wallet", interceptor.JwtVerify, createWallet)
	}
}

// func getTransaction(c *gin.Context) {
// 	var transactions []model.Transaction
// 	db.GetDB().Find(&transactions)
// 	c.JSON(200, transactions)
// }

// https://gorm.io/docs/query.html

// type TransactionResult struct {
// 	ID uint
// 	Total float64
// 	Paid float64
// 	Change float64
// 	PaymentType string
// 	PaymentDetail string
// 	OrderList string
// 	Staff string
// 	CreatedAt time.Time
// }
// func getWallet(c *gin.Context) {
// 	var result []TransactionResult
// 	db.GetDB().Debug().Raw("SELECT transactions.id, total, paid, change, payment_type, payment_detail, order_list, users.username as Staff, transactions.created_at FROM transactions join users on transactions.staff_id = users.id", nil).Scan(&result)
// 	fmt.Printf("%v", result)
// 	c.JSON(200, result)
// 	c.JSON(401, gin.H{"result": "get wallet"})
// }
func filePaymentExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
func savePaymentImage(image *multipart.FileHeader, wallet *model.Wallet, c *gin.Context) {
	if image != nil {
		runningDir, _ := os.Getwd()
		wallet.PaymentImg = image.Filename
		extension := filepath.Ext(image.Filename)
		fileName := fmt.Sprintf("slip-%d%s", wallet.ID, extension)
		filePath := fmt.Sprintf("%s/uploaded/images/%s", runningDir, fileName)

		if filePaymentExists(filePath) {
			os.Remove(filePath)
		}
		c.SaveUploadedFile(image, filePath)
		db.GetDB().Model(&wallet).Update("payment_img", fileName)
	}
}

func createWallet(c *gin.Context) {
	var wallet model.Wallet
	wallet.Value, _ = strconv.ParseFloat(c.PostForm("value"), 64)
	wallet.UserID = c.GetString("jwt_user_id")
	wallet.CreatedAt = time.Now()
	db.GetDB().Create(&wallet)
	image, _ := c.FormFile("paymentImg")
	savePaymentImage(image, &wallet, c)
	c.JSON(http.StatusOK, gin.H{"result": "ok", "data": wallet})
}

type WalletBalance struct {
	Balance float64
}

func getBalance(c *gin.Context) {

	userID := c.GetString("jwt_user_id")
	var result []WalletBalance
	db.GetDB().Debug().Raw("SELECT (SELECT sum (w.value) FROM wallets w WHERE w.user_id = ? ) - (SELECT sum (o.price) FROM orders o WHERE o.user_id = ? ) AS balance", userID, userID).Scan(&result)
	fmt.Printf("%v", result)
	c.JSON(200, result)
}
