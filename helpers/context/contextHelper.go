package contextHelper

import (
	"github.com/gin-gonic/gin"
)

type UserContext struct {
	ClientUuid string
	Uuid       string
	Username   string
	Name       string
	Email      string
}

func GetUserFromContext(c *gin.Context) *UserContext {
	userContext, exists := c.Get("user")
	if exists != true {
		return nil
	}
	user := userContext.(UserContext)
	return &user
}
