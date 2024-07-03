package main

import (
	"automatizacao/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()
	routes.Route(app)
	// set port
	app.Run(":8080")
}
