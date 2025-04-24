package middleware

import (
	"github.com/gin-gonic/gin"
)

func ExtractClaims(ctx *gin.Context) (map[string]interface{}, bool) {
	claims, exists := ctx.Get("claims")
	if !exists {
		return nil, false
	}
	claimMap, ok := claims.(map[string]interface{})
	return claimMap, ok
}
