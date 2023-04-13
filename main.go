package main

import (
	"ae86-auth/models"
	"ae86-auth/routes"
)

func main() {
	router := routes.SetupRouter()

	models.ConnectDB()

	router.Run(":8084")
}
