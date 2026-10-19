package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequiredAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Access denied,token was not found",
			})

			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("UserID", int(claims["sub"].(float64)))
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token invalid!",
			})
		}
	}
}
