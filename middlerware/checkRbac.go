package middleware

import (
	"ktp/common"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckTeacherRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, _ := c.Get("username")
		has, err := common.E.Enforce(username.(string), "/teacher", "POST")
		if err != nil {
			fmt.Println("Enforce failed, err:" , err)
			c.Abort()
		}
		if !has {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "你没有权限",
				"status": 403,
			})
			c.Abort()
		}
		c.Next()
	}
}

