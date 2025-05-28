package middleware

import (
	"github.com/anhhuy1010/DATN-cms-customer/helpers/util"
	"github.com/gin-gonic/gin"
)

func ExtractClaims(ctx *gin.Context) (*util.Claims, bool) {
	claims, exists := ctx.Get("claims")
	if !exists {
		return nil, false
	}
	claimStruct, ok := claims.(*util.Claims)
	return claimStruct, ok
}
