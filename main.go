package main

import (
	"automatizacao/handler"
	"automatizacao/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()
	// Set routes
	routes.Route(app)
	go handler.VerifyAndRun()
	// Init endpoints
	app.Run(":8080")
}
