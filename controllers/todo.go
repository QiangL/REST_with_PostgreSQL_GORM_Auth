package controllers

import (
	"TODO_rest/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CreateTodoInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

func GetAllTodos(c *gin.Context) {
	var todos []models.Todo
	userId := c.GetUint("userID")

	//models.DB.Raw("SELECT * FROM todos WHERE user_id = ?", c.GetString("userID")).Find(&todos)
	models.DB.Where("user_id = ?", userId).Find(&todos)

	c.IndentedJSON(http.StatusOK, gin.H{"todos": todos})
}

func GetActiveTodos(c *gin.Context) {
	var todos []models.Todo
	userId := c.GetUint("userID")

	//models.DB.Raw("SELECT * FROM todos WHERE user_id = ? AND done = false", userId).Find(&todos)
	models.DB.Where("user_id = ? AND done = false", userId).Find(&todos)

	c.IndentedJSON(http.StatusOK, gin.H{"active_todos": todos})
}

func GetDoneTodos(c *gin.Context) {
	var todos []models.Todo
	userId := c.GetUint("userID")

	//models.DB.Raw("SELECT * FROM todos WHERE user_id = ? AND done = true", userId).Find(&todos)
	models.DB.Where("user_id = ? AND done = true", userId).Find(&todos)

	c.IndentedJSON(http.StatusOK, gin.H{"closed_todos": todos})
}

func CreateTodo(c *gin.Context) {
	var input CreateTodoInput
	userID := c.GetUint("userID")

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo := models.Todo{Name: input.Name, Description: input.Description, UserID: userID}
	models.DB.Create(&todo)
	c.IndentedJSON(http.StatusOK, gin.H{"created": todo})
}

func CloseTodo(c *gin.Context) {
	userId := c.GetUint("userID")
	todoId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	var todo models.Todo
	models.DB.Where("id = ?", todoId).Find(&todo)
	if todo.UserID == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Todo doesn't exist"})
		return
	}
	if todo.UserID != userId {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "You have no access to this todo"})
		return
	}

	todo.Done = true
	models.DB.Save(&todo)

	c.Status(http.StatusOK)
}
