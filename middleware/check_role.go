package middleware

import (
	"net/http"
	"time"

	"github.com/anhhuy1010/DATN-cms-customer/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Uuid     string `json:"uuid"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

func RoleMiddleware() gin.HandlerFunc {
	secretKey := config.GetConfig().GetString("server.secret_token")

	return func(c *gin.Context) {
		tokenString := c.GetHeader("x-token")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Set("customer_uuid", claims.Uuid)
		c.Set("customer_name", claims.UserName)

		c.Next()
	}
}
