package middleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var JWT_SECRET = "SECRET_POWER"

func IsAuth() gin.HandlerFunc {
	return checkJWT()
}

func checkJWT() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.Request.Header.Get("Authorization")
		bearerToken := strings.Split(authHeader, " ")
		
		if len(bearerToken) == 2 {

			token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
				
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(os.Getenv("JWT_SECRET")), nil
			})
			
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				fmt.Println(claims["user_id"], claims["exp"])
			} else {
				c.JSON(422, gin.H{
					"Status" : "Token Invalid",
					"Error" : err,
				})
				c.Abort()
				return
				
			}
		} else {
			c.JSON(422, gin.H{
				"Status" : "Authorization needed",
			})
			c.Abort()
			return
		}

	}
}