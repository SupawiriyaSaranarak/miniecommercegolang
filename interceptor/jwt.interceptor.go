package interceptor

import (
	"fmt"
	"main/model"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var secretKey = "ao[wifjwi[ajf[aw"

func JwtSign(payload model.User) string {
	atClaims := jwt.MapClaims{}
	//Payload begin
	atClaims["id"] = payload.ID
	atClaims["email"] = payload.Email
	atClaims["status"] = payload.Status
	atClaims["exp"] = time.Now().Add(time.Minute * 2000).Unix()
	//Payload end
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, _ := at.SignedString([]byte(secretKey))
	return token
}

func JwtVerify(c *gin.Context) {
	tokenString := strings.Split(c.Request.Header["Authorization"][0], " ")[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method : %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
		userID := fmt.Sprintf("%v", claims["id"])
		email := fmt.Sprintf("%v", claims["email"])
		status := fmt.Sprintf("%v", claims["status"])

		c.Set("jwt_user_id", userID)
		c.Set("jwt_email", email)
		c.Set("jwt_status", status)

		c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"result": "nok", "message": "invalid token", "error": err})
		c.Abort()
	}
}
