package routes

import (
	"TODO_rest/controllers"
	"TODO_rest/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	todo := router.Group("/todo")
	todo.Use(middleware.BasicAuth)
	todo.GET("/", controllers.GetAllTodos)
	todo.POST("/", controllers.CreateTodo)
	todo.GET("/active", controllers.GetActiveTodos)
	todo.GET("/done", controllers.GetDoneTodos)
	todo.PATCH("done/:id", controllers.CloseTodo)

	auth := router.Group("/auth")
	auth.POST("/signup", controllers.CreateUser)
	auth.POST("/change-password", middleware.BasicAuth, controllers.ChangePassword)

	return router
}
