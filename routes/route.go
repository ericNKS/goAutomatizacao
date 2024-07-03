package routes

import (
	"automatizacao/automatizacao"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Route(app *gin.Engine) {
	app.GET("/server/:port", func(c *gin.Context) {
		port := c.Param("port")
		if len(port) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid port"})
		}
		isServerOn := automatizacao.IsServerOn(port)
		c.JSON(200, isServerOn)
	})
	app.GET("/server/run/:port", func(c *gin.Context) {
		port := c.Param("port")
		if len(port) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid port"})
		}
		isServerOn := automatizacao.ExecuteNodeJs(port)
		c.JSON(200, isServerOn)
	})
}
