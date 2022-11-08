package middleware

import (
	"TODO_rest/models"
	b64 "encoding/base64"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func BasicAuth(c *gin.Context) {
	authorization := c.GetHeader("Authorization")
	args := strings.Split(authorization, " ")
	if len(args) != 2 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Authorization error": "Invalid Authorization header"})
		return
	}
	method, token := args[0], args[1]

	if method != "Basic" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Authorization error": "Invalid auth method. Basic required"})
		return
	}

	decoded, err := b64.StdEncoding.DecodeString(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Authorization error": "Invalid token"})
		return
	}

	args = strings.Split(string(decoded), ":")
	if len(args) != 2 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Authorization error": "Invalid credentials"})
		return
	}
	username := args[0]

	var user models.User
	models.DB.Where("username = ?", username).First(&user)
	if b64.StdEncoding.EncodeToString([]byte(user.Username+":"+user.Password)) != token {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Authorization error": "Invalid credentials"})
		return
	}

	c.Set("userID", user.ID)
	c.Next()
}
