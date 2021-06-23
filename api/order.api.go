package api

import (
	"fmt"

	"main/db"
	"main/interceptor"
	"main/model"

	"net/http"

	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetuporderAPI - call this method to setup order group route
func SetupOrderAPI(router *gin.Engine) {
	orderAPI := router.Group("/api")
	{
		orderAPI.GET("/order-buyer", interceptor.JwtVerify, getOrderForBuyer)
		orderAPI.GET("/order-seller", interceptor.JwtVerify, getOrderForSeller)
		orderAPI.GET("/order/:id", interceptor.JwtVerify, getOrderByID)
		orderAPI.POST("/order", interceptor.JwtVerify, createOrder)
		orderAPI.PUT("/order/:id/:status", interceptor.JwtVerify, editOrderStatus)
	}
}

type OrderResultForBuyer struct {
	ID           uint
	ProductID    uint
	ProductName  string
	ProductImage string
	PricePerUnit float64
	Amount       string
	TotalPrice   float64
	Status       string
}

func getOrderForBuyer(c *gin.Context) {
	status := c.Query("status")
	userId := c.GetString("jwt_user_id")
	var result []OrderResultForBuyer
	db.GetDB().Debug().Raw("SELECT o.id as id, p.id as product_id, p.name as product_name, p.product_img as product_image, p.price as price_per_unit, o.amount as Amount, o.price as Total_Price, o.status as Status FROM orders o JOIN products p ON CAST(o.product_id as BIGINT) = p.id WHERE o.user_id = ? AND o.status = ?", userId, status).Scan(&result)
	fmt.Printf("%v", result)
	c.JSON(200, result)
}

type OrderResultForSeller struct {
	ID           uint
	ProductID    uint
	ProductName  string
	ProductImage string
	PricePerUnit float64
	Amount       string
	TotalPrice   float64
	ClientFname  string
	ClientLname  string
	Address      string
	Status       string
}

func getOrderForSeller(c *gin.Context) {
	status := c.Query("status")
	userId := c.GetString("jwt_user_id")
	var result []OrderResultForSeller
	db.GetDB().Debug().Raw("SELECT o.id as ID, p.id as Product_ID, p.name as Product_Name, p.product_img as product_image, p.price as Price_Per_Unit, o.amount as Amount, o.price as TotalPrice, o.status as Status, u.first_name as client_fname, u.last_name as client_lname, u.address as address FROM orders o JOIN products p ON CAST(o.product_id as BIGINT) = p.id JOIN users u ON CAST(o.user_id as BIGINT) = u.id WHERE p.user_id = ? AND o.status = ?", userId, status).Scan(&result)
	fmt.Printf("%v", result)
	c.JSON(200, result)
}
func getOrderByID(c *gin.Context) {
	var order model.Order
	db.GetDB().Where("id=?", c.Param("id")).Find(&order)
	c.JSON(http.StatusOK, order)
}

func createOrder(c *gin.Context) {
	var order model.Order
	var product model.Product
	if c.ShouldBind(&order) == nil {
		buyerUserID := c.GetString("jwt_user_id")
		result := db.GetDB().Where("id = ? AND user_id = ?", order.ProductID, buyerUserID).First(&product)
		if result.RowsAffected == 1 {
			c.JSON(400, gin.H{"message": "Cannot buy your own product."})
		} else {
			db.GetDB().Where("id = ? ", order.ProductID).First(&product)
			order.UserID = c.GetString("jwt_user_id")
			order.CreatedAt = time.Now()
			if err := db.GetDB().Create(&order).Error; err != nil {
				c.JSON(200, gin.H{"result": "nok", "err": err})
			} else {
				db.GetDB().Model(&product).UpdateColumn("stock", gorm.Expr("stock - ?", order.Amount))
				c.JSON(200, gin.H{"result": "ok", "data1": order})
			}
		}

	} else {
		c.JSON(401, gin.H{"status": "unable to bind data"})
	}

	// order := model.Order{}
	// product := model.Product{}
	// order.ProductID = c.PostForm("product_id")
	// produtId,_ := strconv.ParseUint(c.PostForm("product_id"), 10, 32)
	// db.GetDB().Where("id = ?", produtId).First(&product)
	// order.Amount, _ = strconv.ParseInt(c.PostForm("stock"), 10, 64)
	// order.Price, _ = strconv.ParseFloat(c.PostForm("price"), 64)
	// order.UserID = c.PostForm("user_id")
	// order.CreatedAt = time.Now()
	// db.GetDB().Create(&order)
	// db.GetDB().Model(&product).UpdateColumn("stock", gorm.Expr("stock - ?", order.Amount))

	c.JSON(200, gin.H{"result": order})

}

func editOrderStatus(c *gin.Context) {
	var order model.Order
	db.GetDB().Where("id = ?", c.Param("id")).First(&order)
	order.Status = c.Param("status")

	db.GetDB().Save(&order)

	c.JSON(http.StatusOK, gin.H{"result": order})

}

// db.Model(&order).UpdateColumn("quantity", gorm.Expr("quantity - ?", 1))
// db.GetDB().Model(&wallet).Update("payment_img", fileName)

// SELECT transactions.id, total, paid, change, payment_type, payment_detail, order_list, users.username as Staff, transactions.created_at FROM transactions join users on transactions.staff_id = users.id
// SELECT o.id as Order_ID, p.id as Product_ID, p.name as Product_Name, p.price as Price_Per_Unit, o.amount as Amount, o.price as TotalPrice, o.status as Status, u.first_name as client_fname, u.last_name as client_lname, u.address as address FROM orders o JOIN products p ON CAST(o.product_id as BIGINT) = p.id JOIN users u ON CAST(o.user_id as BIGINT) = u.id WHERE p.user_id = '3' AND o.status = 'ordered'
