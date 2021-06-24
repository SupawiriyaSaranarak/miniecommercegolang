package api

import (
	"context"
	"fmt"
	"log"
	"main/db"
	"main/interceptor"
	"main/model"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gin-gonic/gin"
)

// SetupProductAPI - call this method to setup product group route
func SetupProductAPI(router *gin.Engine) {
	productAPI := router.Group("/api")
	{
		productAPI.GET("/product", interceptor.JwtVerify, getProduct)
		productAPI.GET("/my-product", interceptor.JwtVerify, getMyProduct)
		productAPI.GET("/product/:id", interceptor.JwtVerify, getProductByID)
		productAPI.POST("/product", interceptor.JwtVerify, createProduct)
		productAPI.PUT("/product", interceptor.JwtVerify, editProduct)
	}
}
func getProduct(c *gin.Context) {
	var product []model.Product

	keyword := c.Query("keyword")
	if keyword != "" {
		keyword = fmt.Sprintf("%%%s%%", keyword)
		db.GetDB().Where("name like ?", keyword).Order("created_at DESC").Find(&product)
	} else {
		db.GetDB().Find(&product)
	}

	c.JSON(http.StatusOK, product)
}
func getProductByID(c *gin.Context) {
	var product model.Product
	db.GetDB().Where("id=?", c.Param("id")).Order("created_at DESC").Find(&product)
	c.JSON(http.StatusOK, product)
}
func getMyProduct(c *gin.Context) {
	var product []model.Product
	userID := c.GetString("jwt_user_id")
	db.GetDB().Where("user_id = ?", userID).Order("created_at DESC").Find(&product)
	c.JSON(http.StatusOK, product)
}
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
func saveImage(image *multipart.FileHeader, product *model.Product, c *gin.Context) {
	if image != nil {
		runningDir, _ := os.Getwd()
		product.ProductImg = image.Filename
		extension := filepath.Ext(image.Filename)
		fileName := fmt.Sprintf("%d%s", product.ID, extension)
		filePath := fmt.Sprintf("%s/uploaded/images/%s", runningDir, fileName)
		// CLOUDINARY_URL=cloudinary://971856587169986:DcVriGPSWQzAGSMFujdMQQ_qGhU@dvdl7oduu
		if fileExists(filePath) {
			os.Remove(filePath)
		}
		c.SaveUploadedFile(image, filePath)

		//save file to cloudinary
		var cld, err = cloudinary.NewFromURL("cloudinary://971856587169986:DcVriGPSWQzAGSMFujdMQQ_qGhU@dvdl7oduu")
		if err != nil {
			log.Fatalf("Failed to intialize Cloudinary, %v", err)
		}
		var ctx = context.Background()
		// fileDir := fmt.Sprintf("/Users/Glao/Desktop/CodeCamp/mini-e-commerce/mini/uploaded/images/%s", fileName)
		resp, err := cld.Upload.Upload(ctx, filePath, uploader.UploadParams{})
		//filePath
		if err != nil {
			log.Fatalf("Failed to upload file, %v\n", err)
		}
		log.Println(*resp)
		//save file to cloudinary

		db.GetDB().Model(&product).Update("product_img", resp.SecureURL)
	}
}
func createProduct(c *gin.Context) {
	product := model.Product{}
	product.Name = c.PostForm("name")
	product.Stock, _ = strconv.ParseInt(c.PostForm("stock"), 10, 64)
	product.Price, _ = strconv.ParseFloat(c.PostForm("price"), 64)
	product.UserID = c.GetString("jwt_user_id")
	product.CreatedAt = time.Now()
	db.GetDB().Create(&product)

	image, _ := c.FormFile("productImg")
	saveImage(image, &product, c)

	c.JSON(200, gin.H{"result": product})

}

func editProduct(c *gin.Context) {
	var product model.Product
	id, _ := strconv.ParseInt(c.PostForm("id"), 10, 32)
	productID := uint(id)
	userID := c.GetString("jwt_user_id")
	result := db.GetDB().Where("id = ? AND user_id = ?", productID, userID).First(&product)
	fmt.Println(result.RowsAffected)
	fmt.Println(result.Error)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Cannot edit other's product."})
	} else {

		product.Stock, _ = strconv.ParseInt(c.PostForm("stock"), 10, 64)

		// db.GetDB().Model(&product).Update("stock", product.Stock)
		db.GetDB().Save(&product)

		image, _ := c.FormFile("productImg")
		saveImage(image, &product, c)

		c.JSON(http.StatusOK, gin.H{"result": product})
	}
}

// db.Model(&product).UpdateColumn("quantity", gorm.Expr("quantity - ?", 1))
// db.GetDB().Model(&wallet).Update("payment_img", fileName)
