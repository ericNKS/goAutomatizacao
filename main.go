package main

import (
	"automatizacao/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	// Function to load .env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := gin.Default()
	routes.Route(app)
	// set port
	app.Run(":8080")
}
