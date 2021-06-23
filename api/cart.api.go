package api

import (
	"fmt"
	"main/db"
	"main/interceptor"
	"main/model"

	"time"

	"github.com/gin-gonic/gin"
)

// SetupcartAPI - call this method to setup order group route
func SetupCartAPI(router *gin.Engine) {
	cartAPI := router.Group("/api")
	{
		cartAPI.GET("/cart", interceptor.JwtVerify, getCart)
		cartAPI.POST("/cart", interceptor.JwtVerify, createCart)
		cartAPI.DELETE("/cart/:id", interceptor.JwtVerify, deleteCart)
	}
}

// func getOrder(c *gin.Context) {
// 	var order []model.Order

// 	keyword := c.Query("keyword")
// 	if keyword != "" {
// 		keyword = fmt.Sprintf("%%%s%%", keyword)
// 		db.GetDB().Where("name like ?", keyword).Find(&order)
// 	} else {
// 		db.GetDB().Find(&order)
// 	}

// 	c.JSON(http.StatusOK, order)
// }
type CartResult struct {
	ID           uint
	ProductID    uint
	ProductName  string
	ProductImage string
	PricePerUnit float64
	Stock        int
	Amount       string
}

func getCart(c *gin.Context) {
	userId := c.GetString("jwt_user_id")
	var result []CartResult
	db.GetDB().Debug().Raw("SELECT c.id as id, p.id as product_id, p.name as product_name, p.product_img as product_image, p.price as price_per_unit, p.stock as Stock, c.amount as amount FROM carts c JOIN products p ON CAST(c.product_id as BIGINT) = p.id WHERE c.user_id = ?", userId).Scan(&result)
	fmt.Printf("%v", result)
	c.JSON(200, result)
}

func createCart(c *gin.Context) {
	var cart model.Cart
	var product model.Product
	if c.ShouldBind(&cart) == nil {
		buyerUserID := c.GetString("jwt_user_id")
		result := db.GetDB().Where("id = ? AND user_id = ?", cart.ProductID, buyerUserID).First(&product)
		if result.RowsAffected == 1 {
			c.JSON(400, gin.H{"message": "Cannot add your own product to cart."})
		} else {
			db.GetDB().Where("id = ? ", cart.ProductID).First(&product)
			cart.UserID = c.GetString("jwt_user_id")
			cart.CreatedAt = time.Now()
			if err := db.GetDB().Create(&cart).Error; err != nil {
				c.JSON(200, gin.H{"result": "nok", "err": err})
			} else {
				c.JSON(200, gin.H{"result": "ok", "data1": cart})
			}
		}

	} else {
		c.JSON(401, gin.H{"status": "unable to bind data"})
	}

}

func deleteCart(c *gin.Context) {
	var cart model.Cart
	buyerUserID := c.GetString("jwt_user_id")
	result := db.GetDB().Where("id = ? AND user_id = ?", c.Param("id"), buyerUserID).Delete(&cart)

	if result.RowsAffected == 1 {
		c.JSON(200, gin.H{"message": "The product has already been remove."})
	} else {
		c.JSON(400, gin.H{"message": "Cannot remove product from other's cart."})
	}
}

// db.Model(&order).UpdateColumn("quantity", gorm.Expr("quantity - ?", 1))
// db.GetDB().Model(&wallet).Update("payment_img", fileName)
