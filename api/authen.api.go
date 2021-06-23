package api

import (
"main/db"
"main/model"
"main/interceptor"
"net/http"
"time" 
"golang.org/x/crypto/bcrypt"
"github.com/gin-gonic/gin"
)

func SetupAuthenAPI(router *gin.Engine) {
	authenAPI := router.Group("/api")
	{
		authenAPI.POST("/login", login)
		authenAPI.POST("/register", register)
	}	
}
func checkPasswordHash(password,hash string) bool {
	err := bcrypt.CompareHashAndPassword( []byte(hash), []byte(password))
	return err == nil
}
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword( []byte(password), 14)
	return string(bytes), err
}
func login(c *gin.Context) {
	var user model.User 
	if c.ShouldBind(&user) == nil {
		var queryUser model.User 
		if err := db.GetDB().First(&queryUser, "email = ?", user.Email).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"result": "nok", "error": err})
		} else if checkPasswordHash(user.Password, queryUser.Password) == false {
			c.JSON(http.StatusUnauthorized, gin.H{"result": "nok", "error": "invalid password"})
		} else {
			token := interceptor.JwtSign(queryUser)

			c.JSON(200, gin.H{"result": "login", "token": token})
		}
	} else {
		c.JSON(401, gin.H{"status": "unable to bind data"})
	}
	
}
func register(c *gin.Context) {
	var user model.User
	if c.ShouldBind(&user) == nil {
		user.Password, _ = hashPassword(user.Password)
		user.CreatedAt = time.Now()
		if err := db.GetDB().Create(&user).Error; err != nil {
			c.JSON(200, gin.H{"result": "nok", "err": err})
		} else {
			c.JSON(200, gin.H{"result": "ok", "data": user})
		}
	} else {
		c.JSON(401, gin.H{"status": "unable to bind data"})
	}	
}
