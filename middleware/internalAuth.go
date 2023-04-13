package middleware

import (
	"ae86-auth/models"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var internalAuthToken = os.Getenv("INTERNAL_AUTH_TOKEN")

func InternalBasicAuth(c *gin.Context) {
	authorization := c.GetHeader("Authorization")
	args := strings.Split(authorization, " ")
	if len(args) != 2 {
		models.ErrorLog.Printf("input token error: %s", authorization)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Authorization error": "传入token错误"})
		return
	}
	method, token := args[0], args[1]

	if method != "Basic" {
		models.ErrorLog.Printf("input token error, not basic")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Authorization error": "token格式错误_1"})
		return
	}
	if token != internalAuthToken {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Authorization error": "token解析错误"})
		return
	}

	c.Next()
}
