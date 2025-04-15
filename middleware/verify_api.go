package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/anhhuy1010/DATN-cms-customer/config"
	"github.com/anhhuy1010/DATN-cms-customer/helpers/respond"
)

func VerifyApi() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := config.GetConfig()

		apiKey := c.Request.Header.Get("X-API-KEY")
		secretKey := cfg.GetString("server.secret_key")

		if apiKey != secretKey {
			c.JSON(http.StatusForbidden, respond.Forbidden())
			c.Abort()
			return
		}

		c.Next()
	}
}
