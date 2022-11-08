package controllers

import (
	"TODO_rest/models"
	b64 "encoding/base64"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CreateUserInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func CreateUser(c *gin.Context) {
	var input CreateUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{Username: input.Username, Password: input.Password}
	models.DB.Create(&user)

	var token = b64.StdEncoding.EncodeToString([]byte(user.Username + ":" + user.Password))
	c.IndentedJSON(http.StatusCreated, gin.H{"token": token})
}

func ChangePassword(c *gin.Context) {
	userId := c.GetUint("userID")
	newPassword := c.Query("new_password")

	var user models.User
	models.DB.Where("id = ?", userId).Scan(&user)
	user.Password = newPassword
	models.DB.Save(&user)

	c.Status(http.StatusOK)
}
