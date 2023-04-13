package routes

import (
	"ae86-auth/controllers"
	"ae86-auth/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	todo := router.Group("/ae86")
	todo.Use(middleware.BasicAuth)
	todo.POST("/", controllers.RedirectRequest)
	todo.GET("/", controllers.RedirectRequest)

	ping := router.Group("/ping")
	ping.Use(middleware.InternalBasicAuth)
	ping.GET("/", controllers.Ping)
	ping.GET("/clearcache", controllers.ClearCache)
	ping.GET("/setredisratelimit", controllers.SetRedisRateLimit)
	ping.GET("/adduser", controllers.AddUser)
	ping.GET("/changepassword", controllers.ChangePassword)
	ping.GET("/changeratelimit", controllers.ChangeRateLimit)
	ping.GET("/changeexpiredate", controllers.ChangeExpireDate)

	return router
}
