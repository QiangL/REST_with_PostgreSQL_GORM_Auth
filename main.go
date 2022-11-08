package main

import (
	"TODO_rest/models"
	"TODO_rest/routes"
)

func main() {
	router := routes.SetupRouter()

	models.ConnectDB()

	router.Run(":8080")
}
