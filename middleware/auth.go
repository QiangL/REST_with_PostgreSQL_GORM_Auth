package middleware

import (
	"ae86-auth/models"
	"net/http"
	"strings"
	"time"

	"encoding/base64"

	"github.com/gin-gonic/gin"
)

func BasicAuth(c *gin.Context) {
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
	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		models.ErrorLog.Printf("input token error: decode fail, %s, %v", token, err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Authorization error": "token解析错误"})
		return
	}

	args = strings.Split(string(decoded), ":")
	if len(args) != 2 {
		models.ErrorLog.Printf("input token error, decode not two part: %s", decoded)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Authorization error": "token格式错误_2"})
		return
	}
	username := args[0]

	var user models.Ae86user
	var key = models.BuildAuthKey(username)
	err = models.MemoryStore.Get(key, &user)
	nowTime := time.Now()
	if err != nil {
		models.DB.Where("username = ?", username).First(&user)
		models.InfoLog.Printf("get user from db: %s, %s, ratelimit: %d, expiredate: %v, now: %v", user.Username, user.Password, user.RateLimit, user.ExpireDate, nowTime)
		if user.ID != 0 {
			models.MemoryStore.Set(key, user, 3*time.Hour)
		}
	}

	models.InfoLog.Printf("get user: %s, %s, ratelimit: %d, expiredate: %v, now: %v", user.Username, user.Password, user.RateLimit, user.ExpireDate, nowTime)
	var real_token = base64.StdEncoding.EncodeToString([]byte(user.Username + ":" + user.Password))
	if real_token != token {
		models.ErrorLog.Printf("input token error, invalid token: real_token: %s, input token: %s", real_token, token)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Authorization error": "token不正确"})
		return
	}

	if nowTime.Unix() > user.ExpireDate.Unix() {
		models.ErrorLog.Printf("user expired: %s, now: %v, user exporedate: %v", username, nowTime, user.ExpireDate)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Authorization error": "你的账号已过有效期"})
		return
	}

	var rateLimitKey = models.BuildRateLimitKey(username, nowTime)
	var cnt int64
	err = models.RedisStore.Get(rateLimitKey, &cnt)
	models.InfoLog.Printf("get rate from cache: %s, %d", rateLimitKey, cnt)
	if err == nil && cnt > user.RateLimit {
		models.ErrorLog.Printf("exceed rate limit: %s, %d", username, cnt)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Authorization error": "你被限流了,请调大阈值"})
		return
	}

	if err != nil {
		cnt = 0
	}

	cnt++
	models.RedisStore.Set(rateLimitKey, cnt, 25*time.Hour)

	c.Next()
}
