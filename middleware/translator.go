package middleware

import (
	"github.com/anhhuy1010/DATN-cms-customer/config"
	"github.com/anhhuy1010/DATN-cms-customer/helpers/translator"
	"github.com/gin-gonic/gin"
)

func Translator() gin.HandlerFunc {
	return func(c *gin.Context) {
		locale := c.Request.Header.Get("X-LANG")
		if locale != "" {
			if !translator.IsLocaleSupported(locale) {
				cfg := config.GetConfig()
				locale = cfg.GetString("server.locale")
			}
		} else {
			cfg := config.GetConfig()
			locale = cfg.GetString("server.locale")
		}
		c.Set("locale", locale)

		c.Next()
	}
}
