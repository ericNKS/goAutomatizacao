package main

import (
	"automatizacao/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// Function that execute node
	// Param it's the directory where are your node project
	//ExecuteNodeJs("C:\\Diretorio\\do\\projeto", "8000")

	app := gin.Default()
	routes.Route(app)
	app.Run(":8080")
}
