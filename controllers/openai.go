package controllers

import (
	"ae86-auth/models"
	"encoding/base64"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var remote, _ = url.Parse(os.Getenv("AZURE_OPENAI_ENDPOINT"))
var proxy = httputil.NewSingleHostReverseProxy(remote)
var apiKey = os.Getenv("AZURE_OPENAI_API_KEY")

func RedirectRequest(c *gin.Context) {
	c.Request.Header.Add("api-key", apiKey)
	c.Request.Host = remote.Host
	c.Request.URL.Path = ""
	proxy.ServeHTTP(c.Writer, c.Request)
}

func Ping(c *gin.Context) {
	c.Status(http.StatusOK)
}

func ClearCache(c *gin.Context) {
	key := models.BuildAuthKey(c.Request.URL.Query()["username"][0])
	models.InfoLog.Printf("clear mem cache:" + key)
	var user models.Ae86user
	err := models.MemoryStore.Get(key, &user)
	if err == nil {
		models.ErrorLog.Printf("before clear cache: %s, %s", key, user.ExpireDate.GoString())
	} else {
		models.ErrorLog.Printf("before clear cache: %s not exist", key)
		c.Status(http.StatusOK)
		return
	}
	models.MemoryStore.Delete(models.BuildAuthKey(c.Request.URL.Query()["username"][0]))
	err = models.MemoryStore.Get(key, &user)
	if err != nil {
		models.ErrorLog.Printf("mem cache success:" + key)
	} else {
		models.ErrorLog.Printf("mem cache after clear: %s", user.ExpireDate.GoString())
	}
	c.Status(http.StatusOK)
}

func SetRedisRateLimit(c *gin.Context) {
	username := c.Request.URL.Query()["username"][0]
	key := models.BuildRateLimitKey(username, time.Now())
	models.ErrorLog.Printf("clear rate limit:" + key)
	var cnt int64
	err := models.RedisStore.Get(key, &cnt)
	if err == nil {
		models.ErrorLog.Printf("before clear rate limit: %s, %d", key, cnt)
	} else {
		models.ErrorLog.Printf("before clear rate limit: %s not exist", key)
		c.Status(http.StatusOK)
		return
	}
	cnt, err = strconv.ParseInt(c.Request.URL.Query()["ratelimit"][0], 10, 0)
	if err == nil {
		models.RedisStore.Set(models.BuildRateLimitKey(username, time.Now()), cnt, 25*time.Hour)
	} else {
		models.ErrorLog.Printf("cnt parse fail:" + err.Error())
	}

	err = models.RedisStore.Get(key, &cnt)
	if err != nil {
		models.ErrorLog.Printf("get rate limit fail:" + err.Error())
	} else {
		models.ErrorLog.Printf("new rate limit value: %d", cnt)
	}
	models.MemoryStore.Delete(models.BuildAuthKey(username))
	c.Status(http.StatusOK)
}

func AddUser(c *gin.Context) {
	models.ErrorLog.Printf("AddUser invoke")
	username := c.Request.URL.Query()["username"][0]
	password := c.Request.URL.Query()["password"][0]
	ratelimit, err := strconv.Atoi(c.Request.URL.Query()["ratelimit"][0])
	if err != nil {
		models.ErrorLog.Printf("ratelimit parse fail:" + err.Error())
		c.Status(http.StatusServiceUnavailable)
		return
	}
	expireDate, err := time.Parse("20060102", c.Request.URL.Query()["expiredate"][0])
	if err != nil {
		models.ErrorLog.Printf("expiredate parse fail:" + err.Error())
		c.Status(http.StatusServiceUnavailable)
		return
	}
	models.ErrorLog.Printf("AddUser invoke: %s, %s, %d, %s", username, password, ratelimit, expireDate.Format("20060102"))

	user := models.Ae86user{
		Username:   username,
		Password:   password,
		RateLimit:  int64(ratelimit),
		ExpireDate: expireDate,
	}

	models.DB.Create(&user)
	var real_token = base64.StdEncoding.EncodeToString([]byte(user.Username + ":" + user.Password))
	c.JSON(http.StatusOK, gin.H{"token": real_token})
}

func ChangeRateLimit(c *gin.Context) {
	models.ErrorLog.Printf("ChangeRateLimit invoke")
	username := c.Request.URL.Query()["username"][0]
	ratelimit, err := strconv.Atoi(c.Request.URL.Query()["ratelimit"][0])
	models.ErrorLog.Printf("ChangeRateLimit invoke: %s, %d", username, ratelimit)
	if err != nil {
		models.ErrorLog.Printf("ratelimit parse fail:" + err.Error())
		c.Status(http.StatusServiceUnavailable)
		return
	}

	user := models.Ae86user{
		Username:  username,
		RateLimit: int64(ratelimit),
	}

	models.DB.Model(&user).Where("username = ?", username).Update(&user)
	models.MemoryStore.Delete(models.BuildAuthKey(username))
	c.Status(http.StatusOK)
}

func ChangePassword(c *gin.Context) {
	models.ErrorLog.Printf("ChangePassword invoke")
	username := c.Request.URL.Query()["username"][0]
	password := c.Request.URL.Query()["password"][0]

	models.ErrorLog.Printf("ChangePassword invoke: %s, %s", username, password)

	user := models.Ae86user{
		Username: username,
		Password: password,
	}

	models.DB.Model(&user).Where("username = ?", username).Update(&user)
	models.MemoryStore.Delete(models.BuildAuthKey(username))
	var real_token = base64.StdEncoding.EncodeToString([]byte(user.Username + ":" + user.Password))
	c.JSON(http.StatusOK, gin.H{"token": real_token})
}

func ChangeExpireDate(c *gin.Context) {
	models.ErrorLog.Printf("ChangeExpireDate invoke")
	username := c.Request.URL.Query()["username"][0]
	expireDate, err := time.Parse("20060102", c.Request.URL.Query()["expiredate"][0])
	if err != nil {
		models.ErrorLog.Printf("expiredate parse fail:" + err.Error())
		c.Status(http.StatusServiceUnavailable)
		return
	}

	models.ErrorLog.Printf("ChangePassword invoke: %s, %s", username, expireDate.Format("20060102"))

	user := models.Ae86user{
		Username:   username,
		ExpireDate: expireDate,
	}

	models.DB.Model(&user).Where("username = ?", username).Update(&user)
	models.MemoryStore.Delete(models.BuildAuthKey(username))
	c.Status(http.StatusOK)
}
